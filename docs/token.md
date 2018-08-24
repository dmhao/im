
**简要描述：**

- token获取(临时接口,后面会改动)
<font color=red>
请求其它api接口时，Header中设置此参数，用作验证
Token : xxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxxxxx
</font>

**请求URL：**
- 格式： ` ip:port/token?secretId=:secretId&secretKey=:secretKey&userId=:userId `
- 示例： /token?secretId=asddsa&secretKey=asddsa&userId=1

**请求方式：**
- GET

**URI参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|secretId    |是| string | 应用的secretId  测试：asddsa |
|secretKey   |是| string | 应用的secretKey  测试：asddsa |
|userId   |是| string | 用户的id   |


 **返回示例**

 ```
 {
     "Code" : "0",
     "Msg" : "成功",
     "Data" : {
	 	  "AppId" : 1,
          "Token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzIzMTkzODIsImlzcyI6ImhqaF9pbSJ9.gfc-Bk0QWNwyYk3qW-IMBM09SWkd_8CZ7tFtB-YEo8Q"
     }
 }

 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|token |string   |请求im登录时填写的token  |



