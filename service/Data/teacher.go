package data

import (
	"context"
	"encoding/json"
	"newaoe/dao"

	"github.com/go-redis/redis/v8"
)

func TeacherGetAll() ([]dao.Teacher, error) {
	key := "teacherAll"
	//先从redis缓存拿数据
	ctx := context.Background()
	res, exist := dao.RedisGet(ctx, key)
	if exist {
		var dataRet []dao.Teacher
		json.Unmarshal([]byte(res), &dataRet)
		return dataRet, nil
	}
	//从数据库查询数据
	teacher, err := dao.TeacherGetAll()
	if err != nil {
		return nil, err
	}
	//缓存到redis
	js, _ := json.Marshal(teacher)
	dao.RedisSet(ctx, key, js, redis.KeepTTL)
	//返回数据
	return teacher, nil
}

func TeacherExist(name string) bool {
	teacher, err := TeacherGetAll()
	if err != nil {
		return false
	}
	//
	for _, id := range teacher {
		if id.Name == name {
			return true
		}
	}
	//
	return false
}
