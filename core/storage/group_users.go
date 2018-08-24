package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type GroupUser struct {
	Id       int64 `gorm:"primary_key"`
	AppId    int
	GroupId  int
	UserId   int64
	UserRole int8
	JoinTime int64
	Status   int8 `gorm:"default:1"`
}

const DefaultRole = 0
const MasterRole = 1
const ManagerRole = 2
const ShowGroupUser = 1
const HideGroupUser = 0
const DeleteGroupUser = -1

func GetUserExistsCount(appId int, groupId int, userIds []string, status int8) int {
	var count int
	dbClient.Model(&GroupUser{}).
		Where("app_id=? and group_id=? and user_id in (?) and status=?", appId, groupId, userIds, status).
		Count(&count)
	return count
}

func GetGroupUser(appId int, groupId int, userId int64, status int8) *GroupUser {
	groupUser := &GroupUser{}
	dbClient.Model(groupUser).
		Where("app_id=? and group_id=? and user_id=? and status=?", appId, groupId, userId, status).
		First(groupUser)
	return groupUser
}

//获取群组用户分页列表
func GetGroupUsers(appId int, groupId int, status int8, page int) []GroupUser {
	offset := (page - 1) * ApiRequestLimit

	var groupUsers []GroupUser
	dbClient.Model(&GroupUser{}).
		Where("app_id=? and group_id=? and status=?", appId, groupId, status).
		Offset(offset).Limit(ApiRequestLimit).Find(&groupUsers)
	return groupUsers
}

//获取群组所有用户列表
func GetAllGroupUser(appId int, groupId int, status int8) []*GroupUser {
	var groupUsers []*GroupUser
	dbClient.Model(&GroupUser{}).
		Where("app_id=? and group_id=? and status=?", appId, groupId, status).Find(&groupUsers)
	return groupUsers
}

//获取用户加入的群组
func GetAllUserJoinGroupId(appId int, userId int64, status int8) []int {
	groupIdsCache, err := GetUserJoinGroupIdCache(appId, userId)
	if err == nil && len(groupIdsCache) > 0 {
		return groupIdsCache
	}
	var userGroups []GroupUser
	dbClient.Model(&GroupUser{}).
		Where("app_id=? and user_id=? and status=?", appId, userId, status).Find(&userGroups)

	var groupIds []int
	if len(userGroups) > 0 {
		for _, groupUser := range userGroups {
			groupIds = append(groupIds, groupUser.GroupId)
		}
		SetUserJoinGroupIdCache(appId, userId, groupIds)
	}
	return groupIds
}

func CreateGroupUsers(appId int, groupId int, groupUsers []*GroupUser) (err error) {
	tx := dbClient.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
		tx.Close()
	}()
	for _, tmpGroupUser := range groupUsers {
		err = tx.Model(&GroupUser{}).Create(tmpGroupUser).Error
		if err != nil {
			return err
		}
	}

	err = tx.Model(&Group{}).
		Where("app_id=? and group_id=?", appId, groupId).
		Update("user_count", gorm.Expr("user_count+?", len(groupUsers))).Error
	if err != nil {
		return err
	}
	return nil
}

func GetGroupUserByRole(appId int, groupId int, status int8, userRoles []int8) []*GroupUser {
	var groupUsers []*GroupUser
	dbClient.Model(&GroupUser{}).
		Where("app_id=? and group_id=? and status=? and user_role in (?)", appId, groupId, status, userRoles).
		Find(&groupUsers)
	return groupUsers
}

func UpdateGroupUserDeleteData(appId int, groupId int, userId int64) (err error) {
	tx := dbClient.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
		tx.Close()
	}()

	err = dbClient.Model(&GroupUser{}).
		Where("app_id=? and group_id=? and user_id=? ", appId, groupId, userId).
		Updates(map[string]interface{}{"status": DeleteGroupUser, "update_time": time.Now().Unix()}).Error
	if err != nil {
		return err
	}

	err = tx.Model(&Group{}).
		Where("app_id=? and group_id=?", appId, groupId).
		Update("user_count", gorm.Expr("user_count-?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func SetManager(appId int, groupId int, userId int64) (err error) {
	err = dbClient.Model(&GroupUser{}).
		Where("app_id=? and group_id=? and user_id=?", appId, groupId, userId).
		Updates(map[string]interface{}{"user_role": ManagerRole, "update_time": time.Now().Unix()}).Error
	return err
}

func UnsetManager(appId int, groupId int, userId int64) (err error) {
	err = dbClient.Model(&GroupUser{}).
		Where("app_id=? and group_id=? and user_id=?", appId, groupId, userId).
		Updates(map[string]interface{}{"user_role": DefaultRole, "update_time": time.Now().Unix()}).Error
	return err
}
