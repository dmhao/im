package core

import (
	"github.com/golang/protobuf/proto"
	"im/core/log"
	"im/mp"
	"im/mp/common"
)

const (
	Ping   = 1
	Pong   = 2
	Login  = 3
	Logout = 4
	Auth   = 5

	Msg    = 6
	MsgAck = 7

	GroupMsg    = 8
	GroupMsgAck = 9

	SyncOffline       = 10
	SyncOfflineMsg    = 11
	SyncOfflineMsgAck = 12
)

const (
	RouteLogin        = 101
	RouteLogout       = 102
	TransportMsg      = 103
	TransportGroupMsg = 104
)

const (
	SystemMessageType = 2
	UserMessageType   = 1
	ImChartType = 1
	GroupImChartType = 2
)

type Message struct {
	Seq     int
	Cmd     int8
	Version int8
	Data    proto.Message
}

type RouteMessage struct {
	Seq       int
	Cmd       int8
	Version   int
	RouteData proto.Message
}

type ApiMessage struct {
	AppId int
	Data  *common.MsgIm
}

var cmdMsgType = make(map[int8]func() proto.Message)
var cmdRouteMsgType = make(map[int8]func() proto.Message)

func init() {
	cmdMsgType[Ping] = func() proto.Message { return &common.Empty{} }
	cmdMsgType[Pong] = func() proto.Message { return &common.Empty{} }

	cmdMsgType[Login] = func() proto.Message { return &common.LoginIm{} }
	cmdMsgType[Logout] = func() proto.Message { return &common.LogoutIm{} }
	cmdMsgType[Auth] = func() proto.Message { return &common.AuthIm{} }

	cmdMsgType[Msg] = func() proto.Message { return &common.MsgIm{} }
	cmdMsgType[MsgAck] = func() proto.Message { return &common.MsgImAck{} }

	cmdMsgType[GroupMsg] = func() proto.Message { return &common.MsgIm{} }
	cmdMsgType[GroupMsgAck] = func() proto.Message { return &common.MsgImAck{} }

	cmdMsgType[SyncOffline] = func() proto.Message { return &common.SyncOffline{} }
	cmdMsgType[SyncOfflineMsg] = func() proto.Message { return &common.SyncOfflineMsg{} }
	cmdMsgType[SyncOfflineMsgAck] = func() proto.Message { return &common.MsgImAck{} }

	cmdRouteMsgType[RouteLogout] = func() proto.Message { return &mp.RouteLogoutIm{} }
	cmdRouteMsgType[RouteLogin] = func() proto.Message { return &mp.RouteLoginIm{} }

	cmdRouteMsgType[TransportMsg] = func() proto.Message { return &mp.RouteMsgIm{} }
	cmdRouteMsgType[TransportGroupMsg] = func() proto.Message { return &mp.RouteMsgIm{} }
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
		log.Warnln("msg包体解码失败", err)
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
		log.Warnln("路由包体解码失败", err)
	}
	msg.RouteData = routeData
}
