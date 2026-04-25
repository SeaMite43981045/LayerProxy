// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"LayerProxy/database"
	"LayerProxy/logger"
	"LayerProxy/models"
	"LayerProxy/utils"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Key string `json:"key" binding:"required"`
}

//go:embed .dist/*
var staticFiles embed.FS

func saveConfigToFile(cfg models.ConfigFile) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("./config/config.json", data, 0644)
}

func StartAPI(ctx context.Context, cfg models.ConfigFile) {
	gin.DefaultWriter = io.MultiWriter(io.Discard)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	subFS, _ := fs.Sub(staticFiles, ".dist")
	staticHandler := http.FS(subFS)

	r.Use(func(c *gin.Context) {
		c.Next()
		logger.LogRequest(c)
	})

	a := r.Group("/api")
	{
		a.POST("/setup", func(c *gin.Context) {
			if cfg.Server.Key != "" {
				c.JSON(http.StatusFound, gin.H{"message": "已经初始化，请前往登录", "redirect": "/login"})
				return
			}

			var req AuthRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "密钥不能为空"})
				return
			}

			hashedKey, err := bcrypt.GenerateFromPassword([]byte(req.Key), bcrypt.DefaultCost)
			if err != nil {
				logger.Error("密钥加密失败: " + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "内部错误"})
				return
			}

			cfg.Server.Key = string(hashedKey)
			if err := saveConfigToFile(cfg); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存配置失败"})
				return
			}

			logger.Info("管理员密钥设置成功")
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "初始化成功"})
		})

		a.POST("/login", func(c *gin.Context) {
			if cfg.Server.Key == "" {
				c.JSON(http.StatusForbidden, gin.H{"error": "请先完成初始化", "redirect": "/setup"})
				return
			}

			var req AuthRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "请输入密钥"})
				return
			}

			err := bcrypt.CompareHashAndPassword([]byte(cfg.Server.Key), []byte(req.Key))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "密钥错误"})
				return
			}

			token, err := utils.GenerateToken("admin")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失败"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"token":   token,
				"message": "登录成功",
			})
		})
	}

	v1 := r.Group("/api/v1")
	v1.Use(utils.JWTAuthMiddleware())
	{
		v1.GET("/servers", func(c *gin.Context) {
			servers := database.GetServersFromDB()
			c.JSON(http.StatusOK, servers)
		})

		v1.POST("/servers", func(c *gin.Context) {
			var s models.ProxyInstance
			if err := c.ShouldBindJSON(&s); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
				return
			}

			if s.Name == "" || s.BackendIP == "" || s.Subdomain == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "配置项 Name/IP/域名 不能为空"})
				return
			}

			_, err := database.DB.Exec("INSERT OR REPLACE INTO servers (name, backend_ip, subdomain) VALUES (?, ?, ?)",
				s.Name, s.BackendIP, s.Subdomain)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库操作失败: " + err.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "实例已保存"})
		})

		v1.DELETE("/servers/:name", func(c *gin.Context) {
			name := c.Param("name")
			result, err := database.DB.Exec("DELETE FROM servers WHERE name = ?", name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			rowsAffected, _ := result.RowsAffected()
			c.JSON(http.StatusOK, gin.H{"status": "deleted", "rows_affected": rowsAffected})
		})
	}

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		f, err := subFS.Open(path[1:])
		if err == nil {
			f.Close()
			http.FileServer(staticHandler).ServeHTTP(c.Writer, c.Request)
			return
		}

		index, _ := fs.ReadFile(subFS, "index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Server.WebPort,
		Handler: r,
	}

	go func() {
		logger.Info(fmt.Sprintf("LayerProxy Web 服务已启动: http://localhost:%s", cfg.Server.WebPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Web 服务异常: " + err.Error())
		}
	}()

	<-ctx.Done()
	logger.Info("正在关闭 Web 管理接口...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Web 服务关闭失败: " + err.Error())
	}
}
