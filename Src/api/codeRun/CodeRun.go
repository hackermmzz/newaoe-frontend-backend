package api_codeRun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	"new-aoe-judge/Src/util"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func getExitReason(code int) string {
	m := map[int]string{0: "Normal Exit", 1: "General Error", 125: "Docker Run Error", 126: "Command Cannot Execute", 127: "Command Not Found", 130: "SIGINT", 132: "SIGILL (Illegal Instruction)", 134: "SIGABRT (Abort / Assertion Failed)", 136: "SIGFPE (Arithmetic Exception / Divide by Zero)", 137: "SIGKILL (Possibly OOM / Memory Limit)", 139: "SIGSEGV (Segmentation Fault)", 143: "SIGTERM"}
	if x, ok := m[code]; ok {
		return x
	}
	return fmt.Sprintf("Unknown Exit Code (%d)", code)
}

type CodeRunRetInfo struct {
	ExitReason string
	ExitCode   int
	Err        error
	TimeCost   int64 //秒数
}

func CodeRun(
	cpuResource global.CPUResourceAllocateInfo,
	id string,
	indices int64,
	rundir string,
	buildDir string,
	crashFile *os.File,
	logFile *os.File,
	recordFile *os.File,
	resultFile *os.File,
	debugFile *os.File,
) *CodeRunRetInfo {
	workdir := "/tmp/project"
	container := fmt.Sprintf("coderun_%s_%d_%s", id, indices, util.UUID())
	// 挂载文件
	mounts := []string{
		"-v", util.JoinPath(config.Conf.NewAOEFolder, "res.rcc") + ":" + workdir + "/res.rcc",
		"-v", util.JoinPath(config.Conf.NewAOEFolder, "config.json") + ":" + workdir + "/config.json:ro",
		"-v", util.JoinPath(buildDir, "newAOE") + ":" + workdir + "/newAOE:ro",
		"-v", util.JoinPath(resultFile.Name()) + ":" + workdir + "/" + config.Conf.RunResultFileName,
		"-v", util.JoinPath(recordFile.Name()) + ":" + workdir + "/" + config.Conf.RecordFileName,
		"-v", util.JoinPath(debugFile.Name()) + ":" + workdir + "/" + config.Conf.RunDebugLogOutputFileName,
		"-v", util.JoinPath(crashFile.Name()) + ":" + workdir + "/" + config.Conf.CrashLogFileName,
	}
	// 挂载地图
	for _, x := range global.MapFiles {
		mounts = append(mounts, "-v", util.JoinPath(config.Conf.NewAOEFolder, x)+":"+workdir+"/"+x+":ro")
	}
	// 生成脚本
	script, e := util.FormatFile(util.JoinPath("assets", "bash", "coderun.sh"), map[string]interface{}{
		"id":                        id,
		"indices":                   strconv.FormatInt(indices, 10),
		"workdir":                   workdir,
		"RunResultFileName":         config.Conf.RunResultFileName,
		"RecordOutputFileName":      config.Conf.RecordFileName,
		"RunDebugLogOutputFileName": config.Conf.RunDebugLogOutputFileName,
		"CrashLogFileName":          config.Conf.CrashLogFileName,
		"AOERunSpeed":               config.Conf.AOERunSpeed,
	},
	)

	if e != nil {
		return &CodeRunRetInfo{
			Err: fmt.Errorf("CodeRun生成脚本失败: %w", e),
		}
	}
	// 构建命令
	args := []string{
		"run",
		"--cpuset-cpus", cpuResource.GetCores(),
		"--cpus", cpuResource.GetCPUS(),
		"--name", container,
		"--label", "newaoe-judge",
		"-m", config.Conf.RunMemoryLimit,
		"--memory-swap", config.Conf.RunMemoryLimit,
		"--tmpfs", "/tmp:rw,size=" + config.Conf.RunDiskLimit,
	}
	args = append(args, mounts...)
	args = append(args, "-w", workdir, config.Conf.NewAOEDockerImg, "bash", "-c", script)
	//
	var dockerStateCMD *exec.Cmd
	var dockerErr error
	var outBuf, errBuf bytes.Buffer
	var ret CodeRunRetInfo
	var exitCode int
	// 记录开始时间
	startTime := time.Now()
	dockerDone := make(chan struct{}, 1)
	go func() {
		defer func() {
			dockerDone <- struct{}{}
		}()
		// 运行命令
		cmd := exec.Command("docker", args...)
		cmd.Stdout = logFile
		cmd.Stderr = logFile
		if e = cmd.Run(); e != nil {
			_, ok := e.(*exec.ExitError)
			if !ok {
				dockerErr = fmt.Errorf("CodeRun运行容器失败: %w", e)
				return
			}
		}
		//获取容器状态码
		dockerStateCMD = exec.Command(
			"docker",
			"inspect",
			container,
		)
		dockerStateCMD.Stdout = &outBuf
		dockerStateCMD.Stderr = &errBuf
		if e = dockerStateCMD.Run(); e != nil {
			dockerErr = fmt.Errorf("CodeRun获取容器状态失败: %w", e)
			return
		}
		exitCode = dockerStateCMD.ProcessState.ExitCode()
	}()
	//移除容器
	defer func() {
		//移除容器
		if e = exec.Command("docker", "rm", container).Run(); e != nil {
			global.LogError(fmt.Errorf("CodeRun移除容器失败: %w", e))
		}
	}()
	//
	select {
	case <-dockerDone:
		if dockerErr != nil {
			ret = CodeRunRetInfo{
				Err: dockerErr,
			}
		} else {
			//判断是否是被OOM了
			var info []map[string]interface{}
			if e = json.Unmarshal(outBuf.Bytes(), &info); e != nil {
				ret = CodeRunRetInfo{
					Err: fmt.Errorf("CodeRun解析容器状态失败: %w", e),
				}
			}
			oom := info[0]["State"].(map[string]interface{})["OOMKilled"]
			if oom.(bool) {
				exitCode = 137
			}
			//获取退出原因
			reason := getExitReason(exitCode)
			ret = CodeRunRetInfo{
				ExitReason: reason,
				ExitCode:   exitCode,
			}
		}
	case <-time.After(time.Duration(config.Conf.RunTimeout) * time.Second):
		//停止容器
		if e = exec.Command("docker", "kill", container).Run(); e != nil {
			ret = CodeRunRetInfo{
				Err: fmt.Errorf("CodeRun停止容器失败: %w", e),
			}
		}
		//返回超时
		ret = CodeRunRetInfo{
			Err: fmt.Errorf("CodeRun获取容器状态超时"),
		}
	}
	//返回结果
	ret.TimeCost = int64(time.Since(startTime).Seconds())
	return &ret
}
