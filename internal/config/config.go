package config

import (
	"log"

	"github.com/spf13/viper"
)

// 以项目根节点作为路径的起点
func LoadConfig(path string) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}
}

func MySQLDSN() string {
	return viper.GetString("mysql.dsn")
}

func ServerPort() int {
	return viper.GetInt("server.port")
}

func JWTSecret() string {
	return viper.GetString("jwt.secret")
}
