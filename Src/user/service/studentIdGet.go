package service

import (
	"errors"
	"newaoe/Src/user/dao"
	"newaoe/Src/user/model"
)

func StudentInfoGetByEmailOrID(id_or_emial string) (*model.Student, error) {
	info := dao.UserGetByIdOrEmail(nil, id_or_emial)
	if info == nil {
		return nil, errors.New("查无此账号!")
	}
	return info, nil
}
