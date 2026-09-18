package model

import "time"

type FeedbackInfo struct {
	Indices    int       `json:"indices" xorm:"pk autoincr 'indices'"`
	SubmitTime time.Time `json:"submittime" xorm:"notnull 'submittime'"`
	BaseFolder string    `json:"basefolder" xorm:"basefolder"`
}

func (c FeedbackInfo) TableName() string {
	return "Feedback"
}
