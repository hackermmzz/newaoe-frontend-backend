package global

var (
	CodeRunType_ReleaseRun = 1000 //release版本运行
	CodeRunType_DebugRun   = 2000 //重度调试模式运行
)

const (
	Code_Status_Unknown         = 0    //未知状态
	Code_Status_Error           = 1000 //服务器异常
	Code_Status_Wait            = 2000 //在等待队列里面
	Code_Status_OJRecevied      = 2050 //OJ接收成功
	Code_Status_Compile         = 3000 //编译中
	Code_Status_Compile_Fail    = 4000 //编译错误
	Code_Status_Compile_Success = 5000 //编译成功
	Code_Status_Running         = 6000 //正在运行出结果
	Code_Status_Crash           = 7000 //游戏崩溃
	Code_Status_Fail            = 8000 //游戏失败
	Code_Status_Success         = 9000 //运行胜利
)
