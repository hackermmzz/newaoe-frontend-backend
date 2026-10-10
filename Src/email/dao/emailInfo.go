package dao

import (
	database "newaoe/Src/databse"
	"newaoe/Src/email/model"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func EmailInfoInsert(session *xorm.Session, info model.EmailInfo) (int, error) {
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

// 前一个为是否成功，后一个为是否有改变
func EmailInfoUpdateSendStatus(session *xorm.Session, indices int, send bool) (bool, int64, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	affected, err := session.ID(indices).
		Cols("send"). // 明确指定只更新send字段
		Update(&model.EmailInfo{
			Send: send,
		})
	if err != nil {
		return false, 0, err
	}
	return true, affected, nil
}

func EmailInfoGetForUpdate(session *xorm.Session, indices int) (*model.EmailInfo, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ret model.EmailInfo
	has, err := session.ForUpdate().Where("indices = ?", indices).Get(&ret)
	if err != nil {
		util.DebugError("EmailInfoGet:", err)
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &ret, nil
}
