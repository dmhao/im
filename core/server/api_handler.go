package server

import (
	"im/core"
)

//处理api发送的消息，写入路由服务推送消息
func (im *imServer) handlerApiMessages(apiMessages []*core.ApiMessage) {
	for _, apiMessage := range apiMessages {
		appId := apiMessage.AppId
		if bindImAndSave(appId, apiMessage.Data) {
			if apiMessage.Data.ChartType == core.ImChartType {
				im.sendRouteImMsg(appId, core.TransportMsg, apiMessage.Data)
			} else if apiMessage.Data.ChartType == core.GroupImChartType {
				im.sendRouteImMsg(appId, core.TransportGroupMsg, apiMessage.Data)
			}
		}
	}
}
