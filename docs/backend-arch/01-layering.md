# 后端分层约定

**版本**：v0.1 · Phase 0

## 目录与职责

```
cmd/                    # 入口（api、migrate）
internal/
  domain/               # 实体与领域类型，无外部依赖
  application/          # 用例编排（Phase 1+ 引入）
  infrastructure/       # DB、Casbin、外部适配
  interfaces/
    http/               # Gin Handlers
    middleware/         # Auth、Tenant、RBAC
  pkg/                  # 跨层工具（response、tenant、rbac）
migrations/             # Goose SQL（Schema 权威来源）
```

## Phase 0 现状

- Handler 可直接使用 `*gorm.DB`（MVP 骨架）
- Phase 1 起逐步引入 `repository` + `application` service

## 多租户查询规范

```go
db := persistence.DBFromContext(globalDB, c.Request.Context())
db.Find(&leads) // 自动 WHERE tenant_id = ?
```

禁止在业务 Handler 中手写未过滤的 `tenant_id` 查询。
