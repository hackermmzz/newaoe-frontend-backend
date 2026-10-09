package system

import (
	"errors"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	"os"
	"strconv"
	"strings"
	"time"
)

func systemUpdate(info *global.OJSystemInfo) error {
	//切割版本号
	currentVersionMajor, currentVersionMinor, currentVersionPatch, err := splitVersion(config.Conf.CurrentVersion)
	if err != nil {
		return err
	}
	newVersionMajor, newVersionMinor, newVersionPatch, err := splitVersion(info.Version)
	if err != nil {
		return err
	}
	//比较版本号
	if compareVersion(
		[]int{currentVersionMajor, currentVersionMinor, currentVersionPatch},
		[]int{newVersionMajor, newVersionMinor, newVersionPatch},
	) {
		updateSystem()
	}
	return nil
}

func compareVersion(currentVersion, newVersion []int) bool {
	//比较版本号
	if currentVersion[0] != newVersion[0] {
		return currentVersion[0] < newVersion[0]
	}
	if currentVersion[1] != newVersion[1] {
		return currentVersion[1] < newVersion[1]
	}
	if currentVersion[2] != newVersion[2] {
		return currentVersion[2] < newVersion[2]
	}
	return false
}

func splitVersion(version string) (int, int, int, error) {
	versionStr := strings.Split(version, ".")
	if len(versionStr) != 3 {
		return 0, 0, 0, errors.New("版本号格式错误!")
	}
	major, err := strconv.Atoi(versionStr[0])
	if err != nil {
		return 0, 0, 0, err
	}
	minor, err := strconv.Atoi(versionStr[1])
	if err != nil {
		return 0, 0, 0, err
	}
	patch, err := strconv.Atoi(versionStr[2])
	if err != nil {
		return 0, 0, 0, err
	}
	return major, minor, patch, nil
}

// 目前更新系统只要确保当前所有的任务完成了就可以了
func updateSystem() {
	//阻塞等待所有任务完成
	global.LogInfo("正在等待所有任务完成.......")
	global.Profiler.TaskRefuseMoreTask()
	defer global.Profiler.TaskAcceptMoreTask()
	for global.Profiler.GetTaskSem() > 0 {
		//等待所有任务完成
		time.Sleep(time.Second * time.Duration(1))
	}
	global.LogInfo("所有任务完成，开始更新系统.......")
	//更新系统(直接退出，会有run.py来更新的)
	/*
		err := util.GitPull("./", "judge", 3)
		if err != nil {
			global.LogError("更新系统失败: " + err.Error())
			return
		}
	*/
	global.LogSuccess("更新系统成功!")
	//退出进程
	os.Exit(0)
}
