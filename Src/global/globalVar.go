package global

import (
	"net/http"
	"new-aoe-judge/Src/config"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var GRPCAuth GRPCAuthInfo
var Profiler *ProfilerInfo
var MapFiles []string
var JudgeMaxCore int
var CompileCPULimit, RunCPULimit int
var LogFile *os.File
var logMu sync.Mutex
var authMu sync.RWMutex
var CompileWaitQueue chan interface{}
var RunWaitQueue chan interface{}
var httpClient *resty.Client

func Init() error {
	//初始化连接池
	httpClient = resty.New().
		SetTimeout(30 * time.Second).
		SetTransport(&http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     60 * time.Second,
		})
	//初始化核心资源
	initCoreResource()
	//读取地图文件
	entries, _ := os.ReadDir(config.Conf.NewAOEFolder)
	for _, x := range entries {
		if filepath.Ext(x.Name()) == ".njust" {
			MapFiles = append(MapFiles, x.Name())
		}
	}
	//打开日志文件
	var e error
	LogFile, e = os.OpenFile(config.Conf.ProcessLogFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0644)
	if e != nil {
		return e
	}
	//创建队列
	CompileWaitQueue = make(chan interface{}, config.Conf.JudgeCodeFetchQueueMaxPayload)
	RunWaitQueue = make(chan interface{}, config.Conf.JudgeCodeFetchQueueMaxPayload/2)
	//
	return e
}

// initCoreResource 初始化核心资源
func initCoreResource() {
	//获取核数(减去1个核心,0核心给系统使用)
	JudgeMaxCore = max(runtime.NumCPU()-1, 2)
	compileOccurCore := max(JudgeMaxCore/3, 1)
	runOccurCore := JudgeMaxCore - compileOccurCore
	//配置编译和运行的CPU资源限制
	CompileCPULimit = 100
	RunCPULimit = 100
	//创建性能分析信息
	compileCoreInfo := make([]SingleCoreInfo, compileOccurCore)
	runCoreInfo := make([]SingleCoreInfo, runOccurCore)
	for i := range compileCoreInfo {
		compileCoreInfo[i].CoreID = i + 1
	}
	for i := range runCoreInfo {
		runCoreInfo[i].CoreID = i + compileOccurCore + 1
	}
	Profiler = NewProfilerInfo(compileCoreInfo, runCoreInfo)
}

//
