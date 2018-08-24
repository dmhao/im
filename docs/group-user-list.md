
**简要描述：**

- 用户群组列表

**请求URL：**
- 格式 ` ip:port/v1/users/:userId/groups?appId=:appId&page=:page `
- 示例 /v1/users/1/groups?appId=1&page=1

**请求方式：**
- Get

**URI参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|userId |是  |string/int | 用户id    |
|appId     |是  |string/int |应用id    |
|page     |是  |string/int |页数 默认1    |


 **返回示例**

 ```
 {
     "Code" : "0",
     "Msg" : "成功",
     "Data" : [
          {
               "GroupId" : "1",
               "AppId" : "1",
               "GroupName" : "测试群组1",
               "GroupDes" : "测试群组",
               "GroupIcon" : "xxxxxxxxxx",
               "MasterUserId" : "1",
               "UserCount" : "2",
               "MaxUserCount" : "200",
               "JoinNeedExamine" : "1",
               "CreateTime" : "1532077897",
               "UpdateTime" : "1532077897",
               "Status" : "1"
          },
     ]
 }

 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |



