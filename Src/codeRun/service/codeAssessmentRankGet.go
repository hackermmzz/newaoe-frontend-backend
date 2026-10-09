package service

import (
	"errors"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	UserDao "newaoe/Src/user/dao"
	"newaoe/Src/util"
)

// 处理一下数据
type AssessmentRankFetchInfo struct {
	RankInfo model.AssessmentRankInfo    `json:"rankinfo"`
	Msg      model.AssessmentRankMsgInfo `json:"msg"`
	Avatar   string                      `json:"avatar"`
}

func AssessmentRankFetch(beg int, end int) ([]AssessmentRankFetchInfo, error) {
	session := database.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return nil, util.NewError("服务器异常")
	}
	defer session.Rollback()
	//获取RankInfo
	info := dao.AssessmentRankGetByRange(session, beg, end)
	if info == nil {
		info = make([]model.AssessmentRankInfo, 0)
	}

	//获取人物头像
	ids := make([]string, len(info))
	for i, d := range info {
		ids[i] = d.ID
	}
	avatars := UserDao.UserGetByIDs(nil, ids)
	if len(avatars) != len(info) {
		return nil, errors.New("服务器异常!")
	}
	idToAvatar := make(map[string]string)
	for _, d := range avatars {
		idToAvatar[d.Id] = d.Avatar
	}

	finaldata := make([]AssessmentRankFetchInfo, len(info))

	for i, d := range info {
		var msgInfo model.AssessmentRankMsgInfo
		msgInfo.Unmarshal([]byte(d.Msg))
		finaldata[i].Avatar = idToAvatar[d.ID]
		finaldata[i].RankInfo = d
		finaldata[i].Msg = msgInfo
	}

	return finaldata, nil
}
