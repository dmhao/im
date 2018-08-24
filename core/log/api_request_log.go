package log

import (
	"github.com/gin-gonic/gin"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"im/core/config"
	"im/core/tools"
	"os"
	"path"
	"time"
)

func Logger() gin.HandlerFunc {
	logger := logrus.New()
	setNull(logger)
	logger.SetLevel(logrus.InfoLevel)

	//获取日志路径 并创建
	logDir := tools.ResolveRealDirPath(config.GetCommonConf().LogDir)
	os.MkdirAll(logDir, os.ModePerm)
	logPath := path.Join(logDir, "api-request.log")

	logWriter, err := rotatelogs.New(
		logPath+"."+getLogFormat(config.GetCommonConf().RotationFormat)+".log",
		rotatelogs.WithLinkName(logPath),                                           // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(config.GetCommonConf().LogExpire*time.Hour),          // 文件最大保存时间
		rotatelogs.WithRotationTime(config.GetCommonConf().RotationTime*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		panic(err)
	}

	writeMap := lfshook.WriterMap{
		logrus.WarnLevel: logWriter,
		logrus.InfoLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	logger.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logger.Infof("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}
