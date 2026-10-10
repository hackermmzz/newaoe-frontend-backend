package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"

	"xorm.io/xorm"
)

var AssessmentRankVersionLua = `
local version = redis.call("GET", "AssessmentRankVersion")
if version then
    return tonumber(version)
end
redis.call("SET", "AssessmentRankVersion", 0)

return 0`

var AssessmentRankVersionIncrLua = `
return redis.call("INCR", "AssessmentRankVersion")
`

func getVersion() (int, error) {
	resCMD := redis.RedisLua(context.Background(), AssessmentRankVersionLua, []string{})
	err := resCMD.Err()
	if err != nil {
		return -1, err
	}
	result, err := resCMD.Int()
	if err != nil {
		return -1, err
	}
	return result, nil
}

func incrVersion() error {
	resCMD := redis.RedisLua(context.Background(), AssessmentRankVersionLua, []string{})
	if err := resCMD.Err(); err != nil {
		return err
	}
	return nil
}

// end不包括
func AssessmentRankGetByRange(session *xorm.Session, beg int, end int) (ret []model.AssessmentRankInfo, err error) {
	//
	//先查询缓存
	version, err := getVersion()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("AssessmentRankGetByRange_%d/%d/%d", version, beg, end)
	dataBytes, exist, err := redis.RedisGet(context.Background(), key)
	if err != nil {
		return nil, err
	}
	if exist {
		var ranks []model.AssessmentRankInfo
		err := json.Unmarshal(dataBytes, &ranks)
		if err == nil {
			return ranks, nil
		}
		return nil, err
	}
	//写入缓存
	defer func() {
		data, _ := json.Marshal(ret)
		if _, err := redis.RedisSet(context.Background(), key, data, time.Duration(5)*time.Minute); err != nil {
			util.DebugError("AssessmentRank cache set:", err)
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var ranks []model.AssessmentRankInfo
	//
	if beg < 0 || end <= beg {
		return ranks, nil
	}
	//先默认查询都是胜利的
	err = session.Desc("status").Asc("frame").Desc("score").Asc("id").Limit(end-beg, beg).Find(&ranks)
	if err != nil {
		return nil, err
	}
	if len(ranks) == 0 {
		return nil, nil
	}
	//判断里面是否存在不是胜利的情况
	l, r := 0, len(ranks)-1
	for l <= r {
		mid := (l + r) / 2
		status := ranks[mid].Status
		if status >= model.AssessmenrRankStatus_Win {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	if l >= len(ranks) {
		return ranks, nil //结果都是胜利的
	}
	//如果查询里面存在不是胜利的，则需考虑其他的查询
	var ranks1 []model.AssessmentRankInfo
	err = session.Desc("status").Desc("score").Asc("frame").Asc("id").Limit(end-beg-l, beg+l).Find(&ranks1)
	if err != nil {
		return nil, err
	}
	//合并
	ranks = append(ranks[:l], ranks1...)
	//
	return ranks, nil
}

func AssessmentRankBatchUpdateOrInsertIfBetter(session *xorm.Session, ranks []model.AssessmentRankInfo) (err error) {
	defer func() {
		//提高版本
		if err == nil {
			if err := incrVersion(); err != nil {
				util.DebugError("AssessmentRank incrVersion:", err)
			}
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if len(ranks) == 0 {
		return nil
	}
	// 1. 收集 ID
	ids := make([]string, 0, len(ranks))

	for _, rank := range ranks {
		ids = append(ids, rank.ID)
	}

	// 2. 一次性查询数据库中已有的数据
	var existRanks []model.AssessmentRankInfo

	err = session.In("id", ids).Find(&existRanks)
	if err != nil {
		return err
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
			return err
		}
	}

	// 6. 更新更好的成绩
	for i := range updateRanks {
		rank := &updateRanks[i]

		_, err = session.ID(rank.ID).AllCols().Update(rank)

		if err != nil {
			return err
		}
	}
	return nil
}

func AssessmentRankIsBetter(newRank *model.AssessmentRankInfo, oldRank *model.AssessmentRankInfo) bool {
	// win 优先：true > false
	if newRank.Status != oldRank.Status {
		return newRank.Status > oldRank.Status
	}
	if newRank.Status == model.AssessmenrRankStatus_Win {
		// frame 越小越好
		if newRank.Frame != oldRank.Frame {
			return newRank.Frame < oldRank.Frame
		}
		return newRank.Score > oldRank.Score
	} else {
		// score 越大越好
		if newRank.Score != oldRank.Score {
			return newRank.Score > oldRank.Score
		}
		return newRank.Frame < oldRank.Frame
	}
}
