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
)

func getRankVersion() (int, error) {
	return redis.RedisGetVersion(context.Background(), "RankVersion")
}

func incrRankVersion() {
	if _, err := redis.RedisIncrease(context.Background(), "RankVersion"); err != nil {
		util.DebugError("Rank incrVersion:", err)
	}
}

// end不包括
func RankGetByRange(session *database.Session, beg int, end int) (ret []model.RankInfo, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		version, cacheErr := getRankVersion()
		if cacheErr != nil {
			return nil, cacheErr
		}
		key = fmt.Sprintf("RankGetByRange:%d:%d:%d", version, beg, end)
		data, exist, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if exist {
			if err := json.Unmarshal(data, &ret); err != nil {
				return nil, err
			}
			return ret, nil
		}
	}
	defer func() {
		if err != nil || !useCache {
			return
		}
		session.PostTransactionAdd(func() {
			data, marshalErr := json.Marshal(ret)
			if marshalErr != nil {
				util.DebugError("RankGetByRange cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("RankGetByRange cache set:", cacheErr)
			}
		})
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if beg < 0 || end <= beg {
		return ret, nil
	}
	err = session.Desc("win").Asc("frame").Desc("score").Asc("id").Limit(end-beg, beg).Find(&ret)
	if err != nil {
		return nil, err
	}
	//
	return ret, nil
}

func RankUpdateOrInsert(session *database.Session, rank *model.RankInfo) (ok bool, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(incrRankVersion)
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	affected, err := session.Where("id = ?", rank.ID).AllCols().Update(rank)

	if err != nil {
		return false, err
	}

	// 没有对应 ID，则插入
	if affected == 0 {
		_, err = session.Insert(rank)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func RankBatchUpdateOrInsertIfBetter(session *database.Session, ranks []model.RankInfo) (ok bool, err error) {
	defer func() {
		if err == nil && len(ranks) > 0 {
			session.PostTransactionAdd(incrRankVersion)
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if len(ranks) == 0 {
		return true, nil
	}
	// 1. 收集 ID
	ids := make([]string, 0, len(ranks))

	for _, rank := range ranks {
		ids = append(ids, rank.ID)
	}

	// 2. 一次性查询数据库中已有的数据
	var existRanks []model.RankInfo

	err = session.In("id", ids).ForUpdate().Find(&existRanks)
	if err != nil {
		return false, err
	}

	// 3. ID -> 旧 Rank
	existMap := make(map[string]model.RankInfo, len(existRanks))

	for _, rank := range existRanks {
		existMap[rank.ID] = rank
	}

	insertRanks := make([]model.RankInfo, 0)
	updateRanks := make([]model.RankInfo, 0)

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
			return false, err
		}
	}

	// 6. 更新更好的成绩
	for i := range updateRanks {
		rank := &updateRanks[i]

		_, err = session.ID(rank.ID).AllCols().Update(rank)

		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func RankIsBetter(newRank *model.RankInfo, oldRank *model.RankInfo) bool {
	//反序列化
	var newinfo, oldinfo model.RankMsgInfo
	err := json.Unmarshal([]byte(newRank.Msg), &newinfo)
	if err != nil {
		util.DebugError("RankIsBetter newRank.Msg unmarshal:", err)
		return false
	}
	err = json.Unmarshal([]byte(oldRank.Msg), &oldinfo)
	if err != nil {
		util.DebugError("RankIsBetter oldRank.Msg unmarshal:", err)
		return true
	}
	return newinfo.Status.IsBetter(oldinfo.Status)
}
