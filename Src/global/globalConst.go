package global

const (
	CodeRunTypeRelease CodeRunType = iota
	CodeRunTypeDebug
)

const (
	Code_Status_Error           = 0 //服务器异常
	Code_Status_Wait            = 1 //在等待队列里面
	Code_Status_Compile         = 2 //编译中
	Code_Status_Compile_Success = 3 //编译成功
	Code_Status_Compile_Fail    = 4 //编译错误
	Code_Status_Running         = 5 //正在运行出结果
	Code_Status_Success         = 6 //运行胜利
	Code_Status_Fail            = 7 //游戏失败
	Code_Status_Crash           = 8 //游戏崩溃
)
