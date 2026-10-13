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

const forumAnnouncementCacheVersionKey = "ForumAnnouncementCacheVersion"

func getForumAnnouncementCacheVersion() (int, error) {
	return redis.RedisGetVersion(context.Background(), forumAnnouncementCacheVersionKey)
}

func incrForumAnnouncementCacheVersion() {
	if _, err := redis.RedisIncrVersion(context.Background(), forumAnnouncementCacheVersionKey); err != nil {
		util.DebugError("ForumAnnouncement cache incrVersion:", err)
	}
}

func ForumAnnouncementInsert(session *database.Session, announcement *model.ForumAnnouncement) (indices int64, err error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	defer func() {
		if err != nil {
			return
		}
		session.PostTransactionAdd(incrForumAnnouncementCacheVersion)
	}()

	_, err = session.Insert(announcement)
	if err != nil {
		return 0, err
	}
	return int64(announcement.Indices), nil
}

// ForumAnnouncementGetByRange 按索引倒序查询公告，不包含 end。
func ForumAnnouncementGetByRange(session *database.Session, beg int, end int) (ret []model.ForumAnnouncement, err error) {
	if beg < 0 || beg >= end {
		return nil, util.NewError("beg<0 or beg>=end!")
	}

	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		version, cacheErr := getForumAnnouncementCacheVersion()
		if cacheErr != nil {
			return nil, cacheErr
		}
		key = fmt.Sprintf("ForumAnnouncementGetByRange:%d:%d:%d", version, beg, end)
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
				util.DebugError("ForumAnnouncementGetByRange cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("ForumAnnouncementGetByRange cache set:", cacheErr)
			}
		})
	}()

	err = session.Desc("indices").Limit(end-beg, beg).Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
