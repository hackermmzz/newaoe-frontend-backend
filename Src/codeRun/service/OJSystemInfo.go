package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"newaoe/Src/codeRun/grpc/grpc_api"
	MapDao "newaoe/Src/map/dao"
	"newaoe/Src/util"
	"strconv"
	"strings"
)

type OJSystemInfo struct {
	Version string `json:"version"` //最新的版本号
}

func GetOJSystemInfo() (*grpc_api.OJSystemInfoReply, error) {
	var msg OJSystemInfo
	//
	value, err := MapDao.MapGet(nil, "OJSystemInfo")
	if err != nil {
		return nil, err
	} else {
		json.Unmarshal([]byte(value), &msg)
	}
	//
	finalMsg, _ := json.Marshal(msg)
	return &grpc_api.OJSystemInfoReply{
		Data: string(finalMsg),
	}, nil
}

func SetOJSystemInfo(info OJSystemInfo) error {
	//获取旧的version
	oldData, err := GetOJSystemInfo()
	if err != nil {
		//判断key是否存在，不存在则创建
		exist, err := MapDao.MapExist(nil, "OJSystemInfo")
		if err != nil {
			return err
		}
		if !exist {
			data, _ := json.Marshal(OJSystemInfo{
				Version: "1.0.0",
			})
			err := MapDao.MapInsert(nil, "OJSystemInfo", string(data))
			if err != nil {
				return err
			}
			//再次查询数据
			oldData, err = GetOJSystemInfo()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	var oldStatus OJSystemInfo
	json.Unmarshal([]byte(oldData.Data), &oldStatus)
	//切割旧版本和新版本的version
	oa, ob, oc, err := splitVersion(oldStatus.Version)
	if err != nil {
		return err
	}
	na, nb, nc, err := splitVersion(info.Version)
	if err != nil {
		return err
	}
	if !compareVersion([]int{oa, ob, oc}, []int{na, nb, nc}) {
		return util.NewError(fmt.Sprintf("版本比之前的版本还旧:%v:%v", oldStatus.Version, info.Version))
	}
	//更新数据
	databytes, _ := json.Marshal(info)
	err = MapDao.MapUpdate(nil, "OJSystemInfo", string(databytes))
	if err != nil {
		return err
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
