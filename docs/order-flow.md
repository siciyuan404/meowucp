# 订单链路说明

本文件用于梳理当前订单链路实现与待完善点，便于后续补齐。

## 当前实现概览

- 订单创建入口：`internal/service/order_service.go` 的 `CreateOrder`
- 购物车校验：`internal/service/cart_service.go`（库存仅在加购时校验）
- 商品与库存能力：`internal/service/product_service.go`、`internal/service/inventory_service.go`
- 库存日志写入：`internal/service/order_service.go`（当前仅写日志）

## 已知待完善项

- 订单项的 `ProductName`/`SKU` 仍是占位字段：`internal/service/order_service.go`
- 库存扣减未走库存服务，只写了库存日志：`internal/service/order_service.go`
- 缺少事务包裹，可能出现部分成功（订单已建、库存未扣等）：`internal/service/order_service.go`
- 订单创建时未再次校验库存或做并发安全处理：`internal/service/cart_service.go`

## 推荐补齐方向

- 在创建订单时批量读取商品信息，填充订单项字段：`internal/repository/product_repository.go`
- 使用 `InventoryService.AdjustStock` 完成扣减与日志记录：`internal/service/inventory_service.go`
- 将「创建订单 + 创建订单项 + 扣减库存 + 清空购物车」放入同一事务
- 若需要防超卖，库存更新需改为原子更新或锁机制
