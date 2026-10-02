package api_codeCompile

import (
	"fmt"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	"new-aoe-judge/Src/util"
	"os"
	"os/exec"
)

type CodeCompileResult struct {
	OK    bool
	Msg   string
	Error error
}

func CodeCompile(
	cpuResource global.CPUResourceAllocateInfo,
	id string,
	indices int64,
	buildDir string,
	logfile *os.File,
	runtype int64,
) CodeCompileResult {
	//格式化编译脚本
	script, e := util.FormatFile(util.JoinPath("assets", "bash", "codecompile.sh"), map[string]interface{}{
		"DebugMode": runtype == int64(global.CodeRunTypeDebug),
	})
	if e != nil {
		return CodeCompileResult{OK: false, Error: fmt.Errorf("format失败:%w", e)}
	}
	//设置编译CPU限制
	args := []string{
		"run",
		"--cpuset-cpus", cpuResource.GetCores(),
		"--cpus", cpuResource.GetCPUS(),
		"--name", fmt.Sprintf("codecompile_%s_%d_%s", id, indices, util.UUID()),
		"--label", "newaoe-judge",
		"--rm",
		"-v", config.Conf.NewAOEFolder + ":/app/project:ro",
		"-v", buildDir + ":/app/build",
		"-w", "/app/", config.Conf.NewAOEDockerImg,
		"bash", "-c", script,
	}
	//执行编译命令
	cmd := exec.Command("docker", args...)
	cmd.Stdout = logfile
	cmd.Stderr = logfile
	e = cmd.Run()
	if e != nil {
		if _, ok := e.(*exec.ExitError); !ok {
			return CodeCompileResult{OK: false, Error: fmt.Errorf("CodeCompile运行容器失败:%w", e)}
		}
	}
	//拿到状态码
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 || e != nil {
		msg := util.ReadFileAnyBytes(logfile, 1024)
		return CodeCompileResult{OK: false, Msg: string(msg)}
	}
	//编译成功
	return CodeCompileResult{OK: true, Msg: ""}
}
