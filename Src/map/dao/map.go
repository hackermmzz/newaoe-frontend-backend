package dao

import (
	"context"
	"encoding/json"
	"fmt"
	database "newaoe/Src/databse"
	"newaoe/Src/map/model"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"
)

type mapCacheData struct {
	Exist bool   `json:"exist"`
	Value string `json:"value"`
}

func mapCacheKey(key string) string {
	return fmt.Sprintf("MapGet:%s", key)
}

func getMapCache(key string) (mapCacheData, bool, error) {
	data, exist, err := redis.RedisGet(context.Background(), mapCacheKey(key))
	if err != nil || !exist {
		return mapCacheData{}, exist, err
	}
	var cache mapCacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		return mapCacheData{}, false, err
	}
	return cache, true, nil
}

func setMapCache(key string, data mapCacheData) {
	value, err := json.Marshal(data)
	if err != nil {
		util.DebugError("Map cache marshal:", err)
		return
	}
	if _, err := redis.RedisSet(context.Background(), mapCacheKey(key), value, 5*time.Minute); err != nil {
		util.DebugError("Map cache set:", err)
	}
}

func deleteMapCache(key string) {
	if _, err := redis.RedisDel(context.Background(), mapCacheKey(key)); err != nil {
		util.DebugError("Map cache del:", err)
	}
}

func MapInsert(session *database.Session, key string, value string) (err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteMapCache(key)
			})
		}
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	info := &model.MapInfo{
		Key:   key,
		Value: value,
	}

	_, err = session.Insert(info)
	if err != nil {
		return err
	}

	return nil
}

func MapDelete(session *database.Session, key string) (err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteMapCache(key)
			})
		}
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	_, err = session.
		Where("`key` = ?", key).
		Delete(&model.MapInfo{})

	return err
}

func MapUpdate(session *database.Session, key string, value string) (err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteMapCache(key)
			})
		}
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	_, err = session.
		Where("`key` = ?", key).
		Cols("value").
		Update(&model.MapInfo{
			Value: value,
		})

	return err
}

func MapExist(session *database.Session, key string) (bool, error) {
	useCache := database.CacheAvailable(session)
	if useCache {
		cache, cached, err := getMapCache(key)
		if err != nil {
			return false, err
		}
		if cached {
			return cache.Exist, nil
		}
	}
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	var info model.MapInfo
	exist, err := session.
		Where("`key` = ?", key).
		Get(&info)

	if err != nil {
		return false, err
	}
	if useCache {
		session.PostTransactionAdd(func() {
			setMapCache(key, mapCacheData{Exist: exist, Value: info.Value})
		})
	}
	return exist, nil
}

func MapGet(session *database.Session, key string) (string, error) {
	useCache := database.CacheAvailable(session)
	if useCache {
		cache, cached, err := getMapCache(key)
		if err != nil {
			return "", err
		}
		if cached {
			if !cache.Exist {
				return "", fmt.Errorf("%v不存在!", key)
			}
			return cache.Value, nil
		}
	}
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	info := &model.MapInfo{}
	exist, err := session.ID(key).Get(info)
	if err != nil {
		return "", err
	}

	if !exist {
		if useCache {
			session.PostTransactionAdd(func() {
				setMapCache(key, mapCacheData{Exist: false})
			})
		}
		return "", fmt.Errorf("%v不存在!", key)
	}

	if useCache {
		session.PostTransactionAdd(func() {
			setMapCache(key, mapCacheData{Exist: true, Value: info.Value})
		})
	}
	return info.Value, nil
}

func MapUpdateOrInsert(session *database.Session, key string, value string) (err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteMapCache(key)
			})
		}
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	exist, err := session.ID(key).Exist(&model.MapInfo{})
	if err != nil {
		return err
	}

	if exist {
		_, err = session.
			ID(key).
			Cols("value").
			Update(&model.MapInfo{
				Value: value,
			})
		return err
	}

	_, err = session.Insert(&model.MapInfo{
		Key:   key,
		Value: value,
	})
	return err
}
