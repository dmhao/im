```
	//心跳   请求ping   返回pong
	Ping    	= 1
	Pong    	= 2

	//登录  请求Login   返回Auth
	//退出请求 Logout
	Login   	= 3
	Logout  	= 4
	Auth    	= 5

	//单聊  请求Msg    返回MsgAck
	Msg     	= 6
	MsgAck     	= 7

	//群聊  请求GroupMsg   返回GroupMsgAck
	GroupMsg 	= 8
	GroupMsgAck = 9

	//离线消息同步   请求SyncOffline       返回SyncOfflineMsg   离线消息的ack  SyncOfflineMsgAck
	SyncOffline = 10
	SyncOfflineMsg = 11
	SyncOfflineMsgAck = 12

```