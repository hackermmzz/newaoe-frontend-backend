package service

import (
	"newaoe/Src/forum/dao"
	"newaoe/Src/forum/model"
	"newaoe/Src/util"
)

func ForumCommentFetch(beg int, end int) ([]model.ForumComment, error) {
	// 校验范围
	if e := util.ValidateInclusiveRange(beg, end); e != nil {
		return nil, e
	}
	res, err := dao.ForumCommentGetByRange(nil, beg, end)
	if err != nil {
		return nil, err
	}
	return res, nil
}
