package api_codeGet

import (
	"context"
	"fmt"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/util"
	"os"
)

type StudentCode struct {
	OK       bool
	Msg      string
	Header   string
	Source   string
	ID       string
	Indices  int64
	RunType  int64
	RunDir   string
	BuildDir string
}

func GetOneStudentCode(ctx context.Context, server grpc_api.CodeClient) StudentCode {
	resp, e := server.GetCode(ctx, &grpc_api.CodeRequest{Auth: global.AuthString()})
	if e != nil {
		return StudentCode{Msg: fmt.Sprintf("服务器异常: %v", e)}
	}
	if !resp.GetOk() {
		return StudentCode{Msg: resp.GetMsg()}
	}
	//创建运行目录和编译目录
	runDir, buildDir, e := createRunDir(resp.GetId(), resp.GetIndices())
	if e != nil {
		return StudentCode{Msg: fmt.Sprintf("创建运行目录失败: %v", e)}
	}
	//下载头文件
	headerPath := util.JoinPath(buildDir, "UsrAI.h")
	_, e = global.DownloadFile(resp.GetHeaderUrl(), headerPath)
	if e != nil {
		return StudentCode{Msg: "下载头文件失败!"}
	}
	//下载源文件
	sourcePath := util.JoinPath(buildDir, "UsrAI.cpp")
	_, e = global.DownloadFile(resp.GetSourceUrl(), sourcePath)
	if e != nil {
		return StudentCode{Msg: "下载源文件失败!"}
	}
	global.UpdateAuth(resp.GetAuth())
	return StudentCode{
		OK:       true,
		Msg:      "获取代码成功!",
		Header:   headerPath,
		Source:   sourcePath,
		ID:       resp.GetId(),
		Indices:  resp.GetIndices(),
		RunType:  resp.GetRuntype(),
		RunDir:   runDir,
		BuildDir: buildDir,
	}
}

func createRunDir(id string, indices int64) (string, string, error) {
	//创建运行目录
	dirName := fmt.Sprintf("%s_%d_%s", id, indices, util.UUID())
	rundir := util.JoinPath(config.Conf.RunDir, dirName)
	if e := os.MkdirAll(rundir, 0755); e != nil {
		return "", "", fmt.Errorf("创建运行目录失败: %s, %s", rundir, e.Error())
	}
	//创建编译目录
	buildDir := util.JoinPath(rundir, "build")
	if e := os.MkdirAll(buildDir, 0755); e != nil {
		return "", "", fmt.Errorf("创建编译目录失败: %s, %s", buildDir, e.Error())
	}
	return rundir, buildDir, nil
}
