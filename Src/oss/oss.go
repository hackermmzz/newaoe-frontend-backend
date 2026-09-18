package oss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"newaoe/Src/config"
	"newaoe/Src/util"

	"path"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var OssClient *minio.Client

func OssInit() {
	var err error
	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
	OssClient, err = minio.New(config.Conf.OSS.Host, &minio.Options{
		Creds:     credentials.NewStaticV4(config.Conf.OSS.AccessKey, config.Conf.OSS.SecretKey, ""),
		Secure:    false,
		Transport: tr,
	})
	if err != nil {
		panic("OSS连接失败!" + err.Error())
	}
	// 检查桶是否存在，如果不存在则创建
	exists, err := OssClient.BucketExists(context.Background(), config.Conf.OSS.BucketName)
	if err != nil {
		panic("检查OSS桶失败!" + err.Error())
	}
	if !exists {
		panic(config.Conf.OSS.BucketName + "桶不存在!")
	}
	//获取公共目录下所有用户默认头像
	avatarPath := path.Join(config.Conf.OSS.PublicBaseFolder, config.Conf.User.UserAvatarFolder)
	allAvatar := OssClient.ListObjects(context.Background(), config.Conf.OSS.BucketName, minio.ListObjectsOptions{
		Prefix:    avatarPath,
		Recursive: true,
	})
	//将默认头像路径添加到配置中（实际使用时会随机选择一个）
	for obj := range allAvatar {
		if obj.Err != nil {
			panic("获取默认头像失败:" + obj.Err.Error())
		} else {
			config.Conf.User.UserDefalutAvatar = append(config.Conf.User.UserDefalutAvatar, obj.Key)
		}
	}
	if len(config.Conf.User.UserDefalutAvatar) == 0 {
		panic(fmt.Sprintf("获取默认头像失败: %v桶的路径%v内没有默认头像资源!", config.Conf.OSS.BucketName, avatarPath))
	}
	//
	util.DebugSuccess("OSS连接成功!")
}

// 判断错误是否值得重试而不是无脑重试
func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	errResp := minio.ToErrorResponse(err)

	switch errResp.Code {
	case "NoSuchKey", "NoSuchBucket", "AccessDenied":
		return false
	}
	return true
}

// 简单指数退避
func backoff(i int) {
	time.Sleep(time.Duration(1<<i) * 100 * time.Millisecond)
}

// OssCheckFileExist 判断文件是否存在
func OssCheckFileExist(filePath string) *minio.ObjectInfo {
	var lastErr error

	for i := 0; i < config.Conf.OSS.MaxRetry; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		info, err := OssClient.StatObject(ctx, config.Conf.OSS.BucketName, filePath, minio.StatObjectOptions{})
		cancel()

		if err == nil {
			return &info
		}

		lastErr = err
		errResp := minio.ToErrorResponse(err)

		// 文件不存在属于“正常情况”，不需要重试
		if errResp.Code == "NoSuchKey" {
			return nil
		}

		// 不可重试错误，直接返回
		if !isRetryable(err) {
			return nil
		}

		backoff(i)
	}

	util.DebugError("OssCheckFileExist:", lastErr.Error())
	return nil
}

// 批量并发检查 OSS 文件是否存在，并返回文件信息
func OssCheckFilesExist(filePaths []string) map[string]*minio.ObjectInfo {
	result := make(map[string]*minio.ObjectInfo, len(filePaths))
	var mu sync.Mutex
	var wg sync.WaitGroup
	// 限制并发数量
	sem := make(chan struct{}, 20)
	for _, filePath := range filePaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() {
				<-sem
			}()
			var info *minio.ObjectInfo
			for i := 0; i < config.Conf.OSS.MaxRetry; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				objInfo, err := OssClient.StatObject(ctx, config.Conf.OSS.BucketName, path, minio.StatObjectOptions{})
				cancel()
				if err == nil {
					info = &objInfo
					break
				}
				errResp := minio.ToErrorResponse(err)
				// 文件不存在，不需要重试
				if errResp.Code == "NoSuchKey" {
					break
				}
				// 不可重试错误
				if !isRetryable(err) {
					break
				}
				backoff(i)
			}
			mu.Lock()
			result[path] = info
			mu.Unlock()
		}(filePath)
	}
	wg.Wait()
	return result
}

// 我们这里全量读取，所以要注意文件不能很大
func OssGetFileData(filePath string) []byte {
	var lastErr error

	for i := 0; i < config.Conf.OSS.MaxRetry; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		reader, err := OssClient.GetObject(ctx, config.Conf.OSS.BucketName, filePath, minio.GetObjectOptions{})
		if err != nil {
			cancel()
			lastErr = err

			if !isRetryable(err) {
				return nil
			}
			backoff(i)
			continue
		}
		// ⚠️ 注意：GetObject 是懒加载，这里才真正发请求
		data, err := io.ReadAll(reader)
		reader.Close()
		cancel()
		if err == nil {
			return data
		}
		lastErr = err
		// 判断是否文件不存在
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return nil
		}
		if !isRetryable(err) {
			return nil
		}
		backoff(i)
	}
	//无法获取
	util.DebugError("OssGetFileData:", lastErr.Error())
	return nil
}

// 获取下载链接
func OssGetDownloadFileUrl(filePath string, expireDuration time.Duration, attachment bool) string {
	var lastErr error

	var reqParams url.Values
	if attachment {
		reqParams = make(url.Values)
		reqParams.Set("response-content-disposition", "attachment")
	}
	//
	for i := 0; i < config.Conf.OSS.MaxRetry; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		u, err := OssClient.PresignedGetObject(
			ctx,
			config.Conf.OSS.BucketName,
			filePath,
			expireDuration,
			reqParams,
		)
		cancel()

		if err == nil {
			return u.String()
		}

		lastErr = err

		if !isRetryable(err) {
			return ""
		}

		backoff(i)
	}

	util.DebugError("OssGetDownloadFileUrl:", lastErr.Error())
	return ""
}

// 文件上传接口(返回上传链接) 批量获取上传链接（并发）
func GetUploadFileUrls(filePath []string, expireDuration []time.Duration) []string {
	if len(filePath) != len(expireDuration) {
		return nil
	}
	result := make([]string, len(filePath))
	var wg sync.WaitGroup
	var mu sync.Mutex
	// 控制并发数
	sem := make(chan struct{}, 20)
	success := true
	for i, path := range filePath {
		wg.Add(1)
		go func(index int, objectPath string, expire time.Duration) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() {
				<-sem
			}()
			var lastErr error
			for j := 0; j < config.Conf.OSS.MaxRetry; j++ {
				ctx, cancel := context.WithTimeout(
					context.Background(),
					2*time.Second,
				)
				u, err := OssClient.PresignedPutObject(
					ctx,
					config.Conf.OSS.BucketName,
					objectPath,
					expire,
				)
				cancel()
				if err == nil {
					mu.Lock()
					result[index] = u.String()
					mu.Unlock()
					return
				}
				lastErr = err
				if !isRetryable(err) {
					break
				}
				backoff(j)
			}
			util.DebugError("GetUploadFileUrls:", lastErr.Error())
			mu.Lock()
			success = false
			mu.Unlock()
		}(i, path, expireDuration[i])
	}
	wg.Wait()
	if !success {
		return nil
	}
	return result
}

// 获取一个上传链接
func GetUploadFileUrl(filePath string, expireDuration time.Duration) string {
	for j := 0; j < config.Conf.OSS.MaxRetry; j++ {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			2*time.Second,
		)
		u, err := OssClient.PresignedPutObject(
			ctx,
			config.Conf.OSS.BucketName,
			filePath,
			expireDuration,
		)
		cancel()
		if err == nil {
			return u.String()
		}
		if !isRetryable(err) {
			break
		}
		backoff(j)
	}
	return ""
}

// 批量获取下载链接（并发）
func OssGetDownloadFileUrls(filePaths []string, expireDuration time.Duration, attachment bool) []string {

	urls := make([]string, len(filePaths))
	var reqParams url.Values
	if attachment {
		reqParams = make(url.Values)
		reqParams.Set("response-content-disposition", "attachment")
	}
	var wg sync.WaitGroup
	// 限制并发数量
	sem := make(chan struct{}, 20)
	for index, filePath := range filePaths {
		wg.Add(1)
		go func(i int, path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() {
				<-sem
			}()
			var lastErr error
			for retry := 0; retry < config.Conf.OSS.MaxRetry; retry++ {
				ctx, cancel := context.WithTimeout(
					context.Background(),
					2*time.Second,
				)
				u, err := OssClient.PresignedGetObject(
					ctx,
					config.Conf.OSS.BucketName,
					path,
					expireDuration,
					reqParams,
				)
				cancel()
				if err == nil {
					urls[i] = u.String()
					return
				}
				lastErr = err
				if !isRetryable(err) {
					break
				}
				backoff(retry)
			}
			util.DebugError(
				"OssGetDownloadFileUrl:",
				path,
				lastErr,
			)
			// 失败保持空字符串
			urls[i] = ""
		}(index, filePath)
	}
	wg.Wait()
	return urls
}

// OssUploadFileData 上传文件数据到OSS
// filePath: OSS存储路径
// data: 文件内容
func OssUploadFileData(filePath string, data []byte, contentType string) bool {
	var lastErr error
	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}
	for i := 0; i < config.Conf.OSS.MaxRetry; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := OssClient.PutObject(
			ctx,
			config.Conf.OSS.BucketName,
			filePath,
			bytes.NewReader(data),
			int64(len(data)),
			opts,
		)
		cancel()
		if err == nil {
			return true
		}
		lastErr = err
		if !isRetryable(err) {
			break
		}
		backoff(i)
	}
	util.DebugError("OssUploadFileData:", filePath, lastErr)
	return false
}
