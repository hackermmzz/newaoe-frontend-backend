package dao

import (
	"newaoe/util"
	"time"

	"xorm.io/xorm"
)

type RankInfo struct {
	ID         string    `json:"id" xorm:"id pk"`              //用户学号
	Win        bool      `json:"win" xorm:"win"`               //胜利或者失败
	SubmitTime time.Time `json:"submittime" xorm:"submittime"` //提交日期
	Score      int       `json:"score" xorm:"score"`           //分数
	Frame      int       `json:"frame" xorm:"frame"`           //运行时间
	Msg        string    `json:"msg" xorm:"msg"`               //运行结果
}

func (r RankInfo) TableName() string {
	return "Rank"
}

func RankGetByRange(session *xorm.Session, beg int, end int) []RankInfo {
	var ranks []RankInfo
	if beg < 0 || end <= beg {
		return ranks
	}
	err := session.Desc("win").Desc("score").Asc("frame").Asc("id").Limit(end-beg, beg).Find(&ranks)
	if err != nil {
		util.Debug("RankGetByRange:", err)
		return nil
	}
	return ranks
}

func RankUpdateOrInsert(session *xorm.Session, rank *RankInfo) bool {
	affected, err := session.Where("id = ?", rank.ID).AllCols().Update(rank)

	if err != nil {
		util.Debug("RankUpdateOrInsert Update:", err)
		return false
	}

	// 没有对应 ID，则插入
	if affected == 0 {
		_, err = session.Insert(rank)
		if err != nil {
			util.Debug("RankUpdateOrInsert Insert:", err)
			return false
		}
	}

	return true
}

func RankBatchUpdateOrInsert(session *xorm.Session, ranks []RankInfo) bool {
	if len(ranks) == 0 {
		return true
	}
	// 1. 收集所有 ID
	ids := make([]string, 0, len(ranks))
	for _, rank := range ranks {
		ids = append(ids, rank.ID)
	}
	// 2. 一次查询已经存在的 Rank
	var existRanks []RankInfo
	err := session.In("id", ids).Cols("id").Find(&existRanks)
	if err != nil {
		util.Debug("RankBatchUpdateOrInsert Find:", err)
		return false
	}
	// 3. 建立存在 ID 的集合
	existMap := make(map[string]bool, len(existRanks))
	for _, rank := range existRanks {
		existMap[rank.ID] = true
	}
	// 4. 分成 Insert 和 Update
	insertRanks := make([]RankInfo, 0)
	updateRanks := make([]RankInfo, 0)
	for _, rank := range ranks {
		if existMap[rank.ID] {
			updateRanks = append(updateRanks, rank)
		} else {
			insertRanks = append(insertRanks, rank)
		}
	}
	// 5. 批量 Insert
	if len(insertRanks) > 0 {
		_, err = session.Insert(&insertRanks)
		if err != nil {
			util.Debug("RankBatchUpdateOrInsert Insert:", err)
			return false
		}
	}
	// 6. Update
	for i := range updateRanks {
		rank := &updateRanks[i]
		_, err = session.ID(rank.ID).AllCols().Update(rank)
		if err != nil {
			util.Debug("RankBatchUpdateOrInsert Update:", err)
			return false
		}
	}
	return true
}

func RankGetCount(session *xorm.Session) int {
	count, err := session.Count(&RankInfo{})
	if err != nil {
		util.Debug("RankGetCount:", err)
		return 0
	}
	return int(count)
}
