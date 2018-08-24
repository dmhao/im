package main

import (
	"im/core/server"
	"im/core/log"
	"im/core/storage"
	"im/core/config"
	"im/core/api"
	"net"
	"os"
	"os/signal"
	"syscall"
	"math/rand"
	"time"
	"runtime"
)



func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	rand.Seed(time.Now().UnixNano())
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	log.Init(log.ClientServer)

	commonConf := config.GetCommonConf()
	err = storage.NewRedis(commonConf.RedisAddr, commonConf.RedisPassWord)
	if err != nil {
		log.Fatalln("redis连接失败", err)
	}

	err = storage.NewDB(commonConf.DBUser, commonConf.DBPassword, commonConf.DBHost, commonConf.DBPort, commonConf.DBDatabase)
	if err != nil {
		log.Fatalln("mysql连接失败", err)
	}

	tcpAddr,err := net.ResolveTCPAddr("tcp", config.GetClientConf().ClientServerAddr)
	if err != nil {
		log.Fatalln("连接服务地址解析失败",  err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln("连接服务监听失败",  err)
	}

	//启动api
	go api.ServerStart()

	cs := server.NewClientServer(tcpListener)
	go func() {
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		log.Infoln("收到服务关闭信号，服务准备关闭")
		cs.Stop()
	}()
	cs.Start()
	storage.StopDB()
	storage.StopRedis()
	log.Infoln("连接服务器已经关闭")
	os.Exit(0)
}






