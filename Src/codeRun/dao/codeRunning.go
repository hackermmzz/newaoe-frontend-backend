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
	affected, err := session.ID(info.Indices).AllCols().Update(info)
	if err != nil {
		return false, err
	}
	return affected > 0, err
}

// 批量更新状态
func CodeRunningBatchUpdateIfBetter(session *database.Session, data []model.CodeRunningInfo) (better []model.CodeRunningInfo, err error) {
	if len(data) == 0 {
		return nil, nil
	}
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//先获取之前的状态
	indices := make([]int, 0)
	for _, item := range data {
		indices = append(indices, item.Indices)
	}
	oldStatus := make([]model.CodeRunningInfo, 0)
	if err := session.In("indices", indices).ForUpdate().Find(&oldStatus); err != nil {
		return nil, err
	}
	oldStatusMap := make(map[int]model.CodeRunningInfo)
	for _, item := range oldStatus {
		oldStatusMap[item.Indices] = item
	}
	//对比数据，只更新更好的数据
	finalBetter := make([]model.CodeRunningInfo, 0)
	for _, item := range data {
		_, ok := oldStatusMap[item.Indices]
		if ok && item.Status > oldStatusMap[item.Indices].Status {
			finalBetter = append(finalBetter, item)
		}
	}
	//更新数据
	for _, item := range finalBetter {
		_, err = session.Where("indices=?", item.Indices).Cols("status").Update(item)
		if err != nil {
			return nil, err
		}
	}
	//
	session.PostTransactionAdd(func() {
		for _, item := range finalBetter {
			deleteCodeRunningExistCache(item.Indices)
		}
		if len(finalBetter) > 0 {
			incrCodeRunningCacheVersion()
		}
	})
	return finalBetter, nil
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
