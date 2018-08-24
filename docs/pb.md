```
syntax = "proto3";
package common;

/*
    用于登录
    Uid             用户uid
    Token           应用的token
    AppId           当前应用ID
    DeviceType      设备类型
*/
message LoginIm {
    int64       Uid = 1;
    string      Token = 2;
    int32       AppId = 3;
    int32       DeviceType = 4;
}


/*
    退出登录
    Uid             用户uid
    Token           应用的token
*/
message LogoutIm {
    int64   Uid = 1;
    string  Token = 2;
}

/*
    登录的返回消息
    Status  登录状态
    Code   返回码
    Data   返回的内容
*/
message AuthIm {
    int32   Status = 1;
    int32   Code = 2;
    string  Data = 3;
}


/*
    单聊群聊消息体
    SenderId            发送方的id
    ReceiverId          接收方的id
    ChartType           (1 单聊  2群聊)
    MsgType             消息类型 (1 用户消息 2系统消息)
    MsgId               消息ID  由服务端填充
    TalkId              会话ID  第一次发送消息，会话ID由服务端填充。
    TraceId             消息调试ID
    Timestamp           消息的发送时间，  毫秒级时间戳
    Content             消息内容
*/
message MsgIm {
    int64   SenderId = 1;
    int64   ReceiverId = 2;
    int32   ChartType = 3;
    int32   MsgType = 4;
    string  MsgId = 5;
    string  TalkId = 6;
    string  TraceId = 7;
    int64   Timestamp = 8;
    string  Content = 9;
}

/*
    消息的Ack结构
    TraceId     消息调试ID
    MsgId       消息的id
    TalkId      会话的id
    Timestamp   消息发送时间
*/
message  MsgImAck {
    string  TraceId = 1;
    string  MsgId = 2;
    string  TalkId = 3;
    int64   Timestamp = 4;
}

/*
    用于ping  pong 心跳时的空结构
*/
message Empty {

}


/*
    同步离线消息
    Limit       同步条数，  Limit为0  同步所有离线消息
*/
message SyncOffline {
    int32   Limit = 1;
}


/*
    离线消息返回
    MsgList         返回的消息数组
    SurplusCount    剩余未同步的消息个数
*/
message SyncOfflineMsg {
    repeated    MsgIm   MsgList = 1;
    int32   SurplusCount = 2;
}
```