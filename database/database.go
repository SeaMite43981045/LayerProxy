// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package database

import (
	"database/sql"
	"os"
	"path/filepath"

	"LayerProxy/logger"
	"LayerProxy/models"

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
			name TEXT UNIQUE NOT NULL CHECK(name <> ''),
			backend_ip TEXT NOT NULL CHECK(backend_ip <> ''),
			subdomain TEXT NOT NULL CHECK(subdomain <> '')
		);`
	_, err = DB.Exec(query)
	if err != nil {
		logger.Error("初始化数据库表失败: " + err.Error())
	}

	_, err = DB.Exec("DELETE FROM servers WHERE name = '' OR name IS NULL")
	if err != nil {
		logger.Error("清理无效数据库记录失败: " + err.Error())
	}
}

func GetServersFromDB() []models.ProxyInstance {
	var servers []models.ProxyInstance

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
		var s models.ProxyInstance
		err := rows.Scan(&s.Name, &s.BackendIP, &s.Subdomain)
		if err != nil {
			logger.Error("解析数据库行失败: " + err.Error())
			continue
		}
		servers = append(servers, s)
	}

	if err = rows.Err(); err != nil {
		logger.Error("读取数据库过程出错: " + err.Error())
	}

	return servers
}
