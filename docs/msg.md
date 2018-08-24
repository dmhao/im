#### 文本消息
```
{
        "type" : "txt",
        "msg" : "hello from rest" //消息内容
}
```


#### 图片消息
```
{
        "type" : "img",   // 消息类型
        "url": "https://a1.easemob.com/easemob-demo/chatdemou.jpg",  //url
        "bounds" : {
          "width" : 480,
          "height" : 720
		}
}
```

#### 语音消息
```
{
        "type": "audio",  // 消息类型
        "length": 10, //时长
		"size": 1233,
		"url": "https://a1.easemob.com/easemob-demo/chatdemoui/chatfiles.mp3",  //返回的url
}
```


#### 视频消息
```
{
        "type": "video",// 消息类型
        "thumb": "https://a1.easemob.com/easemob-demo/chatdemoui.jpg",//成功上传视频缩略图返回的url
        "length": 10,//视频播放长度
        "size": 58103,//视频文件大小
        "url": "https://a1.easemob.com/easemob-demo/chatdemoui/chatfil.mp4"//成功上传视频文件返回的url
}
```
#### 系统级透传消息
```
{
	"type":"system_transport",
	"content": "{\"name\":\"李浩\"}" //自定义的json数据，业务自行处理
}
```

#### 应用级透传消息
```
{
	"type":"transport",
	"content": "{\"name\":\"李浩\"}" //自定义的json数据，业务自行处理
}
```




