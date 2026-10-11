package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/codeAssessmentSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"
)

const codeAssessmentCacheVersionKey = "CodeAssessmentCacheVersion"

func getCodeAssessmentCacheVersion() (int, error) {
	return redis.RedisGetVersion(context.Background(), codeAssessmentCacheVersionKey)
}

func incrCodeAssessmentCacheVersion() {
	if _, err := redis.RedisIncrVersion(context.Background(), codeAssessmentCacheVersionKey); err != nil {
		util.DebugError("CodeAssessment cache incrVersion:", err)
	}
}

func CodeAssessmentAdd(session *database.Session, data model.CodeAssessmentInfo) (indices int64, err error) {
	defer func() {
		if err != nil {
			return
		}
		session.PostTransactionAdd(func() {
			keys := []string{
				fmt.Sprintf("CodeAssessmentGetByIndices:%d", data.Indices),
				fmt.Sprintf("CodeAssessmentGetByID:%s", data.ID),
			}
			for _, key := range keys {
				if _, cacheErr := redis.RedisDel(context.Background(), key); cacheErr != nil {
					util.DebugError("CodeAssessmentAdd cache del:", cacheErr)
				}
			}
			incrCodeAssessmentCacheVersion()
		})
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	_, err = session.Insert(&data)
	if err != nil {
		return 0, err
	}
	return int64(data.Indices), nil
}

func CodeAssessmentGetByID(session *database.Session, id string) (ret []model.CodeAssessmentInfo, err error) {
	key := fmt.Sprintf("CodeAssessmentGetByID:%s", id)
	useCache := database.CacheAvailable(session)
	if useCache {
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
				util.DebugError("CodeAssessmentGetByID cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("CodeAssessmentGetByID cache set:", cacheErr)
			}
		})
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	err = session.Where("id = ?", id).Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func CodeAssessmentGetByIndices(session *database.Session, indices int) (ret *model.CodeAssessmentInfo, err error) {
	key := fmt.Sprintf("CodeAssessmentGetByIndices:%d", indices)
	useCache := database.CacheAvailable(session)
	if useCache {
		data, exist, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if exist {
			var info model.CodeAssessmentInfo
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
				util.DebugError("CodeAssessmentGetByIndices cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("CodeAssessmentGetByIndices cache set:", cacheErr)
			}
		})
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var info model.CodeAssessmentInfo
	has, err := session.Where("indices = ?", indices).Get(&info)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &info, nil
}

// 不包含end
func CodeAssessmentGetByRange(session *database.Session, beg int, end int) (ret []model.CodeAssessmentInfo, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		version, cacheErr := getCodeAssessmentCacheVersion()
		if cacheErr != nil {
			return nil, cacheErr
		}
		key = fmt.Sprintf("CodeAssessmentGetByRange:%d:%d:%d", version, beg, end)
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
				util.DebugError("CodeAssessmentGetByRange cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("CodeAssessmentGetByRange cache set:", cacheErr)
			}
		})
	}()
	//

	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	if beg < 0 || beg >= end {
		return nil, util.NewError("beg<0 or beg>=end!")
	}

	err = session.
		Limit(end-beg, beg).
		Find(&ret)
	if err != nil {
		return nil, util.NewError("CodeAssessmentGetByRange:", err)
	}
	return ret, nil
}
