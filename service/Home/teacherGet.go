package Home

import (
	data "newaoe/service/Data"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

func StudentGetTeacher(ctx *gin.Context) {
	teacher, err := data.TeacherGetAll()
	if err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常", "")
		return
	}
	//
	util.ResponseACK_MSG(ctx, "获取成功", teacher)
}
