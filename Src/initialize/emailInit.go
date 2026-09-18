package initialize

import (
	"newaoe/Src/email/service"
)

func EmailInit() {
	//初始化邮箱连接
	service.EmailSendServiceInit()
}
