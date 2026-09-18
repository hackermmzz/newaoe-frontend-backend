package dao

import (
	"context"
	"encoding/json"
	"newaoe/Src/codeAssessmentSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/util"

	RD "github.com/go-redis/redis/v8"
)

func TeacherGetAll() []model.Teacher {
	key := "teacherAll"
	//先从redis缓存拿数据
	ctx := context.Background()
	res, exist := redis.RedisGet(ctx, key)
	if exist {
		var dataRet []model.Teacher
		json.Unmarshal([]byte(res), &dataRet)
		return dataRet
	}
	//从数据库查询数据
	var teachers []model.Teacher
	err := database.DB.Find(&teachers)
	if err != nil {
		util.DebugError("TeacherGetAll:", err)
		return nil
	}
	//缓存到redis
	js, _ := json.Marshal(teachers)
	redis.RedisSet(ctx, key, js, RD.KeepTTL)
	//返回数据
	return teachers
}

func TeacherExist(name string) bool {
	teacher := TeacherGetAll()
	if len(teacher) == 0 {
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
