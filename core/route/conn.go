package route

import (
	"context"
	"im/core"
	"im/core/log"
	"im/core/tools"
	"net"
	"sync"
	"time"
	"github.com/golang/protobuf/proto"
)

type imServers struct {
	rw   *sync.RWMutex
	data map[*imServerConn]*apps
}

const writeDeadline = 60 * time.Second

type imServerConn struct {
	rs      	*routeServer
	tcpConn 	*net.TCPConn
	allApp  	*apps

	writeCh   	chan *core.RouteMessage
	readPacket	*core.PacketData
	writePacket []byte

	once   		*sync.Once
	wg     		*sync.WaitGroup
	ctx    		context.Context
	cancel 		context.CancelFunc
}

func NewImServerConn(rs *routeServer, conn *net.TCPConn) *imServerConn {
	csc := &imServerConn{
		rs:      rs,
		tcpConn: conn,
		allApp: &apps{
			rw:   &sync.RWMutex{},
			data: make(map[int]*appInfo),
		},
		writeCh:   make(chan *core.RouteMessage, 100),
		readPacket:	&core.PacketData{
			PacketH: make([]byte, core.HeaderLen),
			PacketB: make([]byte, core.BaseBodyLen),
		},
		writePacket: make([]byte, core.HeaderLen + core.BaseBodyLen),
		once:      &sync.Once{},
		wg:        &sync.WaitGroup{},
	}

	csc.ctx, csc.cancel = context.WithCancel(rs.ctx)

	return csc
}

//捕获协程中的程序崩溃
func recoverFunc() {
	err := recover()
	if err != nil {
		stack := tools.Stack(3)
		log.Warnln("协程崩溃", err, string(stack))
	}
}

func (isc *imServerConn) readLoop() {
	defer func() {
		isc.Close()
		recoverFunc()
	}()

	for {
		select {
		case <-isc.ctx.Done():
			return
		case <-isc.rs.ctx.Done():
			return
		default:
			msg, err := core.ReadRoutePacket(isc.tcpConn, isc.readPacket)
			if err != nil ||  msg == nil {
				return
			}
			isc.dealMessage(msg)
		}
	}
}

func (isc *imServerConn) sendMessage(msg *core.RouteMessage) bool {
	if msg == nil {
		return false
	}
	t := time.NewTimer(writeDeadline)
	select {
	case isc.writeCh <- msg:
		t.Stop()
		return true
	case <-t.C:
		log.Warnln("发送消息超时", isc.tcpConn.RemoteAddr())
		return false
	}
}

func (isc *imServerConn) writeLoop() {
	defer func() {
		isc.Close()
		recoverFunc()
		isc.wg.Done()
	}()
	for {
		select {
		case <-isc.ctx.Done():
			return
		case <-isc.rs.ctx.Done():
			return
		case msg := <-isc.writeCh:
			if msg == nil {
				return
			}
			dataBytes, err := proto.Marshal(msg.RouteData)
			if err != nil {
				log.Infoln("消息msg编码失败", err, msg)
				continue
			}
			bodyLen := len(dataBytes)
			if bodyLen > core.BaseBodyLen {
				isc.writePacket = make([]byte, bodyLen)
			}
			core.SetPacketHeader(byte(msg.Cmd), 1, 1, bodyLen, isc.writePacket)
			copy(isc.writePacket[core.HeaderLen:], dataBytes)
			_, err = isc.tcpConn.Write(isc.writePacket[:core.HeaderLen+bodyLen])
			if err != nil {
				log.Warnln("会话写入错误", err, msg)
			}
		}
	}
}

func (isc *imServerConn) handlerConn() {
	log.Infoln("im服务", isc.tcpConn.RemoteAddr(), "已连接")

	//启动写消息协程
	isc.wg.Add(1)
	go isc.writeLoop()

	//消息读取
	isc.readLoop()

	//等待其他 写协程  消息处理协程关闭
	isc.wg.Wait()

	isc.rs.removeImServerConn(isc)
	log.Infoln("im服务", isc.tcpConn.RemoteAddr(), "断开连接")

	close(isc.writeCh)
	isc.tcpConn.Close()
	isc.rs.wg.Done()
}

func (isc *imServerConn) Close() {
	isc.once.Do(func() {
		isc.tcpConn.CloseRead()
		isc.tcpConn.Close()
		isc.cancel()
	})
}
