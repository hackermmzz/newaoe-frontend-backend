package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MSG struct {
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func GetResponseACK_MSG(msg string, data interface{}) MSG {
	message := MSG{
		Status: true,
		Msg:    msg,
		Data:   data,
	}
	return message
}

func GetResponseNAK_MSG(msg string, data interface{}) MSG {
	message := MSG{
		Status: false,
		Msg:    msg,
		Data:   data,
	}
	return message
}

func ResponseACK_MSG(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, GetResponseACK_MSG(msg, data))
}

func ResponseNAK_MSG(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, GetResponseNAK_MSG(msg, data))
}
