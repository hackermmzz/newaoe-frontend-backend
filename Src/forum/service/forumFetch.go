package service

import (
	"newaoe/Src/forum/dao"
	"newaoe/Src/forum/model"
	userDao "newaoe/Src/user/dao"
	"newaoe/Src/util"
)

type ForumAnnouncementInfo struct {
	model.ForumAnnouncement
	Avatar string `json:"avatar"`
}

// 不包含end
func ForumFetch(beg int, end int) ([]ForumAnnouncementInfo, error) {
	// 校验范围
	if e := util.ValidateInclusiveRange(beg, end); e != nil {
		return nil, e
	}
	// 查询数据
	res, err := dao.ForumAnnouncementGetByRange(nil, beg, end)
	if err != nil {
		return nil, err
	}
	// 处理数据
	ret := make([]ForumAnnouncementInfo, len(res))
	for i, v := range res {
		//获取用户头像
		avatar, err := userDao.UserGetAvatar(nil, v.ID)
		if err != nil {
			avatar = ""
			util.DebugError("为什么获取用户头像失败了呢?", err)
		}
		ret[i] = ForumAnnouncementInfo{
			ForumAnnouncement: v,
			Avatar:            avatar,
		}
	}
	return ret, nil
}
