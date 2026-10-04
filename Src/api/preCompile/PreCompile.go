package api_preCompile

import (
	"errors"
	"fmt"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	"new-aoe-judge/Src/util"
	"os/exec"
)

func PreCompile() error {
	// 更新new-aoe目录
	if err := pullNewAOE(); err != nil {
		return err
	}
	// 读取预编脚本
	script := util.ReadAnyText(util.JoinPath("assets", "bash", "precompile.sh"))
	if script == "" {
		return errors.New("预编脚本读取失败")
	}
	cmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"--name", fmt.Sprintf("PreCompile_%s", util.UUID()),
		"--label", "newaoe-judge",
		"-v", config.Conf.NewAOEFolder+":/app/newaoe",
		"-w", "/app",
		config.Conf.NewAOEDockerImg,
		"bash", "-c", string(script),
	)
	cmd.Stdout = global.LogFile
	cmd.Stderr = global.LogFile
	err := cmd.Run()
	if err != nil {
		return err
	}
	//判断退出码
	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("预编失败，退出码：%d", cmd.ProcessState.ExitCode())
	}
	//预编成功
	return nil
}

// 更新new-aoe目录
func pullNewAOE() error {
	return util.GitPull(config.Conf.NewAOEFolder, "main", 3)
}
