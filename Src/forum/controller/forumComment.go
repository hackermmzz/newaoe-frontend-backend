package controller

import (
	"newaoe/Src/forum/model"
	"newaoe/Src/forum/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

type commentRecevied struct {
	Content                    string `json:"content"`
	IndicesForForumAnnoucement int64  `json:"indicesForForumAnnoucement"`
}

func ForumComment(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//解析请求体
	var req commentRecevied
	if ok := util.JsonCtx(ctx, &req); !ok {
		util.ResponseNAK_MSG(ctx, "请求体解析失败", nil)
		return
	}
	//校验参数
	if len(req.Content) == 0 {
		util.ResponseNAK_MSG(ctx, "内容不能为空", nil)
		return
	}
	if len(req.Content) >= 100 {
		util.ResponseNAK_MSG(ctx, "内容不能超过100个字符", nil)
		return
	}
	if req.IndicesForForumAnnoucement <= 0 {
		util.ResponseNAK_MSG(ctx, "公告ID不能为空", nil)
		return
	}
	//插入评论
	err := service.CommentAdd(&model.ForumComment{
		ID:                         id,
		Content:                    req.Content,
		IndicesForForumAnnoucement: req.IndicesForForumAnnoucement,
		SubmitTime:                 util.UTC_Time(),
	})
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "评论成功", nil)
}
