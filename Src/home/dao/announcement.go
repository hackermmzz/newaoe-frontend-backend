package dao

import (
	database "newaoe/Src/databse"
	"newaoe/Src/home/model"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func AnnouncementAdd(session *xorm.Session, data model.AnnouncementInfo) int64 {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	// Insert 返回影响行数 + error
	_, err := session.Insert(&data)
	if err != nil {
		util.DebugError("AnnouncementAdd:", err)
		return 0
	}
	return data.Indices
}

func AnnouncementGetLatest(session *xorm.Session) *model.AnnouncementInfo {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	var data model.AnnouncementInfo

	has, err := session.Where("enabled = ?", true).Desc("indices").Get(&data)

	if err != nil {
		util.DebugError("AnnouncementGetLatest:", err)
		return nil
	}

	if !has {
		return nil
	}

	return &data
}
