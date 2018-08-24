
**简要描述：**

- 所有群组列表

**请求URL：**
- 格式 ` ip:port/v1/groups?appId=:appId&page=:page `

- 示例 /v1/groups?appId=1

**请求方式：**
- Get

**URI参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|appId |是  |string/int |应用Id   |
|page |是  |string/int |页数  默认1   |


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



