package route

import (
	"sync"
)

//存储群组中所有的用户
type groupUsers struct {
	rw   *sync.RWMutex
	data map[int64]bool
}

//应用的数据,users中存储当前所有在线状态的用户集合 groups中存储所有群组的集合
type appInfo struct {
	rw    *sync.RWMutex
	appId int
	//app中用户集合
	users map[int64]bool
	//app中群组集合
	groups map[int]*groupUsers
}

//app数据集合
type apps struct {
	rw   *sync.RWMutex
	data map[int]*appInfo
}

//判断app和该app下是否存在某个用户
func (allApp *apps) appAndUserIdExists(appId int, userId int64) bool {
	ai := allApp.hasAppId(appId)
	if ai == nil {
		return false
	}
	return ai.hasUserId(userId)
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

//所有的app中是否存在某个app
func (allApp *apps) hasAppId(appId int) *appInfo {
	allApp.rw.RLock()
	defer allApp.rw.RUnlock()

	ai, ok := allApp.data[appId]
	if !ok {
		return nil
	}
	return ai
}

//创建app
func newAppInfo(appId int) *appInfo {
	ai := &appInfo{}
	ai.rw = &sync.RWMutex{}
	ai.appId = appId
	ai.users = make(map[int64]bool)
	ai.groups = make(map[int]*groupUsers)
	return ai
}

//app中添加一个在线用户
func (ai *appInfo) insertUserId(userId int64) {
	ai.rw.Lock()
	defer ai.rw.Unlock()

	ai.users[userId] = true
}

//判断一个app中某个用户是否在线
func (ai *appInfo) hasUserId(userId int64) bool {
	ai.rw.RLock()
	defer ai.rw.RUnlock()

	online, ok := ai.users[userId]
	if !ok {
		return false
	}
	return online
}

//app中删除一个用户
func (ai *appInfo) removeUserId(userId int64) {
	ai.rw.Lock()
	defer ai.rw.Unlock()

	delete(ai.users, userId)
}

//查找app中是否存在某个群组不存在创建一个群组，返回此群组
func (ai *appInfo) findOrCreateGroupUsers(groupId int) *groupUsers {
	ai.rw.Lock()
	defer ai.rw.Unlock()

	if gu, ok := ai.groups[groupId]; ok {
		return gu
	}

	gu := &groupUsers{}
	gu.rw = &sync.RWMutex{}
	gu.data = make(map[int64]bool)

	ai.groups[groupId] = gu
	return gu
}

//群组中添加一个在线用户
func (gu *groupUsers) insertUserId(userId int64) {
	gu.rw.Lock()
	defer gu.rw.Unlock()

	gu.data[userId] = true
}

//樽俎中删除一个用户
func (gu *groupUsers) removeUserId(userId int64) {
	gu.rw.Lock()
	defer gu.rw.Unlock()

	delete(gu.data, userId)
}
