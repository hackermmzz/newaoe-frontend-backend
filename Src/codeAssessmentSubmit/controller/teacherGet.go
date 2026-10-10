package controller

import (
	"newaoe/Src/codeAssessmentSubmit/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func StudentGetTeacher(ctx *gin.Context) {
	teacher, err := dao.TeacherGetAll(nil)
	if err != nil {
		util.DebugError("TeacherGetAll:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "获取成功", teacher)
}
