# 侧边栏、导航栏与多页面实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建包含顶部导航栏、左侧侧边栏、仪表盘/日志/设置三个页面的完整后台管理界面 Shell，并完成所有后端 API。

**Architecture:** 前端使用 Naive UI 的 NLayout 组合实现全屏布局，vue-i18n 做多语言；后端使用 Gin 提供 REST API 和 SSE 日志流，gopsutil 获取系统信息，GitHub API 检查更新。

**Tech Stack:** Vue 3 + TypeScript + Naive UI + vue-i18n / Go + Gin + gopsutil + SQLite

---

## 文件结构

### 后端新增/修改

| 文件 | 操作 | 说明 |
|------|------|------|
| `utils/version.go` | 创建 | 版本号常量 |
| `models/models.go` | 修改 | 新增 SystemInfo、UserPreferences 结构体 |
| `database/database.go` | 修改 | 新增 preferences 表初始化 |
| `logger/logger.go` | 修改 | 新增 LogBroadcaster 支持 SSE |
| `http/server.go` | 修改 | 注册所有新路由 |
| `http/system.go` | 创建 | 系统信息 handler |
| `http/config.go` | 创建 | 配置管理 handler |
| `http/logs.go` | 创建 | 日志文件和 SSE handler |
| `http/update.go` | 创建 | 检查更新 handler |
| `http/preferences.go` | 创建 | 偏好设置 handler |

### 前端新增/修改

| 文件 | 操作 | 说明 |
|------|------|------|
| `frontend/src/i18n.ts` | 创建 | vue-i18n 配置 |
| `frontend/src/locales/zh.ts` | 创建 | 中文语言包 |
| `frontend/src/locales/en.ts` | 创建 | 英文语言包 |
| `frontend/src/router/index.ts` | 修改 | 补充 Dashboard/Logging/Settings 路由 |
| `frontend/src/http/httpRequest.ts` | 修改 | 添加新 API 方法 |
| `frontend/src/components/AppLayout.vue` | 创建 | 布局 Shell |
| `frontend/src/components/Sidebar.vue` | 创建 | 侧边栏 |
| `frontend/src/components/Navbar.vue` | 创建 | 导航栏 |
| `frontend/src/views/DashboardView.vue` | 创建 | 仪表盘页面 |
| `frontend/src/views/LoggingView.vue` | 创建 | 日志页面 |
| `frontend/src/views/SettingsView.vue` | 创建 | 设置页面 |
| `frontend/src/main.ts` | 修改 | 引入 i18n |
| `frontend/src/App.vue` | 修改 | 使用 AppLayout |

---

## Task 1: 后端基础 — 版本号与数据模型

**Files:**
- Create: `utils/version.go`
- Modify: `models/models.go`

- [ ] **Step 1: 创建版本号文件**

Create `utils/version.go`:

```go
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

const Version = "1.0.0"
const RepoURL = "SeaMite43981045/LayerProxy"
```

- [ ] **Step 2: 在 models.go 中新增结构体**

Modify `models/models.go`，在文件末尾追加：

```go
type SystemInfo struct {
	CPUModel    string  `json:"cpu_model"`
	CPUCores    int     `json:"cpu_cores"`
	CPUThreads  int     `json:"cpu_threads"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryTotal uint64  `json:"memory_total"`
	MemoryUsed  uint64  `json:"memory_used"`
	MemoryFree  uint64  `json:"memory_free"`
	OSName      string  `json:"os_name"`
	OSVersion   string  `json:"os_version"`
	Uptime      uint64  `json:"uptime"`
}

type UserPreferences struct {
	Language string `json:"language"`
	Theme    string `json:"theme"`
}
```

- [ ] **Step 3: Commit**

```bash
git add utils/version.go models/models.go
git commit -m "feat: add version constant and new data models

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 2: 后端 — 数据库偏好设置表

**Files:**
- Modify: `database/database.go`

- [ ] **Step 1: 在数据库初始化中添加 preferences 表**

Modify `database/database.go`，在 `InitDB()` 函数中的 `CREATE TABLE servers` 之后添加：

```go
_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS preferences (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	language TEXT DEFAULT 'zh',
	theme TEXT DEFAULT 'dark'
)`)
if err != nil {
	logger.Error("创建 preferences 表失败: " + err.Error())
	return
}
```

- [ ] **Step 2: Commit**

```bash
git add database/database.go
git commit -m "feat: add preferences table to database

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 3: 后端 — Logger 广播器（SSE 支持）

**Files:**
- Modify: `logger/logger.go`

- [ ] **Step 1: 在 logger.go 中新增广播器逻辑**

Modify `logger/logger.go`。在 `var` 块中替换为：

```go
var (
	LogChan       = make(chan string, 100)
	broadcastChan = make(chan LogEntry, 100)
)

type LogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type Broadcaster struct {
	clients map[chan LogEntry]bool
	mu      sync.RWMutex
}

var LogBroadcaster = &Broadcaster{
	clients: make(map[chan LogEntry]bool),
}

func (b *Broadcaster) Subscribe() chan LogEntry {
	ch := make(chan LogEntry, 10)
	b.mu.Lock()
	b.clients[ch] = true
	b.mu.Unlock()
	return ch
}

func (b *Broadcaster) Unsubscribe(ch chan LogEntry) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *Broadcaster) Broadcast(entry LogEntry) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.clients {
		select {
		case ch <- entry:
		default:
		}
	}
}
```

并在文件头部添加 `sync` 到 import 中：

```go
import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
)
```

- [ ] **Step 2: 修改 Info/Warning/Error 以支持广播**

Replace `Info` function:

```go
func Info(message ...string) {
	msg := strings.Join(message, " ")
	color.Printf("<fg=gray;>[%s]</> <bg=blue;> INFO </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][INFO] - %s\n", GetFormatTime(), msg))
	LogBroadcaster.Broadcast(LogEntry{Time: GetFormatTime(), Level: "INFO", Message: msg})
}
```

Replace `Warning` function:

```go
func Warning(message ...string) {
	msg := strings.Join(message, " ")
	color.Printf("<fg=gray;>[%s]</> <bg=yellow;> WARN </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][WARN] - %s\n", GetFormatTime(), msg))
	LogBroadcaster.Broadcast(LogEntry{Time: GetFormatTime(), Level: "WARN", Message: msg})
}
```

Replace `Error` function:

```go
func Error(message ...string) {
	msg := strings.Join(message, " ")
	color.Printf("<fg=gray;>[%s]</> <bg=red;> ERROR </> %s\n", GetFormatTime(), msg)
	WriteToFile(fmt.Sprintf("[%s][ERROR] - %s\n", GetFormatTime(), msg))
	LogBroadcaster.Broadcast(LogEntry{Time: GetFormatTime(), Level: "ERROR", Message: msg})
}
```

- [ ] **Step 3: Commit**

```bash
git add logger/logger.go
git commit -m "feat: add log broadcaster for SSE streaming

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 4: 后端 — 系统信息 API

**Files:**
- Create: `http/system.go`

- [ ] **Step 1: 安装 gopsutil 依赖**

```bash
cd "E:/go/LayerProxy"
go get github.com/shirou/gopsutil/v4
```

- [ ] **Step 2: 创建系统信息 handler**

Create `http/system.go`:

```go
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"net/http"
	"runtime"
	"time"

	"LayerProxy/logger"
	"LayerProxy/models"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func HandleSystemInfo(c *gin.Context) {
	info := models.SystemInfo{}

	// OS info
	info.OSName = runtime.GOOS
	hostInfo, err := host.Info()
	if err == nil {
		info.OSVersion = hostInfo.PlatformVersion
		info.Uptime = hostInfo.Uptime
	}

	// CPU info
	cpuInfos, err := cpu.Info()
	if err == nil && len(cpuInfos) > 0 {
		info.CPUModel = cpuInfos[0].ModelName
		info.CPUCores = int(cpuInfos[0].Cores)
	}
	cpuCounts, err := cpu.Counts(true)
	if err == nil {
		info.CPUThreads = cpuCounts
	}
	cpuPercents, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercents) > 0 {
		info.CPUUsage = cpuPercents[0]
	}

	// Memory info
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.MemoryTotal = memInfo.Total
		info.MemoryUsed = memInfo.Used
		info.MemoryFree = memInfo.Free
	}

	c.JSON(http.StatusOK, info)
}
```

- [ ] **Step 3: Commit**

```bash
git add http/system.go go.mod go.sum
git commit -m "feat: add system info API with gopsutil

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 5: 后端 — 配置管理 API

**Files:**
- Create: `http/config.go`

- [ ] **Step 1: 创建配置 handler**

Create `http/config.go`:

```go
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"LayerProxy/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigUpdateRequest struct {
	WebPort          string `json:"web_port"`
	PortStartAt      int    `json:"port_start_at"`
	WildcardDomain   string `json:"wildcard_domain"`
	WildcardMainPort string `json:"wildcard_main_port"`
}

func HandleGetConfig(c *gin.Context) {
	cfg, exists := c.Get("config")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "配置未加载"})
		return
	}
	config := cfg.(models.ConfigFile)

	c.JSON(http.StatusOK, gin.H{
		"web_port":            config.Server.WebPort,
		"port_start_at":       config.Port.PortStartAt,
		"wildcard_domain":     config.Wildcard.WildcardDomain,
		"wildcard_main_port":  config.Wildcard.WildcardMainPort,
		"enable_wildcard":     config.Wildcard.EnableWildcard,
	})
}

func HandleUpdateConfig(c *gin.Context) {
	var req ConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	cfg, exists := c.Get("config")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "配置未加载"})
		return
	}
	config := cfg.(models.ConfigFile)

	config.Server.WebPort = req.WebPort
	config.Port.PortStartAt = req.PortStartAt
	config.Wildcard.WildcardDomain = req.WildcardDomain
	config.Wildcard.WildcardMainPort = req.WildcardMainPort

	if err := saveConfigToFile(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "配置已保存，部分设置需重启后生效"})
}
```

- [ ] **Step 2: Commit**

```bash
git add http/config.go
git commit -m "feat: add config read/update API

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 6: 后端 — 日志文件与 SSE API

**Files:**
- Create: `http/logs.go`

- [ ] **Step 1: 创建日志 handler**

Create `http/logs.go`:

```go
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
```

- [ ] **Step 2: Commit**

```bash
git add http/logs.go
git commit -m "feat: add log file list, download, delete and SSE stream API

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 7: 后端 — 检查更新与偏好设置 API

**Files:**
- Create: `http/update.go`
- Create: `http/preferences.go`

- [ ] **Step 1: 创建检查更新 handler**

Create `http/update.go`:

```go
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"LayerProxy/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func HandleCheckUpdate(c *gin.Context) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", utils.RepoURL)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "无法连接到 GitHub"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "GitHub API 返回错误"})
		return
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析响应失败"})
		return
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := utils.Version

	hasUpdate := latestVersion != currentVersion

	c.JSON(http.StatusOK, gin.H{
		"current_version": currentVersion,
		"latest_version":  latestVersion,
		"has_update":      hasUpdate,
		"release_url":     release.HTMLURL,
	})
}
```

- [ ] **Step 2: 创建偏好设置 handler**

Create `http/preferences.go`:

```go
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package http

import (
	"LayerProxy/database"
	"LayerProxy/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetPreferences(c *gin.Context) {
	var prefs models.UserPreferences
	err := database.DB.QueryRow("SELECT language, theme FROM preferences WHERE id = 1").Scan(&prefs.Language, &prefs.Theme)
	if err != nil {
		prefs.Language = "zh"
		prefs.Theme = "dark"
	}
	c.JSON(http.StatusOK, prefs)
}

func HandleUpdatePreferences(c *gin.Context) {
	var prefs models.UserPreferences
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if prefs.Language != "zh" && prefs.Language != "en" {
		prefs.Language = "zh"
	}
	if prefs.Theme != "dark" && prefs.Theme != "light" {
		prefs.Theme = "dark"
	}

	_, err := database.DB.Exec(
		"INSERT INTO preferences (id, language, theme) VALUES (1, ?, ?) ON CONFLICT(id) DO UPDATE SET language=excluded.language, theme=excluded.theme",
		prefs.Language, prefs.Theme,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "偏好设置已保存"})
}
```

- [ ] **Step 3: Commit**

```bash
git add http/update.go http/preferences.go
git commit -m "feat: add update check and preferences API

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 8: 后端 — 注册所有新路由

**Files:**
- Modify: `http/server.go`

- [ ] **Step 1: 在 server.go 中注册新路由**

Modify `http/server.go`，在 `v1.DELETE("/servers/:name", ...)` 之后、`r.NoRoute` 之前添加：

```go
		v1.GET("/system/info", HandleSystemInfo)
		v1.GET("/config", HandleGetConfig)
		v1.POST("/config", HandleUpdateConfig)
		v1.GET("/logs/stream", HandleLogStream)
		v1.GET("/logs/files", HandleListLogFiles)
		v1.GET("/logs/files/:name", HandleDownloadLogFile)
		v1.DELETE("/logs/files/:name", HandleDeleteLogFile)
		v1.GET("/update", HandleCheckUpdate)
		v1.GET("/preferences", HandleGetPreferences)
		v1.POST("/preferences", HandleUpdatePreferences)
```

- [ ] **Step 2: 在 StartAPI 中将 config 存入 gin context**

在 `r := gin.Default()` 之后添加：

```go
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})
```

- [ ] **Step 3: 编译验证后端**

```bash
cd "E:/go/LayerProxy"
go build .
```

Expected: 编译成功，无错误。

- [ ] **Step 4: Commit**

```bash
git add http/server.go
git commit -m "feat: register all new API routes

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 9: 前端 — I18n 配置与语言包

**Files:**
- Create: `frontend/src/i18n.ts`
- Create: `frontend/src/locales/zh.ts`
- Create: `frontend/src/locales/en.ts`
- Modify: `frontend/src/main.ts`

- [ ] **Step 1: 安装 vue-i18n**

```bash
cd "E:/go/LayerProxy/frontend"
npm install vue-i18n
```

- [ ] **Step 2: 创建中文语言包**

Create `frontend/src/locales/zh.ts`:

```typescript
export default {
  sidebar: {
    dashboard: '仪表盘',
    logging: '日志',
    settings: '设置',
  },
  navbar: {
    docs: '官方文档',
    checkUpdate: '检查更新',
    theme: '主题',
    user: '管理员',
  },
  dashboard: {
    title: '仪表盘',
    cpu: 'CPU',
    memory: '内存',
    system: '系统',
    os: '操作系统',
    version: '版本',
    uptime: '运行时间',
    cores: '核心数',
    threads: '线程数',
    total: '总量',
    used: '已用',
    free: '可用',
    usage: '使用率',
    instances: '代理实例',
    addInstance: '添加实例',
    editInstance: '编辑实例',
    deleteInstance: '删除实例',
    name: '名称',
    backendIP: '后端 IP',
    subdomain: '子域名',
    status: '状态',
    actions: '操作',
  },
  logging: {
    title: '日志',
    liveStream: '实时日志流',
    logFiles: '日志文件',
    connected: '已连接',
    disconnected: '未连接',
    clear: '清空',
    download: '下载',
    delete: '删除',
    noLogs: '暂无日志',
  },
  settings: {
    title: '设置',
    systemConfig: '系统配置',
    preferences: '偏好设置',
    webPort: 'Web 端口',
    webPortDesc: 'Web 管理后台监听的端口号',
    portStartAt: '起始端口',
    portStartAtDesc: '为每个代理实例分配端口的起始值',
    wildcardDomain: '泛域名',
    wildcardDomainDesc: '用于子域名路由的通配符域名',
    wildcardMainPort: '泛域名主端口',
    wildcardMainPortDesc: '泛域名模式下的主监听端口',
    language: '语言',
    theme: '主题',
    dark: '深色',
    light: '浅色',
    save: '保存更改',
    saved: '保存成功',
  },
  common: {
    confirm: '确认',
    cancel: '取消',
    deleteConfirm: '确定要删除吗？',
    loading: '加载中...',
    error: '错误',
    success: '成功',
  },
}
```

- [ ] **Step 3: 创建英文语言包**

Create `frontend/src/locales/en.ts`:

```typescript
export default {
  sidebar: {
    dashboard: 'Dashboard',
    logging: 'Logs',
    settings: 'Settings',
  },
  navbar: {
    docs: 'Documentation',
    checkUpdate: 'Check Update',
    theme: 'Theme',
    user: 'Admin',
  },
  dashboard: {
    title: 'Dashboard',
    cpu: 'CPU',
    memory: 'Memory',
    system: 'System',
    os: 'OS',
    version: 'Version',
    uptime: 'Uptime',
    cores: 'Cores',
    threads: 'Threads',
    total: 'Total',
    used: 'Used',
    free: 'Free',
    usage: 'Usage',
    instances: 'Proxy Instances',
    addInstance: 'Add Instance',
    editInstance: 'Edit Instance',
    deleteInstance: 'Delete Instance',
    name: 'Name',
    backendIP: 'Backend IP',
    subdomain: 'Subdomain',
    status: 'Status',
    actions: 'Actions',
  },
  logging: {
    title: 'Logs',
    liveStream: 'Live Stream',
    logFiles: 'Log Files',
    connected: 'Connected',
    disconnected: 'Disconnected',
    clear: 'Clear',
    download: 'Download',
    delete: 'Delete',
    noLogs: 'No logs available',
  },
  settings: {
    title: 'Settings',
    systemConfig: 'System Config',
    preferences: 'Preferences',
    webPort: 'Web Port',
    webPortDesc: 'Port for the web management interface',
    portStartAt: 'Port Start At',
    portStartAtDesc: 'Starting port number for proxy instances',
    wildcardDomain: 'Wildcard Domain',
    wildcardDomainDesc: 'Wildcard domain for subdomain routing',
    wildcardMainPort: 'Wildcard Main Port',
    wildcardMainPortDesc: 'Main listening port in wildcard mode',
    language: 'Language',
    theme: 'Theme',
    dark: 'Dark',
    light: 'Light',
    save: 'Save Changes',
    saved: 'Saved successfully',
  },
  common: {
    confirm: 'Confirm',
    cancel: 'Cancel',
    deleteConfirm: 'Are you sure you want to delete?',
    loading: 'Loading...',
    error: 'Error',
    success: 'Success',
  },
}
```

- [ ] **Step 4: 创建 i18n 配置**

Create `frontend/src/i18n.ts`:

```typescript
import { createI18n } from 'vue-i18n'
import zh from './locales/zh'
import en from './locales/en'

const savedLang = localStorage.getItem('lp_language') || 'zh'

const i18n = createI18n({
  legacy: false,
  locale: savedLang,
  fallbackLocale: 'zh',
  messages: {
    zh,
    en,
  },
})

export default i18n
```

- [ ] **Step 5: 在 main.ts 中引入 i18n**

Modify `frontend/src/main.ts`，在 `createApp(App)` 之前引入 i18n，并在 use 链中添加：

```typescript
import i18n from './i18n'

// ...
const app = createApp(App)
app.use(router)
app.use(i18n)
app.mount('#app')
```

- [ ] **Step 6: Commit**

```bash
git add frontend/src/i18n.ts frontend/src/locales/zh.ts frontend/src/locales/en.ts frontend/src/main.ts frontend/package.json frontend/package-lock.json
git commit -m "feat: add vue-i18n with zh/en language packs

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 10: 前端 — 路由与布局组件

**Files:**
- Modify: `frontend/src/router/index.ts`
- Create: `frontend/src/components/AppLayout.vue`
- Create: `frontend/src/components/Sidebar.vue`
- Create: `frontend/src/components/Navbar.vue`
- Modify: `frontend/src/App.vue`

- [ ] **Step 1: 更新路由配置**

Modify `frontend/src/router/index.ts`:

```typescript
import AppLayout from '@/components/AppLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import LoginView from '@/views/LoginView.vue'
import LoggingView from '@/views/LoggingView.vue'
import SettingsView from '@/views/SettingsView.vue'
import SetupView from '@/views/SetupView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      name: 'Login',
      path: '/login',
      component: LoginView,
    },
    {
      name: 'Setup',
      path: '/setup',
      component: SetupView,
    },
    {
      path: '/',
      component: AppLayout,
      redirect: { name: 'Dashboard' },
      children: [
        {
          name: 'Dashboard',
          path: 'dashboard',
          component: DashboardView,
        },
        {
          name: 'Logging',
          path: 'logging',
          component: LoggingView,
        },
        {
          name: 'Settings',
          path: 'settings',
          component: SettingsView,
        },
      ],
    },
  ],
})

export default router
```

- [ ] **Step 2: 创建 Sidebar 组件**

Create `frontend/src/components/Sidebar.vue`:

```vue
<template>
  <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="180" show-trigger>
    <n-menu
      :value="activeKey"
      :collapsed-width="64"
      :collapsed-icon-size="22"
      :options="menuOptions"
      @update:value="handleMenuSelect"
    />
  </n-layout-sider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NLayoutSider, NMenu } from 'naive-ui'
import { LayoutDashboard, ScrollText, Settings } from 'lucide-vue-next'
import { h } from 'vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const activeKey = computed(() => route.name as string)

const menuOptions = [
  {
    label: t('sidebar.dashboard'),
    key: 'Dashboard',
    icon: () => h(LayoutDashboard),
  },
  {
    label: t('sidebar.logging'),
    key: 'Logging',
    icon: () => h(ScrollText),
  },
  {
    label: t('sidebar.settings'),
    key: 'Settings',
    icon: () => h(Settings),
  },
]

function handleMenuSelect(key: string) {
  router.push({ name: key })
}
</script>
```

- [ ] **Step 3: 创建 Navbar 组件**

Create `frontend/src/components/Navbar.vue`:

```vue
<template>
  <n-layout-header bordered style="height: 48px; padding: 0 20px; display: flex; align-items: center; justify-content: space-between;">
    <div style="font-size: 16px; font-weight: bold;">
      LayerProxy
    </div>
    <n-space align="center">
      <n-button text tag="a" href="https://github.com/SeaMite43981045/LayerProxy" target="_blank">
        {{ t('navbar.docs') }}
      </n-button>
      <n-button text @click="handleCheckUpdate">
        {{ t('navbar.checkUpdate') }}
      </n-button>
      <n-button text @click="toggleTheme">
        <n-icon>
          <component :is="isDark ? Moon : Sun" />
        </n-icon>
      </n-button>
      <n-tag size="small">{{ t('navbar.user') }}</n-tag>
    </n-space>
  </n-layout-header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { NLayoutHeader, NSpace, NButton, NIcon, NTag } from 'naive-ui'
import { Moon, Sun } from 'lucide-vue-next'
import httpRequest from '@/http/httpRequest'

const { t } = useI18n()
const isDark = ref(true)

function toggleTheme() {
  isDark.value = !isDark.value
}

async function handleCheckUpdate() {
  try {
    const res = await httpRequest.checkUpdate()
    if (res.data.has_update) {
      window.$message?.info(`发现新版本: ${res.data.latest_version}`)
    } else {
      window.$message?.success('当前已是最新版本')
    }
  } catch {
    window.$message?.error('检查更新失败')
  }
}
</script>
```

- [ ] **Step 4: 创建 AppLayout 组件**

Create `frontend/src/components/AppLayout.vue`:

```vue
<template>
  <n-layout style="height: 100vh;">
    <Navbar />
    <n-layout has-sider style="height: calc(100vh - 48px);">
      <Sidebar />
      <n-layout-content style="padding: 20px; overflow: auto;">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { NLayout, NLayoutContent } from 'naive-ui'
import Navbar from './Navbar.vue'
import Sidebar from './Sidebar.vue'
</script>
```

- [ ] **Step 5: 修改 App.vue**

Modify `frontend/src/App.vue`:

```vue
<template>
  <router-view />
</template>
```

- [ ] **Step 6: Commit**

```bash
git add frontend/src/router/index.ts frontend/src/components/AppLayout.vue frontend/src/components/Sidebar.vue frontend/src/components/Navbar.vue frontend/src/App.vue
git commit -m "feat: add app layout with sidebar and navbar

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 11: 前端 — HTTP 请求封装扩展

**Files:**
- Modify: `frontend/src/http/httpRequest.ts`

- [ ] **Step 1: 在 HttpRequest 类中添加新方法**

Modify `frontend/src/http/httpRequest.ts`，在类中添加以下方法：

```typescript
  async systemInfo() {
    return this._makeRequest('/v1/system/info', 'get')
  }

  async getConfig() {
    return this._makeRequest('/v1/config', 'get')
  }

  async updateConfig(data: unknown) {
    return this._makeRequest('/v1/config', 'post', data)
  }

  async listLogFiles() {
    return this._makeRequest('/v1/logs/files', 'get')
  }

  async downloadLogFile(name: string) {
    return this._makeRequest(`/v1/logs/files/${name}`, 'get')
  }

  async deleteLogFile(name: string) {
    return this._makeRequest(`/v1/logs/files/${name}`, 'delete')
  }

  async checkUpdate() {
    return this._makeRequest('/v1/update', 'get')
  }

  async getPreferences() {
    return this._makeRequest('/v1/preferences', 'get')
  }

  async updatePreferences(data: unknown) {
    return this._makeRequest('/v1/preferences', 'post', data)
  }
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/http/httpRequest.ts
git commit -m "feat: extend http client with new APIs

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 12: 前端 — Dashboard 页面

**Files:**
- Create: `frontend/src/views/DashboardView.vue`
- Delete: `frontend/src/views/HomeView.vue` (after confirming DashboardView works)

- [ ] **Step 1: 创建 DashboardView**

Create `frontend/src/views/DashboardView.vue`:

```vue
<template>
  <div style="display: flex; flex-direction: column; gap: 16px; height: 100%;">
    <div style="display: flex; gap: 16px;">
      <n-card style="flex: 1;">
        <n-statistic :label="t('dashboard.cpu')" :value="systemInfo.cpu_model">
          <template #prefix>
            <n-icon><Cpu /></n-icon>
          </template>
        </n-statistic>
        <n-space vertical style="margin-top: 8px;">
          <n-text depth="3">{{ systemInfo.cpu_cores }} {{ t('dashboard.cores') }} / {{ systemInfo.cpu_threads }} {{ t('dashboard.threads') }}</n-text>
          <n-progress type="line" :percentage="Math.round(systemInfo.cpu_usage)" indicator-placement="inside" />
        </n-space>
      </n-card>
      <n-card style="flex: 1;">
        <n-statistic :label="t('dashboard.memory')" :value="formatBytes(systemInfo.memory_total)">
          <template #prefix>
            <n-icon><MemoryStick /></n-icon>
          </template>
        </n-statistic>
        <n-space vertical style="margin-top: 8px;">
          <n-text depth="3">{{ t('dashboard.used') }}: {{ formatBytes(systemInfo.memory_used) }} / {{ t('dashboard.free') }}: {{ formatBytes(systemInfo.memory_free) }}</n-text>
          <n-progress type="line" :percentage="memoryUsagePercent" indicator-placement="inside" />
        </n-space>
      </n-card>
      <n-card style="flex: 1;">
        <n-statistic :label="t('dashboard.system')" :value="systemInfo.os_name">
          <template #prefix>
            <n-icon><Monitor /></n-icon>
          </template>
        </n-statistic>
        <n-space vertical style="margin-top: 8px;">
          <n-text depth="3">{{ t('dashboard.version') }}: {{ systemInfo.os_version }}</n-text>
          <n-text depth="3">{{ t('dashboard.uptime') }}: {{ formatUptime(systemInfo.uptime) }}</n-text>
        </n-space>
      </n-card>
    </div>
    <n-card style="flex: 1; overflow: auto;">
      <template #header>
        <n-space align="center">
          <span>{{ t('dashboard.instances') }}</span>
          <n-button size="small" @click="showModal = true">{{ t('dashboard.addInstance') }}</n-button>
        </n-space>
      </template>
      <n-data-table :columns="columns" :data="instances" :bordered="false" size="small" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard, NStatistic, NIcon, NProgress, NSpace, NText,
  NDataTable, NButton,
} from 'naive-ui'
import { Cpu, MemoryStick, Monitor } from 'lucide-vue-next'
import httpRequest from '@/http/httpRequest'

const { t } = useI18n()

const systemInfo = ref({
  cpu_model: '', cpu_cores: 0, cpu_threads: 0, cpu_usage: 0,
  memory_total: 0, memory_used: 0, memory_free: 0,
  os_name: '', os_version: '', uptime: 0,
})

const instances = ref([])
const showModal = ref(false)

const memoryUsagePercent = computed(() => {
  if (!systemInfo.value.memory_total) return 0
  return Math.round((systemInfo.value.memory_used / systemInfo.value.memory_total) * 100)
})

const columns = [
  { title: t('dashboard.name'), key: 'name' },
  { title: t('dashboard.backendIP'), key: 'backend_ip' },
  { title: t('dashboard.subdomain'), key: 'subdomain' },
  { title: t('dashboard.status'), key: 'status' },
  { title: t('dashboard.actions'), key: 'actions' },
]

function formatBytes(bytes: number) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatUptime(seconds: number) {
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  return `${d}d ${h}h ${m}m`
}

async function loadData() {
  try {
    const sysRes = await httpRequest.systemInfo()
    systemInfo.value = sysRes.data
    const instRes = await httpRequest.getServers()
    instances.value = instRes.data
  } catch {
    window.$message?.error('加载数据失败')
  }
}

onMounted(loadData)
</script>
```

- [ ] **Step 2: 在 httpRequest.ts 中补充 getServers 方法（如果不存在）**

检查 `frontend/src/http/httpRequest.ts`，确保有 `getServers` 方法：

```typescript
  async getServers() {
    return this._makeRequest('/v1/servers', 'get')
  }
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/views/DashboardView.vue
git commit -m "feat: add Dashboard view with system info cards

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 13: 前端 — Logging 页面

**Files:**
- Create: `frontend/src/views/LoggingView.vue`

- [ ] **Step 1: 创建 LoggingView**

Create `frontend/src/views/LoggingView.vue`:

```vue
<template>
  <div style="display: flex; gap: 16px; height: 100%;">
    <n-card style="flex: 1; display: flex; flex-direction: column;" content-style="flex: 1; display: flex; flex-direction: column; padding: 0;">
      <template #header>
        <n-space align="center" justify="space-between" style="width: 100%;">
          <span>{{ t('logging.liveStream') }}</span>
          <n-space>
            <n-tag :type="connected ? 'success' : 'error'" size="small">
              {{ connected ? t('logging.connected') : t('logging.disconnected') }}
            </n-tag>
            <n-button size="small" @click="clearLogs">{{ t('logging.clear') }}</n-button>
          </n-space>
        </n-space>
      </template>
      <div
        ref="logContainer"
        style="flex: 1; background: var(--n-color); font-family: monospace; padding: 12px; overflow-y: auto; white-space: pre-wrap; font-size: 13px;"
      >
        <div v-if="logs.length === 0" style="color: var(--n-text-color-3);">{{ t('logging.noLogs') }}</div>
        <div v-for="(log, index) in logs" :key="index">
          <span style="color: var(--n-text-color-3);">[{{ log.time }}]</span>
          <n-tag :type="logLevelColor(log.level)" size="tiny" style="margin: 0 4px;">{{ log.level }}</n-tag>
          <span>{{ log.message }}</span>
        </div>
      </div>
    </n-card>
    <n-card style="width: 240px;" content-style="padding: 0;">
      <template #header>{{ t('logging.logFiles') }}</template>
      <n-list hoverable clickable>
        <n-list-item v-for="file in logFiles" :key="file.name">
          <n-thing :title="file.name" :description="formatBytes(file.size)" />
          <template #suffix>
            <n-space>
              <n-button text @click="downloadFile(file.name)">
                <n-icon><Download /></n-icon>
              </n-button>
              <n-button text type="error" @click="deleteFile(file.name)">
                <n-icon><Trash2 /></n-icon>
              </n-button>
            </n-space>
          </template>
        </n-list-item>
      </n-list>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard, NSpace, NTag, NButton, NList, NListItem,
  NThing, NIcon,
} from 'naive-ui'
import { Download, Trash2 } from 'lucide-vue-next'
import httpRequest from '@/http/httpRequest'

const { t } = useI18n()

const logs = ref<{ time: string; level: string; message: string }[]>([])
const logFiles = ref<{ name: string; size: number }[]>([])
const connected = ref(false)
const logContainer = ref<HTMLDivElement>()
let eventSource: EventSource | null = null

function logLevelColor(level: string) {
  switch (level) {
    case 'ERROR': return 'error'
    case 'WARN': return 'warning'
    default: return 'info'
  }
}

function formatBytes(bytes: number) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function clearLogs() {
  logs.value = []
}

function scrollToBottom() {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

function connectSSE() {
  const token = localStorage.getItem('lp_token')
  eventSource = new EventSource(`/api/v1/logs/stream?token=${token}`)
  eventSource.addEventListener('log', (e) => {
    try {
      const data = JSON.parse(e.data)
      logs.value.push(data)
      if (logs.value.length > 500) {
        logs.value.shift()
      }
      scrollToBottom()
    } catch {}
  })
  eventSource.onopen = () => { connected.value = true }
  eventSource.onerror = () => { connected.value = false }
}

async function loadLogFiles() {
  try {
    const res = await httpRequest.listLogFiles()
    logFiles.value = res.data
  } catch {}
}

async function downloadFile(name: string) {
  try {
    const res = await httpRequest.downloadLogFile(name)
    const blob = new Blob([res.data])
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = name
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    window.$message?.error('下载失败')
  }
}

async function deleteFile(name: string) {
  try {
    await httpRequest.deleteLogFile(name)
    window.$message?.success('删除成功')
    loadLogFiles()
  } catch {
    window.$message?.error('删除失败')
  }
}

onMounted(() => {
  connectSSE()
  loadLogFiles()
})

onUnmounted(() => {
  eventSource?.close()
})
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/LoggingView.vue
git commit -m "feat: add Logging view with SSE stream and file management

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 14: 前端 — Settings 页面

**Files:**
- Create: `frontend/src/views/SettingsView.vue`

- [ ] **Step 1: 创建 SettingsView**

Create `frontend/src/views/SettingsView.vue`:

```vue
<template>
  <n-card style="max-width: 800px; margin: 0 auto;">
    <n-tabs type="line" animated>
      <n-tab-pane :name="t('settings.systemConfig')" :tab="t('settings.systemConfig')">
        <n-form label-placement="left" label-width="160px" style="margin-top: 16px;">
          <n-form-item :label="t('settings.webPort')">
            <n-input v-model:value="config.web_port" />
            <template #feedback>{{ t('settings.webPortDesc') }}</template>
          </n-form-item>
          <n-form-item :label="t('settings.portStartAt')">
            <n-input-number v-model:value="config.port_start_at" />
            <template #feedback>{{ t('settings.portStartAtDesc') }}</template>
          </n-form-item>
          <n-form-item :label="t('settings.wildcardDomain')">
            <n-input v-model:value="config.wildcard_domain" />
            <template #feedback>{{ t('settings.wildcardDomainDesc') }}</template>
          </n-form-item>
          <n-form-item :label="t('settings.wildcardMainPort')">
            <n-input v-model:value="config.wildcard_main_port" />
            <template #feedback>{{ t('settings.wildcardMainPortDesc') }}</template>
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="saveConfig">{{ t('settings.save') }}</n-button>
          </n-form-item>
        </n-form>
      </n-tab-pane>
      <n-tab-pane :name="t('settings.preferences')" :tab="t('settings.preferences')">
        <n-form label-placement="left" label-width="160px" style="margin-top: 16px;">
          <n-form-item :label="t('settings.language')">
            <n-select v-model:value="prefs.language" :options="languageOptions" />
          </n-form-item>
          <n-form-item :label="t('settings.theme')">
            <n-select v-model:value="prefs.theme" :options="themeOptions" />
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="savePrefs">{{ t('settings.save') }}</n-button>
          </n-form-item>
        </n-form>
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard, NTabs, NTabPane, NForm, NFormItem,
  NInput, NInputNumber, NSelect, NButton,
} from 'naive-ui'
import httpRequest from '@/http/httpRequest'

const { t, locale } = useI18n()

const config = ref({
  web_port: '',
  port_start_at: 25565,
  wildcard_domain: '',
  wildcard_main_port: '',
})

const prefs = ref({
  language: 'zh',
  theme: 'dark',
})

const languageOptions = [
  { label: '中文', value: 'zh' },
  { label: 'English', value: 'en' },
]

const themeOptions = [
  { label: t('settings.dark'), value: 'dark' },
  { label: t('settings.light'), value: 'light' },
]

async function loadConfig() {
  try {
    const res = await httpRequest.getConfig()
    config.value = res.data
  } catch {}
}

async function saveConfig() {
  try {
    await httpRequest.updateConfig(config.value)
    window.$message?.success(t('settings.saved'))
  } catch {
    window.$message?.error(t('common.error'))
  }
}

async function loadPrefs() {
  try {
    const res = await httpRequest.getPreferences()
    prefs.value = res.data
  } catch {}
}

async function savePrefs() {
  try {
    await httpRequest.updatePreferences(prefs.value)
    localStorage.setItem('lp_language', prefs.value.language)
    locale.value = prefs.value.language
    window.$message?.success(t('settings.saved'))
  } catch {
    window.$message?.error(t('common.error'))
  }
}

onMounted(() => {
  loadConfig()
  loadPrefs()
})
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/SettingsView.vue
git commit -m "feat: add Settings view with config and preferences tabs

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## Task 15: 集成测试与收尾

**Files:**
- Modify: `http/server.go` (SSE token 传递)
- Delete: `frontend/src/views/HomeView.vue`

- [ ] **Step 1: 修改 SSE token 传递方式**

SSE 无法自定义 headers，需要通过 query param 传递 token。修改 `http/server.go` 中的 JWT middleware 或创建一个 SSE 专用的 token 解析逻辑。

在 `HandleLogStream` 中，在函数开头添加 token 解析：

Modify `http/logs.go` 的 `HandleLogStream`，在函数开头添加：

```go
	token := c.Query("token")
	if token == "" {
		token = c.GetHeader("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}
	if token == "" || !utils.ValidateToken(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
```

并在 `utils/jwt.go` 中暴露 `ValidateToken` 函数（如果不存在）：

检查 `utils/jwt.go`，确保有 ValidateToken：

```go
func ValidateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return err == nil && token.Valid
}
```

- [ ] **Step 2: 编译验证前后端**

```bash
cd "E:/go/LayerProxy"
go build .
```

```bash
cd "E:/go/LayerProxy/frontend"
npm run build
```

Expected: 两者都编译成功。

- [ ] **Step 3: 删除 HomeView**

```bash
git rm frontend/src/views/HomeView.vue
```

- [ ] **Step 4: 最终 Commit**

```bash
git add -A
git commit -m "feat: complete sidebar, navbar, dashboard, logging and settings

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```

---

## 自审清单

### 1. Spec 覆盖度

| 需求 | 对应 Task |
|------|----------|
| 顶部导航栏 + 左侧边栏 + 内容区布局 | Task 10 |
| 仪表盘页面（CPU/内存/系统信息） | Task 12 |
| 日志页面（SSE 实时流 + 文件管理） | Task 13 |
| 设置页面（系统配置 + 偏好设置） | Task 14 |
| 路由补充（Dashboard/Logging/Settings） | Task 10 |
| 后端 API（系统信息/配置/日志/更新/偏好） | Task 4-8 |
| I18n 多语言 | Task 9 |
| 填满屏幕、Naive UI 默认颜色 | Task 10-14 中使用 Naive UI 组件 |

### 2. 占位符扫描

计划中无 TBD、TODO、"implement later" 等占位符。所有代码均为可直接使用的完整代码。

### 3. 类型一致性

- `models.SystemInfo`、`models.UserPreferences` 在 Task 1 定义，后续 Task 4、7 中使用，字段名一致。
- `ConfigUpdateRequest` 字段名与 `config.json` 结构一致。
- API 路由路径前后端一致。
