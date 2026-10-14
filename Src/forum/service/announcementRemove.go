package service

import (
	"newaoe/Src/forum/dao"
	"newaoe/Src/user/model"
	"newaoe/Src/util"
)

// 删除帖子indices，当且仅当id用户为vip或者他是帖子发布者
func ForumAnnouncementRemove(indices int64, id string, vip int) error {
	bNeedRemove := vip >= model.VIP_SUPER
	if !bNeedRemove {
		//检测用户是否为帖子发布者
		if bPublisher, err := dao.ForumAnnouncementIsPublisher(nil, indices, id); !bPublisher || err != nil {
			if err == nil {
				err = util.NewError(id, "权限错误!")
			}
			return err
		}
		bNeedRemove = true
	}
	//删除帖子
	if bNeedRemove {
		ok, err := dao.ForumAnnouncementRemove(nil, indices)
		if err != nil {
			return err
		}
		if !ok {
			return util.NewError(id, "查无此帖子!")
		}
		return nil
	}
	return util.NewError(id, "权限错误!")
}
