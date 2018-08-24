package log

import (
	"github.com/gin-gonic/gin"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"im/core/config"
	"im/core/tools"
	"net/http/httputil"
	"os"
	"path"
	"time"
)

func Recovery() gin.HandlerFunc {
	logger := logrus.New()
	setNull(logger)
	logger.SetLevel(logrus.ErrorLevel)

	//获取日志路径 并创建
	logDir := tools.ResolveRealDirPath(config.GetCommonConf().LogDir)
	os.MkdirAll(logDir, os.ModePerm)
	logPath := path.Join(logDir, "api-error.log")

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
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	logger.AddHook(lfHook)

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := tools.Stack(3)
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				logger.Errorf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httpRequest), err, stack)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
