package global

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

type CodeRunType int64

// ///////////////////////// OJSystemInfo OJ系统信息
type OJSystemInfo struct {
	Version string `json:"version"` //最新的版本号
}

// //////////////////////// GRPCAuthInfo gRPC 认证信息
type GRPCAuthInfo struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

func (a GRPCAuthInfo) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}

// //////////////////////// CodeRunStatusInfo 代码运行状态信息
type CodeRunStatusInfo struct {
	Status   int32  `json:"status"`
	Food     int    `json:"food"`
	Wood     int    `json:"wood"`
	Gold     int    `json:"gold"`
	Stone    int    `json:"stone"`
	Frame    int    `json:"frame"`
	Win      bool   `json:"win"`
	Score    int    `json:"score"`
	Data     string `json:"data"`
	TimeCost int64  `json:"time_cost"` //秒
}

func (s CodeRunStatusInfo) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// //////////////////////// ProfilerInfo 性能分析信息
type SingleCoreInfo struct {
	CoreID int `json:"coreID"`
	Usage  int `json:"usage"`
}

type CoreInfo struct {
	mu        sync.Mutex
	cond      *sync.Cond
	CoreInfos []SingleCoreInfo `json:"cores"`
}

type CPUResourceAllocateInfo struct {
	CoreIDs     []int `json:"coreIDs"`
	CPUSPerCore []int `json:"cpusPerCore"`
	CPUS        int   `json:"cpus"`
}

type ProfilerInfo struct {
	//统计性能
	mu       sync.Mutex
	taskWait int
	compile  int
	running  int
	//任务信号量
	taskMu        sync.Mutex
	taskAcquireMu sync.Mutex //锁住后任务只能减少不能增加
	taskSem       int64
	//分配资源
	compileCore *CoreInfo
	runCore     *CoreInfo
}

// GetCores 获取核心列表
func (c CPUResourceAllocateInfo) GetCores() string {
	ret := ""
	for i, v := range c.CoreIDs {
		if i != 0 {
			ret += ","
		}
		ret += fmt.Sprintf("%d", v)
	}
	return ret
}

// GetCPUS 获取cpu资源
func (c CPUResourceAllocateInfo) GetCPUS() string {
	ret := ""
	ret += fmt.Sprintf("%.2f", float64(c.CPUS)/100.0)
	return ret
}

// newCoreInfo 创建核心信息
func newCoreInfo(singleCoreInfos []SingleCoreInfo) *CoreInfo {
	r := &CoreInfo{
		CoreInfos: singleCoreInfos,
	}
	r.cond = sync.NewCond(&r.mu)
	return r
}

// NewProfilerInfo 创建性能分析信息
func NewProfilerInfo(compileCores []SingleCoreInfo, runCores []SingleCoreInfo) *ProfilerInfo {
	return &ProfilerInfo{
		taskWait:    0,
		compile:     0,
		running:     0,
		compileCore: newCoreInfo(compileCores),
		runCore:     newCoreInfo(runCores),
	}
}

// 分配cpu资源
func (p *CoreInfo) allocateResource(cpus int) CPUResourceAllocateInfo {
	p.mu.Lock()
	defer p.mu.Unlock()
	//尝试分配资源函数
	tryAllocate := func() *CPUResourceAllocateInfo {
		//分配资源(简单分配)
		tmp := make([]SingleCoreInfo, len(p.CoreInfos))
		copy(tmp, p.CoreInfos)
		//排序核心按使用率从低到高
		sort.Slice(tmp, func(i, j int) bool {
			return tmp[i].Usage < tmp[j].Usage
		})
		var ret CPUResourceAllocateInfo
		ret.CPUS = cpus
		for i := 0; i < len(tmp) && cpus > 0 && tmp[i].Usage < 100; i += 1 {
			minus := min(cpus, 100-tmp[i].Usage)
			tmp[i].Usage += minus
			cpus -= minus
			ret.CoreIDs = append(ret.CoreIDs, tmp[i].CoreID)
			ret.CPUSPerCore = append(ret.CPUSPerCore, minus)
		}
		//返回分配结果
		if cpus == 0 {
			p.CoreInfos = tmp
			return &ret
		}
		return nil
	}
	//尝试分配资源
	var res *CPUResourceAllocateInfo
	for res = tryAllocate(); res == nil; {
		//等待核心可用
		p.cond.Wait()
	}
	//返回分配结果
	return *res
}

// 回收资源
func (p *CoreInfo) recycleResource(cpuResource CPUResourceAllocateInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//创建映射表
	mp := make(map[int]int)
	for _, v := range p.CoreInfos {
		mp[v.CoreID] = v.Usage
	}
	//
	for i, x := range cpuResource.CoreIDs {
		mp[x] -= cpuResource.CPUSPerCore[i]
	}
	//更新核心信息
	finalRes := make([]SingleCoreInfo, 0, len(p.CoreInfos))
	for k, v := range mp {
		finalRes = append(finalRes, SingleCoreInfo{
			CoreID: k,
			Usage:  v,
		})
	}
	//更新核心信息
	p.CoreInfos = finalRes
	//通知其他核心可用
	p.cond.Broadcast()
}

// 获取资源快照
func (p *ProfilerInfo) GetResourceSnapshot() string {
	stats := make(map[string]interface{})
	//记录性能
	stats["TaskWait"] = p.taskWait
	stats["Compile"] = p.compile
	stats["Running"] = p.running
	//记录资源
	stats["CoreForCompile"] = make([]SingleCoreInfo, 0)
	stats["CoreForCompile"] = p.compileCore.CoreInfos

	//记录资源
	stats["CoreForRun"] = make([]SingleCoreInfo, 0)
	stats["CoreForRun"] = p.runCore.CoreInfos
	//返回快照
	dataBytes, _ := json.MarshalIndent(stats, "", "  ")
	return string(dataBytes)
}

// 给代码编译分配资源
func (p *ProfilerInfo) AllocateCodeCompileResource(cpus int) CPUResourceAllocateInfo {
	res := p.compileCore.allocateResource(cpus)
	//Log资源情况
	LogInfo(fmt.Sprintf("ProfilerInfo分配资源成功!当前资源情况: %s", p.GetResourceSnapshot()))
	return res
}

// 回收代码编译资源
func (p *ProfilerInfo) RecycleCodeCompileResource(cpuResource CPUResourceAllocateInfo) {
	p.compileCore.recycleResource(cpuResource)
	//Log资源情况
	LogInfo(fmt.Sprintf("ProfilerInfo回收资源成功!当前资源情况: %s", p.GetResourceSnapshot()))
}

// 给代码运行分配资源
func (p *ProfilerInfo) AllocateCodeRunResource(cpus int) CPUResourceAllocateInfo {
	res := p.runCore.allocateResource(cpus)
	//Log资源情况
	LogInfo(fmt.Sprintf("ProfilerInfo分配资源成功!当前资源情况: %s", p.GetResourceSnapshot()))
	return res
}

// 回收代码运行资源
func (p *ProfilerInfo) RecycleCodeRunResource(cpuResource CPUResourceAllocateInfo) {
	p.runCore.recycleResource(cpuResource)
	//Log资源情况
	LogInfo(fmt.Sprintf("ProfilerInfo回收资源成功!当前资源情况: %s", p.GetResourceSnapshot()))
}

func (p *ProfilerInfo) inc(field *int) {
	p.mu.Lock()
	*field++
	p.mu.Unlock()
}
func (p *ProfilerInfo) dec(field *int) {
	p.mu.Lock()
	if *field > 0 {
		*field--
	}
	p.mu.Unlock()
}

func (p *ProfilerInfo) GetTaskWait() int {
	return p.taskWait
}

func (p *ProfilerInfo) IncreaseTaskWait() {
	p.inc(&p.taskWait)
}
func (p *ProfilerInfo) DecreaseTaskWait() {

	p.dec(&p.taskWait)
}

func (p *ProfilerInfo) GetCompile() int {
	return p.compile
}

func (p *ProfilerInfo) IncreaseCompile() {
	p.inc(&p.compile)
}
func (p *ProfilerInfo) DecreaseCompile() {
	p.dec(&p.compile)
}

func (p *ProfilerInfo) GetRunning() int {
	return p.running
}

func (p *ProfilerInfo) IncreaseRunning() {
	p.inc(&p.running)
}
func (p *ProfilerInfo) DecreaseRunning() {

	p.dec(&p.running)
}

// 拒绝从函数调用开始的所有任务
func (p *ProfilerInfo) TaskRefuseMoreTask() {
	p.taskAcquireMu.Lock()
}

func (p *ProfilerInfo) TaskAcceptMoreTask() {
	p.taskAcquireMu.Unlock()
}

func (p *ProfilerInfo) TaskProcess() {
	//
	p.taskAcquireMu.Lock()
	defer p.taskAcquireMu.Unlock()
	//获取任务锁
	p.taskMu.Lock()
	defer p.taskMu.Unlock()
	p.taskSem++
}

func (p *ProfilerInfo) TaskComplete() {
	p.taskMu.Lock()
	defer p.taskMu.Unlock()
	p.taskSem--
}

func (p *ProfilerInfo) TaskAllDone() bool {
	p.taskMu.Lock()
	defer p.taskMu.Unlock()
	return p.taskSem == 0
}

func (p *ProfilerInfo) GetTaskSem() int64 {
	return p.taskSem
}
