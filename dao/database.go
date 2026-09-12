package dao

import (
	"fmt"
	"newaoe/config"
	"newaoe/util"
	"time"

	_ "github.com/go-sql-driver/mysql" // 导入驱动（下划线表示只执行init函数）
	"xorm.io/xorm"
)

var DB *xorm.Engine

func ConnectDatabase() {
	//连接指定的数据库
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=UTC",
		config.Conf.Mysql.Username,
		config.Conf.Mysql.Password,
		config.Conf.Mysql.Host,
		config.Conf.Mysql.DBName,
	)
	DB, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}
	// 设置最大空闲连接数（推荐：2~10）
	DB.SetMaxIdleConns(config.Conf.Mysql.DBMaxIdleConns)
	// 设置最大打开连接数（根据MySQL配置调整，一般不超过50）
	DB.SetMaxOpenConns(config.Conf.Mysql.DBMaxOpenConns)
	// 设置连接的最大生命周期（避免长时间占用连接）
	DB.SetConnMaxLifetime(time.Duration(config.Conf.Mysql.DBConnMaxLifetime) * time.Second)
	// 设置连接的最大空闲时间
	DB.SetConnMaxIdleTime(time.Duration(config.Conf.Mysql.DBConnMaxIdleTime) * time.Second)
	//多ping几次防止网络不好导致误判
	for i := 0; i < 10; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
	}
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}
	util.Debug("数据库连接成功")

}
