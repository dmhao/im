package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"im/core"
	"im/core/server"
	"im/core/storage"
	"im/core/tools"
	"im/mp/common"
	"time"
)

const (
	agreeOp  = "agree"
	refuseOp = "refuse"
)

func CreateExamineUser(c *gin.Context) {
	appId := c.GetInt("appId")
	userId := c.GetInt64("userId")
	groupId := tools.StrToi(c.Param("groupId"))
	groupName := c.PostForm("GroupName")
	userName := c.PostForm("UserName")
	if groupId == 0 || userId == 0 || groupName == "" || userName == "" {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}

	//判断群组是否存在
	groupData := storage.GetGroup(appId, groupId, storage.ShowGroupStatus)
	if groupData.GroupId == 0 {
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}

	//查看群组用户是否存在
	groupUser := storage.GetGroupUser(appId, groupId, userId, storage.ShowGroupUser)
	if groupUser.Id > 0 {
		mContext{c}.ErrorResponse(OperateRepeat, OperateRepeatMsg)
		return
	}
	groupExamineUser := &storage.GroupExamineUser{
		AppId: appId,
		GroupId: groupId,
		GroupName: groupName,
		UserId: userId,
		UserName: userName,
		CreateTime: time.Now().Unix(),
		ExamineType: storage.ApplyForExamineType,
	}

	storage.CreateGroupExamineUser(groupExamineUser)
	if groupExamineUser.ExamineId != 0 {
		//获取群管理列表  推送审核消息
		groupUsers := storage.GetGroupUserByRole(appId, groupId, storage.ShowGroupUser,
			[]int8{storage.MasterRole, storage.ManagerRole})
		msgContent := examineMsgContent(groupExamineUser)

		var msgData []*core.ApiMessage
		for _, groupUser := range groupUsers {
			apiMsg := makeApiMessage(groupUser, msgContent)
			msgData = append(msgData, apiMsg)
		}
		if len(msgData) > 0 {
			t := time.NewTimer(20 * time.Second)
			select {
			case server.ApiMessagesCh <- msgData:
				t.Stop()
				mContext{c}.SuccessResponse(nil)
				return
			case <-t.C:
				mContext{c}.ErrorResponse(ServiceTimeOut, ServiceTimeOutMsg)
				return
			}
		}
	}
	mContext{c}.ErrorResponse(CreateError, CreateErrorMsg)
}

func examineMsgContent(examineUser *storage.GroupExamineUser) string {
	agreeUrl := fmt.Sprintf("/v1/groups/%v/examineUsers/%v?appId=%v&examineId=%v&op=%s",
		examineUser.GroupId, examineUser.UserId, examineUser.AppId, examineUser.ExamineId, agreeOp)
	refuseUrl := fmt.Sprintf("/v1/groups/%v/examineUsers/%v?appId=%v&examineId=%v&op=%s",
		examineUser.GroupId, examineUser.UserId, examineUser.AppId, examineUser.ExamineId, refuseOp)

	agreeButton := &Button{
		ClickType: "url",
		Text:      "同意",
		JumpData:  agreeUrl,
	}
	refuseButton := &Button{
		ClickType: "url",
		Text:      "拒绝",
		JumpData:  refuseUrl,
	}

	buttonMsg := ExamineButtonMsg{
		ApplyUserId:   examineUser.UserId,
		ApplyUserName: examineUser.UserName,
		GroupId:       examineUser.GroupId,
		GroupName:     examineUser.GroupName,
		ButtonMsg: ButtonMsg{
			Intro:   examineUser.UserName + "-申请加入群-" + examineUser.GroupName,
			Buttons: []*Button{agreeButton, refuseButton},
		},
	}

	content := struct {
		Type     string      `json:"type"`
		DataType string      `json:"data_type"`
		Data     interface{} `json:"data"`
	}{
		Type:     "join_group_examine",
		DataType: "button",
		Data:     buttonMsg,
	}

	transPortMsg := TransPort{
		Type:    "system_transport",
		Content: tools.JsonMarshal(content),
	}
	return tools.JsonMarshal(transPortMsg)
}

func makeApiMessage(groupUser *storage.GroupUser, msgContent string) *core.ApiMessage {
	uuidBytes, _ := uuid.NewV4()
	msgIm := &common.MsgIm{
		TraceId: "im-api-" + uuidBytes.String(),
		ChartType: 1,
		MsgType: core.SystemMessageType,
		ReceiverId: groupUser.UserId,
		Content: msgContent,
	}


	apiMsg := &core.ApiMessage{
		AppId: groupUser.AppId,
		Data: msgIm,
	}
	return apiMsg
}

func ExamineUserOp(c *gin.Context) {
	appId := c.GetInt("appId")
	opUserId := c.GetInt64("userId")
	op := c.Query("op")
	examineId := tools.StrToInt64(c.Query("examineId"))
	groupId := tools.StrToi(c.Param("groupId"))
	examineUserId := tools.StrToInt64(c.Param("examineUserId"))
	if examineId == 0 || groupId == 0 || examineUserId == 0 || (op != agreeOp && op != refuseOp) {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}

	opUserData := storage.GetGroupUser(appId, groupId, opUserId, storage.ShowGroupUser)
	if opUserData.UserRole != storage.MasterRole && opUserData.UserRole != storage.ManagerRole {
		mContext{c}.ErrorResponse(ForbiddenError, NoPowerUpdate)
		return
	}

	examineData := storage.GetExamineUserByExamineId(examineId)
	if examineData.Status != storage.DefaultExamineStatus {
		mContext{c}.ErrorResponse(OperateRepeat, OperateRepeatMsg)
		return
	}
	var err error
	if op == agreeOp {
		err = storage.AgreeExamineUser(appId, groupId, examineId, examineUserId, opUserId)
	} else if op == refuseOp {
		err = storage.RefuseExamineUser(appId, examineId, examineUserId, opUserId)
	}

	if err != nil {
		mContext{c}.ErrorResponse(OperateError, OperateErrorMsg)
		return
	}
	mContext{c}.SuccessResponse(nil)
}
