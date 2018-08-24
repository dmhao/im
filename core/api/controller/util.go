package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Button struct {
	ClickType string `json:"click_type"`
	Text      string `json:"text"`
	JumpData  string `json:"jump_data"`
}

type ButtonMsg struct {
	Intro   string      `json:"intro"`
	Buttons interface{} `json:"buttons"`
}

type ExamineButtonMsg struct {
	ApplyUserId   int64  `json:"apply_user_id"`
	ApplyUserName string `json:"apply_user_name"`
	GroupId       int    `json:"group_id"`
	GroupName     string `json:"group_name"`
	ButtonMsg
}

type TransPort struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// GetPage returns paging parameter.
func GetPage(c *gin.Context) int {
	ret, _ := strconv.Atoi(c.Query("page"))
	if 1 > ret {
		ret = 1
	}

	return ret
}

type mContext struct {
	*gin.Context
}

func (c mContext) ErrorResponse(code int, msg string) {
	result := NewResult()
	result.Code = code
	result.Msg = msg
	result.Data = nil
	c.AbortWithStatusJSON(200, result)
}

func (c mContext) SuccessResponse(data interface{}) {
	result := NewResult()
	result.Code = Success
	result.Msg = SuccessMsg
	result.Data = data
	c.JSON(200, result)
	c.Next()
}
