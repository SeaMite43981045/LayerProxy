# 侧边栏、导航栏与多页面设计文档

## 概述

为 LayerProxy 增加一个完整的后台管理界面 shell，包含顶部导航栏、左侧侧边栏、以及三个功能页面（仪表盘、日志、设置）。同时完成后端 API 以支持实时日志流、系统信息获取、配置管理和更新检查。

## 布局架构

采用经典的 **顶部导航栏 + 左侧边栏 + 右侧内容区** 布局：

- **NLayout** 根容器（`height: 100vh`，填满整个屏幕）
- **NLayoutHeader**（顶部）：全宽导航栏，左侧显示 Logo，右侧放置功能按钮
- **NLayout**（header 下方）：
  - **NLayoutSider**（左侧）：侧边栏导航，3 个菜单项
  - **NLayoutContent**（右侧）：路由视图 `router-view`，页面内容填满剩余空间

**约束：** 使用 Naive UI 默认颜色，不自定义 `NLayout`、`NLayoutSider` 等组件的背景色。

## 组件设计

### AppLayout.vue

Shell 组件，包裹整个管理后台界面：
- 渲染 `Navbar` + `Sidebar` + `router-view`
- 使用 `NLayout` 组合实现整体布局

### Sidebar.vue

左侧导航栏组件：
- 3 个菜单项：仪表盘、日志、设置
- 使用 `NMenu`，搭配 Lucide 图标：`LayoutDashboard`、`ScrollText`、`Settings`
- 当前激活项高亮，点击切换路由

### Navbar.vue

顶部导航栏组件：
- 左侧：应用 Logo / 名称
- 右侧功能项：
  - **官方文档**：外链跳转按钮（`https://github.com/SeaMite43981045/LayerProxy`）
  - **检查更新**：调用后端 API 检查 GitHub Release 版本
  - **主题切换**：仅显示当前主题图标（Dark → `Moon`，Light → `Sun`，来自 lucide-vue）
  - **用户头像/名称**：占位展示

## 页面设计

### Dashboard（仪表盘）

路由：`/dashboard`，Name: `Dashboard`

内容：
- **系统信息卡片**（顶部横向排列，使用 `NCard`）：
  - CPU：型号、核心数/线程数、使用率进度条（`NProgress`），图标 `Cpu`
  - 内存：总容量、已用/可用、使用率进度条，图标 `MemoryStick`
  - 系统：OS 名称、版本号、运行时间，图标 `Monitor`
- **代理实例列表**（下方填满）：使用 `NDataTable`，保留原 HomeView 的表格样式和功能（增删改查）

### Logging（日志）

路由：`/logging`，Name: `Logging`

内容：
- **左侧：实时日志流**（占主要宽度）
  - 等宽字体日志显示区域（`NCode` 风格或自定义 div）
  - SSE 连接到 `/api/v1/logs/stream`
  - 顶部工具栏：连接状态指示、清空按钮（使用 `NButton` 默认样式）
  - 自动滚动到底部，用户手动上滑时暂停自动滚动
- **右侧：日志文件列表**（固定宽度约 240px）
  - 列出 `logs/` 目录下所有日志文件
  - 每项显示文件名、下载按钮（`Download` 图标）、删除按钮（`Trash2` 图标）

### Settings（设置）

路由：`/settings`，Name: `Settings`

内容：
- 使用 `NTabs` 分为两个标签页：

**系统配置标签页：**
- 使用 `NForm` + `NFormItem` 左右布局
- 每个配置项包含：标签名称 + 解释说明（subtitle） + 输入框
- 字段：
  - Web 端口 — "Web 管理后台监听的端口号"
  - 起始端口 — "为每个代理实例分配端口的起始值"
  - 泛域名 — "用于子域名路由的通配符域名"
  - 泛域名主端口 — "泛域名模式下的主监听端口"
- 保存按钮在右下角，点击后调用 `/api/v1/config`

**偏好设置标签页：**
- 语言选择器（`NSelect`）：中文 / English
- 默认主题偏好（`NSelect`）：深色 / 浅色
- 保存到后端 `/api/v1/preferences`，同时也缓存在 localStorage

## 路由变更

在 `frontend/src/router/index.ts` 中补充以下路由：

```typescript
{
  name: 'Dashboard',
  path: '/dashboard',
  component: DashboardView,
},
{
  name: 'Logging',
  path: '/logging',
  component: LoggingView,
},
{
  name: 'Settings',
  path: '/settings',
  component: SettingsView,
},
```

其中 `/` 重定向到 `/dashboard`。

## API 设计

| 方法 | 路由 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/system/info` | 获取系统信息（CPU/内存/OS/运行时间） | JWT |
| GET | `/api/v1/config` | 获取当前系统配置（对应 config.json） | JWT |
| POST | `/api/v1/config` | 更新系统配置，写入 config.json | JWT |
| GET | `/api/v1/logs/stream` | SSE 实时日志流 | JWT |
| GET | `/api/v1/logs/files` | 获取日志文件列表 | JWT |
| GET | `/api/v1/logs/files/:name` | 下载指定日志文件 | JWT |
| DELETE | `/api/v1/logs/files/:name` | 删除指定日志文件 | JWT |
| GET | `/api/v1/update` | 检查 GitHub Release 更新 | JWT |
| GET | `/api/v1/preferences` | 获取用户偏好设置 | JWT |
| POST | `/api/v1/preferences` | 保存用户偏好设置 | JWT |

### 接口详情

**GET /api/v1/system/info**

返回当前服务器的系统信息：

```json
{
  "cpu_model": "Intel(R) Core(TM) i7-13700K",
  "cpu_cores": 16,
  "cpu_threads": 24,
  "cpu_usage": 60.5,
  "memory_total": 34359738368,
  "memory_used": 13743895347,
  "memory_free": 20615843021,
  "os_name": "windows",
  "os_version": "10.0.26300",
  "uptime": 302400
}
```

**GET /api/v1/config**

返回当前 `config.json` 的内容（脱敏，不包含 key 字段）。

**POST /api/v1/config**

请求体：

```json
{
  "web_port": "23754",
  "port_start_at": 25565,
  "wildcard_domain": "*.example.com",
  "wildcard_main_port": "25565"
}
```

更新后写入 `config.json`，部分配置需要重启生效，后端返回提示信息。

**GET /api/v1/logs/stream**

SSE 端点，Content-Type: `text/event-stream`。

实时推送通过 `logger` 包输出的每一条日志。格式：

```
event: log
data: {"time":"2026-04-27T10:23:01Z","level":"INFO","message":"Server started"}

```

**GET /api/v1/logs/files**

返回日志目录下的文件列表：

```json
[
  {"name": "app.log", "size": 1024000, "modified": "2026-04-27T10:00:00Z"},
  {"name": "app-20260426.log", "size": 512000, "modified": "2026-04-26T23:59:00Z"}
]
```

**GET /api/v1/logs/files/:name**

直接返回指定日志文件的二进制内容，Content-Type: `text/plain`，支持浏览器下载。

**DELETE /api/v1/logs/files/:name**

删除指定的日志文件（禁止删除正在被写入的当前日志文件）。

**GET /api/v1/update**

查询 GitHub Release API，对比当前版本号：

```json
{
  "current_version": "1.0.0",
  "latest_version": "1.1.0",
  "has_update": true,
  "release_url": "https://github.com/SeaMite43981045/LayerProxy/releases/tag/v1.1.0"
}
```

当前版本号硬编码在独立的 `version.go` 中（例如 `const Version = "1.0.0"`）。

**GET /api/v1/preferences**

```json
{
  "language": "zh",
  "theme": "dark"
}
```

**POST /api/v1/preferences**

保存用户偏好，存入 SQLite 数据库（新增 `preferences` 表）。

## I18n 设计

- 使用 `vue-i18n` 实现前端多语言
- 默认语言：中文（`zh`）
- 支持语言：中文（`zh`）、英文（`en`）
- 语言文件存放在 `frontend/src/locales/` 目录下
- 语言切换后即时生效，偏好保存到后端和 localStorage
- **后端日志保持中文输出，不做国际化**

## 数据模型

### SystemInfo

```go
type SystemInfo struct {
    CPUModel     string  `json:"cpu_model"`
    CPUCores     int     `json:"cpu_cores"`
    CPUThreads   int     `json:"cpu_threads"`
    CPUUsage     float64 `json:"cpu_usage"`
    MemoryTotal  uint64  `json:"memory_total"`
    MemoryUsed   uint64  `json:"memory_used"`
    MemoryFree   uint64  `json:"memory_free"`
    OSName       string  `json:"os_name"`
    OSVersion    string  `json:"os_version"`
    Uptime       uint64  `json:"uptime"`
}
```

### UserPreferences

```go
type UserPreferences struct {
    Language string `json:"language"`
    Theme    string `json:"theme"`
}
```

## 依赖

前端新增依赖：
- `vue-i18n`（多语言）

后端新增依赖：
- `github.com/shirou/gopsutil/v4`（跨平台获取 CPU、内存、系统详细信息）

SSE 使用 Go 标准库 `net/http` 实现。
