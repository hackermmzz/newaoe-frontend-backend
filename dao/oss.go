package dao

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"newaoe/config"
	"newaoe/util"
	"path"
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
	util.Debug("OSS连接成功!")
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

	util.Debug("OssCheckFileExist:", lastErr.Error())
	return nil
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
	util.Debug("OssGetFileData:", lastErr.Error())
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

	util.Debug("OssGetDownloadFileUrl:", lastErr.Error())
	return ""
}

// 文件上传接口(返回上传链接)
func GetUploadFileUrls(filePath []string, expireDuration []time.Duration) []string {
	var result []string

	for i, path := range filePath {
		var lastErr error
		success := false

		for j := 0; j < config.Conf.OSS.MaxRetry; j++ {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

			u, err := OssClient.PresignedPutObject(
				ctx,
				config.Conf.OSS.BucketName,
				path,
				expireDuration[i],
			)
			cancel()

			if err == nil {
				result = append(result, u.String())
				success = true
				break
			}

			lastErr = err

			if !isRetryable(err) {
				return nil
			}

			backoff(j)
		}

		if !success {
			util.Debug("GetUploadFileUrls:", lastErr.Error())
			return nil
		}
	}

	return result
}
