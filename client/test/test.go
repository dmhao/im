package test

import (
"net"
"im/client/core"
"fmt"
"im/mp/common"
"time"
"os"
"os/signal"
"syscall"
"im/core/tools"
"github.com/tidwall/gjson"
"bufio"
"strings"
)

/**
普通聊天
call-{"receiver_uid": 321321, "content":"你好"}
*/


var operation chan string
var operateData chan interface{}
//var tcpAddr = "172.16.100.11:7777"
var tcpAddr = "127.0.0.1:7777"
//var tcpAddr = "192.168.139.143:7777"
//var tcpAddr = "192.168.139.139:7777"

var uid int64 = 1
var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxIiwiZXhwIjoxNTM1MTY5MTgxLCJpc3MiOiJoamhfaW0iLCJzdWIiOiIxIn0.74XZh3jg2qyxZXPA4g3LTx5T1LWMu6WiixOfzP6mt4M"

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", tcpAddr)
	conn,err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("连接失败", err)
		return
	}
	operateData = make(chan interface{})
	operation = make(chan string)


	//登录写入
	msg := &core.Message{}
	msg.Cmd = core.Login
	data := &common.LoginIm{}
	data.Uid = uid
	data.Token = token
	data.AppId = 1
	data.DeviceType = 1
	msg.Data = data
	dataBytes := core.MakePacket(msg)
	conn.Write(dataBytes)
	//读取用户信息
	authMsg := core.ReadPacket(conn)
	authData := authMsg.Data.(*common.AuthIm)
	if authData.Status == 1 {
		fmt.Println("uid:", uid, string(authData.Data))
	} else {
		fmt.Println("uid:", uid, string(authData.Data))
	}
	syncData := &common.SyncOffline{}
	syncData.Limit = 0
	msg.Cmd = core.SyncOffline
	msg.Data = syncData
	dataBytes = core.MakePacket(msg)
	conn.Write(dataBytes)
	syncOfflineMsg := core.ReadPacket(conn)
	fmt.Println(syncOfflineMsg)
	if syncOfflineMsg.Data != nil {
		syncOfflineData := syncOfflineMsg.Data.(*common.SyncOfflineMsg)
		if len(syncOfflineData.MsgList) > 0 {
			length := len(syncOfflineData.MsgList)
			lastMsg := syncOfflineData.MsgList[length-1]
			ackMsg := &common.MsgImAck{}
			ackMsg.MsgId = lastMsg.MsgId
			ackMsg.Timestamp = lastMsg.Timestamp
			ackMsg.TraceId = lastMsg.TraceId
			ackMsg.TalkId = lastMsg.TalkId

			msg := &core.Message{}
			msg.Cmd = core.SyncOfflineMsgAck
			msg.Data = ackMsg
			msgBytes := core.MakePacket(msg)
			conn.Write(msgBytes)
		}
	}

	client := new(Client)
	client.uid = uid
	client.conn = conn

	go ping(client)
	go run(client)

	go hook(client)
	for {
		inputReader := bufio.NewReader(os.Stdin)
		input, _ := inputReader.ReadString('\n')
		inputSlice := strings.Split(input, "-")
		if len(inputSlice) == 2 {
			operation <- inputSlice[0]
			operateData <- inputSlice[1]
		}
	}
}

func main1() {
	start := 70000
	/*	if len(os.Args)!=0{
			getStart,_ := strconv.Atoi(os.Args[0])
			start = start + getStart
		}*/
	end := start + 10000
	for i:= start; i < end; i++ {
		go test(int64(i))
		if i % 300 == 0 {
			time.Sleep(3*time.Second)
		}
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case s := <-sig:
		fmt.Println( "路由服务器关闭 关闭信号-", s)
	}
}

var cover = 0
//捕获协程中的程序崩溃
func recoverFunc() {
	err := recover()
	if err != nil {
		stack := tools.Stack(3)
		fmt.Println("协程崩溃", err, string(stack))
	}
	stack := tools.Stack(6)
	fmt.Println("协程崩溃", err, string(stack))
	cover++
}

func test(uid int64) {
	defer recoverFunc()
	addr, _ := net.ResolveTCPAddr("tcp", tcpAddr)

	conn,err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("连接失败", err)
		return
	}
	operateData = make(chan interface{})
	operation = make(chan string)
	//登录写入
	msg := new(core.Message)
	msg.Cmd = core.Login
	data := new(common.LoginIm)
	data.Uid = uid
	data.Token = token
	data.AppId = 1
	data.DeviceType = 1
	msg.Data = data
	dataBytes := core.MakePacket(msg)
	conn.Write(dataBytes)
	//读取用户信息
	authMsg := core.ReadPacket(conn)
	authData := authMsg.Data.(*common.AuthIm)
	if authData.Status == 1 {
		fmt.Println("uid:", uid, string(authData.Data))
	} else {
		fmt.Println("uid:", uid, string(authData.Data))
	}
	syncData := new(common.SyncOffline)
	syncData.Limit = 0
	msg.Cmd = core.SyncOffline
	msg.Data = syncData
	dataBytes = core.MakePacket(msg)
	conn.Write(dataBytes)
	syncOfflineMsg := core.ReadPacket(conn)
	if syncOfflineMsg.Data != nil {
		syncOfflineData := syncOfflineMsg.Data.(*common.SyncOfflineMsg)
		if len(syncOfflineData.MsgList) > 0 {
			length := len(syncOfflineData.MsgList)
			lastMsg := syncOfflineData.MsgList[length-1]
			ackMsg := new(common.MsgImAck)
			ackMsg.MsgId = lastMsg.MsgId
			ackMsg.Timestamp = lastMsg.Timestamp
			ackMsg.TraceId = lastMsg.TraceId
			ackMsg.TalkId = lastMsg.TalkId

			msg := new(core.Message)
			msg.Cmd = core.SyncOfflineMsgAck
			msg.Data = ackMsg
			msgBytes := core.MakePacket(msg)
			conn.Write(msgBytes)
		}
	}

	client := new(Client)
	client.uid = uid
	client.conn = conn

	go ping(client)
	run(client)
}
func ping(client *Client) {
	for {
		time.Sleep(20 * time.Second)
		msg := new(core.Message)
		msg.Cmd = core.Ping
		msg.Version = 1
		msg.Data = new(common.Empty)
		msgData := core.MakePacket(msg)
		client.conn.Write(msgData)
	}
}


func hook(client *Client) {
	for {
		select {
		case cmd := <-operation:
			opData := <-operateData
			if cmd == "call" {
				chartType := gjson.Get(opData.(string), "chartType").Int()
				receiverId := gjson.Get(opData.(string), "receiverId").Int()
				opStr := opData.(string)
				content := gjson.Get(opStr, "content")
				msg := new(core.Message)
				msg.Version = 1
				if chartType == 2 {
					msg.Cmd = core.GroupMsg
				} else if chartType == 1 {
					msg.Cmd = core.Msg
				}
				data := new(common.MsgIm)
				data.SenderId = client.uid
				data.ReceiverId = receiverId
				data.Content = content.String()
				data.ChartType = int32(chartType)
				msg.Data = data
				byteData := core.MakePacket(msg)
				_, err := client.conn.Write(byteData)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}


func run(client *Client) {

	for {
		msg := core.ReadPacket(client.conn)
		if msg == nil {
			fmt.Println(" conn close ")
			client.conn.Close()
			return
		}
		switch msg.Cmd {
		case core.Pong:
			/*			fmt.Println("ping-pong uid:", server.uid)
						fmt.Println(cover)*/
		case core.GroupMsgAck:
			fmt.Println(msg.Data)
		case core.MsgAck:
			fmt.Println(msg.Data)
		case core.Msg:
			msgData := msg.Data.(*common.MsgIm)
			fmt.Println(msgData)
			contentType := gjson.Get(msgData.Content,"type").String()
			if contentType == "system_transport" {
				fmt.Println(gjson.Parse(msgData.Content))
			}
			msgAck := new(common.MsgImAck)
			msgAck.MsgId = msgData.MsgId
			msg.Cmd = core.MsgAck
			msg.Data = msgAck
			msgBytes := core.MakePacket(msg)
			client.conn.Write(msgBytes)
		case core.GroupMsg:
			msgData := msg.Data.(*common.MsgIm)
			fmt.Println(msgData)
			msgAck := new(common.MsgImAck)
			msgAck.MsgId = msgData.MsgId
			msg.Cmd = core.GroupMsgAck
			msg.Data = msgAck
			msgBytes := core.MakePacket(msg)
			client.conn.Write(msgBytes)
		}
	}
}

type Client struct {
	conn	*net.TCPConn
	uid		int64
}
