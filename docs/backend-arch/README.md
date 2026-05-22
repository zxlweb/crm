# 后端架构文档

Go 后端实现约定，与 `backend/` 代码目录对应。

## 文件索引

| 文件 | 内容 |
|------|------|
| [00-database-schema.md](./00-database-schema.md) | 核心表结构、索引、GORM 示例 |

## 代码分层（Clean Architecture）

```
backend/
├── cmd/api/              # HTTP 入口
├── cmd/migrate/          # Goose 迁移 CLI
├── migrations/           # SQL 迁移（权威 Schema）
├── internal/domain/      # 实体模型
├── internal/infrastructure/persistence/  # DB、Casbin
├── internal/interfaces/http/             # Handlers
├── internal/interfaces/middleware/       # Auth / Tenant / RBAC
└── internal/pkg/           # response、tenant、rbac
```

## 关键约定

| 主题 | 约定 |
|------|------|
| 多租户 | 业务表含 `tenant_id`；查询用 `DBFromContext` + `TenantScope` |
| 迁移 | 生产用 `make migrate-up`；`AUTO_MIGRATE=true` 仅本地 |
| RBAC | `resource:action`；启动时 `SyncCasbinPolicies` |
| API 响应 | `internal/pkg/response` 统一结构 |

## 相关 ADR

- [0001 多租户](../architecture/adr/0001-shared-db-multi-tenancy.md)
- [0002 Casbin](../architecture/adr/0002-casbin-rbac.md)
- [0003 Goose](../architecture/adr/0003-goose-migrations.md)
