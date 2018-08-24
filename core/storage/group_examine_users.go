package storage

import (
	"github.com/jinzhu/gorm"
	"time"
)

const ApplyForExamineType = 0
const InviteExamineType = 1

const (
	DefaultExamineStatus = 0
	RefuseExamineStatus  = -1
	AgreeExamineStatus   = 1
)

type GroupExamineUser struct {
	ExamineId   int64 `gorm:"primary_key"`
	AppId       int
	GroupId     int
	GroupName   string
	UserId      int64
	UserName    string
	OpUserId    int64
	CreateTime  int64
	ExamineType int8 `gorm:"default:0"`
	Status      int8 `gorm:"default:0"`
	ExamineTime int64
}

func GetExamineUserByExamineId(examineId int64) *GroupExamineUser {
	examineData := &GroupExamineUser{}
	err := dbClient.Model(examineData).Where("examine_id=?", examineId).First(examineData).Error
	if err != nil {
		return nil
	}
	return examineData
}

func CreateGroupExamineUser(groupExamineUser *GroupExamineUser) error {
	err := dbClient.Model(&GroupExamineUser{}).Create(groupExamineUser).Error
	if err != nil {
		return err
	}
	return nil
}

//审核通过  添加用户至群组成员表并修改群组信息
func AgreeExamineUser(appId int, groupId int, examineId int64, examineUserId int64, opUserId int64) (err error) {
	timeNow := time.Now().Unix()
	setVal := map[string]interface{}{
		"status":       AgreeExamineStatus,
		"examine_time": timeNow,
		"op_user_id":   opUserId,
	}
	groupUser := &GroupUser{
		AppId: appId,
		GroupId: groupId,
		UserId: examineUserId,
		JoinTime: timeNow,
	}

	tx := dbClient.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
		tx.Close()
	}()

	err = dbClient.Model(&GroupExamineUser{}).
		Where("app_id=? and examine_id=? and user_id=? and status = 0", appId, examineId, examineUserId).
		Updates(setVal).Error
	if err != nil {
		return err
	}

	err = tx.Model(&GroupUser{}).Create(groupUser).Error
	if err != nil {
		return err
	}

	err = tx.Model(&Group{}).
		Where("app_id=? and group_id=?", appId, groupId).
		Update("user_count", gorm.Expr("user_count+?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

//拒绝加入
func RefuseExamineUser(appId int, examineId int64, examineUserId int64, opUserId int64) error {
	setVal := map[string]interface{}{
		"status":       RefuseExamineStatus,
		"examine_time": time.Now().Unix(),
		"op_user_id":   opUserId,
	}
	err := dbClient.Model(&GroupExamineUser{}).
		Where("app_id=? and examine_id=? and user_id=? and status = 0", appId, examineId, examineUserId).
		Updates(setVal).Error

	if err != nil {
		return err
	}
	return nil
}
