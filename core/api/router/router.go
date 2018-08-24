package router

import (
	"github.com/gin-gonic/gin"
	"im/core/api/controller"
	"im/core/config"
)

func Init(e *gin.Engine) {
	e.LoadHTMLGlob("../core/api/templates/*")
	{
		v1 := e.Group("/v1")
		v1.Use(controller.CheckRequest)
		//群组列表
		v1.GET("/groups", controller.Groups)
		//群组详情
		v1.GET("/groups/:groupId", controller.GroupInfo)
		//创建群组
		v1.POST("/users/:userId/groups", controller.CreateGroup)
		//修改用户群组
		v1.POST("/users/:userId/groups/:groupId", controller.UpdateGroup)
		//用户的群组列表
		v1.GET("/users/:userId/groups", controller.UserGroups)
		//删除用户群组
		v1.DELETE("/users/:userId/groups/:groupId", controller.DeleteGroup)

		//群组成员
		v1.GET("/groups/:groupId/users", controller.GroupUsers)
		//群组添加成员
		v1.POST("/groups/:groupId/users", controller.CreateGroupUser)
		//群组踢出成员
		v1.DELETE("/groups/:groupId/users/:userId", controller.DeleteGroupUser)
		//设为管理员
		v1.POST("/groups/:groupId/managers", controller.SetManager)
		//取消管理员
		v1.DELETE("/groups/:groupId/managers/:managerId", controller.UnsetManager)

		v1.POST("/groups/:groupId/examineUsers", controller.CreateExamineUser)
		v1.GET("/groups/:groupId/examineUsers/:examineUserId", controller.ExamineUserOp)

		v1.POST("/images", controller.UploadImage)
		v1.POST("/videos", controller.UploadVideo)
		v1.POST("/audios", controller.UploadAudio)

		v1.POST("/messages", controller.CreateMessage)
		//获取token
		e.GET("/token", controller.Token)
	}
	e.GET("/mem", controller.Mem)
	if config.GetClientConf().ApiPProf {
		{
			debug := e.Group("/debug")
			debug.Any("/pprof/", controller.DebugIndex)
			debug.Any("/pprof/heap", controller.DebugIndex)
			debug.Any("/pprof/goroutine", controller.DebugIndex)
			debug.Any("/pprof/block", controller.DebugIndex)
			debug.Any("/pprof/mutex", controller.DebugIndex)
			debug.Any("/pprof/threadcreate", controller.DebugIndex)
			debug.Any("/pprof/profile", controller.DebugProfile)
			debug.Any("/pprof/symbol", controller.DebugSymbol)
			debug.Any("/pprof/trace", controller.DebugTrace)
			debug.Any("/pprof/cmdline", controller.DebugCmdline)
		}
	}

}
