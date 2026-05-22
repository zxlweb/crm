# ADR-0001: 共享数据库 + tenant_id 多租户隔离

**状态**：Accepted  
**日期**：2026-05-21  
**决策者**：Software Architect

## 背景（Context）

CRM 需支持多企业（租户）共用一套部署，数据必须隔离，且 MVP 阶段团队规模与运维资源有限。

## 决策（Decision）

采用 **共享 PostgreSQL 数据库 + 业务表 `tenant_id` 列** 实现行级隔离。

实现要点：
- 业务表强制 `tenant_id UUID NOT NULL`
- HTTP 层 `X-Tenant-ID` + `context.Context` 传递租户
- GORM 查询使用 `TenantScope` 自动附加 `WHERE tenant_id = ?`
- 生产环境可叠加 PostgreSQL RLS（后续）

## 备选方案（Options Considered）

| 方案 | 优点 | 缺点 |
|------|------|------|
| 共享 DB + tenant_id ✅ | 成本低、迭代快 | 隔离依赖应用层纪律 |
| Schema 隔离 | 隔离性更高 | 迁移与连接管理复杂 |
| 独立数据库 | 隔离最强 | 运维与成本高，不适合 MVP |

## 后果（Consequences）

### 正面
- 单库部署简单，适合 SaaS MVP 快速验证
- 与 GORM Scope、中间件模式契合

### 负面 / 风险
- 开发疏漏可能导致跨租户查询；需 Code Review + 集成测试覆盖
- 超大规模单租户可能影响邻居（需后续分片策略）

### 后续行动
- [ ] 生产环境评估启用 RLS
- [ ] Phase 1 登录流程校验用户是否属于目标租户

## 关联文档

- [01-multi-tenancy.md](../01-multi-tenancy.md)
- [backend-arch/00-database-schema.md](../../backend-arch/00-database-schema.md)
- [tasks/00-mvp-task-breakdown.md](../../tasks/00-mvp-task-breakdown.md) Phase 0.4
