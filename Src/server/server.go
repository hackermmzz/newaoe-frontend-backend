package server

import (
	"fmt"
	"log"
	"net/http"
	"new-aoe-judge/Src/config"
	"path/filepath"
)

func ServerInit() {
	for {
		server()
	}
}

func server() {
	// 首页目录
	homeDir := filepath.Clean(config.Conf.ServerHome)

	// 日志目录
	logDir := filepath.Clean(config.Conf.ProcessLogDir)

	mux := http.NewServeMux()

	/*
		日志路由：

		/log/ProcessLog0.html
		映射到：

		logDir/ProcessLog0.html
	*/
	mux.Handle(
		"/log/",
		http.StripPrefix(
			"/log/",
			http.FileServer(http.Dir(logDir)),
		),
	)

	/*
		首页及其他静态文件：

		/
		映射到：

		homeDir/index.html

		例如：

		/style.css
		映射到：

		homeDir/style.css
	*/
	mux.Handle(
		"/",
		http.FileServer(http.Dir(homeDir)),
	)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.Conf.ProcessLogPort),
		Handler: mux,
	}

	fmt.Println("HTTP server started:")
	fmt.Println(fmt.Sprintf("http://0.0.0.0:%d/", config.Conf.ProcessLogPort))

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
