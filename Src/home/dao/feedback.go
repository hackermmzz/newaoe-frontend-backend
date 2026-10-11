package dao

import (
	"context"
	"encoding/json"
	"fmt"
	database "newaoe/Src/databse"
	"newaoe/Src/home/model"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"
)

const feedbackCacheVersionKey = "FeedbackCacheVersion"

func getFeedbackCacheVersion() (int, error) {
	return redis.RedisGetVersion(context.Background(), feedbackCacheVersionKey)
}

func incrFeedbackCacheVersion() {
	if _, err := redis.RedisIncrVersion(context.Background(), feedbackCacheVersionKey); err != nil {
		util.DebugError("Feedback cache incrVersion:", err)
	}
}

func FeedbackGetByIndices(session *database.Session, indices int) (ret *model.FeedbackInfo, err error) {
	key := fmt.Sprintf("FeedbackGetByIndices:%d", indices)
	useCache := database.CacheAvailable(session)
	if useCache {
		data, exist, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if exist {
			var info model.FeedbackInfo
			if err := json.Unmarshal(data, &info); err != nil {
				return nil, err
			}
			return &info, nil
		}
	}
	defer func() {
		if err != nil || ret == nil || !useCache {
			return
		}
		session.PostTransactionAdd(func() {
			data, marshalErr := json.Marshal(ret)
			if marshalErr != nil {
				util.DebugError("FeedbackGetByIndices cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("FeedbackGetByIndices cache set:", cacheErr)
			}
		})
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var info model.FeedbackInfo
	has, err := session.Where("indices=?", indices).Get(&info)
	if err != nil || !has {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &info, nil
}

func FeedbackInsert(session *database.Session, info model.FeedbackInfo) (indices int, err error) {
	defer func() {
		if err != nil {
			return
		}
		session.PostTransactionAdd(func() {
			key := fmt.Sprintf("FeedbackGetByIndices:%d", info.Indices)
			if _, cacheErr := redis.RedisDel(context.Background(), key); cacheErr != nil {
				util.DebugError("FeedbackInsert cache del:", cacheErr)
			}
			incrFeedbackCacheVersion()
		})
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	// Insert 返回影响行数 + error
	_, err = session.Insert(&info)
	if err != nil {
		return 0, err
	}
	return info.Indices, nil
}

// 不包含end
func FeedbackGetRange(session *database.Session, beg int, end int) (ret []model.FeedbackInfo, err error) {
	if beg < 0 || beg >= end {
		return nil, util.NewError("beg<0 or beg>=end!")
	}

	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		version, cacheErr := getFeedbackCacheVersion()
		if cacheErr != nil {
			return nil, cacheErr
		}
		key = fmt.Sprintf("FeedbackGetRange:%d:%d:%d", version, beg, end)
		data, exist, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if exist {
			if err := json.Unmarshal(data, &ret); err != nil {
				return nil, err
			}
			return ret, nil
		}
	}
	defer func() {
		if err != nil || !useCache {
			return
		}
		session.PostTransactionAdd(func() {
			data, marshalErr := json.Marshal(ret)
			if marshalErr != nil {
				util.DebugError("FeedbackGetRange cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("FeedbackGetRange cache set:", cacheErr)
			}
		})
	}()

	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	err = session.OrderBy("indices desc").
		Limit(end-beg, int(beg)).
		Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
