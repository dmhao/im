package route

import (
	"im/core"
	"im/core/log"
	"im/core/storage"
	"im/core/tools"
	"im/mp"
)

//选择合适的handler处理收到的消息
func (isc *imServerConn) dealMessage(msg *core.RouteMessage) {
	switch msg.Cmd {
	case core.RouteLogin:
		isc.hRouteLogin(msg)
	case core.RouteLogout:
		isc.hRouteLogout(msg)
	case core.TransportMsg:
		isc.hTransportMsg(msg)
	case core.TransportGroupMsg:
		isc.hTransportGroupMsg(msg)
	}
}

//路由服务器  注册登录im的用户
func (isc *imServerConn) hRouteLogin(msg *core.RouteMessage) {
	routeData := msg.RouteData.(*mp.RouteLoginIm)
	appId := int(routeData.GetAppId())
	userId := routeData.GetUserId()

	//查找appId对应的app信息
	ai := isc.allApp.findORCreateAppInfo(appId)
	//在app中添加用户userId
	ai.insertUserId(userId)

	//获取用户加入的群组, 并把用户参加的群组中加入用户的连接
	groupIds := storage.GetAllUserJoinGroupId(appId, userId, storage.ShowGroupUser)
	if len(groupIds) > 0 {
		for _, groupId := range groupIds {
			gu := ai.findOrCreateGroupUsers(groupId)
			gu.insertUserId(userId)
		}
	}
	log.Infoln("用户登录", appId, userId)
}

//路由服务器  注销登录im的用户
func (isc *imServerConn) hRouteLogout(msg *core.RouteMessage) {
	routeData := msg.RouteData.(*mp.RouteLogoutIm)
	appId := int(routeData.GetAppId())
	userId := routeData.GetUserId()
	ai := isc.allApp.findORCreateAppInfo(appId)
	ai.removeUserId(userId)
	groupIds := storage.GetAllUserJoinGroupId(appId, userId, storage.ShowGroupUser)
	if len(groupIds) > 0 {
		for _, groupId := range groupIds {
			gu := ai.findOrCreateGroupUsers(groupId)
			gu.removeUserId(userId)
		}
	}
	log.Infoln("用户退出", appId, userId)
}

//路由服务器 转发单聊消息
func (isc *imServerConn) hTransportMsg(msg *core.RouteMessage) {
	routeData := msg.RouteData.(*mp.RouteMsgIm)
	appId := int(routeData.GetAppId())
	receiverUserId := routeData.GetReceiverId()

	var iscs []*imServerConn
	isc.rs.iss.rw.RLock()
	for tmpIsc, _ := range isc.rs.iss.data {
		if tmpIsc.allApp.appAndUserIdExists(appId, receiverUserId) {
			iscs = append(iscs, tmpIsc)
		}
	}
	isc.rs.iss.rw.RUnlock()
	if len(iscs) > 0 {
		for _, is := range iscs {
			if is.allApp.hasAppId(appId) != nil {
				is.sendMessage(msg)
			}
		}
	}
	log.Infoln("单聊转发", tools.JsonMarshal(routeData))
}

//路由服务器 转发群聊的消息
func (isc *imServerConn) hTransportGroupMsg(msg *core.RouteMessage) {
	routeData := msg.RouteData.(*mp.RouteMsgIm)
	appId := int(routeData.GetAppId())
	cscs := isc.rs.getAllImServerConn()
	for _, cs := range cscs {
		if cs.allApp.hasAppId(appId) != nil {
			cs.sendMessage(msg)
			log.Infoln("群聊转发至im服务", cs.tcpConn.RemoteAddr().String(), tools.JsonMarshal(routeData))
		}
	}
	log.Infoln("群聊转发", tools.JsonMarshal(routeData))
}
