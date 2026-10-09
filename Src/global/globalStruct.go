package global

import (
	"encoding/json"
)

type CodeRunType int64

// ///////////////////////// OJSystemInfo OJ系统信息
type OJSystemInfo struct {
	Version string `json:"version"` //最新的版本号
}

// //////////////////////// GRPCAuthInfo gRPC 认证信息
type GRPCAuthInfo struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

func (a GRPCAuthInfo) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}

// //////////////////////// CodeRunStatusInfo 代码运行状态信息
type CodeRunStatusInfo struct {
	Status   int32  `json:"status"`
	Food     int    `json:"food"`
	Wood     int    `json:"wood"`
	Gold     int    `json:"gold"`
	Stone    int    `json:"stone"`
	Frame    int    `json:"frame"`
	Win      bool   `json:"win"`
	Score    int    `json:"score"`
	Data     string `json:"data"`
	TimeCost int64  `json:"time_cost"` //秒
}

func (s CodeRunStatusInfo) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}
