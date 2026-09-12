package Server

import (
	"time"

	"github.com/gin-contrib/cors"
)

func cors_Config() {
	Engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://114.66.62.156:5119"},               // 允许所有源（生产环境建议指定具体域名）
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},                  // 允许的方法，必须包含OPTIONS
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 预检请求的缓存时间，减少OPTIONS请求次数
	}))
}
