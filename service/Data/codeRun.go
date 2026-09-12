package data

import (
	"encoding/json"
	"fmt"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"
	"time"
)

// 根据Redis缓存判断用户是否可以运行/提交代码
func CodeRunRecordExist(id string) bool {
	exist, err := dao.RDB.Exists(dao.RDB.Context(), fmt.Sprintf("CodeRunOrSubmit:%v", id)).Result()
	//
	if err != nil {
		return false
	}
	return exist == 1
}

// 加入Redis缓存判断是否可以运行/提交代码
func CodeRunRecordAdd(id string) bool {
	err := dao.RDB.Set(dao.RDB.Context(), fmt.Sprintf("CodeRunOrSubmit:%v", id), "", time.Duration(config.Conf.Code.CodeSubmitInterval)*time.Second).Err()
	//
	if err != nil {
		return false
	}
	return true
}

// 获取用户代码提交/运行 的TTL
func CodeRunRecordTTL(id string) int64 {
	ttl, err := dao.RDB.TTL(dao.RDB.Context(), fmt.Sprintf("CodeRunOrSubmit:%v", id)).Result()
	if err != nil {
		util.Debug("CodeRunRecordTTL:", err)
		return int64(1e9)
	}
	return int64(ttl.Seconds())
}

// 根据索引获取代码文件
func CodeRunGetByIndices(indices int) *dao.CodeRunInfo {
	ctx := dao.RDB.Context()
	ret := &dao.CodeRunInfo{}
	key := fmt.Sprintf("CodeRunIndices:%v", indices)
	if info, exist := dao.RedisGet(ctx, key); exist {
		json.Unmarshal([]byte(info), ret)
		return ret
	}
	// 从数据库中获取
	duration := time.Duration(config.Conf.Code.CodeSubmitInterval) * time.Second // 缓存时间和提交间隔相同
	if ret = dao.CodeRunGetByIndices(indices); ret != nil {
		// 缓存到Redis
		dt, _ := json.Marshal(*ret)
		// 缓存到Redis
		dao.RedisSet(ctx, key, dt, duration)
		return ret
	}
	return nil
}
