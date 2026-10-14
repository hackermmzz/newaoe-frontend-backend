package subrouter

import (
	"newaoe/Src/filter"
	"newaoe/Src/forum/controller"

	"github.com/gin-gonic/gin"
)

func ForumGroup_Init(router *gin.RouterGroup) {
	//论坛组
	forumGroup := router.Group("/forum")
	forumGroup.Use(filter.FilterCookieCheck())
	//
	forumGroup.POST("/announcement", filter.FilterForumAnnouncement(), controller.ForumAnnouncement)
	forumGroup.POST("/announcementACK", controller.ForumAnnouncementACK)
	forumGroup.GET("/forumFetch", controller.ForumFetch)
	forumGroup.POST("/comment", filter.FilterForumComment(), controller.ForumComment)
	forumGroup.GET("/commentFetch", controller.ForumCommentFetch)
	forumGroup.GET("/removeForumAnnouncement", controller.ForumAnnouncementRemove)
}
