package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"strings"
	"time"
)

const codeRunCacheVersionKey = "CodeRunCacheVersion"

func codeRunCacheKey(indices int) string {
	return fmt.Sprintf("CodeRunGetByIndices:%d", indices)
}

func CodeRunCacheDelete(indices int) {
	if _, err := redis.RedisDel(context.Background(), codeRunCacheKey(indices)); err != nil {
		util.DebugError("CodeRun cache del:", err)
	}
}

func getCodeRunCacheVersion() (int, error) {
	return redis.RedisGetVersion(context.Background(), codeRunCacheVersionKey)
}

func CodeRunCacheVersionIncrease() {
	if _, err := redis.RedisIncrVersion(context.Background(), codeRunCacheVersionKey); err != nil {
		util.DebugError("CodeRun cache incrVersion:", err)
	}
}

func codeRunVersionedCacheKey(name string, args ...any) (string, error) {
	version, err := getCodeRunCacheVersion()
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d:%s", name, version, util.SHA256(string(data))), nil
}

func getCodeRunVersionedCache(key string, ret any) (bool, error) {
	data, exist, err := redis.RedisGet(context.Background(), key)
	if err != nil || !exist {
		return exist, err
	}
	if err := json.Unmarshal(data, ret); err != nil {
		return false, err
	}
	return true, nil
}

func setCodeRunVersionedCache(key string, value any, name string) {
	data, err := json.Marshal(value)
	if err != nil {
		util.DebugError(name+" cache marshal:", err)
		return
	}
	if _, err := redis.RedisSet(context.Background(), key, data, 5*time.Minute); err != nil {
		util.DebugError(name+" cache set:", err)
	}
}

func CodeRunAdd(session *database.Session, c model.CodeRunInfo) (indices int, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				CodeRunCacheDelete(c.Indices)
				CodeRunCacheVersionIncrease()
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err = session.Insert(&c)
	if err != nil {
		return 0, err
	}
	return c.Indices, nil
}

func CodeRunExist(session *database.Session, indices int) (bool, error) {
	info, err := CodeRunGetByIndices(session, indices)
	if err != nil {
		return false, err
	}
	return info != nil, nil
}

func CodeRunUpdateStatusAsync(session *database.Session, indices int, status string) (ok bool, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				CodeRunCacheDelete(indices)
				CodeRunCacheVersionIncrease()
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var info model.CodeRunInfo
	info.Status = status
	_, err = session.Where("indices=?", indices).
		Update(info)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CodeRunUpdateStatusSync(session *database.Session, indices int, status string) (ok bool, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				CodeRunCacheDelete(indices)
				CodeRunCacheVersionIncrease()
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var info model.CodeRunInfo
	info.Status = status
	_, err = session.Where("indices=?", indices).
		Update(info)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CodeRunBatchUpdateStatusIfBetter 使用临时表批量更新运行状态和版本。
// 事务由调用方控制，缓存应在事务提交成功后失效。
func CodeRunBatchUpdateStatusIfBetter(session *database.Session, data []model.CodeRunInfo) (err error) {
	if len(data) == 0 {
		return nil
	}
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//先获取旧数据
	indices := make([]int, 0)
	for _, item := range data {
		indices = append(indices, item.Indices)
	}
	var oldStatus []model.CodeRunInfo
	if err := session.In("indices", indices).ForUpdate().Find(&oldStatus); err != nil {
		return err
	}
	oldStatusMap := make(map[int]model.CodeRunInfo)
	for _, item := range oldStatus {
		oldStatusMap[item.Indices] = item
	}
	//对比数据，只更新更好的数据
	finalBetter := make([]model.CodeRunInfo, 0)
	for _, item := range data {
		_, ok := oldStatusMap[item.Indices]
		if ok && statusIsBetter(item.Status, oldStatusMap[item.Indices].Status) {
			finalBetter = append(finalBetter, item)
		}
	}
	//更新数据
	for _, item := range finalBetter {
		_, err = session.Where("indices=?", item.Indices).Cols("status").Update(item)
		if err != nil {
			return err
		}
	}
	// 删除缓存
	session.PostTransactionAdd(func() {
		for _, item := range finalBetter {
			CodeRunCacheDelete(item.Indices)
		}
		CodeRunCacheVersionIncrease()
	})
	return nil
}

func CodeRunGetByIndices(session *database.Session, indices int) (ret *model.CodeRunInfo, err error) {
	key := codeRunCacheKey(indices)
	useCache := database.CacheAvailable(session)
	if useCache {
		data, exist, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if exist {
			var info model.CodeRunInfo
			if err := json.Unmarshal(data, &info); err != nil {
				return nil, err
			}
			return &info, nil
		}
	}
	defer func() {
		if err != nil || ret == nil || !useCache {
			return
		}
		session.PostTransactionAdd(func() {
			data, marshalErr := json.Marshal(ret)
			if marshalErr != nil {
				util.DebugError("CodeRunGetByIndices cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("CodeRunGetByIndices cache set:", cacheErr)
			}
		})
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var info model.CodeRunInfo
	has, err := session.Where("indices=?", indices).Get(&info)
	if err != nil || !has {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &info, nil
}

func CodeRunBatchGetByIndices(session *database.Session, indices []int) (ret []model.CodeRunInfo, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = codeRunVersionedCacheKey("CodeRunBatchGetByIndices", indices)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getCodeRunVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setCodeRunVersionedCache(key, ret, "CodeRunBatchGetByIndices")
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if len(indices) == 0 {
		return ret, nil
	}
	err = session.In("indices", indices).
		Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func CodeRunGetById(session *database.Session, id string) (ret []model.CodeRunInfo, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = codeRunVersionedCacheKey("CodeRunGetById", id)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getCodeRunVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setCodeRunVersionedCache(key, ret, "CodeRunGetById")
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	err = session.Where("id = ?", id).Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func CodeRunCountById(session *database.Session, id string) (ret int64, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = codeRunVersionedCacheKey("CodeRunCountById", id)
		if err != nil {
			return 0, err
		}
		if exist, cacheErr := getCodeRunVersionedCache(key, &ret); cacheErr != nil {
			return 0, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setCodeRunVersionedCache(key, ret, "CodeRunCountById")
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	ret, err = session.Where("id = ?", id).Count(&model.CodeRunInfo{})
	if err != nil {
		return 0, err
	}
	return ret, nil
}

// 不包含end
func CodeRunGetRangeById(session *database.Session, id string, beg int, end int, class ...int) (ret []model.CodeRunInfo, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = codeRunVersionedCacheKey("CodeRunGetRangeById", id, beg, end, class)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getCodeRunVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setCodeRunVersionedCache(key, ret, "CodeRunGetRangeById")
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if beg < 0 || beg >= end {
		return nil, util.NewError("beg<0 or beg>=end!")
	}
	session = getExtraClassInCondition(session, class...)
	err = session.Where("id = ?", id).Desc("submittime").
		Limit(end-beg, int(beg)).
		Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// 获取指定范围的提交记录(不包含end)
func CodeRunGetByRange(session *database.Session, beg int, end int, class ...int) (ret []model.CodeRunInfo, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = codeRunVersionedCacheKey("CodeRunGetByRange", beg, end, class)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getCodeRunVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setCodeRunVersionedCache(key, ret, "CodeRunGetByRange")
			})
		}
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	session = getExtraClassInCondition(session, class...)
	err = session.Desc("indices").Limit(end-beg, beg).Find(&ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// 生成IN class的额外条件
func getExtraClassInCondition(session *database.Session, class ...int) *database.Session {
	if len(class) == 0 {
		return session
	}
	if len(class) == 1 {
		session.Where("class = ?", class[0])
		return session
	}
	placeholders := make([]string, len(class))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	cond := "class IN (" + strings.Join(placeholders, ",") + ")"
	args := make([]interface{}, len(class))
	for i, v := range class {
		args[i] = v
	}
	session.And(cond, args...)
	return session
}

// 比对数据，只更新更好的数据
func statusIsBetter(newStatus string, oldStatus string) bool {
	//序列化
	newStatusInfo := model.CodeRunStatusInfo{}
	err := newStatusInfo.Unmarshal([]byte(newStatus))
	if err != nil {
		return false
	}
	oldStatusInfo := model.CodeRunStatusInfo{}
	err = oldStatusInfo.Unmarshal([]byte(oldStatus))
	if err != nil {
		return true
	}
	return newStatusInfo.IsBetter(oldStatusInfo)
}
