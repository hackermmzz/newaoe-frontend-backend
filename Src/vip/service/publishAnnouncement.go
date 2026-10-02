package service

import (
	"newaoe/Src/home/dao"
	"newaoe/Src/home/model"
	"newaoe/Src/util"
)

func PublishAnnouncement(title string, content string, enabled bool) error {
	curTime := util.UTC_Time()
	data := model.AnnouncementInfo{
		Title:      title,
		Content:    content,
		Version:    1,
		Enabled:    enabled,
		CreateTime: curTime,
		UpdateTime: curTime,
	}
	if 0 == dao.AnnouncementAdd(nil, data) {
		return util.NewError("数据库异常!")
	}
	return nil
}
