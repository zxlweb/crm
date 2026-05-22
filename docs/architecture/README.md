# 系统架构文档

描述系统**如何组织**（How），技术取舍见 [adr/](./adr/README.md)（Why）。

## 文件索引

| 文件 | 主题 |
|------|------|
| [00-main-architecture.md](./00-main-architecture.md) | 系统上下文、技术栈、分层、风险 |
| [01-multi-tenancy.md](./01-multi-tenancy.md) | 多租户方案与实现要点 |
| [02-rbac-design.md](./02-rbac-design.md) | RBAC + ABAC 权限模型 |

## ADR 决策记录

| ADR | 决策 |
|-----|------|
| [0001](./adr/0001-shared-db-multi-tenancy.md) | 共享 DB + tenant_id |
| [0002](./adr/0002-casbin-rbac.md) | Casbin + resource:action |
| [0003](./adr/0003-goose-migrations.md) | Goose 迁移 |

新建决策：复制 [adr/0000-template.md](./adr/0000-template.md)。

## 与代码映射

| 架构主题 | 代码位置 |
|----------|----------|
| 多租户中间件 | `backend/internal/interfaces/middleware/tenant.go` |
| Tenant Context | `backend/internal/pkg/tenant/` |
| RBAC | `backend/internal/infrastructure/persistence/casbin*.go` |
