
**简要描述：**

- 群组创建

**请求URL：**
- 格式 ` ip:port/v1/users/:userId/groups?appId=:appId `
- 示例 /v1/users/1/groups?appId=1

**请求方式：**
- POST

**URI参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|appId |是  |string/int |应用id   |


**POST参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|GroupName |是  |string |群组名   |
|GroupDes |是  |string | 群组介绍    |
|GroupIcon     |是  |string | 群组图标    |
|MaxUserCount     |是  |int | 最大成员数    |
|JoinNeedExamine     |是  |int | 加入是否需要审核 1需要审核  0不需要审核  |

 **返回示例**


 ```
 {
     "Code" : "0",
     "Msg" : "成功",
     "Data" : {
          "GroupId" : "2",
          "AppId" : "1",
          "GroupName" : "测试群组1",
          "GroupDes" : "测试群组",
          "GroupIcon" : "xxxxxxxxxx",
          "MasterUserId" : "1",
          "UserCount" : "1",
          "MaxUserCount" : "200",
          "JoinNeedExamine" : "1",
          "CreateTime" : "1532316246",
          "UpdateTime" : "1532316246",
          "Status" : "1"
     }
 }

 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|GroupId |int   |群组id  |
|AppId |int   |应用ID  |
|GroupName |string   |群组名  |
|GroupDes |string   |群组介绍  |
|GroupIcon |string   |群组图标  |
|MasterUserId |int   |群主  |
|UserCount |int   |当前用户数  |
|MaxUserCount |int   |最大用户数  |
|JoinNeedExamine |int   |加入群组是否需要 审核  1需审核  0无需审核|
|CreateTime |int   |创建日期  |



