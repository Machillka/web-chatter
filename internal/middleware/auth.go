package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/machillka/web-chatter/internal/config"
)

// Claims 结构体用来解析 JWT 的自定义字段
type Claims struct {
	Sub uint `json:"sub"`
	jwt.RegisteredClaims
}

// AuthRequired 验证 Authorization header 中的 Bearer token
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少 Authorization 头"})
			return
		}

		parts := strings.Fields(auth)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的 Authorization 格式"})
			return
		}

		tokenStr := parts[1]
		jwtKey := []byte(config.JWTSecret())

		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token 验证失败"})
			return
		}

		claims := token.Claims.(*Claims)
		// 将用户 ID 注入 Context，以便后续 handler 使用
		c.Set("userID", claims.Sub)
		c.Next()
	}
}
