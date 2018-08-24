package storage

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"im/core/log"
	"im/mp/common"
	"strconv"
	"im/core"
)

var redisClient *redis.Client

func NewRedis(addr, password string) error {
	option := &redis.Options{Addr: addr, DB: 0, Password:password}
	redisClient = redis.NewClient(option)
	err := redisClient.Ping().Err()
	return err
}

func StopRedis() {
	redisClient.Close()
}

//循环群组所有用户 通过redis pipeline批量增加im消息队列数据
func AppendUserGroupImMsg(appId int, groupUsers []*GroupUser, msg *common.MsgIm) {
	msgData, err := json.Marshal(msg)
	if err != nil {
		log.Infoln("消息编码失败", err, msg)
		return
	}
	p := redisClient.Pipeline()
	for _, groupUser := range groupUsers {
		userMessageLK := fmt.Sprintf(userMessageListCacheSet.key, appId, groupUser.UserId)
		userMessageHk := fmt.Sprintf(userMessageHashCacheSet.key, appId, groupUser.UserId)
		p.RPush(userMessageLK, msg.MsgId)
		if userMessageListCacheSet.expire != 0 {
			p.Expire(userMessageLK, userMessageListCacheSet.expire)
		}
		p.HSet(userMessageHk, msg.MsgId, msgData)
		if userMessageHashCacheSet.expire != 0 {
			p.Expire(userMessageHk, userMessageHashCacheSet.expire)
		}
	}
	p.Exec()
}

//增加im消息队列数据
func AppendUserImMsg(appId int, userId int64, msg *common.MsgIm) {
	msgData, err := json.Marshal(msg)
	if err != nil {
		log.Infoln("消息编码失败", err, msg)
		return
	}

	p := redisClient.Pipeline()
	userMessageLK := fmt.Sprintf(userMessageListCacheSet.key, appId, userId)
	p.RPush(userMessageLK, msg.MsgId).Err()
	if userMessageListCacheSet.expire != 0 {
		p.Expire(userMessageLK, userMessageListCacheSet.expire)
	}

	userMessageHk := fmt.Sprintf(userMessageHashCacheSet.key, appId, userId)
	p.HSet(userMessageHk, msg.MsgId, msgData).Err()
	if userMessageHashCacheSet.expire != 0 {
		p.Expire(userMessageHk, userMessageHashCacheSet.expire)
	}
	p.Exec()
}

//删除用户消息队列某条数据
func RemoveUserImMsg(appId int, userId int64, msgId string) {
	p := redisClient.Pipeline()
	userMessageLK := fmt.Sprintf(userMessageListCacheSet.key, appId, userId)
	p.LRem(userMessageLK, 0, msgId).Err()

	userMessageHk := fmt.Sprintf(userMessageHashCacheSet.key, appId, userId)
	p.HDel(userMessageHk, msgId).Err()
	p.Exec()
}

//不存在就创建会话Id 并 更新会话时间
func UpdateOrCreateTalkId(appId int, msg *common.MsgIm) {
	var talkId string
	if msg.TalkId == "" {
		if msg.ChartType == core.GroupImChartType {
			talkId = "g_" + strconv.FormatInt(msg.ReceiverId, 10)
		} else if msg.ChartType == core.ImChartType {
			var startId, endId int64
			if msg.SenderId > msg.ReceiverId {
				startId = msg.SenderId
				endId = msg.ReceiverId
			} else {
				startId = msg.ReceiverId
				endId = msg.SenderId
			}
			talkId = "c_" + strconv.FormatInt(startId, 10) + "_" + strconv.FormatInt(endId, 10)
		}
	} else {
		talkId = msg.TalkId
	}
	allTalkIdZSetK := fmt.Sprintf(allTalkIdZSetCacheSet.key, appId)
	redisClient.ZAdd(allTalkIdZSetK, redis.Z{Score: float64(msg.Timestamp), Member: talkId}).Err()
	msg.TalkId = talkId
	if allTalkIdZSetCacheSet.expire != 0 {
		redisClient.Expire(allTalkIdZSetK, allTalkIdZSetCacheSet.expire)
	}
}

//更新用户的会话聊表 聊天时间
func UpdateUserTalkIdTime(appId int, msg *common.MsgIm) {
	if msg.ChartType == core.GroupImChartType {
		tUserTalkIdZSetK := fmt.Sprintf(userTalkIdZSetCacheSet.key, appId, msg.SenderId)
		redisClient.ZAdd(tUserTalkIdZSetK, redis.Z{Score: float64(msg.Timestamp), Member: msg.TalkId})
		if userTalkIdZSetCacheSet.expire != 0 {
			redisClient.Expire(tUserTalkIdZSetK, userTalkIdZSetCacheSet.expire)
		}

	} else if msg.ChartType == core.ImChartType {
		rUserTalkIdZSetK := fmt.Sprintf(userTalkIdZSetCacheSet.key, appId, msg.ReceiverId)
		redisClient.ZAdd(rUserTalkIdZSetK, redis.Z{Score: float64(msg.Timestamp), Member: msg.TalkId})

		sUserTalkIdZSetK := fmt.Sprintf(userTalkIdZSetCacheSet.key, appId, msg.SenderId)
		redisClient.ZAdd(sUserTalkIdZSetK, redis.Z{Score: float64(msg.Timestamp), Member: msg.TalkId})

		if userTalkIdZSetCacheSet.expire != 0 {
			redisClient.Expire(rUserTalkIdZSetK, userTalkIdZSetCacheSet.expire)
			redisClient.Expire(sUserTalkIdZSetK, userTalkIdZSetCacheSet.expire)
		}
	}
}

func RemoveUserOfflineMsg(appId int, userId int64, ackMsgId string) {
	//获取所有的消息倒序列表数据
	userMessageLK := fmt.Sprintf(userMessageListCacheSet.key, appId, userId)
	allMsgIds, err := redisClient.LRange(userMessageLK, 0, -1).Result()
	if err != nil {
		log.Warnln("用户离线消息队列读取失败", err, "userId", userId)
		return
	}
	msgIndex := 0
	findIndex := false
	//寻找消息  找到消息时记录对应的下标
	msgIdsCount := len(allMsgIds)
	if msgIdsCount > 0 {
		for i := msgIdsCount; i >= 0; i -- {
			if allMsgIds[i-1] == ackMsgId {
				msgIndex = i
				findIndex = true
				break
			}
		}
		//找到ack消息的下标  删除队列  和  消息hash实体中的数据
		if findIndex {
			userMessageHk := fmt.Sprintf(userMessageHashCacheSet.key, appId, userId)

			p := redisClient.Pipeline()
			p.LTrim(userMessageLK, int64(msgIndex), -1)
			p.HDel(userMessageHk, allMsgIds[:msgIndex]...)
			p.Exec()
		}
	}
}

//获取用户的离线消息
func GetUserOfflineMsg(appId int, userId int64, limit int) ([]*common.MsgIm, int) {
	userMessageLK := fmt.Sprintf(userMessageListCacheSet.key, appId, userId)
	offlineMsgCount := redisClient.LLen(userMessageLK).Val()

	var surplusCount = 0
	var validMsgIds []string
	var allOfflineMsg []*common.MsgIm

	//离线消息 大于 0
	if offlineMsgCount > 0 {
		//获取所有的消息列表数据
		allMsgIds, err := redisClient.LRange(userMessageLK, 0, offlineMsgCount).Result()
		if err != nil {
			log.Warnln("用户离线消息队列读取失败", err, "userId", userId)
			return allOfflineMsg, surplusCount
		}
		var readNum = 0
		for _, msgId := range allMsgIds {
			//limit 不等于0时，读到的个数  和  limit相等时跳出
			if readNum == limit && limit != 0 {
				break
			}
			validMsgIds = append(validMsgIds, msgId)
			readNum++
		}

		if len(validMsgIds) > 0 {
			userMessageHk := fmt.Sprintf(userMessageHashCacheSet.key, appId, userId)
			msgSlice := redisClient.HMGet(userMessageHk, validMsgIds...).Val()
			if err != nil {
				log.Warnln("通过消息队列读取消息实体失败", err, "userId", userId)
				return allOfflineMsg, surplusCount
			}

			for _, msgData := range msgSlice {
				msgIm := &common.MsgIm{}
				if msgData != nil {
					err := json.Unmarshal([]byte(msgData.(string)), msgIm)
					if err == nil {
						allOfflineMsg = append(allOfflineMsg, msgIm)
					}
				}
			}
		}
		// redis列表中总数  -  有效读取的个数  - 读取的条数
		surplusCount = int(offlineMsgCount) - int(limit)
	}

	if surplusCount < 0 {
		surplusCount = 0
	}
	return allOfflineMsg, surplusCount
}

func GetAppInfoCacheByAppId(appId int) (*App, error) {
	app := &App{}
	appInfoCk := fmt.Sprintf(appInfoKvCacheSet.key, appId)
	appInfoStr := redisClient.Get(appInfoCk).Val()
	err := json.Unmarshal([]byte(appInfoStr), app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func SetAppInfoCacheByAppId(appId int, app *App) {
	appInfoCk := fmt.Sprintf(appInfoKvCacheSet.key, appId)
	appBytes, err := json.Marshal(app)
	if err == nil {
		redisClient.Set(appInfoCk, string(appBytes), appInfoKvCacheSet.expire)
	}
}

func GetAppRouteServiceCacheByAppId(appId int) ([]*RouteService, error) {
	var routeServices []*RouteService
	routeServiceCk := fmt.Sprintf(appRouteServiceKvCacheSet.key, appId)
	routeServiceStr := redisClient.Get(routeServiceCk).Val()
	err := json.Unmarshal([]byte(routeServiceStr), routeServices)
	if err != nil {
		return routeServices, err
	}
	return routeServices, nil
}

func SetAppRouteServiceCacheByAppId(appId int, routeServices []*RouteService) {
	routeServiceCk := fmt.Sprintf(appRouteServiceKvCacheSet.key, appId)
	routeServiceBytes, err := json.Marshal(routeServices)
	if err == nil {
		redisClient.Set(routeServiceCk, string(routeServiceBytes), appRouteServiceKvCacheSet.expire)
	}
}


func GetUserJoinGroupIdCache(appId int, userId int64) ([]int, error) {
	joinGroupIdCk := fmt.Sprintf(userJoinGroupIdKvCacheSet.key, appId, userId)
	joinGroupIdCkStr := redisClient.Get(joinGroupIdCk).Val()

	var groupIds []int
	err := json.Unmarshal([]byte(joinGroupIdCkStr), &groupIds)
	if err != nil {
		return groupIds, err
	}
	return groupIds, nil
}

func SetUserJoinGroupIdCache(appId int, userId int64, groupIds []int) {
	joinGroupIdCk := fmt.Sprintf(userJoinGroupIdKvCacheSet.key, appId, userId)
	joinGroupIdBytes, err := json.Marshal(groupIds)
	if err == nil {
		redisClient.Set(joinGroupIdCk, string(joinGroupIdBytes), userJoinGroupIdKvCacheSet.expire)
	}
}
