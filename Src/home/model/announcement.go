package model

import "time"

type AnnouncementInfo struct {
	Indices    int64     `json:"indices" xorm:"'indices' pk autoincr BIGINT"`
	Title      string    `json:"title" xorm:"'title' VARCHAR(255) notnull"`
	Content    string    `json:"content" xorm:"'content' TEXT notnull"`
	Version    int       `json:"version" xorm:"'version' INT notnull default 1"`
	Enabled    bool      `json:"enabled" xorm:"'enabled' BOOLEAN notnull default true"`
	CreateTime time.Time `json:"createTime" xorm:"'createTime' DATETIME created"`
	UpdateTime time.Time `json:"updateTime" xorm:"'updateTime' DATETIME updated"`
}

func (c AnnouncementInfo) TableName() string {
	return "Announcement"
}
