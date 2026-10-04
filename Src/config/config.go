package config

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AuthAccount                     string `yaml:"AuthAccount"`
	AuthPassword                    string `yaml:"AuthPassword"`
	BaseIP                          string `yaml:"BaseIP"`
	GRPCPort                        string `yaml:"GRPCPort"`
	HTTPPort                        string `yaml:"HTTPPort"`
	NewAOEFolder                    string `yaml:"NewAOEFolder"`
	NewAOEDockerImg                 string `yaml:"NewAOEDockerImg"`
	CompileLogFileName              string `yaml:"CompileLogFileName"`
	CompileUseMultiThread           bool   `yaml:"CompileUseMultiThread"`
	RunDir                          string `yaml:"RunDir"`
	RunLogFileName                  string `yaml:"RunLogFileName"`
	RunResultFileName               string `yaml:"RunResultFileName"`
	RunDebugLogOutputFileName       string `yaml:"RunDebugLogOutputFileName"`
	RecordFileName                  string `yaml:"RecordFileName"`
	CrashLogFileName                string `yaml:"CrashLogFileName"`
	RunTimeout                      int    `yaml:"RunTimeout"`
	JudgeSleepTimeWhenGetCodeFailed int    `yaml:"JudgeSleepTimeWhenGetCodeFailed"`
	CodeRunStatusUploadInterval     int    `yaml:"CodeRunStatusUploadInterval"`
	RunMemoryLimit                  int    `yaml:"RunMemoryLimit"`
	RunDiskLimit                    int    `yaml:"RunDiskLimit"`
	ProcessLogDir                   string `yaml:"ProcessLogDir"`
	AOERunSpeed                     int    `yaml:"AOERunSpeed"`
	ServerHome                      string `yaml:"ServerHome"`
	SystemInfoFetchInterval         int    `yaml:"SystemInfoFetchInterval"`
	CurrentVersion                  string `yaml:"CurrentVersion"`
	ProcessLogPort                  int    `yaml:"ProcessLogPort"`
}

var Conf Config

func parseSimpleYAML(path string) (map[string]string, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	m := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(strings.SplitN(sc.Text(), "#", 2)[0])
		if line == "" {
			continue
		}
		p := strings.SplitN(line, ":", 2)
		if len(p) != 2 {
			continue
		}
		v := strings.Trim(strings.TrimSpace(p[1]), "\"'")
		m[strings.TrimSpace(p[0])] = v
	}
	return m, sc.Err()
}
func envOr(m map[string]string, key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return m[key]
}
func atoi(v string, def int) int {
	n, e := strconv.Atoi(v)
	if e != nil {
		return def
	}
	return n
}
func abool(v string) bool {
	b, _ := strconv.ParseBool(v)
	return b
}
func LoadConfig() error {
	m, e := parseSimpleYAML("assets/config.yaml")
	if e != nil {
		return e
	}
	Conf = Config{
		AuthAccount:     envOr(m, "AuthAccount"),
		AuthPassword:    envOr(m, "AuthPassword"),
		BaseIP:          envOr(m, "BaseIP"),
		GRPCPort:        envOr(m, "GRPCPort"),
		HTTPPort:        envOr(m, "HTTPPort"),
		NewAOEFolder:    envOr(m, "new_aoe_folder"),
		NewAOEDockerImg: envOr(m, "new_aoe_docker_img"),

		CompileLogFileName:    envOr(m, "CompileLogFileName"),
		CompileUseMultiThread: abool(envOr(m, "CompileUseMultiThread")),
		RunDir:                envOr(m, "RunDir"),
		RunLogFileName:        envOr(m, "RunLogFileName"),

		RunResultFileName:               envOr(m, "RunResultFileName"),
		RunDebugLogOutputFileName:       envOr(m, "RunDebugLogOutputFileName"),
		RecordFileName:                  envOr(m, "RecordFileName"),
		CrashLogFileName:                envOr(m, "CrashLogFileName"),
		RunMemoryLimit:                  atoi(envOr(m, "RunMemoryLimit"), 256),
		RunDiskLimit:                    atoi(envOr(m, "RunDiskLimit"), 20),
		RunTimeout:                      atoi(envOr(m, "RunTimeout"), 1800),
		JudgeSleepTimeWhenGetCodeFailed: atoi(envOr(m, "JudgeSleepTimeWhenGetCodeFailed"), 5),
		CodeRunStatusUploadInterval:     atoi(envOr(m, "CodeRunStatusUploadInterval"), 1),
		ProcessLogDir:                   envOr(m, "ProcessLogDir"),
		AOERunSpeed:                     atoi(envOr(m, "AOERunSpeed"), 8),
		ServerHome:                      envOr(m, "ServerHome"),
		SystemInfoFetchInterval:         atoi(envOr(m, "SystemInfoFetchInterval"), 5),
		CurrentVersion:                  envOr(m, "CurrentVersion"),
		ProcessLogPort:                  atoi(envOr(m, "ProcessLogPort"), 500119),
	}
	if Conf.GRPCPort == "" {
		return errors.New("GRPCPort 不能为空")
	}
	if Conf.BaseIP == "" {
		return errors.New("BaseIP 不能为空")
	}

	return e
}
