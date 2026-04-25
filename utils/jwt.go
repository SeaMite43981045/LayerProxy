// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("HUGS9DasdG3DfG4U690DFGHseF893sdfg")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 有效期 24 小时
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权，请先登录"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 格式错误"})
			c.Abort()
			return
		}

		// 解析 Token
		tokenString := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或已过期的 Token"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文，方便后续路由使用
		c.Set("username", claims.Username)
		c.Next()
	}
}
