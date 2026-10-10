package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/codeAssessmentSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
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

func TeacherGetAll(session *xorm.Session) ([]model.Teacher, error) {
	//查缓存
	dataBytes, exist := redis.RedisGet(context.Background(), "TeacherGetAll")
	if exist {
		var ret []model.Teacher
		err := json.Unmarshal(dataBytes, &ret)
		if err == nil {
			return ret, nil
		}
		return nil, err
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
		return nil, err
	}
	//返回数据写入缓存
	db, _ := json.Marshal(teachers)
	redis.RedisSet(context.Background(), "TeacherGetAll", db, time.Duration(1)*time.Hour)
	//
	return teachers, nil
}

func TeacherExist(session *xorm.Session, name string) (ret bool, err error) {
	//查缓存
	key := fmt.Sprintf("TeacherExist:%v", name)
	v, ext := redis.RedisGet(context.Background(), key)
	if ext {
		err := json.Unmarshal(v, &ret)
		if err == nil {
			return ret, nil
		}
		return false, err
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
