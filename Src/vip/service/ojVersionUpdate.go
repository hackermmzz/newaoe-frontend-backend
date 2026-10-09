package service

import (
	"errors"
	"newaoe/Src/codeRun/service"
	"strconv"
	"strings"
)

func OJVersionUpdate(version string) error {
	//检查参数是否合法
	if err := checkOJVersionAvaliable(version); err != nil {
		return err
	}
	//
	return service.SetOJSystemInfo(service.OJSystemInfo{
		Version: version,
	})
}

func checkOJVersionAvaliable(version string) error {
	versionStr := strings.Split(version, ".")
	if len(versionStr) != 3 {
		return errors.New("版本号格式错误!")
	}
	_, err := strconv.Atoi(versionStr[0])
	if err != nil {
		return err
	}
	_, err = strconv.Atoi(versionStr[1])
	if err != nil {
		return err
	}
	_, err = strconv.Atoi(versionStr[2])
	if err != nil {
		return err
	}
	return nil
}
