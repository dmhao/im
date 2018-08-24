package server

import (
	"im/core"
	"im/core/auth"
	"im/core/log"
	"im/core/storage"
	"im/core/tools"
	"im/mp"
	"im/mp/common"
	"time"
)

//选择合适的handler处理收到的消息
func (cc *clientConn) dealMessage(msg *core.Message) {
	//验证token是否合法，token过期或非法都将退出当前用户
	if msg.Cmd != core.Login && msg.Cmd != core.Ping {
		valid := auth.CheckTokenRightful(cc.appId, cc.userId, cc.token)
		if !valid {
			cc.hLogout()
		}
	}
	switch msg.Cmd {
	case core.Ping:
		cc.hPing(msg)
	case core.Login:
		cc.hLogin(msg)
	case core.Logout:
		cc.hLogout()
	case core.Msg:
		cc.hIm(msg)
	case core.MsgAck:
		cc.hMsgAck(msg)
	case core.GroupMsg:
		cc.hGroupIm(msg)
	case core.GroupMsgAck:
		cc.hGroupMsgAck(msg)
	case core.SyncOffline:
		cc.hSyncOffline(msg)
	case core.SyncOfflineMsgAck:
		cc.hSyncOfflineAck(msg)
	}
}

//推送用户登录消息到所有路由服务，这样所有的路由服务都可以找到用户所在的im服务
func sendRouteLoginMsg(cc *clientConn) {
	registerIm := &mp.RouteLoginIm{
		UserId: cc.userId,
		AppId: int32(cc.appId),
	}

	routeMsg := &core.RouteMessage{
		Cmd: core.RouteLogin,
		RouteData: registerIm,
	}
	cc.im.pushRouteMessage(cc.appId, routeMsg, true)
}

//推送消息到某一台路由服务，由路由服务转发消息
func sendOneRouteLoginMsg(cc *clientConn, routeAddr string) {
	loginIm := &mp.RouteLoginIm{
		UserId: cc.userId,
		AppId: int32(cc.appId),
	}

	routeMsg := &core.RouteMessage{
		Cmd: core.RouteLogin,
		RouteData: loginIm,
	}
	cc.im.pushOneRouteMessage(cc.appId, routeMsg, routeAddr)
}

//推送用户退出消息到所有路由服务，所有的路由服务都应该注销此用户登录状态
func sendRouteLogoutMsg(cc *clientConn) {
	logoutIm := &mp.RouteLogoutIm{
		UserId: cc.userId,
		AppId: int32(cc.appId),
	}

	routeMsg := &core.RouteMessage{
		Cmd: core.RouteLogout,
		RouteData: logoutIm,
	}
	cc.im.pushRouteMessage(cc.appId, routeMsg, true)
}

//绑定数据并保存消息
func bindImAndSave(appId int, msgIm *common.MsgIm) bool {
	msgIm.Timestamp = time.Now().UnixNano()
	storage.UpdateOrCreateTalkId(appId, msgIm)
	storage.MakeMsgId(msgIm)
	if storage.AddMessage(appId, msgIm) {
		//更新用户的会话列表
		storage.UpdateUserTalkIdTime(appId, msgIm)
		//保存离线消息
		storage.AppendUserImMsg(appId, msgIm.ReceiverId, msgIm)
		log.Infoln("消息已存储", tools.JsonMarshal(msgIm))
		return true
	}
	return false
}

//绑定群组数据并保存消息
func bindGroupImAndSave(appId int, msgIm *common.MsgIm) bool {
	msgIm.Timestamp = time.Now().UnixNano()
	storage.UpdateOrCreateTalkId(appId, msgIm)
	storage.MakeMsgId(msgIm)
	if storage.AddMessage(appId, msgIm) {
		//更新用户的会话列表
		storage.UpdateUserTalkIdTime(appId, msgIm)
		log.Infoln("消息已存储", tools.JsonMarshal(msgIm))
		return true
	}
	return false
}

//收到客户端单聊消息，返回给客户端的ack
func rspImAck(cc *clientConn, msgIm *common.MsgIm) {
	ackMsg := &common.MsgImAck{
		MsgId: msgIm.MsgId,
		TraceId: msgIm.TraceId,
		TalkId: msgIm.TalkId,
		Timestamp: msgIm.Timestamp,
	}

	msg := &core.Message{
		Cmd: core.MsgAck,
		Data: ackMsg,
	}
	cc.sendMessage(msg)
}

//收到客户端群聊消息，返回给客户端的ack
func rspGroupImAck(cc *clientConn, msgIm *common.MsgIm) {
	ackMsg := &common.MsgImAck{
		MsgId: msgIm.MsgId,
		TraceId: msgIm.TraceId,
		TalkId: msgIm.TalkId,
		Timestamp: msgIm.Timestamp,
	}

	msg := &core.Message{
		Cmd: core.GroupMsgAck,
		Data: ackMsg,
	}
	cc.sendMessage(msg)
}

//处理收到的ping请求
func (cc *clientConn) hPing(msg *core.Message) {
	msg.Cmd = core.Pong
	msg.Data = &common.Empty{}
	cc.sendMessage(msg)
}

//处理当前会话收到的单聊消息
func (cc *clientConn) hIm(msg *core.Message) {
	if msg.Data == nil {
		return
	}
	msgIm := msg.Data.(*common.MsgIm)
	if msgIm.SenderId == 0 || msgIm.ReceiverId == 0 || msgIm.Content == "" {
		log.Warnln("消息不完整", tools.JsonMarshal(msgIm))
		return
	}
	if bindImAndSave(cc.appId, msgIm) {
		cc.im.sendRouteImMsg(cc.appId, core.TransportMsg, msgIm)
		//返回发送者的ack
		rspImAck(cc, msgIm)
	}
}

//处理推送单聊消息后，客户端返回的ack
func (cc *clientConn) hMsgAck(msg *core.Message) {
	ackMck := msg.Data.(*common.MsgImAck)
	msgId := ackMck.MsgId
	storage.RemoveUserImMsg(cc.appId, cc.userId, msgId)
}

//处理推送群聊消息后，客户端返回的ack
func (cc *clientConn) hGroupMsgAck(msg *core.Message) {
	ackMck := msg.Data.(*common.MsgImAck)
	msgId := ackMck.MsgId
	storage.RemoveUserImMsg(cc.appId, cc.userId, msgId)
}

//处理当前会话收到的群组消息
func (cc *clientConn) hGroupIm(msg *core.Message) {
	if msg.Data == nil {
		return
	}
	msgIm := msg.Data.(*common.MsgIm)
	if msgIm.SenderId == 0 || msgIm.ReceiverId == 0 || msgIm.Content == "" {
		log.Warnln("消息不完整", tools.JsonMarshal(msgIm))
		return
	}
	msgIm.Timestamp = time.Now().UnixNano()
	storage.UpdateOrCreateTalkId(cc.appId, msgIm)
	storage.MakeMsgId(msgIm)

	if bindGroupImAndSave(cc.appId, msgIm) {
		//读取所有的群成员， 循环保存离线消息数据
		groupId := int(msgIm.GetReceiverId())
		groupUsers := storage.GetAllGroupUser(cc.appId, groupId, storage.ShowGroupUser)
		if len(groupUsers) > 0 {
			var gus []*storage.GroupUser
			for _, groupUser := range groupUsers {
				if groupUser.UserId != cc.userId {
					gus = append(gus, groupUser)
				}
			}
			if len(gus) > 0 {
				storage.AppendUserGroupImMsg(cc.appId, gus, msgIm)
			}
		}
		//上报路由服务器，转发群消息
		cc.im.sendRouteImMsg(cc.appId, core.TransportGroupMsg, msgIm)
		rspGroupImAck(cc, msgIm)
	}
}

//处理当前会话收到的登录消息
func (cc *clientConn) hLogin(msg *core.Message) {
	if msg.Data == nil {
		return
	}
	//读取登录消息的token
	msgData := msg.Data.(*common.LoginIm)
	tokenStr := msgData.GetToken()
	appId := int(msgData.GetAppId())

	retData := &common.AuthIm{
		Status: 0,
		Code: 101,
		Data: "不是个有效token",
	}
	if tokenStr != "" {
		//检测token合法性
		token, err := auth.CheckToken(appId, tokenStr)
		if token != nil && err == nil {
			if token.Valid {
				deviceType := int(msgData.GetDeviceType())
				userId := msgData.GetUid()
				cc.appId = appId
				cc.userId = userId
				cc.deviceType = deviceType
				cc.token = token

				//appInfo中存储此用户的连接
				ai := cc.im.allApp.findORCreateAppInfo(cc.appId)
				ai.InsertConn(cc)
				//获取用户加入的所有群组, 循环群组列表加入此用户的会话
				groupIds := storage.GetAllUserJoinGroupId(cc.appId, cc.userId, storage.ShowGroupUser)
				if len(groupIds) > 0 {
					for _, groupId := range groupIds {
						gu := ai.findOrCreateGroupUsers(groupId)
						gu.InsertConn(cc)
					}
				}
				sendRouteLoginMsg(cc)

				retData.Status = 1
				retData.Code = 100
				retData.Data = "登录成功"
			}
		}
	}
	msg.Cmd = core.Auth
	msg.Data = retData
	cc.sendMessage(msg)
}

//处理当前会话收到的退出消息
func (cc *clientConn) hLogout() {
	if cc.appId != 0 && cc.userId != 0 {
		//移除app中此用户的会话
		ai := cc.im.allApp.findORCreateAppInfo(cc.appId)
		ai.RemoveConn(cc)
		//获取用户加入的所有群组, 循环群组列表移除此用户的会话
		groupIds := storage.GetAllUserJoinGroupId(cc.appId, cc.userId, storage.ShowGroupUser)
		if len(groupIds) > 0 {
			for _, groupId := range groupIds {
				gu := ai.findOrCreateGroupUsers(groupId)
				gu.RemoveConn(cc)
			}
		}
		sendRouteLogoutMsg(cc)
	}
	cc.tcpConn.CloseRead()
	cc.tcpConn.Close()
}

//处理当前会话收到的离线消息
func (cc *clientConn) hSyncOffline(msg *core.Message) {
	so := msg.Data.(*common.SyncOffline)
	limit := int(so.Limit)
	msgIms, surplusCount := storage.GetUserOfflineMsg(cc.appId, cc.userId, limit)
	offlineMsg := &common.SyncOfflineMsg{
		MsgList: msgIms,
		SurplusCount: int32(surplusCount),
	}

	pushMsg := &core.Message{
		Cmd: core.SyncOfflineMsg,
		Data: offlineMsg,
	}
	cc.sendMessage(pushMsg)
}

//处理推送离线消息后，客户端返回的离线消息ack
func (cc *clientConn) hSyncOfflineAck(msg *core.Message) {
	msgAck := msg.Data.(*common.MsgImAck)
	msgId := msgAck.MsgId
	storage.RemoveUserOfflineMsg(cc.appId, cc.userId, msgId)
}
