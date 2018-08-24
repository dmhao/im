package core

import (
	"im/mp/common"
	"im/mp"
	"github.com/golang/protobuf/proto"
	"fmt"
)

const (
	Ping    	= 1
	Pong    	= 2
	Login   	= 3
	Logout  	= 4
	Auth    	= 5

	Msg     	= 6
	MsgAck     	= 7

	GroupMsg 	= 8
	GroupMsgAck = 9

	SyncOffline = 10
	SyncOfflineMsg = 11
	SyncOfflineMsgAck = 12
)

const (
	RouteLogin = 101
	RouteLogout = 102
	TransportMsg = 103
	TransportGroupMsg = 104
)


type Message struct {
	Seq				int
	Cmd				int8
	Version			int8
	Data			proto.Message
}


type RouteMessage struct {
	Seq				int
	Cmd				int8
	Version			int
	RouteData 		proto.Message
}


var cmdMsgType = make(map[int8]func() proto.Message)
var cmdRouteMsgType = make(map[int8]func() proto.Message)


func init() {
	cmdMsgType[Ping] = func() proto.Message { return new(common.Empty) }
	cmdMsgType[Pong] = func() proto.Message { return new(common.Empty) }

	cmdMsgType[Login] = func() proto.Message { return new(common.LoginIm) }
	cmdMsgType[Logout] = func() proto.Message { return new(common.LogoutIm)}
	cmdMsgType[Auth] = func() proto.Message { return new(common.AuthIm) }

	cmdMsgType[Msg] = func() proto.Message { return new(common.MsgIm) }
	cmdMsgType[MsgAck] = func() proto.Message { return new(common.MsgImAck) }

	cmdMsgType[GroupMsg] = func() proto.Message { return new(common.MsgIm) }
	cmdMsgType[GroupMsgAck] = func() proto.Message { return new(common.MsgImAck) }

	cmdMsgType[SyncOffline] = func() proto.Message { return new(common.SyncOffline) }
	cmdMsgType[SyncOfflineMsg] = func() proto.Message { return new(common.SyncOfflineMsg) }
	cmdMsgType[SyncOfflineMsgAck] = func() proto.Message { return new(common.MsgImAck) }


	cmdRouteMsgType[RouteLogout] = func() proto.Message { return new(mp.RouteLogoutIm) }
	cmdRouteMsgType[RouteLogin] = func() proto.Message { return new(mp.RouteLoginIm) }

	cmdRouteMsgType[TransportMsg] = func() proto.Message { return new(mp.RouteMsgIm) }
	cmdRouteMsgType[TransportGroupMsg] = func() proto.Message { return new(mp.RouteMsgIm) }
}


//pb序列化消息中data
func (msg *Message) formatData(dataBytes []byte) {
	cmd := msg.Cmd
	structFunc, ok := cmdMsgType[cmd]
	if !ok {
		msg.Data = nil
		return
	}
	data := structFunc()
	err := proto.Unmarshal(dataBytes, data)
	if err != nil {
		fmt.Println( "msg包体解码失败", err)
	}
	msg.Data = data
}


func (msg *RouteMessage) formatData(dataBytes []byte) {
	cmd := msg.Cmd
	structFunc, ok := cmdRouteMsgType[cmd]
	if !ok {
		msg.RouteData = nil
		return
	}
	routeData := structFunc()
	err := proto.Unmarshal(dataBytes, routeData)
	if err != nil {
		fmt.Println( "路由包体解码失败", err)
	}
	msg.RouteData = routeData
}


