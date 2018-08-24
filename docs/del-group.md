
**简要描述：**

- 删除群组

**请求URL：**
- 格式 ` ip:port/v1/users/:userId/groups/:groupId?appId=:appId `
- 示例 /v1/users/1/groups/3?appId=1

**请求方式：**
- DELETE

**URI参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|userId |是  |string/int |用户id  |
|groupId |是  |string/int |群组id   |
|appId |是  |string/int |应用id   |

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



