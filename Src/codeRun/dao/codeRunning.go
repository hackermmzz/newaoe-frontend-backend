package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"
)

const codeRunningCacheVersionKey = "CodeRunningCacheVersion"

func codeRunningExistCacheKey(indices int) string {
	return fmt.Sprintf("CodeRunningExist:%d", indices)
}

func deleteCodeRunningExistCache(indices int) {
	if _, err := redis.RedisDel(context.Background(), codeRunningExistCacheKey(indices)); err != nil {
		util.DebugError("CodeRunningExist cache del:", err)
	}
}

func getCodeRunningCacheVersion() (int, error) {
	return redis.RedisGetVersion(context.Background(), codeRunningCacheVersionKey)
}

func incrCodeRunningCacheVersion() {
	if _, err := redis.RedisIncrVersion(context.Background(), codeRunningCacheVersionKey); err != nil {
		util.DebugError("CodeRunning cache incrVersion:", err)
	}
}

func CodeRunningExist(session *database.Session, indices int) (ret bool, err error) {
	key := codeRunningExistCacheKey(indices)
	useCache := database.CacheAvailable(session)
	if useCache {
		data, exist, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return false, cacheErr
		}
		if exist {
			if err := json.Unmarshal(data, &ret); err != nil {
				return false, err
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
				util.DebugError("CodeRunningExist cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("CodeRunningExist cache set:", cacheErr)
			}
		})
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var info model.CodeRunningInfo
	has, err := session.ID(indices).Get(&info)
	if err != nil {
		return false, err
	}
	return has, nil
}

func CodeRunningInsert(session *database.Session, data model.CodeRunningInfo) (ok bool, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteCodeRunningExistCache(data.Indices)
				incrCodeRunningCacheVersion()
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err = session.Insert(data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CodeRunningUpdate(session *database.Session, info model.CodeRunningInfo) (ok bool, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteCodeRunningExistCache(info.Indices)
				incrCodeRunningCacheVersion()
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	affected, err := session.ID(info.Indices).Update(info)
	if err != nil {
		return false, err
	}
	return affected > 0, err
}

func CodeRunningRemove(session *database.Session, indices int) (ok bool, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteCodeRunningExistCache(indices)
				incrCodeRunningCacheVersion()
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	d := &model.CodeRunningInfo{
		Indices: indices,
	}
	_, err = session.Delete(d)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CodeRunningBatchRemove(session *database.Session, indices []int) (ok bool, err error) {
	cacheIndices := append([]int(nil), indices...)
	defer func() {
		if err != nil {
			return
		}
		session.PostTransactionAdd(func() {
			for _, index := range cacheIndices {
				deleteCodeRunningExistCache(index)
			}
			if len(cacheIndices) > 0 {
				incrCodeRunningCacheVersion()
			}
		})
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if len(indices) == 0 {
		return true, nil
	}
	//
	_, err = session.In("indices", indices).Delete(&model.CodeRunningInfo{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func CodeRunningGetExpireTime(session *database.Session, expireDuration time.Duration, number int) ([]model.CodeRunningInfo, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var result []model.CodeRunningInfo
	// 当前时间减去超时时间
	expireTime := time.Now().Add(-expireDuration)
	err := session.
		Where("submittime <= ?", expireTime).
		Limit(number). // 最多返回 number 条
		Find(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 不包含end
func CodeRunningGetByRange(session *database.Session, beg int, end int) (ret []model.CodeRunningInfo, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		version, cacheErr := getCodeRunningCacheVersion()
		if cacheErr != nil {
			return nil, cacheErr
		}
		key = fmt.Sprintf("CodeRunningGetByRange:%d:%d:%d", version, beg, end)
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
				util.DebugError("CodeRunningGetByRange cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("CodeRunningGetByRange cache set:", cacheErr)
			}
		})
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	if beg >= end || beg < 0 {
		return nil, nil
	}

	err = session.
		Asc("indices").
		Limit(end-beg, beg).
		Find(&ret)

	if err != nil {
		return nil, err
	}
	return ret, nil
}
