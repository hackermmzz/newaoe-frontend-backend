package service

import (
	"errors"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	UserDao "newaoe/Src/user/dao"
	"newaoe/Src/util"
)

// 处理一下数据
type RankFetchInfo struct {
	RankInfo model.RankInfo    `json:"rankinfo"`
	Msg      model.RankMsgInfo `json:"msg"`
	Avatar   string            `json:"avatar"`
}

func RankFetch(beg int, end int) ([]RankFetchInfo, error) {
	//获取RankInfo
	info, err := dao.RankGetByRange(nil, beg, end)
	if err != nil {
		return nil, err
	}
	if info == nil {
		info = make([]model.RankInfo, 0)
	}

	//获取人物头像
	ids := make([]string, len(info))
	for i, d := range info {
		ids[i] = d.ID
	}
	avatars, err := UserDao.UserGetByIDs(nil, ids)
	if err != nil {
		return nil, err
	}
	if len(avatars) != len(info) {
		return nil, errors.New("服务器异常!")
	}
	idToAvatar := make(map[string]string)
	for _, d := range avatars {
		idToAvatar[d.Id] = d.Avatar
	}

	finaldata := make([]RankFetchInfo, 0, len(info))

	for _, d := range info {
		var msgInfo model.RankMsgInfo
		if err := msgInfo.Unmarshal([]byte(d.Msg)); err != nil {
			util.DebugError("RankFetch: rank msg unmarshal failed:", d.ID, err)
			continue
		}
		finaldata = append(finaldata, RankFetchInfo{
			Avatar:   idToAvatar[d.ID],
			RankInfo: d,
			Msg:      msgInfo,
		})
	}

	return finaldata, nil
}
