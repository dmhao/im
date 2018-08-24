package server

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"im/core"
	"im/core/log"
	"im/core/tools"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"github.com/golang/protobuf/proto"
)

const readDeadline = 30 * time.Second
const writeDeadline = 60 * time.Second

type clientConn struct {
	im      	*imServer
	tcpConn 	*net.TCPConn

	deviceType 	int
	appId      	int
	userId     	int64
	token      	*jwt.Token

	writeCh    	chan *core.Message
	readPacket	*core.PacketData
	writePacket []byte

	once   		*sync.Once
	wg     		*sync.WaitGroup
	ctx    		context.Context
	cancel 		context.CancelFunc
}

func newClientConn(im *imServer, tcpConn *net.TCPConn) *clientConn {
	cc := &clientConn{
		im:        	im,
		tcpConn:   	tcpConn,
		writeCh:   	make(chan *core.Message, 10),
		readPacket:	&core.PacketData{
			PacketH: make([]byte, core.HeaderLen),
			PacketB: make([]byte, core.BaseBodyLen),
		},
		writePacket: make([]byte, core.HeaderLen + core.BaseBodyLen),
		once:      	&sync.Once{},
		wg:        	&sync.WaitGroup{},
	}
	cc.ctx, cc.cancel = context.WithCancel(im.ctx)
	return cc
}

//捕获协程中的程序崩溃
func recoverFunc() {
	err := recover()
	if err != nil {
		stack := tools.Stack(3)
		log.Warnln("协程崩溃", err, string(stack))
	}
}

//循环读取客户端发送的消息
func (cc *clientConn) readLoop() {
	defer func() {
		cc.Close()
		recoverFunc()
	}()
	for {
		select {
		case <-cc.ctx.Done():
			return
		case <-cc.im.ctx.Done():
			return
		default:
			//读取连接中的packet数据
			cc.tcpConn.SetReadDeadline(time.Now().Add(readDeadline))
			msg, err := core.ReadPacket(cc.tcpConn, cc.readPacket)
			if err != nil || msg == nil {
				return
			}
			cc.dealMessage(msg)
		}
	}
}


//循环推送到客户端消息
func (cc *clientConn) writeLoop() {
	defer func() {
		cc.Close()
		recoverFunc()
		cc.wg.Done()
	}()
	for {
		select {
		case <-cc.ctx.Done():
			return
		case <-cc.im.ctx.Done():
			return
		case msg := <-cc.writeCh:
			if msg == nil {
				return
			}
			dataBytes, err := proto.Marshal(msg.Data)
			if err != nil {
				log.Infoln("消息msg编码失败", err, msg)
				continue
			}
			bodyLen := len(dataBytes)
			if bodyLen > cap(cc.writePacket) {
				cc.writePacket = make([]byte, bodyLen)
			}
			core.SetPacketHeader(byte(msg.Cmd), 1, 1, bodyLen, cc.writePacket)
			copy(cc.writePacket[core.HeaderLen:], dataBytes)
			_, err = cc.tcpConn.Write(cc.writePacket[:core.HeaderLen+bodyLen])
			if err != nil {
				log.Warnln("会话写入错误", err)
				return
			}
		}
	}
}

//推送消息到客户端
func (cc *clientConn) sendMessage(msg *core.Message) bool {
	if msg == nil {
		return false
	}
	t := time.NewTimer(writeDeadline)
	select {
	case cc.writeCh <- msg:
		t.Stop()
		return true
	case <-t.C:
		log.Warnln("会话写入超时", cc.userId)
		return false
	}
}

func (cc *clientConn) handlerConn() {
	atomic.AddInt32(cc.im.connCount, 1)
	//启动写消息协程
	cc.wg.Add(1)
	go cc.writeLoop()

	//消息读取
	cc.readLoop()

	//等待其他 写协程  消息处理协程关闭
	cc.wg.Wait()
	atomic.AddInt32(cc.im.connCount, -1)

	close(cc.writeCh)
	cc.tcpConn.Close()
	cc.im.wg.Done()
}

func (cc *clientConn) Close() {
	cc.once.Do(func() {
		cc.cancel()
		cc.hLogout()
	})
}
