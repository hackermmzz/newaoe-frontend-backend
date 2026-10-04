package global

import (
	"fmt"
	"net/http"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/util"
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
	err := initCoreResource()
	if err != nil {
		return err
	}
	//读取地图文件
	entries, _ := os.ReadDir(config.Conf.NewAOEFolder)
	for _, x := range entries {
		if filepath.Ext(x.Name()) == ".njust" {
			MapFiles = append(MapFiles, x.Name())
		}
	}
	//初始化日志文件
	e := logInit()
	if e != nil {
		return e
	}
	//创建队列
	CompileWaitQueue = make(chan interface{}, JudgeMaxCore*2)
	RunWaitQueue = make(chan interface{}, JudgeMaxCore*4)
	//
	return e
}

// initCoreResource 初始化核心资源
func initCoreResource() error {
	if runtime.NumCPU() < 2 {
		return fmt.Errorf("CPU核心数过少, 需要至少2个核心, 当前核心数:%d", runtime.NumCPU())
	}
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
	//
	LogInfo(fmt.Sprintf("初始化核心资源完成, 总核数:%d, 编译核数:%d, 运行核数:%d, 编译CPU限制:%d, 运行CPU限制:%d", JudgeMaxCore, compileOccurCore, runOccurCore, CompileCPULimit, RunCPULimit))
	if JudgeMaxCore == 2 {
		LogInfo("警告: 核心数过少, 可能会导致性能下降")
	}
	if compileOccurCore == 0 {
		return fmt.Errorf("警告: 编译核数为0, 可能会导致编译失败")
	}
	if runOccurCore == 0 {
		return fmt.Errorf("警告: 运行核数为0, 可能会导致运行失败")
	}
	if CompileCPULimit <= 0 {
		return fmt.Errorf("警告: 编译CPU限制为0, 可能会导致编译失败")
	}
	if RunCPULimit <= 0 {
		return fmt.Errorf("警告: 运行CPU限制为0, 可能会导致运行失败")
	}
	return nil
}

func logInit() error {
	//清空文件夹
	err := os.RemoveAll(util.JoinPath(config.Conf.ProcessLogDir))
	if err != nil {
		return err
	}
	//创建文件夹
	err = os.MkdirAll(config.Conf.ProcessLogDir, 0755)
	if err != nil {
		return err
	}
	//创建日志文件
	LogFile, err = os.OpenFile(util.JoinPath(config.Conf.ProcessLogDir, "ProcessLog0.html"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	return nil
}
