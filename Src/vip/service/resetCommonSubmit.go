package service

import (
	"context"
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/redis"
	"newaoe/Src/util"
)

// 如果ids为nil，表示重置所有人的次数
func ResetCommonSubmit(ids []string) error {
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
	return nil
}
