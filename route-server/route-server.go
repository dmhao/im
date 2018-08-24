package main

import (
	"os"
	"os/signal"
	"syscall"
	"im/core/route"
	"im/core/log"
	"im/core/storage"
	"im/core/config"
	"net"
	"runtime"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	log.Init(log.RouteServer)

	commonConf := config.GetCommonConf()
	err = storage.NewRedis(commonConf.RedisAddr, commonConf.RedisPassWord)
	if err != nil {
		log.Fatalln("redis连接失败", err)
	}

	err = storage.NewDB(commonConf.DBUser, commonConf.DBPassword, commonConf.DBHost, commonConf.DBPort, commonConf.DBDatabase)
	if err != nil {
		log.Fatalln("mysql连接失败", err)
	}

	tcpAddr,err := net.ResolveTCPAddr("tcp", config.GetRouteConf().RouteServerAddr)
	if err != nil {
		log.Fatalln("路由监听地址解析错误",  err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln("路由服务器监听失败",  err)
	}

	rs := route.NewRouteServer(listener)
	go func() {
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		log.Infoln("收到服务关闭信号，服务准备关闭")
		rs.Stop()
	}()
	rs.Start()
	storage.StopDB()
	storage.StopRedis()
	log.Infoln("路由服务器已经关闭")
	os.Exit(0)
}