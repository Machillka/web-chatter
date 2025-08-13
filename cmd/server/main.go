package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/machillka/web-chatter/internal/config"
	"github.com/machillka/web-chatter/internal/db"
	"github.com/machillka/web-chatter/internal/handler"
	"github.com/machillka/web-chatter/internal/middleware"
)

func main() {
	config.LoadConfig("config/config.yaml")

	// 打开 MySQL 连接
	err := db.Init()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("✅ MySQL 已连接")

	// 初始化 Gin
	// r := gin.New()
	// r.Use(gin.Logger(), gin.Recovery())

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// TODO: 挂载 auth、ws、http 路由
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	// 受保护的路由
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/profile", handler.Profile)
		auth.PUT("/profile", handler.UpdateProfile)
	}

	wsGroup := r.Group("/ws")
	// wsGroup.Use(middleware.AuthRequired())
	wsGroup.GET("/chat", handler.ChatHandler)

	// 启动服务
	addr := fmt.Sprintf(":%d", config.ServerPort())
	log.Printf("服务器启动: %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
