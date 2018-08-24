package api

import (
	"github.com/gin-gonic/gin"
	"im/core/api/router"
	"im/core/config"
	"im/core/log"
)

func ServerStart() {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(log.Logger(), log.Recovery())
	router.Init(e)
	err := e.Run(config.GetClientConf().ApiPort)
	if err != nil {
		log.Infoln("apiServer启动失败", err)
	} else {
		log.Infoln("apiServer启动完成")
	}
}
