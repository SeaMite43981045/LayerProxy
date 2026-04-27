// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"LayerProxy/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LogFileInfo struct {
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
}

func HandleListLogFiles(c *gin.Context) {
	files, err := os.ReadDir("logs")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取日志目录"})
		return
	}

	var result []LogFileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		result = append(result, LogFileInfo{
			Name:     info.Name(),
			Size:     info.Size(),
			Modified: info.ModTime(),
		})
	}

	c.JSON(http.StatusOK, result)
}

func HandleDownloadLogFile(c *gin.Context) {
	name := c.Param("name")
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件名"})
		return
	}

	path := filepath.Join("logs", name)
	data, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	c.Data(http.StatusOK, "text/plain; charset=utf-8", data)
}

func HandleDeleteLogFile(c *gin.Context) {
	name := c.Param("name")
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件名"})
		return
	}

	currentLog := logger.GetLogFileName()
	if filepath.Base(currentLog) == name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除正在写入的当前日志文件"})
		return
	}

	path := filepath.Join("logs", name)
	if err := os.Remove(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "日志文件已删除"})
}

func HandleLogStream(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ch := logger.LogBroadcaster.Subscribe()
	defer logger.LogBroadcaster.Unsubscribe(ch)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	for {
		select {
		case entry, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(entry)
			fmt.Fprintf(c.Writer, "event: log\ndata: %s\n\n", data)
			flusher.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
