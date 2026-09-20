package controller

import (
	"fmt"
	"newaoe/Src/user/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

// 数据结构体
type touristRegistDataInfo struct {
	Password   string `json:"password"`
	Email      string `json:"email"`
	Verifycode string `json:"verifycode"`
}

func TouristUserRegist(c *gin.Context) {
	//
	var dt touristRegistDataInfo
	if !util.JsonCtx(c, &dt) {
		util.ResponseNAK_MSG(c, "数据报文错误", "")
		return
	}
	//
	//生成id(这里直随机生成，不使用唯一键来生成有序的)
	rid := util.TruncateStringForEmail(dt.Email)
	if rid == "" {
		util.ResponseNAK_MSG(c, "请输入正确的邮箱!", nil)
		return
	}
	id := fmt.Sprintf("tourist_20070128_%v", rid)
	//注册
	err := service.UserRegistTouristUser(id, dt.Password, dt.Email, dt.Verifycode)
	if err != nil {
		util.DebugError("TouristUserRegist:", err)
		util.ResponseNAK_MSG(c, err.Error(), nil)
		return
	}
	//注册成功
	util.ResponseACK_MSG(c, "注册成功", map[string]interface{}{"id": id})
}
