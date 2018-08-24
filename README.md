# im
Golang写的简单Im目前支持单聊、群聊、api推送消息(文本、图片、语音、视频、自定义透传等格式)
![png](https://img.hacpai.com/pipe/450370050/450370050/450370050/20dd99c5d64043768455e7c8ebef5b0b.png)

![png](https://img.hacpai.com/pipe/450370050/450370050/450370050/c9c0055aa7f2458db61b043ef7906f42.png)


### 演示
目前暂无客户端供演示，以服务端功能开发优化为主
![testpng](https://img.hacpai.com/pipe/450370050/450370050/450370050/eb8c680d86d046a5a7490d40bedadc98.png)

### 导入数据库
导入core文件夹下im.sql

### Route-Server
路由服务用于消息转发

windows 运行 route-server下的exe文件

linux 可以通过 go build route-server.go编译

config.yaml 配置文件
默认配置
```go
dbUser: root    ###数据库账号
dbPassword:     ###数据库密码
dbHost: 127.0.0.1   ###数据库地址
dbPort: 3306    ###数据库端口
dbDatabase: im  ###默认数据库
redisAddr: 127.0.0.1:6379 ###redis地址+端口
redisPassWord:  ###redis密码
routeServerAddr: :3333  ###监听地址
logDir: log/   ###log路径
logExpire: 168  ###log保存周期
rotationFormat: Day  ###log拆分规则格式  Day = %Y-%m-%d
rotationTime: 24  ###log拆分周期  单位Hour
```


### Im-Server 
客户端连接的Im服务

windows 运行 im-server下的exe文件

linux 可以通过 go build im-server.go编译

config.yaml配置文件
```go
dbUser: root    ###数据库账号
dbPassword:     ###数据库密码
dbHost: 127.0.0.1   ###数据库地址
dbPort: 3306    ###数据库端口
dbDatabase: im  ###默认数据库
redisAddr: 127.0.0.1:6379 ###redis地址+端口
redisPassWord:  ###redis密码
clientServerAddr: :7777  ###监听的地址
apiPort: :4444  ###api监听的地址
apiPProf : true  ###是否开启pprof
logDir: log/   ###log路径
logExpire: 168  ###log保存周期
rotationFormat: Day  ###log拆分规则格式  Day = %Y-%m-%d
rotationTime: 24  ###log拆分周期  单位Hour
```

### client
简单测试客户端
windows 运行 client 下的exe文件

linux 可以通过 go build main.go编译

单聊输入call-{"receiverId":2,"chartType":1,"content":"你好用户2"}

群聊输入call-{"receiverId":2,"chartType":2,"content":"你好群组2"}

receiverId为接收者Id 单聊receiverId = userId 群聊receiverId = groupId

chartType为聊天类型  chartType=1单聊  chartType=2群聊


