package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/codeCommonSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"
)

func CodeCommonInfoAdd(session *database.Session, data model.CodeCommonInfo) (indices int64, err error) {
	defer func() {
		if err != nil {
			return
		}
		session.PostTransactionAdd(func() {
			keys := []string{
				fmt.Sprintf("CodeCommonGetByIndices:%d", data.Indices),
				fmt.Sprintf("CodeCommonGetByID:%s", data.ID),
			}
			for _, key := range keys {
				if _, cacheErr := redis.RedisDel(context.Background(), key); cacheErr != nil {
					util.DebugError("CodeCommonInfoAdd cache del:", cacheErr)
				}
			}
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

func CodeCommonGetByIndices(session *database.Session, indices int) (ret *model.CodeCommonInfo, err error) {
	key := fmt.Sprintf("CodeCommonGetByIndices:%d", indices)
	useCache := database.CacheAvailable(session)
	if useCache {
		data, exist, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if exist {
			var info model.CodeCommonInfo
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
				util.DebugError("CodeCommonGetByIndices cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("CodeCommonGetByIndices cache set:", cacheErr)
			}
		})
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var info model.CodeCommonInfo
	has, err := session.Where("indices=?", indices).Get(&info)
	if err != nil || !has {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &info, nil
}
func CodeCommonGetByID(session *database.Session, id string) (ret []model.CodeCommonInfo, err error) {
	key := fmt.Sprintf("CodeCommonGetByID:%s", id)
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
				util.DebugError("CodeCommonGetByID cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("CodeCommonGetByID cache set:", cacheErr)
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
