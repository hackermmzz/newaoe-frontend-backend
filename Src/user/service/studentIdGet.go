package service

import (
	"errors"
	"newaoe/Src/user/dao"
)

func StudentIDGetByEmailOrID(id_or_emial string) (string, error) {
	info := dao.UserGetByIdOrEmail(nil, id_or_emial)
	if info == nil {
		return "", errors.New("查无此账号!")
	}
	return info.Id, nil
}
