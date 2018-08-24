package storage

import (
	"im/mp/common"
	"strconv"
)

type Message struct {
	Id         int64
	SenderId   int64
	ReceiverId int64
	ChartType  int32
	MsgType    int32
	MsgId      string
	TalkId     string
	TraceId    string
	Timestamp  int64
	Content    string
	AppId      int
}

func AddMessage(appId int, msgIm *common.MsgIm) bool {
	msg := &Message{
		SenderId: msgIm.SenderId,
		ReceiverId: msgIm.ReceiverId,
		ChartType: msgIm.ChartType,
		MsgType: msgIm.MsgType,
		MsgId: msgIm.MsgId,
		TalkId: msgIm.TalkId,
		TraceId: msgIm.TraceId,
		Timestamp: msgIm.Timestamp,
		Content: msgIm.Content,
		AppId: appId,
	}

	dbClient.Model(msg).Create(msg)
	if msg.Id != 0 {
		return true
	} else {
		return false
	}
}

func MakeMsgId(msgIm *common.MsgIm) {
	msgIm.MsgId = msgIm.TalkId + "-" + strconv.FormatInt(msgIm.Timestamp, 10)
}
