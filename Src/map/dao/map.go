package dao

import (
	"fmt"
	database "newaoe/Src/databse"
	"newaoe/Src/map/model"

	"xorm.io/xorm"
)

func MapInsert(session *xorm.Session, key string, value string) error {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	info := &model.MapInfo{
		Key:   key,
		Value: value,
	}

	_, err := session.Insert(info)
	if err != nil {
		return err
	}

	return nil
}

func MapDelete(session *xorm.Session, key string) error {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	_, err := session.
		Where("`key` = ?", key).
		Delete(&model.MapInfo{})

	return err
}

func MapUpdate(session *xorm.Session, key string, value string) error {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	_, err := session.
		Where("`key` = ?", key).
		Cols("value").
		Update(&model.MapInfo{
			Value: value,
		})

	return err
}

func MapExist(session *xorm.Session, key string) (bool, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	exist, err := session.
		Where("`key` = ?", key).
		Exist(&model.MapInfo{})

	return exist, err
}

func MapGet(session *xorm.Session, key string) (string, error) {
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
		return "", fmt.Errorf("%v不存在!", key)
	}

	return info.Value, nil
}

func MapUpdateOrInsert(session *xorm.Session, key string, value string) error {
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
