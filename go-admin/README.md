# base

一个基础的 Go Web 项目脚手架，包含以下技术栈：

- `gin`
- `ent`
- `postgres`
- `redis`
- `zap`
- `jwt`

## 匿名文章接口

无需登录或传入 Authorization 请求头，仅返回已发布（`status=1`）的文章。

- `GET /api/articles?page=1&pageSize=10&keyword=关键词&category=cdn帮助`：分页列表：cdn新闻按 `sort` 降序，cdn帮助按创建时间 `createdAt` 降序，其他分类按 `sort` 升序；同值按 ID 降序。默认每页 20 条，最多 200 条。`keyword` 搜索标题，`category` 精确匹配分类（可选），数量统计及分页均按分类过滤；`status` 参数不会改变公开范围。
- `GET /api/articles/:id`：文章详情（包括 HTML 正文）；草稿或不存在的文章返回 HTTP 404，非法 ID 返回 HTTP 400。

成功响应为 `{ "code": 0, "msg": "成功", "data": ... }`。列表的 `data` 包含 `list`、`total`、`page`、`pageSize`；列表文章仅包含 `id`、`sort`、`title`、`category`、`summary`、`cover`、`status`、`author`、`createdAt`、`updatedAt`，不会查询或返回正文 `content`。详情额外返回 `content`（HTML 正文）。

文章支持自由填写分类 `category`，可留空，最多 80 字。管理端创建、修改请求支持该字段，管理端与匿名列表、详情均返回分类；现有文章分类默认为空。匿名列表仍不查询或返回正文。

管理端继续使用需要登录和权限校验的 `/api/system/articles` 接口。后端重启后自动迁移新增分类字段，新接口生效。

文章排序字段 `sort` 为 0 至 2147483647 的整数，默认 0。cdn新闻数值越大越靠前，新建新闻时自动读取该分类最大排序号并填入加 1；cdn帮助按创建时间从新到旧；其他分类数值越小越靠前。相同排序值或时间按 ID 降序固定排列。编辑旧稿保留原排序号，后台与匿名列表均在分页前排序。

## 目录结构

- `cmd/api`：应用启动入口
- `configs`：YAML 配置文件目录
- `internal/config`：配置加载模块
- `internal/db`：PostgreSQL 连接封装
- `internal/redis`：Redis 连接封装
- `internal/handler`：HTTP 处理器
- `internal/middleware`：中间件
- `internal/router`：路由注册
- `internal/logger`：Zap 日志初始化
- `ent/schema`：Ent 实体定义

## 配置说明

- `configs/config.yaml`：基础配置
- `configs/config-dev.yaml`：开发环境配置覆盖
- `configs/config-local.yaml`：本地环境配置覆盖
- `configs/config-prod.yaml`：生产环境配置覆盖

加载顺序：

1. `configs/config.yaml`
2. `configs/config-${APP_ENV}.yaml`

默认 `APP_ENV=dev`。

应用配置示例：

```yaml
app:
  env: dev
  addr: ":8080"
  shutdown_timeout_seconds: 30
  auto_seed: true
```

说明：

- `shutdown_timeout_seconds` 是优雅关闭等待时间
- `auto_seed` 控制是否在启动时自动初始化系统基础数据
- 主程序会把应用生命周期 `ctx` 透传给 HTTP Server 和数据库初始化
- 路由层会把应用生命周期 `ctx` 写入 Gin 上下文，键名为 `app_ctx`
- 收到退出信号后先停止接收新请求，等待已有请求处理完成，再取消应用 `ctx`

PostgreSQL 配置示例：

```yaml
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  name: "base"
  ssl_mode: "disable"
  time_zone: "Asia/Shanghai"
  max_open_conns: 10
  max_idle_conns: 10
  conn_max_lifetime_minutes: 30
```

说明：

- 数据库连接参数不再使用单个 `url` 字符串，统一拆成独立配置项
- `max_open_conns` 控制最大打开连接数
- `max_idle_conns` 控制最大空闲连接数
- `conn_max_lifetime_minutes` 控制连接最大生命周期，单位：分钟

日志配置示例：

```yaml
log:
  level: "info"
  format: "console"
  output: "both"
  directory: "logs"
  filename: "app.log"
  max_size_mb: 10
  max_backups: 10
  max_age_days: 7
  compress: true
```

说明：

- `format` 支持 `console` 和 `json`
- `output` 支持 `console`、`file` 和 `both`
- 单个日志文件最大 `10MB`
- 日志默认保留 `7` 天
- 自动压缩历史日志文件

Redis 配置示例：

```yaml
redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  dial_timeout_seconds: 5
  read_timeout_seconds: 3
  write_timeout_seconds: 3
  pool_size: 10
```

说明：

- 启动时会执行 Redis `PING` 检查连接
- Redis 初始化使用主程序生命周期 `ctx`
- 关闭服务时会关闭 Redis 客户端连接
- 路由层会把 Redis 客户端写入 Gin 上下文，键名为 `redis`

启动示例：

```powershell
$env:APP_ENV="local"
go run ./cmd/api
```

## 后续建议

1. 按实际环境修改 `configs` 下的数据库和 JWT 配置
2. 执行 `go generate ./...` 生成 Ent 代码
3. 在此基础上补充用户、登录、权限、业务模块等代码
