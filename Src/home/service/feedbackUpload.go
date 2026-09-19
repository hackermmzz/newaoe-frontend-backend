package service

import (
	"fmt"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/home/dao"
	"newaoe/Src/home/model"
	"newaoe/Src/oss"
	"newaoe/Src/util"
	"path"
	"time"
)

func FeedbackGetDownloadLink(indices int, html string) (string, error) {
	//获取反馈的基础路径
	feedbackinfo := dao.FeedbackGetByIndices(nil, indices)
	if feedbackinfo == nil {
		return "", util.NewError("数据库获取反馈记录失败!", indices)
	}
	//上传html数据
	htmlpath := path.Join(feedbackinfo.BaseFolder, "feedback.html")
	if !oss.OssUploadFileData(htmlpath, []byte(html), "text/html; charset=utf-8") {
		return "", util.NewError("上传html失败!")
	}
	//获取链接(设置永不过期)
	url := oss.OssGetDownloadFileUrl(htmlpath, time.Duration(7*24)*time.Hour, false)
	if url == "" {
		return "", util.NewError("获取html下载链接失败!")
	}
	return url, nil
}

type FeedbackUploadLinkInfo struct {
	Indices           int      `json:"indices"`
	UploadImageUrls   []string `json:"uploadimageurls"`
	UploadFileUrls    []string `json:"uploadfileurls"`
	UploadVideoUrls   []string `json:"uploadvideourls"`
	DownloadImageUrls []string `json:"downloadimageurls"`
	DownloadFileUrls  []string `json:"downloadfileurls"`
	DownloadVideoUrls []string `json:"downloadvideourls"`
}

func FeedbackGetUploadUrls(id string, images []string, files []string, videos []string) (*FeedbackUploadLinkInfo, error) {
	var ret FeedbackUploadLinkInfo
	//生成文件名
	dirName := id + "_" + util.UTC_Time().Format("20060102_150405")
	dir := path.Join(config.Conf.OSS.PublicBaseFolder, config.Conf.Feedback.FeedbackFolder, dirName)
	imagespath := make([]string, len(images))
	filespath := make([]string, len(files))
	videospath := make([]string, len(videos))
	for i := range imagespath {
		imagespath[i] = fmt.Sprintf("%v/image_%d_%s", dir, i, images[i])
	}
	for i := range filespath {
		filespath[i] = fmt.Sprintf("%v/file_%d_%s", dir, i, files[i])
	}
	for i := range videospath {
		videospath[i] = fmt.Sprintf("%v/video_%d_%s", dir, i, videos[i])
	}
	//生成上传链接
	expire_time := time.Duration(60) * time.Minute
	imagesurl := oss.GetUploadFileUrls(imagespath, util.NewArray(len(images), expire_time))
	filesurl := oss.GetUploadFileUrls(filespath, util.NewArray(len(files), expire_time))
	videosurl := oss.GetUploadFileUrls(videospath, util.NewArray(len(videos), expire_time))
	//生成下载链接
	downloadExpireTime := time.Duration(7*24) * time.Hour
	imagesDownloadurls := oss.OssGetDownloadFileUrls(imagespath, downloadExpireTime, false)
	filesDownloadUrls := oss.OssGetDownloadFileUrls(filespath, downloadExpireTime, true)
	videosDownloadUrls := oss.OssGetDownloadFileUrls(videospath, downloadExpireTime, false)
	//
	if len(imagesurl) != len(images) || len(filesurl) != len(files) || len(videosurl) != len(videos) ||
		len(imagesDownloadurls) != len(images) || len(filesDownloadUrls) != len(files) ||
		len(videosDownloadUrls) != len(videos) {
		return nil, util.NewError("上传/下载链接生成失败!")
	}
	//插入记录到数据库
	session := database.NewSession()
	defer session.Close()
	info := model.FeedbackInfo{
		SubmitTime: util.UTC_Time(),
		BaseFolder: dir,
	}
	defer session.Rollback()
	indices := dao.FeedbackInsert(session, info)
	if indices == 0 {

	}
	if err := session.Commit(); err != nil {
		return nil, util.NewError("服务器异常!", err)
	}
	//
	ret = FeedbackUploadLinkInfo{
		Indices:           indices,
		UploadImageUrls:   imagesurl,
		UploadFileUrls:    filesurl,
		UploadVideoUrls:   videosurl,
		DownloadImageUrls: imagesDownloadurls,
		DownloadFileUrls:  filesDownloadUrls,
		DownloadVideoUrls: videosDownloadUrls,
	}
	return &ret, nil
}
