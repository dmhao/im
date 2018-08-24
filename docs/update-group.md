
**简要描述：**

- 修改用户群组

**请求URL：**
- 格式 ` ip:port/v1/users/:userId/groups/:groupId?appId=:appId `
- 示例 /v1/users/1/groups/2?appId=1

**请求方式：**
- POST

**URI参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|userId |是  |string/int |用户id  |
|groupId |是  |string/int |群组id   |
|appId |是  |string/int |应用id   |

**POST参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|GroupName |是  |string |群组名   |
|GroupDes |是  |string | 群组介绍    |
|GroupIcon     |是  |string | 群组图标    |
|MaxUserCount     |是  |int | 最大成员数    |
|JoinNeedExamine     |是  |int | 加入是否需要审核  |
|GroupId     |是  |int | 群组id  |
|AppId     |是  |int | 应用id  |

 **返回示例**


 ```
 {
     "Code" : "0",
     "Msg" : "成功",
     "Data" : {
          "GroupId" : "2",
          "AppId" : "1",
          "GroupName" : "你好",
          "GroupDes" : "你好你好你好",
          "GroupIcon" : "xxxxxxxxxx",
          "MasterUserId" : "0",
          "UserCount" : "0",
          "MaxUserCount" : "200",
          "JoinNeedExamine" : "1",
          "CreateTime" : "0",
          "UpdateTime" : "1532316724",
          "Status" : "0"
     }
 }

 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |



