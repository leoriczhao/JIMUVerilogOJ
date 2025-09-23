package middleware

import (
	"net/http"
	"strings"

	"verilog-oj/backend/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 从配置读取JWT密钥
var jwtSecret = []byte(config.LoadConfig().JWT.Secret)

// AuthRequired JWT认证中间件
func AuthRequired() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "未提供认证Token",
			})
			c.Abort()
			return
		}

		// 检查Bearer格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "认证Token格式错误",
			})
			c.Abort()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Token无效或已过期",
			})
			c.Abort()
			return
		}

		// 提取claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// 将用户信息设置到上下文
			if userID, exists := claims["user_id"]; exists {
				c.Set("user_id", uint(userID.(float64)))
			}
			if username, exists := claims["username"]; exists {
				c.Set("username", username.(string))
			}
			if role, exists := claims["role"]; exists {
				c.Set("role", role.(string))
			}
		}

		c.Next()
	})
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return gin.Logger()
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
