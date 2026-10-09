package service

import (
	"context"
	"newaoe/Src/redis"
)

func SetBlockAllSubmit(block bool) error {
	value := 0
	if block {
		value = 1
	}
	lua := `
	for i, val in pairs(KEYS) do
		redis.call("SET",val,ARGV[1],"KEEPTTL")
	end
	return 1
	`
	resCMD := redis.RedisLua(context.Background(), lua, []string{
		"BlockCodeRun",
		"BlockAssessmentSubmit",
		"BlockCommonSubmit",
	}, value)
	return resCMD.Err()
}
