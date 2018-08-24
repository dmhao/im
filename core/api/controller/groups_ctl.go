package controller

import (
	"github.com/gin-gonic/gin"
	"im/core/auth"
	"im/core/storage"
	"im/core/tools"
	"strconv"
	"strings"
	"time"
)

func Groups(c *gin.Context) {
	appId := c.GetInt("appId")
	page := GetPage(c)
	groups := storage.GetGroups(appId, storage.ShowGroupStatus, page)

	mContext{c}.SuccessResponse(groups)
}

func GroupInfo(c *gin.Context) {
	appId := c.GetInt("appId")
	groupId := tools.StrToi(c.Param("groupId"))
	if groupId == 0 {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}
	group := storage.GetGroup(appId, groupId, storage.ShowGroupStatus)
	mContext{c}.SuccessResponse(group)
}

func CreateGroup(c *gin.Context) {
	appId := c.GetInt("appId")
	userId := tools.StrToInt64(c.Param("userId"))

	group := &storage.Group{}
	bindErr := c.Bind(group)
	if bindErr != nil || userId == 0 || group.MaxUserCount == 0 ||
		group.GroupName == "" || group.GroupIcon == "" {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}
	currentTime := time.Now().Unix()
	group.AppId = appId
	group.MasterUserId = userId
	group.CreateTime = currentTime
	group.UpdateTime = currentTime
	err := storage.CreateGroup(group)
	if err != nil {
		mContext{c}.ErrorResponse(CreateError, CreateErrorMsg)
	} else {
		mContext{c}.SuccessResponse(group)
	}
}

func UpdateGroup(c *gin.Context) {
	appId := c.GetInt("appId")
	groupId := tools.StrToi(c.Param("groupId"))
	userId := tools.StrToInt64(c.Param("userId"))

	group := &storage.Group{}
	bindErr := c.Bind(group)
	if bindErr != nil || groupId == 0 || userId == 0 ||
		group.GroupName == "" || group.GroupIcon == "" || group.MaxUserCount == 0 {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}

	dataGroup := storage.GetGroup(appId, groupId, storage.ShowGroupStatus)
	if dataGroup.GroupId == 0 {
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}
	if dataGroup.MasterUserId != userId {
		mContext{c}.ErrorResponse(ForbiddenError, NoPowerUpdate)
		return
	}

	currentTime := time.Now().Unix()
	group.UpdateTime = currentTime
	group.AppId = appId
	group.GroupId = groupId
	err := storage.UpdateGroupCommonData(group)
	if err != nil {
		mContext{c}.ErrorResponse(UpdateError, UpdateErrorMsg)
	} else {
		mContext{c}.SuccessResponse(group)
	}
}

func DeleteGroup(c *gin.Context) {
	appId := c.GetInt("appId")
	groupId := tools.StrToi(c.Param("groupId"))
	userId := tools.StrToInt64(c.Param("userId"))

	if groupId == 0 || userId == 0 {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}

	dataGroup := storage.GetGroup(appId, groupId, storage.ShowGroupStatus)
	if dataGroup.GroupId == 0 {
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}
	if dataGroup.MasterUserId != userId {
		mContext{c}.ErrorResponse(ForbiddenError, NoPowerUpdate)
		return
	}

	err := storage.UpdateGroupDeleteData(appId, groupId)
	if err != nil {
		mContext{c}.ErrorResponse(DeleteError, DeleteErrorMsg)
	} else {
		mContext{c}.SuccessResponse(nil)
	}
}

func GroupUsers(c *gin.Context) {
	appId := c.GetInt("appId")
	groupId := tools.StrToi(c.Param("groupId"))
	if groupId == 0 {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}
	page := GetPage(c)
	groupUsers := storage.GetGroupUsers(appId, groupId, storage.ShowGroupUser, page)
	mContext{c}.SuccessResponse(groupUsers)
}

func CreateGroupUser(c *gin.Context) {
	appId := c.GetInt("appId")
	opUserId := c.GetInt64("userId")
	groupId := tools.StrToi(c.Param("groupId"))
	userIdsStr, _ := c.GetPostForm("UserIds")
	userIds := strings.Split(userIdsStr, ",")
	if groupId == 0 || len(userIds) == 0 || opUserId == 0 {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}

	//判断群组是否存在
	dataGroup := storage.GetGroup(appId, groupId, storage.ShowGroupStatus)
	if dataGroup.GroupId == 0 {
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}

	//查看群组用户是否存在
	userCount := storage.GetUserExistsCount(appId, groupId, userIds, storage.ShowGroupUser)
	if userCount > 0 {
		mContext{c}.ErrorResponse(OperateRepeat, OperateRepeatMsg)
		return
	}

	//查询操作用户，普通用户无权限操作
	opGroupUser := storage.GetGroupUser(appId, groupId, opUserId, storage.ShowGroupUser)
	if opGroupUser.Id == 0 || opGroupUser.UserRole == 0 {
		mContext{c}.ErrorResponse(ForbiddenError, NoPowerUpdate)
		return
	}

	var groupUsers []*storage.GroupUser
	joinTime := time.Now().Unix()
	for _, userId := range userIds {
		userId64, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			continue
		}
		groupUser := &storage.GroupUser{
			AppId: appId,
			GroupId: groupId,
			UserId: userId64,
			JoinTime: joinTime,
		}
		groupUsers = append(groupUsers, groupUser)
	}
	if len(groupUsers) > 0 {
		err := storage.CreateGroupUsers(appId, groupId, groupUsers)
		if err != nil {
			mContext{c}.ErrorResponse(CreateError, CreateErrorMsg)
		} else {
			mContext{c}.SuccessResponse(groupUsers)
		}
	} else {
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
	}
}

func DeleteGroupUser(c *gin.Context) {
	appId := c.GetInt("appId")
	groupId := tools.StrToi(c.Param("groupId"))
	userId := tools.StrToInt64(c.Param("userId"))
	opUserId := c.GetInt64("userId")

	if groupId == 0 || userId == 0 || opUserId == 0 {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}

	//查看是否存在
	groupUser := storage.GetGroupUser(appId, groupId, userId, storage.ShowGroupUser)
	if groupUser.Id == 0 {
		mContext{c}.ErrorResponse(OperateRepeat, OperateRepeatMsg)
		return
	}

	//查询操作用户，普通用户无权限操作
	opGroupUser := storage.GetGroupUser(appId, groupId, opUserId, storage.ShowGroupUser)
	if opGroupUser.Id == 0 || opGroupUser.UserRole == 0 ||
		(opGroupUser.UserId == userId && opGroupUser.UserRole == storage.MasterRole) {
		mContext{c}.ErrorResponse(ForbiddenError, NoPowerUpdate)
		return
	}

	err := storage.UpdateGroupUserDeleteData(appId, groupId, userId)
	if err != nil {
		mContext{c}.ErrorResponse(DeleteError, DeleteErrorMsg)
	} else {
		mContext{c}.SuccessResponse(nil)
	}
}

//用户群组
func UserGroups(c *gin.Context) {
	appId := c.GetInt("appId")
	userId := tools.StrToInt64(c.Param("userId"))

	if userId == 0 {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}
	page := GetPage(c)
	groups := storage.GetUserJoinGroups(appId, userId, storage.ShowGroupStatus, page)

	mContext{c}.SuccessResponse(groups)
}

//设置为管理员
func SetManager(c *gin.Context) {
	appId := c.GetInt("appId")
	opUserId := c.GetInt64("userId")
	groupId := tools.StrToi(c.Param("groupId"))
	managerId := tools.StrToInt64(c.PostForm("ManagerId"))

	if groupId == 0 || opUserId == 0 || managerId == 0 {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}

	managerUser := storage.GetGroupUser(appId, groupId, managerId, storage.ShowGroupUser)
	if managerUser.Id == 0 {
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}

	if managerUser.UserRole == storage.ManagerRole {
		mContext{c}.ErrorResponse(OperateRepeat, OperateRepeatMsg)
		return
	}

	//非群主没有设置管理员权限
	masterUser := storage.GetGroupUser(appId, groupId, opUserId, storage.ShowGroupUser)
	if masterUser.UserRole != storage.MasterRole {
		mContext{c}.ErrorResponse(ForbiddenError, NoPowerUpdate)
		return
	}

	err := storage.SetManager(appId, groupId, managerId)
	if err != nil {
		mContext{c}.ErrorResponse(OperateError, OperateRepeatMsg)
	} else {
		mContext{c}.SuccessResponse(nil)
	}
}

//取消管理员
func UnsetManager(c *gin.Context) {
	appId := c.GetInt("appId")
	opUserId := c.GetInt64("userId")
	groupId := tools.StrToi(c.Param("groupId"))
	managerId := tools.StrToInt64(c.Param("managerId"))
	if groupId == 0 || opUserId == 0 || managerId == 0 {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}

	managerUser := storage.GetGroupUser(appId, groupId, managerId, storage.ShowGroupUser)
	if managerUser.Id == 0 || managerUser.UserRole != storage.ManagerRole {
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}

	//非群主没有设置管理员权限
	masterUser := storage.GetGroupUser(appId, groupId, opUserId, storage.ShowGroupUser)
	if masterUser.UserRole != storage.MasterRole {
		mContext{c}.ErrorResponse(ForbiddenError, NoPowerUpdate)
	}

	err := storage.UnsetManager(appId, groupId, managerId)
	if err != nil {
		mContext{c}.ErrorResponse(OperateError, OperateErrorMsg)
	} else {
		mContext{c}.SuccessResponse(nil)
	}
}

func Token(c *gin.Context) {
	userId := c.Query("userId")
	secretId := c.Query("secretId")
	secretKey := c.Query("secretKey")

	if userId == "" || secretId == "" || secretKey == "" {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}

	app := storage.GetAppBySecret(secretId, secretKey)
	if app.AppId == 0 {
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}

	ss := auth.GetToken(strconv.Itoa(app.AppId), userId, secretKey)
	mContext{c}.SuccessResponse(gin.H{"AppId": app.AppId, "Token": ss})
}
