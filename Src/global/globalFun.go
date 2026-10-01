package global

import (
	"context"
	"fmt"
	"io"
	"net/http"
	grpc_api "new-aoe-judge/Src/grpc"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func Log(msg ...interface{}) {
	logMu.Lock()
	defer logMu.Unlock()
	if LogFile == nil {
		return
	}
	data := fmt.Sprint(msg...)
	fmt.Fprintf(LogFile, "[%s] %v\n", time.Now().Format("2006-01-02 15:04:05"), data)
	_ = LogFile.Sync()
}
func AuthString() string {
	authMu.RLock()
	defer authMu.RUnlock()
	return GRPCAuth.String()
}
func UpdateAuth(v string) {
	authMu.Lock()
	GRPCAuth.Auth = v
	authMu.Unlock()
}

func PostCodeStatus(
	ctx context.Context,
	server grpc_api.CodeClient,
	id string,
	indices int64,
	status int32,
	data string) (*grpc_api.StatusUpdateReply, error) {
	resp, e := server.CodeStatusUpdate(ctx,
		&grpc_api.CodeStatusUpdateRequest{
			Auth:    AuthString(),
			Indices: indices,
			Id:      id,
			Status:  status,
			Data:    data,
		})
	if e != nil {
		return nil, e
	}
	if resp.GetAuth() != "" {
		UpdateAuth(resp.GetAuth())
	}
	return resp, nil
}

// UploadFile PUT上传本地文件，直接流式读文件，不全部加载进内存
func UploadFile(url string, filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	//获取文件大小
	info, err := f.Stat()
	if err != nil {
		return err
	}
	filesize := info.Size()
	//如果文件大小为0，则使用http.NoBody
	var body io.Reader = f
	if filesize == 0 {
		body = http.NoBody
	}
	//创建PUT请求
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return err
	}
	req.ContentLength = filesize
	//发送请求
	resp, err := httpClient.GetClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf(
			"UploadFile: http status %d: %s",
			resp.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}

	return nil
}

// DownloadFile 下载文件
func DownloadFile(url string, filepath string) (*resty.Response, error) {
	resp, e := httpClient.R().SetOutput(filepath).Get(url)
	if e != nil {
		return resp, e
	}
	if !resp.IsSuccess() {
		return resp, fmt.Errorf("DownloadFile: http status %d", resp.StatusCode())
	}
	return resp, nil
}
