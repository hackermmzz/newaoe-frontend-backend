package dao

import (
	database "newaoe/Src/databse"
	"newaoe/Src/home/model"

	"xorm.io/xorm"
)

func AnnouncementAdd(session *xorm.Session, data model.AnnouncementInfo) (int64, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	// Insert 返回影响行数 + error
	_, err := session.Insert(&data)
	if err != nil {
		return 0, err
	}
	return data.Indices, nil
}

func AnnouncementGetLatest(session *xorm.Session) (*model.AnnouncementInfo, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	var data model.AnnouncementInfo

	has, err := session.Where("enabled = ?", true).Desc("indices").Get(&data)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return &data, nil
}
