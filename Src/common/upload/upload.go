package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/oss"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"
)

// 数据上传
func UploadData(filePath string, data []byte, contentType string) bool {
	return oss.OssUploadFileData(filePath, data, contentType)
}

// 单文件上传接口
func UploadFile(filePath string, expireDuration time.Duration) string {
	url := oss.GetUploadFileUrl(filePath, expireDuration)
	if url == "" {
		util.DebugError("UploadFile:生成上传链接失败")
		return ""
	}
	return url
}

// UploadPostPolicyFile 生成带服务端文件大小校验的 POST Policy 上传凭证。
func UploadPostPolicyFile(filePath string, expireDuration time.Duration, minSize int64, maxSize int64) *oss.UploadPolicyResult {
	result := oss.GetUploadPostPolicyUrl(filePath, expireDuration, minSize, maxSize)
	if result == nil {
		util.DebugError("UploadPostPolicyFile:生成上传凭证失败")
	}
	return result
}

// 文件上传接口(返回上传链接)
func UploadPostPolicyFiles(filePath []string, expireDuration []time.Duration, minSize int64, maxSize int64) []*oss.UploadPolicyResult {
	urls := oss.GetUploadPostPolicyUrls(filePath, expireDuration, minSize, maxSize)
	if urls == nil {
		util.DebugError("UploadFiles:生成上传链接失败")
		return nil
	}
	return urls
}

// 文件上传接口(带Redis缓存文件路径)
func UploadFileWithCachePath(id string, filePath []string, expireDuration time.Duration, minSize int64, maxSize int64) ([]*oss.UploadPolicyResult, string) {
	urls := oss.GetUploadPostPolicyUrls(
		filePath,
		util.NewArray(len(filePath), expireDuration),
		minSize,
		maxSize,
	)
	if urls == nil {
		util.DebugError("UploadFile:生成上传链接失败")
		return nil, ""
	}
	//生成redis key
	key := ""
	for _, v := range filePath {
		key = key + "_" + v
	}
	key = fmt.Sprintf("%v_%v", id, util.SHA256(key)) + util.UUID()
	//redis 设置key
	data_bytes, _ := json.Marshal(filePath)
	ok, err := redis.RedisSet(context.Background(), key, data_bytes, expireDuration)
	if err != nil {
		util.DebugError("设置Redis key失败:", err)
		return nil, ""
	}
	if !ok {
		util.DebugError("设置Redis key失败!")
		return nil, ""
	}
	return urls, key
}

// 根据 key 原子获取并删除对应的文件路径。
func UploadParseAndDelKey(key string) ([]string, error) {
	byte_data, exist, err := redis.RedisGetDel(context.Background(), key)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, util.NewError("UploadParseKey: redis key不存在")
	}
	fileNameArr := make([]string, 0)
	if err := json.Unmarshal(byte_data, &fileNameArr); err != nil {
		return nil, err
	}
	return fileNameArr, nil
}
