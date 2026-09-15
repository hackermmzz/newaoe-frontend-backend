package vip

import (
	"newaoe/service/Home"

	"github.com/gin-gonic/gin"
)

func GetStudentHistory(ctx *gin.Context) {
	Home.StudentHistoryGet(ctx)
}
