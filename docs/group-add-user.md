
**简要描述：**

- 群组批量添加成员

**请求URL：**
- 格式 ` ip:port/v1/groups/:groupId/users?appId=:appId`
- 示例 /v1/groups/1/users?appId=1

**请求方式：**
- POST

**URI参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|groupId |是  |string/int |群组id   |
|appId |是  |string/int |应用id   |

**POST参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|UserIds |是  |string | 用户id,多个半角逗号隔开    |

 **返回示例**


 ```
 {
     "Code" : "0",
     "Msg" : "成功",
     "Data" : [
		 {
			  "Id" : "11",
			  "AppId" : "1",
			  "GroupId" : "2",
			  "UserId" : "3",
			  "UserRole" : "0",
			  "JoinTime" : "1532318426",
			  "Status" : "1"
		 },
	 ]
 }

 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|UserId |int   |用户Id  由业务系统获取用户详情  |
|GroupId |int   |群组ID  |
|UserRole |int   |用户级别  0普通用户  1群主  2管理员  |
|JoinTime |int   |加入时间 |



