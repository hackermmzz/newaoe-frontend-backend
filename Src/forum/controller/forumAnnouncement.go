package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"newaoe/Src/common/upload"
	"newaoe/Src/config"
	"newaoe/Src/oss"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"path/filepath"
	"time"

	"newaoe/Src/forum/dao"
	"newaoe/Src/forum/model"
	"newaoe/Src/forum/service"

	"github.com/gin-gonic/gin"
)

type forumAnnouncement struct {
	Content string   `json:"content"`
	Title   string   `json:"title"`
	Files   []string `json:"files"`
	Videos  []string `json:"videos"`
	Images  []string `json:"images"`
}

type forumAnnouncementPost struct {
	Key       string                    `json:"key"`
	FileUrls  []*oss.UploadPolicyResult `json:"fileUrls"`
	VideoUrls []*oss.UploadPolicyResult `json:"videoUrls"`
	ImageUrls []*oss.UploadPolicyResult `json:"imageUrls"`
}

func ForumAnnouncement(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//解析请求体
	var announcement forumAnnouncement
	if ok := util.JsonCtx(ctx, &announcement); !ok {
		util.ResponseNAK_MSG(ctx, "数据格式错误!", nil)
		return
	}
	//限制标题长度
	if len(announcement.Title) > 100 {
		util.ResponseNAK_MSG(ctx, "标题长度不能超过100个字符!", nil)
		return
	}
	//限制内容长度
	if len(announcement.Content) > 10000 {
		util.ResponseNAK_MSG(ctx, "内容长度不能超过10000个字符!", nil)
		return
	}
	//限制文件数量
	if len(announcement.Files) > 10 {
		util.ResponseNAK_MSG(ctx, "文件数量不能超过10个!", nil)
		return
	}
	//限制视频数量
	if len(announcement.Videos) > 3 {
		util.ResponseNAK_MSG(ctx, "视频数量不能超过10个!", nil)
		return
	}
	//限制图片数量
	if len(announcement.Images) > 10 {
		util.ResponseNAK_MSG(ctx, "图片数量不能超过10个!", nil)
		return
	}
	//生成存储路径
	baseDir := filepath.Join(config.Conf.OSS.PublicBaseFolder, config.Conf.Other.PublicOtherFolder, "forum", id, fmt.Sprintf("%v_%v", util.UTC_Time().String(), util.UUID()))
	contentPath := filepath.Join(baseDir, "content.html")
	filePaths := make([]string, len(announcement.Files))
	videoPaths := make([]string, len(announcement.Videos))
	imagePaths := make([]string, len(announcement.Images))
	fileMap := make(map[string]string)
	imageMap := make(map[string]string)
	videoMap := make(map[string]string)
	for i, v := range announcement.Files {
		filePaths[i] = filepath.Join(baseDir, fmt.Sprintf("file_%d_%v_%s", i, util.UUID(), announcement.Files[i]))
		fileMap[v] = filePaths[i]
	}
	for i, v := range announcement.Videos {
		videoPaths[i] = filepath.Join(baseDir, fmt.Sprintf("video_%d_%v.mp4", i, util.UUID()))
		videoMap[v] = videoPaths[i]
	}
	for i, v := range announcement.Images {
		imagePaths[i] = filepath.Join(baseDir, fmt.Sprintf("image_%d_%v.jpg", i, util.UUID()))
		imageMap[v] = imagePaths[i]
	}
	//过滤前端上传过来的文件并且后端直接上传
	err := service.SafeUploadContent(contentPath, announcement.Content, fileMap, imageMap, videoMap)
	if err != nil {
		util.ResponseNAK_MSG(ctx, "内容上传失败!", nil)
		return
	}
	//生成文件urls
	fileUrls := upload.UploadPostPolicyFiles(announcement.Files, util.NewArray(len(announcement.Files), time.Minute*30), 0, config.Conf.Server.MaxFilsSizeUserUpload)
	//生成视频url
	videoUrls := upload.UploadPostPolicyFiles(announcement.Videos, util.NewArray(len(announcement.Videos), time.Minute*60), 0, 100*1024*1024) //视频最大100MB
	//生成图片url
	imageUrls := upload.UploadPostPolicyFiles(announcement.Images, util.NewArray(len(announcement.Images), time.Minute*30), 0, config.Conf.Server.MaxFilsSizeUserUpload)
	//检查是否生成成功
	if len(fileUrls) != len(announcement.Files) || len(videoUrls) != len(announcement.Videos) || len(imageUrls) != len(announcement.Images) {
		util.ResponseNAK_MSG(ctx, "生成上传链接失败!", nil)
		return
	}
	//生成唯一key
	key := fmt.Sprintf("ForumAnnouncementController:%s", id)
	//上传缓存
	ok, err := redis.RedisSetNX(ctx, key, map[string]interface{}{
		"title":       announcement.Title,
		"contentPath": contentPath,
		"id":          id,
	}, time.Minute*30)
	if !ok || err != nil {
		util.ResponseNAK_MSG(ctx, "缓存上传失败!", nil)
		return
	}
	//返回成功
	util.ResponseACK_MSG(ctx, "成功!", forumAnnouncementPost{
		Key:       key,
		FileUrls:  fileUrls,
		VideoUrls: videoUrls,
		ImageUrls: imageUrls,
	})
}

type forumAnnouncementACK struct {
	Key string `json:"key"`
}

func ForumAnnouncementACK(ctx *gin.Context) {
	var ack forumAnnouncementACK
	if ok := util.JsonCtx(ctx, &ack); !ok {
		util.ResponseNAK_MSG(ctx, "数据格式错误!", nil)
		return
	}
	if ack.Key == "" {
		util.ResponseNAK_MSG(ctx, "key不能为空!", nil)
		return
	}
	//根据key获取内容
	bytes, ok, err := redis.RedisGet(ctx, ack.Key)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	if bytes == nil || !ok {
		util.ResponseNAK_MSG(ctx, "Key不存在!", nil)
		return
	}
	//反序列化数据
	dataCache := make(map[string]interface{})
	err = json.Unmarshal(bytes, &dataCache)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//获取数据
	title, ok0 := dataCache["title"]
	id, ok1 := dataCache["id"]
	contentPath, ok2 := dataCache["contentPath"]
	if !ok0 || !ok1 || !ok2 {
		util.ResponseNAK_MSG(ctx, "数据格式错误!", nil)
		return
	}
	//写入数据库
	indices, err := dao.ForumAnnouncementInsert(nil, &model.ForumAnnouncement{
		ID:          id.(string),
		Title:       title.(string),
		ContentPath: contentPath.(string),
		SubmitTime:  util.UTC_Time(),
	})
	if err != nil || indices == 0 {
		if err == nil {
			err = errors.New("写入数据库失败!")
		}
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//返回成功
	util.ResponseACK_MSG(ctx, "成功!", nil)
}
