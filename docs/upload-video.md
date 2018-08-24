
**简要描述：**

- 上传视频

**请求URL：**
- ` ip:port/v1/videos?appId=:appId`

**请求方式：**
- POST

**参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|VideoType |是  |string |上传图片的类型  message(聊天视频) |
|Video |是  |File | 上传的文件    |

 **返回示例**

 ```
{
    "Code": 0,
    "Msg": "成功",
    "Data": {
        "Thumb": "http:/qn-images.articlechain.cn/video-thumb/dsadsa/2018-08-08/1533713551254579834.jpg",
        "Length": 4.395011,
        "Size": 295906,
        "Url": "http:/qn-videos.articlechain.cn/dsadsa/2018-08-08/1533713551254579834.mp4"
    }
}
 ```



 **返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|Url |string   |上传后返回的地址 |
|Thumb |string   |缩略图地址 |
|Length |float   |时长 |
|Size |int   |文件大小  |



