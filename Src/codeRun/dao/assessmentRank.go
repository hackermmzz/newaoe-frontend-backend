package dao

import (
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

// end不包括
func AssessmentRankGetByRange(session *xorm.Session, beg int, end int) []model.AssessmentRankInfo {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ranks []model.AssessmentRankInfo
	if beg < 0 || end <= beg {
		return ranks
	}
	err := session.Desc("status").Asc("frame").Desc("score").Asc("id").Limit(end-beg, beg).Find(&ranks)
	if err != nil {
		util.DebugError("AssessmentRankGetByRange:", err)
		return nil
	}
	return ranks
}

func AssessmentRankUpdateOrInsert(session *xorm.Session, rank *model.AssessmentRankInfo) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	affected, err := session.Where("id = ?", rank.ID).AllCols().Update(rank)

	if err != nil {
		util.DebugError("AssessmentRankUpdateOrInsert Update:", err)
		return false
	}

	// 没有对应 ID，则插入
	if affected == 0 {
		_, err = session.Insert(rank)
		if err != nil {
			util.DebugError("AssessmentRankUpdateOrInsert Insert:", err)
			return false
		}
	}

	return true
}

func AssessmentRankBatchUpdateOrInsert(session *xorm.Session, ranks []model.AssessmentRankInfo) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if len(ranks) == 0 {
		return true
	}
	// 1. 收集所有 ID
	ids := make([]string, 0, len(ranks))
	for _, rank := range ranks {
		ids = append(ids, rank.ID)
	}
	// 2. 一次查询已经存在的 Rank
	var existRanks []model.AssessmentRankInfo
	err := session.In("id", ids).Cols("id").Find(&existRanks)
	if err != nil {
		util.DebugError("AssessmentRankBatchUpdateOrInsert Find:", err)
		return false
	}
	// 3. 建立存在 ID 的集合
	existMap := make(map[string]bool, len(existRanks))
	for _, rank := range existRanks {
		existMap[rank.ID] = true
	}
	// 4. 分成 Insert 和 Update
	insertRanks := make([]model.AssessmentRankInfo, 0)
	updateRanks := make([]model.AssessmentRankInfo, 0)
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
			util.DebugError("AssessmentRankBatchUpdateOrInsert Insert:", err)
			return false
		}
	}
	// 6. Update
	for i := range updateRanks {
		rank := &updateRanks[i]
		_, err = session.ID(rank.ID).AllCols().Update(rank)
		if err != nil {
			util.DebugError("AssessmentRankBatchUpdateOrInsert Update:", err)
			return false
		}
	}
	return true
}

func AssessmentRankBatchUpdateOrInsertIfBetter(session *xorm.Session, ranks []model.AssessmentRankInfo) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if len(ranks) == 0 {
		return true
	}
	// 1. 收集 ID
	ids := make([]string, 0, len(ranks))

	for _, rank := range ranks {
		ids = append(ids, rank.ID)
	}

	// 2. 一次性查询数据库中已有的数据
	var existRanks []model.AssessmentRankInfo

	err := session.In("id", ids).Find(&existRanks)
	if err != nil {
		util.DebugError("AssessmentRankBatchUpdateOrInsertIfBetter Find:", err)
		return false
	}

	// 3. ID -> 旧 Rank
	existMap := make(map[string]model.AssessmentRankInfo, len(existRanks))

	for _, rank := range existRanks {
		existMap[rank.ID] = rank
	}

	insertRanks := make([]model.AssessmentRankInfo, 0)
	updateRanks := make([]model.AssessmentRankInfo, 0)

	// 4. 比较
	for _, newRank := range ranks {
		oldRank, ok := existMap[newRank.ID]

		// 数据库不存在，直接插入
		if !ok {
			insertRanks = append(insertRanks, newRank)
			continue
		}

		// 数据库存在，只有新成绩更好才更新
		if AssessmentRankIsBetter(&newRank, &oldRank) {
			updateRanks = append(updateRanks, newRank)
		}
	}

	// 5. 批量插入
	if len(insertRanks) > 0 {
		_, err = session.Insert(&insertRanks)

		if err != nil {
			util.DebugError("RankBatchUpdateOrInsertIfBetter Insert:", err)
			return false
		}
	}

	// 6. 更新更好的成绩
	for i := range updateRanks {
		rank := &updateRanks[i]

		_, err = session.ID(rank.ID).AllCols().Update(rank)

		if err != nil {
			util.DebugError("RankBatchUpdateOrInsertIfBetter Update:", err)
			return false
		}
	}
	return true
}

func AssessmentRankIsBetter(newRank *model.AssessmentRankInfo, oldRank *model.AssessmentRankInfo) bool {
	// win 优先：true > false
	if newRank.Status != oldRank.Status {
		return newRank.Status > oldRank.Status
	}
	// win 相同，frame 越小越好
	if newRank.Frame != oldRank.Frame {
		return newRank.Frame < oldRank.Frame
	}
	// frame相同，score 越大越好
	if newRank.Score != oldRank.Score {
		return newRank.Score > oldRank.Score
	}
	// 完全一样，不更新
	return false
}
