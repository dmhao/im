package server

import (
	"im/core/storage"
	"sync"
)

//存储群组中所有的用户
type groupUsers struct {
	rw   *sync.RWMutex
	data map[int64]*clientConn
}

//应用的数据,users中存储当前所有在线状态的用户会话集合 groups中存储所有群组的集合
type appInfo struct {
	rw    *sync.RWMutex
	appId int
	//app中用户集合
	users map[int64]*clientConn
	//app中群组集合
	groups map[int]*groupUsers
}

//app数据集合
type apps struct {
	rw   *sync.RWMutex
	data map[int]*appInfo
}

//查找是否存在某个app不存在创建一个app，返回此app
func (allApp *apps) findORCreateAppInfo(appId int) *appInfo {
	allApp.rw.Lock()
	defer allApp.rw.Unlock()

	if ai, ok := allApp.data[appId]; ok {
		return ai
	}
	ai := newAppInfo(appId)
	allApp.data[appId] = ai
	return ai
}

//app中添加一个用户会话
func (ai *appInfo) InsertConn(cc *clientConn) {
	ai.rw.Lock()
	ai.users[cc.userId] = cc
	ai.rw.Unlock()
}

//app中删除一个用户会话
func (ai *appInfo) RemoveConn(cc *clientConn) {
	ai.rw.Lock()
	delete(ai.users, cc.userId)
	ai.rw.Unlock()
}

//app中查找某个用户的会话
func (ai *appInfo) FindConn(userId int64) *clientConn {
	ai.rw.RLock()
	defer ai.rw.RUnlock()

	if userConn, ok := ai.users[userId]; ok {
		return userConn
	}
	return nil
}

//获取群组中所有用户的会话列表
func (ai *appInfo) GetGroupUsersConn(groupUsers []*storage.GroupUser) map[int64]*clientConn {
	ai.rw.RLock()
	userConnMap := make(map[int64]*clientConn)
	for _, groupUser := range groupUsers {
		if userConn, ok := ai.users[groupUser.UserId]; ok {
			userConnMap[groupUser.UserId] = userConn
		}
	}
	ai.rw.RUnlock()
	return userConnMap
}

//创建app
func newAppInfo(appId int) *appInfo {
	ai := &appInfo{}
	ai.rw = &sync.RWMutex{}
	ai.appId = appId
	ai.users = make(map[int64]*clientConn)
	ai.groups = make(map[int]*groupUsers)
	return ai
}

//查找app中是否存在某个群组不存在创建一个群组，返回此群组
func (ai *appInfo) findOrCreateGroupUsers(groupId int) *groupUsers {
	ai.rw.Lock()
	defer ai.rw.Unlock()

	if gu, ok := ai.groups[groupId]; ok {
		return gu
	}

	gu := &groupUsers{}
	gu.data = make(map[int64]*clientConn)
	gu.rw = &sync.RWMutex{}

	ai.groups[groupId] = gu
	return gu
}

//群组中添加某个用户会话
func (gu *groupUsers) InsertConn(cc *clientConn) {
	gu.rw.Lock()
	gu.data[cc.userId] = cc
	gu.rw.Unlock()
}

//群组中删除某个用户会话
func (gu *groupUsers) RemoveConn(cc *clientConn) {
	gu.rw.Lock()
	delete(gu.data, cc.userId)
	gu.rw.Unlock()
}
