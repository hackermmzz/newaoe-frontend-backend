package controller

import (
	"newaoe/Src/codeAssessmentSubmit/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func StudentGetTeacher(ctx *gin.Context) {
	teacher := dao.TeacherGetAll(nil)
	util.ResponseACK_MSG(ctx, "获取成功", teacher)
}
