#### 申请入群审核消息

 ```
 {
     "type" : "system_transport",
     "content" : "{"type":"group_examine","data_type":"button","data":{"apply_user_id":12,"apply_user_name":"李浩","group_id":1,"group_name":"测试组","msg_des":"李浩-申请加入群-测试组","content":[{"click_type":"url","text":"同意","jump_data":"/groups/1/examineUsers/12?appId=1&examineId=12&op=agree"},{"click_type":"url","text":"拒绝","jump_data":"/groups/1/examineUsers/12?appId=1&examineId=12&op=refuse"}]}}"
 }
 ```

 ```
 {
     "type" : "join_group_examine",
     "data_type" : "button",
     "data" : {
          "apply_user_id" : "12",
          "apply_user_name" : "李浩",
          "group_id" : "1",
          "group_name" : "测试组",
          "intro" : "李浩-申请加入群-测试组",
          "buttons" : [
               {
                    "click_type" : "url",
                    "text" : "同意",
                    "jump_data" : "/groups/1/examineUsers/12?appId=1&examineId=12&op=agree"
               },
               {
                    "click_type" : "url",
                    "text" : "拒绝",
                    "jump_data" : "/groups/1/examineUsers/12?appId=1&examineId=12&op=refuse"
               },
          ]
     }
 }

 ```

 |参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|type |是  |string |join_group_examine 加入群审核   |
|data_type |是  |string |button 按钮类型   |
|apply_user_id |是  |string | 申请人用户id  |
|apply_user_name |是  |string | 申请人用户名   |
|group_id |是  |string | 群组id  |
|group_name |是  |string | 群组名  |、
|intro |是  |string |按钮简介   |
|buttons |是  |array |按钮列表   |
|click_type |是  |string |url链接类型   |
|text |是  |string |按钮文字   |
|jump_data |是  |string |跳转数据   |