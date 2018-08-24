package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"im/core"
	"im/core/server"
	"im/mp/common"
	"time"
)

func CreateMessage(c *gin.Context) {
	appId := c.GetInt("appId")
	msgIm := &common.MsgIm{}
	c.Bind(msgIm)
	uuidBytes, _ := uuid.NewV4()
	msgIm.TraceId = "im-api-" + uuidBytes.String()

	apiMsg := &core.ApiMessage{
		AppId: appId,
		Data: msgIm,
	}
	var msgData []*core.ApiMessage
	msgData = append(msgData, apiMsg)
	t := time.NewTimer(20 * time.Second)
	select {
	case server.ApiMessagesCh <- msgData:
		t.Stop()
		mContext{c}.SuccessResponse(nil)
	case <-t.C:
		mContext{c}.ErrorResponse(ServiceTimeOut, ServiceTimeOutMsg)
	}
}
