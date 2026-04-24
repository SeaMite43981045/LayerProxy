// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package proxy

import (
	"LayerProxy/logger"
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	dbPath := "./data/data.db"
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.MkdirAll(dbDir, 0755)
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Error("无法打开数据库: " + err.Error())
		return
	}

	if err = DB.Ping(); err != nil {
		logger.Error("数据库 Ping 失败: " + err.Error())
		return
	}

	query := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		backend_ip TEXT,
		subdomain TEXT
	);`
	_, err = DB.Exec(query)
	if err != nil {
		logger.Error("初始化数据库表失败: " + err.Error())
	}
}

func GetServersFromDB() []ProxyInstance {
	var servers []ProxyInstance

	if DB == nil {
		logger.Error("数据库对象为 nil，请检查 InitDB 是否执行成功")
		return servers
	}

	rows, err := DB.Query("SELECT name, backend_ip, subdomain FROM servers")
	if err != nil {
		logger.Error("查询服务器列表失败: " + err.Error())
		return servers
	}

	defer rows.Close()

	for rows.Next() {
		var s ProxyInstance
		err := rows.Scan(&s.Name, &s.BackendIP, &s.Subdomain)
		if err != nil {
			logger.Error("解析数据库行失败: " + err.Error())
			continue
		}
		servers = append(servers, s)
	}

	// 检查遍历过程中是否出错
	if err = rows.Err(); err != nil {
		logger.Error("读取数据库过程出错: " + err.Error())
	}

	return servers
}

func StartAPI() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/servers", func(c *gin.Context) {
			servers := GetServersFromDB()
			c.JSON(http.StatusOK, servers)
		})

		v1.POST("/servers", func(c *gin.Context) {
			var s ProxyInstance
			if err := c.ShouldBindJSON(&s); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
				return
			}

			_, err := DB.Exec("INSERT OR REPLACE INTO servers (name, backend_ip, subdomain) VALUES (?, ?, ?)",
				s.Name, s.BackendIP, s.Subdomain)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库操作失败: " + err.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "实例已保存"})
		})

		v1.DELETE("/servers/:name", func(c *gin.Context) {
			name := c.Param("name")
			result, err := DB.Exec("DELETE FROM servers WHERE name = ?", name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			rowsAffected, _ := result.RowsAffected()
			c.JSON(http.StatusOK, gin.H{"status": "deleted", "rows_affected": rowsAffected})
		})
	}

	logger.Info("Gin API 服务已启动: http://localhost:23754/api/v1/")

	if err := r.Run(":23754"); err != nil {
		logger.Error("API 服务启动失败: " + err.Error())
	}
}
