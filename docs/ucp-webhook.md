# UCP Webhook 对接说明

本文件用于说明 UCP webhook 的接收、入队、投递与告警流程。

## 路由与入口

- 接收端点：`POST /ucp/v1/order-webhooks`
  - 路由注册：`cmd/api/main.go`
  - 处理逻辑：`internal/ucp/api/order_webhook_handler.go`
- 管理端点：
  - 审计列表：`GET /api/v1/admin/ucp/webhook-audits`
  - 告警列表：`GET /api/v1/admin/ucp/webhook-alerts`
  - 队列列表：`GET /api/v1/admin/ucp/webhook-jobs`
  - 队列重试：`POST /api/v1/admin/ucp/webhook-jobs/:id/retry`

## 签名校验与防重放

- Header：`UCP-Signature`、`UCP-Key-Id`
- JWK 校验器：`internal/ucp/security`（初始化见 `cmd/api/main.go`）
- 防重放：基于 payload hash 的 Seen/Mark 机制，避免重复处理：`internal/ucp/api/order_webhook_handler.go`

## 入队、投递与告警

- 入队：写入数据库队列表（job）：`internal/service/webhook_queue_service.go`
- 投递：worker 拉取队列并向 delivery_url 推送：`cmd/worker/main.go`
- 告警：按失败次数与去重窗口触发：`cmd/worker/main.go`

## 配置项

配置位于 `configs/config.yaml` / `configs/config.example.yaml` 的 `ucp.webhook` 节点：

- `jwk_set_url`
- `clock_skew_seconds`
- `delivery_url`
- `delivery_timeout_sec`
- `skip_signature_verify`
- `alert_min_attempts`
- `alert_dedupe_seconds`

## 本地联调

- 一键脚本：`scripts/dev-ucp.ps1`
- 脚本内容包含：迁移、启动 mock-jwks/mock-webhook/API/worker、发送 webhook、检查状态
