package dao

import (
	"context"
	"encoding/json"
	"fmt"
	database "newaoe/Src/databse"
	"newaoe/Src/forum/model"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"
)

const forumCommentCacheVersionKey = "ForumCommentCacheVersion"

func getForumCommentCacheVersion() (int, error) {
	return redis.RedisGetVersion(context.Background(), forumCommentCacheVersionKey)
}

func incrForumCommentCacheVersion() {
	if _, err := redis.RedisIncrVersion(context.Background(), forumCommentCacheVersionKey); err != nil {
		util.DebugError("ForumComment cache incrVersion:", err)
	}
}

func ForumCommentInsert(session *database.Session, comment *model.ForumComment) (indices int64, err error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	defer func() {
		if err != nil {
			return
		}
		session.PostTransactionAdd(incrForumCommentCacheVersion)
	}()

	_, err = session.Insert(comment)
	if err != nil {
		return 0, err
	}
	return comment.Indices, nil
}

// ForumCommentGetByRange 按索引倒序查询指定公告的评论，不包含 end。
func ForumCommentGetByRange(session *database.Session, indices int64, beg int, end int) (ret []model.ForumComment, err error) {
	if indices <= 0 {
		return nil, util.NewError("indices<=0!")
	}
	if beg < 0 || beg >= end {
		return nil, util.NewError("beg<0 or beg>=end!")
	}

	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		version, cacheErr := getForumCommentCacheVersion()
		if cacheErr != nil {
			return nil, cacheErr
		}
		key = fmt.Sprintf("ForumCommentGetByRange:%d:%d:%d:%d", version, indices, beg, end)
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

	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	defer func() {
		if err != nil || !useCache {
			return
		}
		session.PostTransactionAdd(func() {
			data, marshalErr := json.Marshal(ret)
			if marshalErr != nil {
				util.DebugError("ForumCommentGetByRange cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("ForumCommentGetByRange cache set:", cacheErr)
			}
		})
	}()

	err = session.Where("indicesForForumAnnoucement = ?", indices).
		Desc("indices").
		Limit(end-beg, beg).
		Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
