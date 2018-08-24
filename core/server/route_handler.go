package server

import (
	"im/core"
	"im/core/log"
	"im/core/storage"
	"im/mp"
)

//选择合适的handler处理收到的消息
func (rc *routeConn) dealMessage(msg *core.RouteMessage) {
	switch msg.Cmd {
	case core.TransportMsg:
		rc.hTransportMsg(msg)
	case core.TransportGroupMsg:
		rc.hTransportGroupMsg(msg)
	}
}

func (rc *routeConn) hTransportMsg(msg *core.RouteMessage) {
	routeMsg := msg.RouteData.(*mp.RouteMsgIm)
	userId := routeMsg.GetReceiverId()
	appId := int(routeMsg.GetAppId())
	appInfo := rc.im.allApp.findORCreateAppInfo(appId)

	userConn, ok := appInfo.users[userId]
	if !ok {
		log.Warnln("单聊转发,没有找到用户连接appId",
			routeMsg.AppId, "userId", routeMsg.ReceiverId)
		return
	}
	pushMsg := &core.Message{
		Cmd: core.Msg,
		Data: routeMsg.TransportData,
	}
	userConn.sendMessage(pushMsg)
	log.Infoln("单聊推送完成", msg)
}

func (rc *routeConn) hTransportGroupMsg(msg *core.RouteMessage) {
	routeMsg := msg.RouteData.(*mp.RouteMsgIm)
	appId := int(routeMsg.GetAppId())
	groupId := int(routeMsg.ReceiverId)

	ai := rc.im.allApp.findORCreateAppInfo(appId)

	//读取所有的群成员
	groupUsers := storage.GetAllGroupUser(appId, groupId, storage.ShowGroupUser)

	groupUserConnMap := ai.GetGroupUsersConn(groupUsers)
	delete(groupUserConnMap, routeMsg.GetSenderId())
	if groupUserConnMap != nil {
		//当前服务器有此群的在线用户连接时 ， 推送群消息
		if len(groupUserConnMap) > 0 {
			pushMsg := &core.Message{
				Cmd: core.GroupMsg,
				Data: routeMsg.TransportData,
			}
			loopGroupUserSendMessage(groupUserConnMap, pushMsg)
		}
	}
	log.Infoln("群聊推送完成", msg)
}

func loopGroupUserSendMessage(gUserMap map[int64]*clientConn, msg *core.Message) {
	//当前服务有此群的在线用户 ， 循环推送群消息
	if len(gUserMap) > 0 {
		for _, userConn := range gUserMap {
			if userConn != nil {
				userConn.sendMessage(msg)
			}
		}
	}
}
