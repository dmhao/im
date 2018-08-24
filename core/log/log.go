package log

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"im/core/config"
	"im/core/tools"
	"os"
	"path"
	"time"
)

var logClient *logrus.Logger

type ServerType int

const (
	RouteServer  ServerType = 1
	ClientServer ServerType = 2
)

var logNameMap = map[ServerType]map[string]string{
	RouteServer:  {"log": "route-server.log", "error_log": "route-server-error.log"},
	ClientServer: {"log": "im-server.log", "error_log": "im-server-error.log"},
}

func setNull(logger *logrus.Logger) {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	logger.Out = src
}

func getLogFormat(format string) string {
	dateFormat := "%Y-%m-%d"
	switch format {
	case "Hour":
		dateFormat = "%Y-%m-%d-%H"
	}

	return dateFormat
}

func Init(serverType ServerType) {
	logClient = logrus.New()
	setNull(logClient)
	logClient.SetLevel(logrus.DebugLevel)

	//获取日志路径 并创建
	logDir := tools.ResolveRealDirPath(config.GetCommonConf().LogDir)
	os.MkdirAll(logDir, os.ModePerm)

	logPath := path.Join(logDir, logNameMap[serverType]["log"])
	errorLogPath := path.Join(logDir, logNameMap[serverType]["error_log"])

	infoWriter, err := rotatelogs.New(
		logPath+"."+getLogFormat(config.GetCommonConf().RotationFormat)+".log",
		rotatelogs.WithLinkName(logPath),                                                          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(config.GetCommonConf().LogExpire)*time.Hour),          // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(config.GetCommonConf().RotationTime)*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		panic(err)
	}

	errorWriter, err := rotatelogs.New(
		errorLogPath+"."+getLogFormat(config.GetCommonConf().RotationFormat)+".log",
		rotatelogs.WithLinkName(errorLogPath),                                                     // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(config.GetCommonConf().LogExpire)*time.Hour),          // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(config.GetCommonConf().RotationTime)*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		panic(err)
	}

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  infoWriter,
		logrus.WarnLevel:  infoWriter,
		logrus.FatalLevel: errorWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	logClient.AddHook(lfHook)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logClient.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logClient.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	logClient.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logClient.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	logClient.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	logClient.Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logClient.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logClient.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logClient.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logClient.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	logClient.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logClient.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	logClient.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	logClient.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logClient.Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	logClient.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	logClient.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	logClient.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	logClient.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	logClient.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	logClient.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	logClient.Fatalln(args...)
}
