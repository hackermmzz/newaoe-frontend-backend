package controller

import (
	"fmt"
	"newaoe/Src/codeAssessmentSubmit/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func TeacherAdd(ctx *gin.Context) {
	teacher := ctx.Query("teacher")
	//这里比较简单，直接走dao
	err := dao.TeacherAdd(nil, teacher)
	if err != nil {
		util.DebugError("TeacherAdd:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, fmt.Sprintf("添加教师: %s 成功", teacher), nil)
}
