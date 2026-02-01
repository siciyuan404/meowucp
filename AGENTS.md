# AGENTS.md

本文件用于指导本仓库的智能体编程工作。
项目基于 Go 1.21，使用 Gin、GORM 与 PostgreSQL。

## 仓库命令

### 构建
- `make build`
- `go build -o bin/api cmd/api/main.go`
- `go build -o bin/admin cmd/admin/main.go`

### 运行
- `make run`（运行 `cmd/api/main.go`）
- `go run cmd/api/main.go`

### 测试
- `make test`（运行 `go test -v ./...`）
- `go test ./...`（全量测试）
- `go test -v ./internal/service`（单个包）
- `go test ./internal/ucp/api -run TestCheckoutCompleteCreatesOrder`
- `go test ./internal/service -run Idempotency`

### 覆盖率
- `make test-coverage`
- `go test -coverprofile=coverage.out ./...`
- `go tool cover -html=coverage.out -o coverage.html`

### 数据库 / 迁移
- `make migrate-up`（执行 `migrations/001_init.sql`）
- `make migrate-down`（删除并重建 schema）

### Docker
- `make docker-up`
- `make docker-down`
- `make docker-logs`

### Lint / 格式化
- 本仓库没有专用 lint 命令。
- 所有 Go 文件使用 `gofmt`。
- 可选：`go vet ./...`

## 代码风格规范

### 通用 Go 风格
- 所有改动必须 `gofmt`。
- 按职责组织文件（api、service、repository、domain、ucp）。
- 函数过长时拆分为小型 helper。
- 只有在现有接口要求时才使用 `context`（大多数代码不需要）。

### Imports
- 统一分组导入：
  - 标准库
  - 第三方库
  - 本地包（`github.com/meowucp/...`）
- 由 `gofmt` 管理顺序和空行。

### 命名
- 遵循 Go 命名规范：
  - 导出符号用 `CamelCase`
  - 非导出符号用 `camelCase`
  - 名称对齐业务语义（如 `OrderService`、`WebhookQueueService`）。
- 统一大写缩写：`ID`、`URL`、`HTTP`、`JSON`、`UCP`。

### 类型与数据模型
- 领域模型位于 `internal/domain`，并使用 GORM 标签。
- 可选字段使用指针（如 `*time.Time`、`*int64`）。
- 尽量使用显式类型，避免 `interface{}`（除非载荷必须）。

### 错误处理
- 遇到错误立即返回，不要静默忽略。
- 静态错误用 `errors.New`，格式化错误用 `fmt.Errorf`。
- 哨兵错误比较用 `errors.Is`（如 `gorm.ErrRecordNotFound`）。
- 错误信息要保留上下文（以操作名开头）。

### API 处理器（Gin）
- 输入校验使用 `ShouldBindJSON` / `ShouldBindQuery` / `PostForm`。
- 统一用 `respondError` 返回错误。
- 返回正确 HTTP 状态码（`400` 表示入参错误，`404` 表示资源不存在）。
- Handler 保持薄逻辑，业务逻辑进入 service 层。

### 服务层
- Service 封装业务逻辑与仓储访问。
- Handler 不直接访问数据库，只通过 service。
- 事务在仓储 helper 中控制（如已有实现）。
- 遵循既有幂等与 webhook 队列模式。

### 仓储层
- 接口定义在 `internal/repository/repository.go`。
- 实现在相邻文件中（如 `order_repository.go`）。
- 使用 GORM 查询风格，未找到返回 `gorm.ErrRecordNotFound`。

### Webhooks
- `WebhookQueueService` 支持：
  - `Enqueue`（原始 payload 入队）
  - `EnqueueOrderEvent`（订单状态变更）
  - `EnqueuePaidEvent`（支付完成）
  - `DeliverOrderEvent`（管理员同步触发）
- 更新订单状态时，需同步触发对应 webhook。

### UCP 专项规范
- UCP 处理器位于 `internal/ucp/api`。
- 使用 `internal/ucp/model` 中的 UCP 类型。
- 状态和消息规则需与 UCP 计划一致（如升级规则）。

### 测试
- 使用 `testing`，适用时用表驱动测试。
- 测试名称清晰且描述性强。
- 优先使用内存 fake 作为依赖。
- Handler 测试必须校验 HTTP 状态码与响应体。
- Webhook 测试允许多条任务入队（可能触发多个事件）。

### 格式与文件
- 保持 CRLF 行尾（Windows）。
- 除非文件已有需求，不要引入非 ASCII 字符。

## 项目结构（高层）
- `cmd/api`：HTTP 服务入口
- `cmd/admin`：管理端服务入口
- `internal/api`：REST API
- `internal/ucp/api`：UCP API
- `internal/service`：业务逻辑
- `internal/repository`：数据库访问
- `internal/domain`：领域模型
- `internal/ucp/worker`：Webhook 投递工作器
- `migrations`：SQL 迁移文件

## 智能体工作说明
- 未发现 Cursor 或 Copilot 指令文件。
- 如新增命令，请同时更新 Makefile 与本文件。
- 保持改动聚焦，避免无关重构。
