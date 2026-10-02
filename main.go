package main

import (
	"bufio"
	"context"
	"fmt"
	api_preCompile "new-aoe-judge/Src/api/preCompile"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/task"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// waitForAnyKey is line-buffered on non-Windows terminals; press a key and Enter.
func waitForAnyKey() error {
	_, err := bufio.NewReader(os.Stdin).ReadByte()
	return err
}

// dialGRPC 连接 gRPC 服务器
func dialGRPC() (*grpc.ClientConn, error) {
	return grpc.Dial(config.Conf.BaseIP+":"+config.Conf.GRPCPort, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100*1024*1024), grpc.MaxCallSendMsgSize(100*1024*1024)), grpc.WithKeepaliveParams(keepaliveParams()))
}

// keepaliveParams 保持连接参数
func keepaliveParams() keepalive.ClientParameters {
	return keepalive.ClientParameters{Time: 10 * time.Second, Timeout: 5 * time.Second, PermitWithoutStream: true}
}

// cleanUp 清理子进程
func cleanUp() {
	global.Log("正在清理子进程...")
	out, e := exec.Command("docker", "ps", "-q", "--filter", "label=newaoe-judge").Output()
	if e == nil {
		ids := ""
		for _, id := range strings.Fields(string(out)) {
			ids += id + " "
		}
		global.Log(fmt.Sprintf("所有容器id为:%s", ids))
		e0 := exec.Command("docker", "kill", ids).Run()
		e1 := exec.Command("docker", "rm", ids).Run()
		if e0 != nil {
			global.Log("清理子进程失败: " + e0.Error())
		}
		if e1 != nil {
			global.Log("清理子进程失败: " + e1.Error())
		}
	} else {
		global.Log("清理子进程失败: " + e.Error())
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

// preProcess 预处理
func preProcess() error {
	// 初始化全局变量
	if e := global.Init(); e != nil {
		return e
	}
	return nil
}

// closeProcess 关闭进程
func closeProcess() {
	// 延迟关闭日志文件
	defer func() {
		if global.LogFile != nil {
			global.LogFile.Close()
		}
	}()
}

// 打印任务情况
func getTaskStatus() string {
	return fmt.Sprintf(
		"当前等待:%d 编译:%d/%d 运行:%d/%d\n",
		global.Profiler.TaskWait,
		global.Profiler.Compile,
		len(global.CompileWaitQueue)+global.Profiler.Compile,
		global.Profiler.Running,
		len(global.RunWaitQueue)+global.Profiler.Running,
	)
}
func main() {
	// 初始化配置
	if e := config.LoadConfig(); e != nil {
		panic(e)
	}
	// 初始化预处理
	if e := preProcess(); e != nil {
		panic(e)
	}
	// 延迟处理一些事务
	defer closeProcess()
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
		waitingConfirm := false
		var keyPressed <-chan error
		clear := func() {
			fmt.Println("再次捕获到 Ctrl+C，正在清理并退出...")
			signal.Stop(sig)
			cancel()
			cleanUp()
		}
		for {
			select {
			case <-sig:
				if waitingConfirm {
					clear()
					return
				}

				waitingConfirm = true
				fmt.Println(
					"捕获到退出信号，按任意其他键继续，按 Ctrl+C 退出。\n" + getTaskStatus(),
				)

				keyDone := make(chan error, 1)
				keyPressed = keyDone
				go func() {
					keyDone <- waitForAnyKey()
				}()

			case err := <-keyPressed:
				keyPressed = nil
				waitingConfirm = false

				if err != nil {
					clear()
				} else {
					fmt.Println("收到其他按键，继续运行...")
				}
			}
		}
	}()
	// 启动任务处理任务
	go task.Task(ctx, server)
	//启动日志记录任务情况
	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Println(getTaskStatus())
		}
	}()
	// 等待上下文取消
	<-ctx.Done()
}
