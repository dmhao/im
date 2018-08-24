package server

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

type router struct {
	rw   *sync.RWMutex
	data map[string]*routeConn
}

type routeConn struct {
	addr    	string
	im      	*imServer
	tcpConn 	*net.TCPConn

	writeCh 	chan *core.RouteMessage
	readPacket	*core.PacketData
	writePacket []byte

	once   		*sync.Once
	wg     		*sync.WaitGroup
	ctx    		context.Context
	cancel 		context.CancelFunc
}

func newRouteConn(im *imServer, conn *net.TCPConn, addr string) *routeConn {
	rConn := &routeConn{
		addr:      	addr,
		im:        	im,
		tcpConn:   	conn,
		writeCh:   	make(chan *core.RouteMessage, 50),
		readPacket:	&core.PacketData{
			PacketH: make([]byte, core.HeaderLen),
			PacketB: make([]byte, core.BaseBodyLen),
		},
		writePacket: make([]byte, core.HeaderLen + core.BaseBodyLen),
		once:      	&sync.Once{},
		wg:        	&sync.WaitGroup{},
	}
	rConn.ctx, rConn.cancel = context.WithCancel(im.ctx)
	return rConn
}

//循环读取路由服务下发的消息
func (rc *routeConn) readLoop() {
	defer func() {
		rc.Close()
		recoverFunc()
	}()
	for {
		select {
		case <-rc.ctx.Done():
			return
		case <-rc.im.ctx.Done():
			return
		default:
			msg, err := core.ReadRoutePacket(rc.tcpConn, rc.readPacket)
			if err != nil || msg == nil {
				return
			}
			rc.dealMessage(msg)
		}
	}
}

func (rc *routeConn) sendMessage(msg *core.RouteMessage) bool {
	if msg == nil {
		return false
	}
	t := time.NewTimer(writeDeadline)
	select {
	case rc.writeCh <- msg:
		t.Stop()
		return true
	case <-t.C:
		log.Warnln("写入数据超时", tools.JsonMarshal(msg))
		return false
	}
}

func (rc *routeConn) writeLoop() {
	defer func() {
		rc.Close()
		recoverFunc()
		rc.wg.Done()
	}()

	for {
		select {
		case <-rc.ctx.Done():
			return
		case <-rc.im.ctx.Done():
			return
		case msg := <-rc.writeCh:
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
				rc.writePacket = make([]byte, bodyLen)
			}
			core.SetPacketHeader(byte(msg.Cmd), 1, 1, bodyLen, rc.writePacket)
			copy(rc.writePacket[core.HeaderLen:], dataBytes)
			_, err = rc.tcpConn.Write(rc.writePacket[:core.HeaderLen+bodyLen])
			if err != nil {
				log.Warnln("路由会话写入", err)
				return
			}
		}

	}
}

func (rc *routeConn) handlerConn() {
	//当前已建立的所有会话同步到新路由服务
	rc.wg.Add(1)
	go func() {
		allUserConn := rc.im.findAllUserConn()
		for _, userConn := range allUserConn {
			sendOneRouteLoginMsg(userConn, rc.addr)
		}
		rc.wg.Done()
	}()

	//定时ping路由服务
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		for {
			t := time.NewTimer(writeDeadline)
			select {
			case <-rc.ctx.Done():
				t.Stop()
				return
			case <-t.C:
				dateBytes := make([]byte, core.HeaderLen)
				core.SetPacketHeader(byte(core.Ping), 1, 1, 0, dateBytes)
				_, err := rc.tcpConn.Write(dateBytes)
				if err != nil {
					return
				}
			}
		}
	}()
	//启动写消息协程
	rc.wg.Add(1)
	go rc.writeLoop()

	//消息读取
	rc.readLoop()

	//等待其他 写协程  消息处理协程关闭
	rc.wg.Wait()

	//删除此会话数据
	rc.im.rt.rw.Lock()
	delete(rc.im.rt.data, rc.addr)
	rc.im.rt.rw.Unlock()

	close(rc.writeCh)
	rc.tcpConn.Close()
	rc.im.wg.Done()
}

func (rc *routeConn) Close() {
	rc.once.Do(func() {
		rc.tcpConn.CloseRead()
		rc.tcpConn.Close()
		rc.cancel()
	})
}
