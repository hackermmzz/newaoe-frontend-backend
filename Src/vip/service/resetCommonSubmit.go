package service

import (
	"context"
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/redis"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
)

// 如果ids为nil，表示重置所有人的次数
func ResetCommonSubmit(ids []string) error {
	resetAll := len(ids) == 0
	curIdx := 0
	for {
		//如果需要重置所有人的次数
		if resetAll {
			//获取id
			res, err := dao.UserGetByRangeOrderByID(nil, curIdx, curIdx+config.Conf.Other.ResetCommonSubmitBatch)
			curIdx += config.Conf.Other.ResetCommonSubmitBatch
			if err != nil {
				return err
			}
			if len(res) == 0 {
				break // 无更多用户
			}
			ids = make([]string, 0)
			for _, info := range res {
				ids = append(ids, info.ID)
			}
		}
		//
		for i := 0; i < len(ids); i += config.Conf.Other.ResetCommonSubmitBatch {
			ctx := context.Background()
			beg := i
			end := min(beg+config.Conf.Other.ResetCommonSubmitBatch, len(ids))
			pipe := redis.NewTxPipeline()
			for ; beg < end; beg += 1 {
				key := fmt.Sprintf("CommonUploadOrRunTimes_%s", ids[beg])
				pipe.Set(ctx, key, 0, util.GetLeftTimeForOneDay())
			}
			_, err := pipe.Exec(ctx)
			if err != nil {
				return err
			}
		}
		//
		if !resetAll {
			break
		}
	}
	return nil
}
