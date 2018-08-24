package storage

import "time"

type Group struct {
	GroupId         int `gorm:"primary_key"`
	AppId           int
	GroupName       string
	GroupDes        string
	GroupIcon       string
	MasterUserId    int64
	UserCount       int16 `gorm:"default:1"`
	MaxUserCount    int16
	JoinNeedExamine int8 `gorm:"default:1"`
	CreateTime      int64
	UpdateTime      int64
	Status          int8 `gorm:"default:1"`
}

const ShowGroupStatus = 1
const HideGroupStatus = 0
const DeleteGroupStatus = -1

func GetGroup(appId int, groupId int, status int8) *Group {
	group := &Group{}
	dbClient.Model(group).Where("app_id=? and group_id=? and status=? ", appId, groupId, status).
		First(group)
	return group
}

//创建群组 并把群主加入群组
func CreateGroup(group *Group) (err error) {
	tx := dbClient.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
		tx.Close()
	}()
	if err = tx.Model(&Group{}).Create(group).Error; err != nil {
		return err
	}
	groupUser := &GroupUser{
		AppId: group.AppId,
		GroupId: group.GroupId,
		UserId: group.MasterUserId,
		JoinTime: time.Now().Unix(),
		UserRole: MasterRole,
	}

	if err = tx.Model(groupUser).Create(groupUser).Error; err != nil {
		return err
	}
	return nil
}

//修改可修改的普通数据
func UpdateGroupCommonData(group *Group) error {
	err := dbClient.Model(&Group{}).
		Omit("app_id", "group_id", "master_user_id", "user_count", "max_user_count", "create_time", "status").
		Where("app_id=? and group_id=?", group.AppId, group.GroupId).
		Update(group).Error
	return err
}

//修改为删除状态
func UpdateGroupDeleteData(appId int, groupId int) (err error) {
	tx := dbClient.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
		tx.Close()
	}()
	//修改群组为删除状态
	err = tx.Model(&Group{}).Where("app_id=? and group_id=?", appId, groupId).
		Updates(map[string]interface{}{"status": DeleteGroupStatus, "update_time": time.Now().Unix()}).Error
	if err != nil {
		return err
	}
	//修改群组的用户都为隐藏状态
	err = tx.Model(&GroupUser{}).Where("app_id=? and group_id=?", appId, groupId).
		Updates(map[string]interface{}{"status": HideGroupUser, "update_time": time.Now().Unix()}).Error
	if err != nil {
		return err
	}

	return err
}

//获取所有的圈子
func GetGroups(appId int, status int8, page int) []Group {
	offset := (page - 1) * ApiRequestLimit

	var groups []Group
	dbClient.Model(&Group{}).
		Where("app_id=? and status=? ", appId, status).
		Offset(offset).Limit(ApiRequestLimit).Find(&groups)

	return groups
}

//获取用户加入的圈子列表
func GetUserJoinGroups(appId int, userId int64, status int8, page int) []Group {
	offset := (page - 1) * ApiRequestLimit
	joinGroupIds := GetAllUserJoinGroupId(appId, userId, ShowGroupUser)
	var groups []Group

	if len(joinGroupIds) > 0 {
		dbClient.Model(&Group{}).
			Where("group_id in (?) and app_id=? and status=? ", joinGroupIds, appId, status).
			Offset(offset).Limit(ApiRequestLimit).Find(&groups)
	}
	return groups
}
