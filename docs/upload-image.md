
**简要描述：**

- 上传图片

**请求URL：**
- ` ip:port/v1/images?appId=:appId`

**请求方式：**
- POST

**参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|ImageType |是  |string |上传图片的类型  message(聊天图片)  group-icon （群组头像）  |
|Image |是  |File | 上传的文件    |

 **返回示例**

 ```
 {
     "Code" : "0",
     "Msg" : "成功",
     "Data" : {
          "Url" : "https:/im-images-oss.oss-cn-beijing.aliyuncs.com/test/2018-08-07/15336338780554167001242x2208.jpg",
          "Width" : "1242",
          "Height" : "2208"
     }
 }

 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|Url |string   |上传后返回的地址 |
|Width |int   |图片宽 |
|Height |int   |图片高  |



