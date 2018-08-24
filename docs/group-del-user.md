
**简要描述：**

- 群组删除成员

**请求URL：**
- 格式 ` ip:port/v1/groups/:groupId/users/:userId?appId=:appId `
- 示例 /v1/groups/1/users/2?appId=1

**请求方式：**
- DELETE

**URI参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|groupId |是  |string/int |群组id   |
|userId |是  |string/int | 用户id    |
|appId     |是  |string/int |应用id    |

 **返回示例**

 ```
 {
     "Code" : "0",
     "Msg" : "成功",
     "Data" : null
 }

 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |



