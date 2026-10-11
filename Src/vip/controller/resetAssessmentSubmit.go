package controller

import (
	"strings"

	"newaoe/Src/codeAssessmentSubmit/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func ResetAssessmentSubmit(ctx *gin.Context) {
	idOrEmail := strings.TrimSpace(ctx.Query("id"))
	if idOrEmail == "" {
		util.ResponseNAK_MSG(ctx, "id或邮箱不能为空", nil)
		return
	}
	err := service.ResetAssessmentSubmitRecord(idOrEmail)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "重置成功!", nil)
}
