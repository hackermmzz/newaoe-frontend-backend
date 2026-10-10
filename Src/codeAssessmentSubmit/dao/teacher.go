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

	"xorm.io/xorm"
)

func TeacherAdd(session *xorm.Session, teacher string) error {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data := model.Teacher{
		Name: teacher,
	}
	_, err := session.Insert(&data)
	if err != nil {
		return err
	}
	//删缓存
	redis.RedisDel(context.Background(), "TeacherGetAll")
	//
	return nil
}

func TeacherGetAll(session *xorm.Session) []model.Teacher {
	//查缓存
	dataBytes, exist := redis.RedisGet(context.Background(), "TeacherGetAll")
	if exist {
		var ret []model.Teacher
		err := json.Unmarshal(dataBytes, &ret)
		if err == nil {
			return ret
		}
		util.DebugError("TeacherGetAll", err)
	}
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//从数据库查询数据
	var teachers []model.Teacher
	err := session.Find(&teachers)
	if err != nil {
		util.DebugError("TeacherGetAll:", err)
		return nil
	}
	//返回数据写入缓存
	db, _ := json.Marshal(teachers)
	redis.RedisSet(context.Background(), "TeacherGetAll", db, time.Duration(1)*time.Hour)
	//
	return teachers
}

func TeacherExist(session *xorm.Session, name string) (ret bool) {
	//查缓存
	key := fmt.Sprintf("TeacherExist:%v", name)
	v, ext := redis.RedisGet(context.Background(), key)
	if ext {
		err := json.Unmarshal(v, &ret)
		if err == nil {
			return ret
		}
		util.DebugError("TeacherExist", err)
	}
	//
	defer func() {
		//写入缓存
		redis.RedisSet(context.Background(), key, ret, time.Duration(5)*time.Minute)
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	teacher := TeacherGetAll(session)
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
