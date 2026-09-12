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

func RankBatchUpdateOrInsertIfBetter(session *xorm.Session, ranks []RankInfo) bool {
	if len(ranks) == 0 {
		return true
	}
	// 1. 收集 ID
	ids := make([]string, 0, len(ranks))

	for _, rank := range ranks {
		ids = append(ids, rank.ID)
	}

	// 2. 一次性查询数据库中已有的数据
	var existRanks []RankInfo

	err := session.In("id", ids).Find(&existRanks)
	if err != nil {
		util.Debug("RankBatchUpdateOrInsertIfBetter Find:", err)
		return false
	}

	// 3. ID -> 旧 Rank
	existMap := make(map[string]RankInfo, len(existRanks))

	for _, rank := range existRanks {
		existMap[rank.ID] = rank
	}

	insertRanks := make([]RankInfo, 0)
	updateRanks := make([]RankInfo, 0)

	// 4. 比较
	for _, newRank := range ranks {
		oldRank, ok := existMap[newRank.ID]

		// 数据库不存在，直接插入
		if !ok {
			insertRanks = append(insertRanks, newRank)
			continue
		}

		// 数据库存在，只有新成绩更好才更新
		if RankIsBetter(&newRank, &oldRank) {
			updateRanks = append(updateRanks, newRank)
		}
	}

	// 5. 批量插入
	if len(insertRanks) > 0 {
		_, err = session.Insert(&insertRanks)

		if err != nil {
			util.Debug("RankBatchUpdateOrInsertIfBetter Insert:", err)
			return false
		}
	}

	// 6. 更新更好的成绩
	for i := range updateRanks {
		rank := &updateRanks[i]

		_, err = session.ID(rank.ID).AllCols().Update(rank)

		if err != nil {
			util.Debug("RankBatchUpdateOrInsertIfBetter Update:", err)
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

func RankIsBetter(newRank *RankInfo, oldRank *RankInfo) bool {
	// win 优先：true > false
	if newRank.Win != oldRank.Win {
		return newRank.Win
	}
	// win 相同，score 越大越好
	if newRank.Score != oldRank.Score {
		return newRank.Score > oldRank.Score
	}
	// score 相同，frame 越小越好
	if newRank.Frame != oldRank.Frame {
		return newRank.Frame < oldRank.Frame
	}
	// 完全一样，不更新
	return false
}
