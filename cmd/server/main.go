package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/machillka/web-chatter/config"
)

func main() {
	// 从环境变量读取 DSN，格式：user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true
	config.LoadConfig("config/config.yaml")
	dsn := config.MySQLDSN()

	if dsn == "" {
		log.Fatal("请设置环境变量 MYSQL_DSN")
	}

	// 打开 MySQL 连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("数据库心跳检测失败: %v", err)
	}
	log.Println("✅ MySQL 已连接")

	// 初始化 Gin
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// TODO: 挂载 auth、ws、http 路由
	// e.g. r.POST("/auth/login", handler.Login)

	// 启动服务
	addr := ":8080"
	log.Printf("服务器启动: %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
