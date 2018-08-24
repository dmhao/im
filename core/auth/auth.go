package auth

import (
	"github.com/dgrijalva/jwt-go"
	"im/core/log"
	"im/core/storage"
	"strconv"
	"time"
)

const Issuer = "hjh_im"
const TokenExpire = 24 * time.Hour

func GetToken(appId string, userId, secretKey string) string {
	signKey := []byte(secretKey)
	claims := &jwt.StandardClaims{
		Subject:   appId,
		Audience:  userId,
		ExpiresAt: time.Now().Add(TokenExpire).Unix(),
		Issuer:    Issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signKey)
	if err != nil {
		log.Warnln("token生成失败", err)
	}
	return ss
}

func CheckToken(appId int, tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		app := storage.GetAppById(appId)
		return []byte(app.SecretKey), nil
	})
	return token, err
}

func CheckTokenRightful(appId int, userId int64, token *jwt.Token) bool {
	if token != nil {
		if token.Valid {
			claims := token.Claims.(*jwt.StandardClaims)
			if strconv.Itoa(appId) == claims.Subject {
				now := time.Now().Unix()
				expire := claims.ExpiresAt
				if expire == 0 {
					return true
				}
				return now <= expire
			}
		}
	}
	return false
}
