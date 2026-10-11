package dao

import (
	"context"
	"encoding/json"
	database "newaoe/Src/databse"
	"newaoe/Src/home/model"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"
)

func AnnouncementAdd(session *database.Session, data model.AnnouncementInfo) (indices int64, err error) {
	defer func() {
		if err != nil {
			return
		}
		session.PostTransactionAdd(func() {
			if _, cacheErr := redis.RedisDel(context.Background(), "AnnouncementGetLatest"); cacheErr != nil {
				util.DebugError("AnnouncementAdd cache del:", cacheErr)
			}
		})
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	// Insert 返回影响行数 + error
	_, err = session.Insert(&data)
	if err != nil {
		return 0, err
	}
	return data.Indices, nil
}

func AnnouncementGetLatest(session *database.Session) (ret *model.AnnouncementInfo, err error) {
	const key = "AnnouncementGetLatest"
	useCache := database.CacheAvailable(session)
	if useCache {
		dataBytes, exist, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if exist {
			var data model.AnnouncementInfo
			if err := json.Unmarshal(dataBytes, &data); err != nil {
				return nil, err
			}
			return &data, nil
		}
	}
	defer func() {
		if err != nil || ret == nil || !useCache {
			return
		}
		session.PostTransactionAdd(func() {
			data, marshalErr := json.Marshal(ret)
			if marshalErr != nil {
				util.DebugError("AnnouncementGetLatest cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("AnnouncementGetLatest cache set:", cacheErr)
			}
		})
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	var data model.AnnouncementInfo

	has, err := session.Where("enabled = ?", true).Desc("indices").Get(&data)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return &data, nil
}
