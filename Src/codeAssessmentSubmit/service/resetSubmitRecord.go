package service

import (
	"errors"
	"newaoe/Src/codeAssessmentSubmit/dao"
	userDao "newaoe/Src/user/dao"
)

func ResetAssessmentSubmitRecord(idOrEmail string) error {
	//获取用户id
	user, err := userDao.UserGetByIdOrEmail(nil, idOrEmail)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}
	//删除提交记录
	err = dao.CodeAssessmentRemoveByID(nil, user.ID)
	if err != nil {
		return err
	}
	return nil
}
