package initialize

import database "newaoe/Src/databse"

func DataBaseInit() {
	//连接数据库
	database.ConnectDatabase()
}
