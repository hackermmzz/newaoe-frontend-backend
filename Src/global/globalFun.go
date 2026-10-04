package global

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"new-aoe-judge/Src/config"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/util"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	//backoff参数
	backOffSleepMS       = 100
	backOffMaxRetryTimes = 6
	backOffMaxBackOffMS  = 1000
	//日志参数
	logFileLine  = 0
	logFileIndex = 0
)

func LogClear() {
	if LogFile == nil {
		return
	}
	_ = LogFile.Truncate(0)
	// Truncate 不一定会自动重置当前文件偏移
	_, _ = LogFile.Seek(0, io.SeekStart)
	//log()
}

func LogSuccess(msg ...interface{}) {
	logColor("success", msg...)
}

func LogError(msg ...interface{}) {
	logColor("error", msg...)
}

func LogInfo(msg ...interface{}) {
	logColor("info", msg...)
}

func logColor(level string, msg ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	message := html.EscapeString(fmt.Sprint(msg...))
	finalMsg := fmt.Sprintf(
		"<div class=\"log %s\"><span class=\"time\">[%s]</span> %s</div>\n",
		level,
		now,
		message,
	)
	log(finalMsg)
}

func log(msg ...interface{}) {
	logMu.Lock()
	defer logMu.Unlock()
	if LogFile == nil {
		return
	}
	//检查是否需要切换日志文件
	if logFileLine >= 500 {
		logFileIndex++
		//打开新的日志文件
		var e error
		var newFile *os.File
		newFile, e = os.OpenFile(util.JoinPath(config.Conf.ProcessLogDir, fmt.Sprintf("ProcessLog%d.html", logFileIndex)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0644)
		if e == nil {
			logFileLine = 0
			//关闭旧文件
			_ = LogFile.Close()
			//设置新的日志文件
			LogFile = newFile
			//fmt.Fprint(LogFile, htmlInfo)
		}
	}
	//
	data := fmt.Sprint(msg...)
	logFileLine += util.GetStringLines(data)
	fmt.Fprint(LogFile, data)
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

// backOff 指数退避重试
// firstSleep：首次休眠毫秒
// retryCnt：最大重试次数
// maxLimit：单次最大休眠毫秒
// fn：待执行业务函数，返回error代表失败
// return：最后一次的错误，全部重试失败返回err
func backOff(sleepMs int, retryCnt int, maxLimit int, fn func() error) error {
	var errs []error = make([]error, 0)
	for i := 0; i < retryCnt; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		errs = append(errs, err)
		// 最后一次，不再休眠
		if i == retryCnt-1 {
			break
		}
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		// 指数增长，不超过maxLimit
		sleepMs = min(sleepMs*2, maxLimit)
	}
	return fmt.Errorf("%v", errs)
}

// PostCodeStatus 更新代码状态
func PostCodeStatus(
	ctx context.Context,
	server grpc_api.CodeClient,
	id string,
	indices int64,
	status int32,
	data string) (*grpc_api.StatusUpdateReply, error) {
	var statusUpdateReply *grpc_api.StatusUpdateReply
	//
	fun := func() error {
		resp, e := server.CodeStatusUpdate(ctx,
			&grpc_api.CodeStatusUpdateRequest{
				Auth:    AuthString(),
				Indices: indices,
				Id:      id,
				Status:  status,
				Data:    data,
			})
		if e != nil {
			return e
		}
		if resp.GetAuth() != "" {
			UpdateAuth(resp.GetAuth())
			statusUpdateReply = resp
		}
		return nil
	}
	//重试6次
	err := backOff(backOffSleepMS, backOffMaxRetryTimes, backOffMaxBackOffMS, fun)
	if err != nil {
		return nil, err
	}
	return statusUpdateReply, nil
}

// 服务器异常的PostCodeStatus
func PostServerErrorStatus(ctx context.Context,
	server grpc_api.CodeClient,
	id string,
	indices int64,
	err error) (*grpc_api.StatusUpdateReply, error) {
	return PostCodeStatus(ctx, server, id, indices, Code_Status_Error, CodeRunStatusInfo{
		Status: Code_Status_Error,
		Data:   fmt.Sprintf("服务器异常:%v", err),
	}.String())
}

// 获取OJ最新的版本等状态
func GetOJSystemInfo(ctx context.Context, server grpc_api.CodeClient) (*OJSystemInfo, error) {
	var reply *grpc_api.OJSystemInfoReply
	fun := func() error {
		var err error
		reply, err = server.OJSystemInfoGet(ctx, &grpc_api.Empty{})
		if err != nil {
			return err
		}
		return nil
	}
	err := backOff(backOffSleepMS, backOffMaxRetryTimes, backOffMaxBackOffMS, fun)
	if err != nil {
		return nil, err
	}
	//解析reply
	var ojSystemInfo OJSystemInfo
	err = json.Unmarshal([]byte(reply.Data), &ojSystemInfo)
	if err != nil {
		return nil, err
	}
	return &ojSystemInfo, nil
}

// UploadFile PUT上传本地文件，直接流式读文件，不全部加载进内存
func UploadFile(url string, filepath string) error {
	//发送请求
	fun := func() error {
		//打开文件
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
	//重试6次
	err := backOff(backOffSleepMS, backOffMaxRetryTimes, backOffMaxBackOffMS, fun)
	if err != nil {
		return err
	}
	return nil
}

// DownloadFile 下载文件
func DownloadFile(url string, filepath string) (*resty.Response, error) {
	var resp *resty.Response
	fun := func() error {
		resp0, e := httpClient.R().SetOutput(filepath).Get(url)
		if e != nil {
			return e
		}
		if !resp0.IsSuccess() {
			return fmt.Errorf("DownloadFile: http status %d", resp0.StatusCode())
		}
		resp = resp0
		return nil
	}
	//重试6次
	err := backOff(backOffSleepMS, backOffMaxRetryTimes, backOffMaxBackOffMS, fun)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
