package dao

import (
	"context"
	"encoding/json"
	"fmt"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/user/model"
	"newaoe/Src/util"
	"time"
)

const userCacheVersionKey = "UserCacheVersion"

func userCacheKey(id string) string {
	return fmt.Sprintf("UserGet:%s", id)
}

func deleteUserCache(id string) {
	if _, err := redis.RedisDel(context.Background(), userCacheKey(id)); err != nil {
		util.DebugError("User cache del:", err)
	}
}

func getUserCacheVersion() (int, error) {
	return redis.RedisGetVersion(context.Background(), userCacheVersionKey)
}

func incrUserCacheVersion() {
	if _, err := redis.RedisIncrVersion(context.Background(), userCacheVersionKey); err != nil {
		util.DebugError("User cache incrVersion:", err)
	}
}

func userVersionedCacheKey(name string, args ...any) (string, error) {
	version, err := getUserCacheVersion()
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d:%s", name, version, util.SHA256(string(data))), nil
}

func getUserVersionedCache(key string, ret any) (bool, error) {
	data, exist, err := redis.RedisGet(context.Background(), key)
	if err != nil || !exist {
		return exist, err
	}
	if err := json.Unmarshal(data, ret); err != nil {
		return false, err
	}
	return true, nil
}

func setUserVersionedCache(key string, value any, name string) {
	data, err := json.Marshal(value)
	if err != nil {
		util.DebugError(name+" cache marshal:", err)
		return
	}
	if _, err := redis.RedisSet(context.Background(), key, data, 5*time.Minute); err != nil {
		util.DebugError(name+" cache set:", err)
	}
}

func UserResetAvatar(session *database.Session, id string, avatarFile string) (bool, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	return UserUpdate(session, id, model.Student{
		Avatar: avatarFile,
	})
}

// 判断用户是否已经注册
func UserExist(session *database.Session, id string) (bool, error) {
	user, err := UserGet(session, id)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

// 判断邮箱是否已经使用
func EmailExist(session *database.Session, email string) (ret bool, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = userVersionedCacheKey("EmailExist", email)
		if err != nil {
			return false, err
		}
		if exist, cacheErr := getUserVersionedCache(key, &ret); cacheErr != nil {
			return false, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setUserVersionedCache(key, ret, "EmailExist")
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	ret, err = session.Where("email = ?", email).Exist(&model.Student{})
	if (!ret) || (err != nil) {
		return false, err
	}
	return true, nil
}

// 添加用户
func UserAdd(session *database.Session, userInfo model.Student) (ok bool, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteUserCache(userInfo.ID)
				incrUserCacheVersion()
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data := &model.Student{
		ID:         userInfo.ID,
		Password:   userInfo.Password,
		Email:      userInfo.Email,
		Vip:        userInfo.Vip,
		RegistDate: util.UTC_Time(),
		Avatar:     util.GetRandomAvatar(),
	}
	_, err = session.Insert(data)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 添加用户
func UserAddByStudentInfo(session *database.Session, info model.Student) (ok bool, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteUserCache(info.ID)
				incrUserCacheVersion()
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err = session.Insert(info)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 更新用户数据
func UserUpdate(session *database.Session, id string, user model.Student) (ok bool, err error) {
	defer func() {
		if err == nil {
			session.PostTransactionAdd(func() {
				deleteUserCache(id)
				incrUserCacheVersion()
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err = session.Where("id = ?", id).Update(&user)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取 id like 的所有用户(不包含end)
func UserGetLikeId(session *database.Session, id_like string, beg int, end int) (ret []model.Student, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = userVersionedCacheKey("UserGetLikeId", id_like, beg, end)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getUserVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setUserVersionedCache(key, ret, "UserGetLikeId")
			})
		}
	}()
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	err = session.Where("id LIKE ?", "%"+id_like+"%").Limit(end-beg, beg).Find(&ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// 获取指定用户的数据
func UserGet(session *database.Session, id string) (ret *model.Student, err error) {
	key := userCacheKey(id)
	useCache := database.CacheAvailable(session)
	if useCache {
		data, exist, cacheErr := redis.RedisGet(context.Background(), key)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if exist {
			var user model.Student
			if err := json.Unmarshal(data, &user); err != nil {
				return nil, err
			}
			return &user, nil
		}
	}
	defer func() {
		if err != nil || ret == nil || !useCache {
			return
		}
		session.PostTransactionAdd(func() {
			data, marshalErr := json.Marshal(ret)
			if marshalErr != nil {
				util.DebugError("UserGet cache marshal:", marshalErr)
				return
			}
			if _, cacheErr := redis.RedisSet(context.Background(), key, data, 5*time.Minute); cacheErr != nil {
				util.DebugError("UserGet cache set:", cacheErr)
			}
		})
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var user model.Student
	has, err := session.Where("id = ? ", id).Get(&user)
	if err != nil || !has {
		return nil, err
	}
	return &user, nil
}

// 获取指定用户的数据
func UserGetByEmail(session *database.Session, email string) (ret *model.Student, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = userVersionedCacheKey("UserGetByEmail", email)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getUserVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setUserVersionedCache(key, ret, "UserGetByEmail")
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var user model.Student
	has, err := session.Where("email = ? ", email).Get(&user)
	if err != nil || !has {
		return nil, err
	}
	return &user, nil
}

// 获取指定用户的数据
func UserGetByIdOrEmail(session *database.Session, id_or_email string) (ret *model.Student, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = userVersionedCacheKey("UserGetByIdOrEmail", id_or_email)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getUserVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setUserVersionedCache(key, ret, "UserGetByIdOrEmail")
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var user model.Student
	has, err := session.Where("id = ? OR email = ?", id_or_email, id_or_email).Get(&user)
	if err != nil || !has {
		return nil, err
	}
	return &user, nil
}

// 根据 ID 升序排序（不包括 end）
func UserGetByRangeOrderByID(session *database.Session, beg int, end int) (ret []model.Student, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = userVersionedCacheKey("UserGetByRangeOrderByID", beg, end)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getUserVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setUserVersionedCache(key, ret, "UserGetByRangeOrderByID")
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
	// LIMIT beg, end-beg  等价于 offset beg limit count
	err = session.Asc("id").Limit(end-beg, beg).Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// 获取指定多个用户的数据
func UserGetByIDs(session *database.Session, ids []string) (ret []model.Student, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = userVersionedCacheKey("UserGetByIDs", ids)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getUserVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setUserVersionedCache(key, ret, "UserGetByIDs")
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//

	if len(ids) == 0 {
		return ret, nil
	}

	err = session.In("id", ids).Find(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// 获取所有学生数据(谨慎使用)
func UserGetAll(session *database.Session) (ret []model.Student, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = userVersionedCacheKey("UserGetAll")
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getUserVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setUserVersionedCache(key, ret, "UserGetAll")
			})
		}
	}()
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	// Find 直接把查询结果填充到切片
	err = session.Find(&ret)
	if err != nil {
		// 出错返回空
		return nil, err
	}
	return ret, nil
}

// 获取用户密码(编码后的密码)
func UserGetPassword(session *database.Session, id string) (string, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data, err := UserGet(session, id)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	return data.Password, nil
}

func UserGetEmail(session *database.Session, id string) (string, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data, err := UserGet(session, id)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	return data.Email, nil
}

func UserGetAvatar(session *database.Session, id string) (string, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data, err := UserGet(session, id)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	return data.Avatar, nil
}

func UserResetPassword(session *database.Session, id string, new_password string) (bool, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	return UserUpdate(session, id, model.Student{
		Password: new_password,
	})
}

// 获取 [beg, end)，不包含 end
func UserGetRangeIDs(session *database.Session, beg int, end int) (ret []string, err error) {
	useCache := database.CacheAvailable(session)
	key := ""
	if useCache {
		key, err = userVersionedCacheKey("UserGetRangeIDs", beg, end)
		if err != nil {
			return nil, err
		}
		if exist, cacheErr := getUserVersionedCache(key, &ret); cacheErr != nil {
			return nil, cacheErr
		} else if exist {
			return ret, nil
		}
	}
	defer func() {
		if err == nil && useCache {
			session.PostTransactionAdd(func() {
				setUserVersionedCache(key, ret, "UserGetRangeIDs")
			})
		}
	}()
	if beg < 0 || end <= beg {
		return nil, util.NewError("beg<0 or beg>=end!")
	}

	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	err = session.
		Table(model.Student{}.TableName()).
		Cols("id").
		OrderBy("id ASC").
		Limit(end-beg, beg).
		Find(&ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func UserIsSuperUser(session *database.Session, id string) (bool, error) {
	//获取用户
	user, err := UserGet(session, id)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return user.Vip == model.VIP_SUPER, nil
}
