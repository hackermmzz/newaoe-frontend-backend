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

// 文件上传接口(返回上传链接)
func UploadFile(filePath []string, expireDuration []time.Duration) bool {
	urls := oss.GetUploadFileUrls(filePath, expireDuration)
	if urls == nil {
		util.DebugError("UploadFile:生成上传链接失败")
		return false
	}
	return true
}

// 文件上传接口(带Redis缓存文件路径)
func UploadFileWithCachePath(id string, filePath []string, expireDuration time.Duration) ([]string, string) {
	urls := oss.GetUploadFileUrls(filePath, util.NewArray(len(filePath), expireDuration))
	if urls == nil {
		util.DebugError("UploadFile:生成上传链接失败")
		return nil, ""
	}
	//生成redis key
	key := ""
	for _, v := range filePath {
		key = key + "_" + v
	}
	key = fmt.Sprintf("%v_%v", id, util.SHA256(key))
	//redis 设置key
	data_bytes, _ := json.Marshal(filePath)
	if !redis.RedisSet(context.Background(), key, data_bytes, expireDuration) {
		util.DebugError("设置Redis key失败!")
		return nil, ""
	}
	return urls, key
}

// 根据ctx字段里面的key解析获取得到文件的数组(不会发送消息)
func UploadParseAndDelKey(body []byte, keyName string, id string) []string {
	//解析Post请求体
	postMap := make(map[string]string)
	if !util.JsonBytes(body, &postMap) {
		return nil
	}
	//从redis里面获取信息
	key, exist := postMap[keyName]
	if !exist {
		util.DebugError("UploadParseKey:数据格式错误!")
		return nil
	}
	byte_data, exist := redis.RedisGet(context.Background(), key)
	if !exist {
		util.DebugError("UploadParseKey:redis key不存在!")
		return nil
	}
	fileNameArr := make([]string, 0)
	err := json.Unmarshal(byte_data, &fileNameArr)
	if err != nil {
		util.DebugError("UploadParseKey:", err)
	}
	if !redis.RedisDel(context.Background(), key) {
		util.DebugError("UploadParseKey:Redis key 删除出问题!")
	}
	return fileNameArr
}
