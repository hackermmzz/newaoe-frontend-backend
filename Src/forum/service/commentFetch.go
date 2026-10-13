package service

import (
	"newaoe/Src/forum/dao"
	"newaoe/Src/forum/model"
	userDao "newaoe/Src/user/dao"
	"newaoe/Src/util"
)

type ForumCommentInfo struct {
	model.ForumComment
	Avatar string `json:"avatar"`
}

func ForumCommentFetch(beg int, end int) ([]ForumCommentInfo, error) {
	// 校验范围
	if e := util.ValidateInclusiveRange(beg, end); e != nil {
		return nil, e
	}
	res, err := dao.ForumCommentGetByRange(nil, beg, end)
	if err != nil {
		return nil, err
	}
	ret := make([]ForumCommentInfo, len(res))
	for i, v := range res {
		//获取用户头像
		avatar, err := userDao.UserGetAvatar(nil, v.ID)
		if err != nil {
			avatar = ""
			util.DebugError("为什么获取用户头像失败了呢?", err)
		}
		ret[i] = ForumCommentInfo{
			ForumComment: v,
			Avatar:       avatar,
		}
	}
	return ret, nil
}
