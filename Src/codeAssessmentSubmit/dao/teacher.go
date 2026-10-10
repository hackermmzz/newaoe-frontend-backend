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

func TeacherAdd(session *xorm.Session, teacher string) (err error) {
	defer func() {
		if err == nil {
			//删缓存
			if _, e := redis.RedisDel(context.Background(), "TeacherGetAll"); e != nil {
				util.DebugError("TeacherGetAll cache del:", e)
			}
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

func TeacherGetAll(session *xorm.Session) (ret []model.Teacher, err error) {
	//查缓存
	dataBytes, exist, err := redis.RedisGet(context.Background(), "TeacherGetAll")
	if err != nil {
		return nil, err
	}
	if exist {
		var ret []model.Teacher
		err := json.Unmarshal(dataBytes, &ret)
		if err == nil {
			return ret, nil
		}
		return nil, err
	}
	//写入缓存
	defer func() {
		db, _ := json.Marshal(ret)
		//写入缓存
		if _, err := redis.RedisSet(context.Background(), "TeacherGetAll", db, time.Duration(1)*time.Hour); err != nil {
			util.DebugError("TeacherGetAll cache set:", err)
		}
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

func TeacherExist(session *xorm.Session, name string) (ret bool, err error) {
	//查缓存
	key := fmt.Sprintf("TeacherExist:%v", name)
	v, ext, err := redis.RedisGet(context.Background(), key)
	if err != nil {
		return false, err
	}
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
		if _, err := redis.RedisSet(context.Background(), key, ret, time.Duration(5)*time.Minute); err != nil {
			util.DebugError("TeacherExist cache set:", err)
		}
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
