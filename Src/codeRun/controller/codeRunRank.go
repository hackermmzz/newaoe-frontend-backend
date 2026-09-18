package controller

import (
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	UserDao "newaoe/Src/user/dao"
	"newaoe/Src/util"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FetchRank(ctx *gin.Context) {
	//解析参数
	range_ := ctx.Query("range")
	parts := strings.Split(range_, ":")
	if len(parts) != 2 {
		util.ResponseNAK_MSG(ctx, "参数传递错误!", "")
		return
	}
	beg, err0 := strconv.Atoi(parts[0])
	end, err1 := strconv.Atoi(parts[1])
	if err0 != nil || err1 != nil {
		util.ResponseNAK_MSG(ctx, "参数传递错误!", "")
		return
	}
	//获取beg到end的排行数据
	session := database.NewSession()
	defer session.Rollback()
	err := session.Begin()
	if err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常", "")
		return
	}
	info := dao.RankGetByRange(session, beg, end+1)
	if info == nil {
		info = make([]model.RankInfo, 0)
	}
	//获取人物头像
	ids := make([]string, len(info))
	for i, d := range info {
		ids[i] = d.ID
	}
	avatars := UserDao.UserGetByIDs(ids)
	if len(avatars) != len(info) {
		util.ResponseNAK_MSG(ctx, "服务器异常", "")
		return
	}
	idToAvatar := make(map[string]string)
	for _, d := range avatars {
		idToAvatar[d.Id] = d.Avatar
	}
	//处理一下数据
	type FinalDataInfo struct {
		model.RankInfo
		Avatar string `json:"avatar"`
	}
	finaldata := make([]FinalDataInfo, len(info))

	for i, d := range info {
		finaldata[i].Avatar = idToAvatar[d.ID]
		finaldata[i].RankInfo = d
	}
	//返回结果
	util.ResponseACK_MSG(ctx, "获取排行成功", finaldata)
}
