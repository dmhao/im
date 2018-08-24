package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"im/core/auth"
	"im/core/tools"
	"strconv"
)

func CheckRequest(c *gin.Context) {
	appIdStr := c.Query("appId")
	appId := tools.StrToi(appIdStr)
	tokenStr := c.GetHeader("Token")
	if appId == 0 || tokenStr == "" {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}
	token, err := auth.CheckToken(appId, tokenStr)
	if token != nil && err == nil {
		if token.Valid {
			claims := token.Claims.(*jwt.StandardClaims)
			if appIdStr == claims.Subject {
				c.Set("appId", appId)
				userId, _ := strconv.ParseInt(claims.Audience, 10, 64)
				c.Set("userId", userId)
				c.Next()
				return
			}
		}
	}
	mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
	return
}
