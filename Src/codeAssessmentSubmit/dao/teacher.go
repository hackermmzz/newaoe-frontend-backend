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

func TeacherAdd(session *database.Session, teacher string) (err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				keys := []string{
					"TeacherGetAll",
					fmt.Sprintf("TeacherExist:%v", teacher),
				}
				for _, key := range keys {
					if _, cacheErr := redis.RedisDel(context.Background(), key); cacheErr != nil {
						util.DebugError("TeacherAdd cache del:", cacheErr)
					}
				}
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data := model.Teacher{
		Name: teacher,
	}
	_, err = session.Insert(&data)
	if err != nil {
		return err
	}
	//
	return nil
}

func TeacherGetAll(session *database.Session) (ret []model.Teacher, err error) {
	//查缓存
	useCache := database.CacheAvailable(session)
	if useCache {
		dataBytes, exist, cacheErr := redis.RedisGet(context.Background(), "TeacherGetAll")
		if cacheErr != nil {
			return nil, cacheErr
		}
		if exist {
			var cached []model.Teacher
			if err := json.Unmarshal(dataBytes, &cached); err == nil {
				return cached, nil
			} else {
				return nil, err
			}
		}
	}
	//写入缓存
	defer func() {
		if err != nil || !useCache {
			return
		}
		session.PostTransactionAdd(func() {
			db, marshalErr := json.Marshal(ret)
			if marshalErr != nil {
				util.DebugError("TeacherGetAll cache marshal:", marshalErr)
				return
			}
			//写入缓存
			if _, err := redis.RedisSet(context.Background(), "TeacherGetAll", db, time.Duration(1)*time.Hour); err != nil {
				util.DebugError("TeacherGetAll cache set:", err)
			}
		})
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//从数据库查询数据
	var teachers []model.Teacher
	err = session.Find(&teachers)
	if err != nil {
		return nil, err
	}
	return teachers, nil
}

func TeacherExist(session *database.Session, name string) (ret bool, err error) {
	//查缓存
	key := fmt.Sprintf("TeacherExist:%v", name)
	useCache := database.CacheAvailable(session)
	if useCache {
		v, ext, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return false, cacheErr
		}
		if ext {
			if err := json.Unmarshal(v, &ret); err == nil {
				return ret, nil
			} else {
				return false, err
			}
		}
	}
	//
	defer func() {
		if err != nil || !useCache {
			return
		}
		session.PostTransactionAdd(func() {
			//写入缓存
			if _, err := redis.RedisSet(context.Background(), key, ret, time.Duration(5)*time.Minute); err != nil {
				util.DebugError("TeacherExist cache set:", err)
			}
		})
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	teacher, err := TeacherGetAll(session)
	if err != nil {
		return false, err
	}
	if len(teacher) == 0 {
		return false, nil
	}
	//
	for _, id := range teacher {
		if id.Name == name {
			return true, nil
		}
	}
	//
	return false, nil
}
