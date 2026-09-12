package User

import (
	"fmt"
	"newaoe/config"
	"newaoe/dao"
	data "newaoe/service/Data"
	"newaoe/service/Email"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

func UserRegist(c *gin.Context) {
	//数据结构体
	type DataInfo struct {
		Id         string `json:"id"`
		Password   string `json:"password"`
		Email      string `json:"email"`
		Verifycode string `json:"verifycode"`
	}
	//
	var dt DataInfo
	if !util.JsonCtx(c, &dt) {
		util.ResponseNAK_MSG(c, "数据报文错误", "")
		return
	}
	//处理邮箱
	dt.Email = userEmailProcess(dt.Id, dt.Email)
	//检查密码是否符合格式
	if !PasswordLegal(dt.Password) {
		util.ResponseNAK_MSG(c, "密码格式错误!", "")
		return
	}
	//检查验证码是否正确
	if !data.RegistVerifyCodeExist(dt.Id, dt.Verifycode) {
		util.ResponseNAK_MSG(c, "验证码错误", "")
		return
	}
	//查询用是否允许注册
	if !data.UserIdLegal(dt.Id) {
		util.ResponseNAK_MSG(c, "该用户不允许注册账号!", "")
		return
	}
	//查询是否已经注册
	if data.UserExist(dt.Id) {
		util.ResponseNAK_MSG(c, "用户已经注册", "")
		return
	}
	//查询邮箱是否已经使用过
	if data.EmailExist(dt.Email) {
		util.ResponseNAK_MSG(c, "邮箱已经使用", "")
		return
	}
	//添加用户
	session := dao.DB.NewSession()
	defer session.Rollback()
	err := session.Begin()
	if err != nil {
		util.ResponseNAK_MSG(c, "服务器异常", "")
		return
	}

	password, err := util.EncodePassword(dt.Password)
	if err != nil {
		util.ResponseNAK_MSG(c, "密码格式不合规则!", "")
		return
	}
	if !data.UserAdd(session, dt.Id, password, dt.Email) {
		util.ResponseNAK_MSG(c, "注册失败", "")
		return
	}

	if err = session.Commit(); err != nil {
		util.ResponseNAK_MSG(c, "服务器异常!", "")
		return
	}
	//注册成功
	util.ResponseACK_MSG(c, "注册成功", "")
}

// 判断密码是否合法
func PasswordLegal(password string) bool {
	return true
}

// 发送注册码
func UserRegistCodeSend(ctx *gin.Context) {
	type DataInfo struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	}
	//
	var dt DataInfo
	if !util.JsonCtx(ctx, &dt) {
		util.ResponseNAK_MSG(ctx, "数据报文错误", "")
		return
	}
	//处理邮箱
	dt.Email = userEmailProcess(dt.Id, dt.Email)
	//判断是否可以发送验证码
	ttl := data.RegistVerifyCodeCanSendTTL(dt.Id)
	if ttl > 0 {
		util.ResponseNAK_MSG(ctx, fmt.Sprintf("请等待%v秒后重试!", ttl), "")
		return
	}
	//发送验证码
	code := util.GenerateVerifyCode(config.Conf.RegistVerifyCode.Length)
	Email.SendTextEmail(dt.Email, "注册验证码", "你的验证码为: "+code)
	//插入Redis
	success := data.RegistVerifyCodeAdd(dt.Id, code)
	//
	if !success {
		util.ResponseNAK_MSG(ctx, "验证码发送失败", "")
		return
	}
	util.ResponseACK_MSG(ctx, "验证码发送成功", "")
}

// 对用户绑定的邮箱进行处理
func userEmailProcess(id string, email string) string {
	//强制邮箱只能是Id@njust.edu.cn
	return id + "@njust.edu.cn"
}
