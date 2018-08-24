package storage

import "time"

type CacheSet struct {
	key    string
	expire time.Duration
}

var allTalkIdZSetCacheSet = &CacheSet{"allTalkIdZset:appId_%v", 0}

var userMessageHashCacheSet = &CacheSet{"userMessageHash:appId_%v:%v", 3600 * 24 * 30 * time.Second}

var userMessageListCacheSet = &CacheSet{"userMessageList:appId_%v:%v", 3600 * 24 * 30 * time.Second}

var userTalkIdZSetCacheSet = &CacheSet{"userTalkIdZSet:appId_%v:%v", 3600 * 24 * 90 * time.Second}

var appInfoKvCacheSet = &CacheSet{"appInfo:appId_%v", 3600 * 24 * time.Second}

var appRouteServiceKvCacheSet = &CacheSet{"routeService:appId_%v", 3600 * time.Second}

var userJoinGroupIdKvCacheSet = &CacheSet{"userJoinGroup:appId_%v:userId_%v", 3600 * 24 * time.Second}