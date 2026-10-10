package dao

import (
	database "newaoe/Src/databse"
	"newaoe/Src/home/model"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func FeedbackGetByIndices(session *xorm.Session, indices int) (*model.FeedbackInfo, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ret model.FeedbackInfo
	has, err := session.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &ret, nil
}

func FeedbackInsert(session *xorm.Session, info model.FeedbackInfo) (int, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	// Insert 返回影响行数 + error
	_, err := session.Insert(&info)
	if err != nil {
		return 0, err
	}
	return info.Indices, nil
}

// 不包含end
func FeedbackGetRange(session *xorm.Session, beg int, end int) ([]model.FeedbackInfo, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if beg < 0 || beg >= end {
		return nil, util.NewError("beg<0 or beg>=end!")
	}
	var ret []model.FeedbackInfo
	err := session.OrderBy("indices desc").
		Limit(end-beg, int(beg)).
		Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
