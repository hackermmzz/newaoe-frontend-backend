package database

import (
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/util"

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
	util.DebugSuccess("数据库连接成功")

}

// 封装一下seesion
type Session struct {
	*xorm.Session
	taskQueue []func()
}

// PostTransactionAdd registers a task that is allowed to run only after the
// database operation is durable. Autocommit operations are already durable
// when this method is called, so their task can run immediately.
func (p *Session) PostTransactionAdd(task func()) {
	if task == nil {
		return
	}
	if !p.IsInTx() {
		task()
		return
	}
	p.taskQueue = append(p.taskQueue, task)
}

// PostTransaction executes and clears all registered tasks exactly once.
func (p *Session) PostTransaction() {
	tasks := p.taskQueue
	p.taskQueue = nil
	for _, task := range tasks {
		task()
	}
}

func (p *Session) Commit() error {
	if err := p.Session.Commit(); err != nil {
		p.taskQueue = nil
		return err
	}
	p.PostTransaction()
	return nil
}

// Rollback discards tasks whose corresponding database changes were not
// committed.
func (p *Session) Rollback() error {
	p.taskQueue = nil
	return p.Session.Rollback()
}

// CacheAvailable is false inside a transaction. Reading from Redis there can
// violate transaction isolation, and writing a transaction snapshot back to
// Redis after commit can overwrite newer cached data.
func CacheAvailable(session *Session) bool {
	return session == nil || !session.IsInTx()
}

func NewSession() *Session {
	return &Session{Session: DB.NewSession()}
}
