# MeowUCP 电商项目

基于 Golang + PostgreSQL + Redis 的轻量级电商系统。

## 技术栈

- **后端**: Golang 1.21
- **Web框架**: Gin
- **数据库**: PostgreSQL
- **缓存**: Redis
- **ORM**: GORM
- **消息队列**: Redis Stream
- **全文搜索**: PostgreSQL pg_trgm

## 功能模块

- 商品系统（商品管理、分类、搜索）
- 购物系统（购物车、结算）
- 订单系统（订单创建、状态管理、退款）
- 用户系统（注册登录、权限管理）
- 支付系统（在线支付集成）
- 库存系统（库存管理）

## 项目结构

```
meowucp/
├── cmd/                    # 应用入口
│   ├── api/               # API 服务
│   └── admin/             # 管理后台
├── internal/              # 私有代码
│   ├── api/              # HTTP handlers
│   ├── domain/           # 领域模型
│   ├── repository/       # 数据访问层
│   ├── service/          # 业务逻辑层
│   └── middleware/      # 中间件
├── pkg/                  # 公共库
├── configs/              # 配置文件
├── migrations/           # 数据库迁移
├── scripts/              # 工具脚本
└── web/                  # 前端资源
```

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 14+
- Redis 6+

### 安装

```bash
# 克隆项目
git clone https://github.com/meowucp/meowucp.git
cd meowucp

# 安装依赖
go mod download

# 复制配置文件
cp configs/config.example.yaml configs/config.yaml

# 启动服务
docker-compose up -d
go run cmd/api/main.go
```

### 数据库初始化

```bash
# 运行迁移
go run scripts/migrate/main.go up
```

### UCP 本地联调一键脚本

```powershell
powershell -ExecutionPolicy Bypass -File scripts/dev-ucp.ps1
```

可选参数：

```powershell
# 清理已有进程并启动
powershell -ExecutionPolicy Bypass -File scripts/dev-ucp.ps1 -StopExisting

# 仅清理，不启动
powershell -ExecutionPolicy Bypass -File scripts/dev-ucp.ps1 -StopExisting -StopOnly

# 启动但不发送 webhook / 不检查状态
powershell -ExecutionPolicy Bypass -File scripts/dev-ucp.ps1 -NoVerify

# 强制 mock webhook 失败（用于重试/告警验证）
powershell -ExecutionPolicy Bypass -File scripts/dev-ucp.ps1 -StopExisting -MockWebhookFailStatus 500
```

## 开发

### 运行 API 服务

```bash
go run cmd/api/main.go
```

### 运行管理后台

```bash
go run cmd/admin/main.go
```

### 运行测试

```bash
go test ./...
```

### 并发库存扣减集成测试

该测试会连接真实数据库，仅在 `integration` tag 下运行：

测试过程中会对 `products` 表执行 `AutoMigrate`，并在结束后删除本次创建的测试数据。

```bash
go test ./internal/repository -tags integration -run TestProductRepositoryUpdateStockWithDeltaIsAtomic
```

可选环境变量：

- `TEST_CONFIG_PATH`（默认 `configs/config.yaml`）
- `TEST_DB_HOST`
- `TEST_DB_PORT`
- `TEST_DB_USER`
- `TEST_DB_PASSWORD`
- `TEST_DB_NAME`
- `TEST_DB_SSLMODE`

## 订单链路注意事项

- 订单创建时应填充真实商品信息（名称、SKU），目前为占位字段：`internal/service/order_service.go`
- 库存扣减建议通过库存服务统一处理（含库存校验与变更日志）：`internal/service/inventory_service.go`
- 订单创建、订单项写入、库存扣减、清空购物车应放在同一事务中，避免部分成功：`internal/service/order_service.go`

## UCP Webhook 对接说明

### 路由与入口

- 接收端点：`POST /ucp/v1/order-webhooks`（`cmd/api/main.go`、`internal/ucp/api/order_webhook_handler.go`）
- 管理端点：
  - 审计列表：`GET /api/v1/admin/ucp/webhook-audits`
  - 告警列表：`GET /api/v1/admin/ucp/webhook-alerts`
  - 队列列表与重试：`GET /api/v1/admin/ucp/webhook-jobs`、`POST /api/v1/admin/ucp/webhook-jobs/:id/retry`

### 签名校验与防重放

- Header：`UCP-Signature`、`UCP-Key-Id`
- JWK 地址：`configs/config.yaml` / `configs/config.example.yaml` 的 `ucp.webhook.jwk_set_url`
- 允许配置跳过校验：`ucp.webhook.skip_signature_verify`
- 防重放逻辑：基于 payload hash 进行去重与限时标记：`internal/ucp/api/order_webhook_handler.go`

### 入队与投递

- 入队为数据库队列表（job），由 worker 拉取并投递：`internal/service/webhook_queue_service.go`、`cmd/worker/main.go`
- 投递地址与超时配置：`configs/config.yaml` / `configs/config.example.yaml` 的 `ucp.webhook.delivery_url`、`ucp.webhook.delivery_timeout_sec`
- 告警策略：`ucp.webhook.alert_min_attempts`、`ucp.webhook.alert_dedupe_seconds`

### 本地联调脚本

- 一键启动与联调：`scripts/dev-ucp.ps1`
  - 自动迁移、启动 mock-jwks/mock-webhook/API/worker，并可发送 webhook 进行验证

## 部署

支持使用 Docker Compose 进行部署。

## License

MIT
