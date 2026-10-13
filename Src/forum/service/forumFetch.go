package service

import (
	"newaoe/Src/forum/dao"
	"newaoe/Src/forum/model"
	"newaoe/Src/util"
)

// 不包含end
func ForumFetch(beg int, end int) ([]model.ForumAnnouncement, error) {
	// 校验范围
	if e := util.ValidateInclusiveRange(beg, end); e != nil {
		return nil, e
	}
	// 查询数据
	res, err := dao.ForumAnnouncementGetByRange(nil, beg, end)
	if err != nil {
		return nil, err
	}
	return res, nil
}
