
**简要描述：**

- 申请入群

**请求URL：**
- ` ip:port/v1/groups/:groupId/examineUsers?appId=:appId `

**请求方式：**
- POST

**URI参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|groupId |是  |string/int |群组id   |
|appId |是  |string/int | 应用id    |


**Post参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|GroupName |是  |string |群组名   |
|UserName |是  |string | 用户名    |

 **返回示例**


 ```
 {
     "Code" : "0",
     "Msg" : "成功",
     "Data" : {
     }
 }

 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |



