
**简要描述：**

- 上传音频

**请求URL：**
- ` ip:port/v1/audios?appId=:appId`

**请求方式：**
- POST

**参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|AudioType |是  |string |上传图片的类型  message(聊天音频) |
|Audio |是  |File | 上传的文件    |

 **返回示例**

 ```
{
    "Code": 0,
    "Msg": "成功",
    "Data": {
        "Length": 45.852281,
        "Size": 706939,
        "Url": "http:/qn-audios.articlechain.cn/ddd/2018-08-08/1533712665379793900.aac"
    }
}
 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|Url |string   |上传后返回的地址 |
|Length |float   |时长 |
|Size |int   |文件大小  |



