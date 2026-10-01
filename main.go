package main

import (
	"context"
	"fmt"
	task "new-aoe-judge/Src/Task"
	api_preCompile "new-aoe-judge/Src/api/preCompile"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// dialGRPC 连接 gRPC 服务器
func dialGRPC() (*grpc.ClientConn, error) {
	return grpc.Dial(config.Conf.BaseIP+":"+config.Conf.GRPCPort, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100*1024*1024), grpc.MaxCallSendMsgSize(100*1024*1024)), grpc.WithKeepaliveParams(keepaliveParams()))
}

// keepaliveParams 保持连接参数
func keepaliveParams() keepalive.ClientParameters {
	return keepalive.ClientParameters{Time: 10 * time.Second, Timeout: 5 * time.Second, PermitWithoutStream: true}
}

// CleanUp 清理子进程
func CleanUp() {
	global.Log("正在清理子进程...")
	out, e := exec.Command("docker", "ps", "-q", "--filter", "label=newaoe-judge").Output()
	if e == nil {
		ids := ""
		for _, id := range strings.Fields(string(out)) {
			ids += id + " "
		}
		_ = exec.Command("docker", "kill", ids).Run()
		_ = exec.Command("docker", "rm", ids).Run()
	}
	global.Log("清理完成!")
}

// getAuthAccount 获取认证账号
func getAuthAccount() {
	global.GRPCAuth.ID = os.Getenv("AUTH_ACCOUNT")
	global.GRPCAuth.Password = os.Getenv("AUTH_PASSWORD")
	if global.GRPCAuth.ID == "" {
		global.GRPCAuth.ID = config.Conf.AuthAccount
	}
	if global.GRPCAuth.Password == "" {
		global.GRPCAuth.Password = config.Conf.AuthPassword
	}
	if global.GRPCAuth.ID == "" {
		fmt.Print("账号(非游客账号): ")
		_, _ = fmt.Scanln(&global.GRPCAuth.ID)
	}
	if global.GRPCAuth.Password == "" {
		fmt.Print("密码: ")
		_, _ = fmt.Scanln(&global.GRPCAuth.Password)
	}
}

// PreProcess 预处理
func PreProcess() error {
	// 初始化全局变量
	if e := global.Init(); e != nil {
		return e
	}
	return nil
}

// CloseProcess 关闭进程
func CloseProcess() {
	// 延迟关闭日志文件
	defer func() {
		if global.LogFile != nil {
			global.LogFile.Close()
		}
	}()
}
func main() {
	// 初始化配置
	if e := config.LoadConfig(); e != nil {
		panic(e)
	}
	// 初始化预处理
	if e := PreProcess(); e != nil {
		panic(e)
	}
	// 延迟处理一些事务
	defer CloseProcess()
	// 初始化认证
	getAuthAccount()
	// 初始化运行目录
	_ = os.MkdirAll(config.Conf.RunDir, 0755)
	// 预编译
	if e := api_preCompile.PreCompile(); e != nil {
		global.Log("预编译失败: " + e.Error())
		return
	}
	global.Log("初始化完成!")
	// 连接 gRPC 服务器
	conn, e := dialGRPC()
	if e != nil {
		global.Log("连接 gRPC 失败: " + e.Error())
		return
	}
	defer conn.Close()
	// 创建 gRPC 客户端
	server := grpc_api.NewCodeClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 注册信号处理函数
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Printf("捕获到退出信号，当前等待:%d 编译:%d 运行:%d\n", global.Profiler.TaskWait, global.Profiler.Compile, global.Profiler.Running)
		cancel()
		CleanUp()
	}()
	// 启动任务处理任务
	go task.Task(ctx, server)
	// 等待上下文取消
	<-ctx.Done()
}
