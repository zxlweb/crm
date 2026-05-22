# CRM 系统架构文档

**版本**：v0.2（MVP 架构）  
**日期**：2026-05-21  
**技术栈**：Nuxt 3 + Go + PostgreSQL

> 浏览器阅读：[00-main-architecture.html](./00-main-architecture.html) · 决策记录：[adr/](./adr/README.md)

## 1. 系统上下文图 (C4 Context)

**系统名称**：Enterprise CRM (SaaS 多租户)

**外部用户/系统**：
- 企业管理员（Tenant Admin）
- 销售/客服/经理等普通用户
- 系统管理员（Super Admin）
- 外部邮件/SMS/支付/第三方登录系统
- 外部数据导入/导出（Excel、API）

**主要交互**：
- 用户通过浏览器访问 Nuxt 3 前端
- 前端通过 RESTful API + JWT 调用 Go Backend
- Backend 访问 PostgreSQL（多租户隔离）

## 2. 多租户实现策略（核心决策）

**推荐方案**：**共享数据库 + tenant_id 隔离**（Row-level Isolation）

**理由**：
- 开发与运维成本最低，适合初期到中型规模
- PostgreSQL Row Level Security (RLS) 可进一步加强隔离
- 易于后续扩展到 Schema 隔离或独立数据库

**实现方式**：
- 所有业务表必须包含 `tenant_id` 字段（UUID 或 BIGINT）
- Backend 中间件自动从 JWT 中提取 `tenant_id`，并注入到所有查询
- 使用 GORM Scope 或 sqlc 封装 tenant-aware 查询
- 关键表（Tenant、User、Role 等）全局共享，其他业务表严格隔离

## 3. RBAC 权限模型设计

**模型**：**RBAC + ABAC 混合**（推荐）

- **核心实体**：
  - Tenant（租户）
  - User（用户，可属于多个租户）
  - Role（角色，支持租户级 + 全局）
  - Permission（资源:操作，例如 `leads:create`, `reports:export`）
  - UserRole、RolePermission 关联表

**实现推荐**：
- **首选**：集成 [Casbin](https://casbin.org/)（支持 RBAC + ABAC + 多租户）
- 备选：自实现 + PostgreSQL 表查询 + 缓存（Redis）
- 数据权限：部分表增加 `owner_id` 或 `department_id` 做进一步过滤

**前端实现**：
- `usePermission()` composable
- `<Permission :required="['deals:delete']">` 组件或指令

## 4. 整体分层架构

### Backend (Go)
```
cmd/
  └─ api/                  # 入口
internal/
  ├── domain/              # 领域模型（Tenant, User, Lead, Deal...）
  ├── application/         # 用例层（Service）
  ├── infrastructure/      # 持久化、外部服务
  ├── interfaces/          # API Handlers + DTO
  ├── middleware/          # Tenant, Auth, RBAC, Logging
  └── config/
```

- Web 框架：**Gin**
- ORM：**GORM**（MVP 阶段快速开发）或 **sqlc + pgx**（生产推荐）
- 迁移工具：**goose**
- 验证：**validator.v10**
- API 文档：**Swagger**

### Frontend (Nuxt 3)
```
layers/                  # 可选分层
src/
  ├── app/               # Nuxt App Router
  ├── components/        # base、ui、feature
  ├── composables/       # useTenant, usePermission, useApi 等
  ├── features/          # 按业务模块（leads、deals、contacts...）
  ├── server/            # Nitro 服务端 API（可选）
  ├── stores/            # Pinia
  └── utils/
```

- 状态管理：**Pinia**
- API 请求：**useFetch / $fetch** + Zod 解析响应
- 表单：**Zod** + Nuxt UI / vee-validate
- UI：**Tailwind CSS + Nuxt UI**（或 Headless UI）

## 5. 认证授权方案

- **认证**：JWT + Refresh Token（HttpOnly Cookie）
- **流程**：
  1. 登录时返回 access_token + refresh_token
  2. 中间件解析 token → 加载 User + Tenant + Roles
  3. 支持多租户切换（用户可在有权限的租户间切换）
- **Super Admin** 特殊通道（独立权限）

## 6. 数据库核心表设计原则

**必须包含的字段**（每张业务表）：
- `id`, `tenant_id`, `created_at`, `updated_at`, `created_by`, `updated_by`
- 软删除：`deleted_at`

**核心共享表**：
- `tenants`
- `users`
- `user_tenants`（用户-租户关联）
- `roles`
- `permissions`
- `role_permissions`
- `audit_logs`（审计）

## 7. 前端架构分层建议

- **composables/**：全局可复用（useTenant, useRBAC, useI18n）
- **features/**：业务模块组织（推荐）
- **i18n**：使用 `nuxt/i18n`，支持 `zh-CN` 和 `en`，动态加载
- **类型安全**：所有 API Response 定义 Zod Schema，前端自动推断类型

## 8. 关键架构决策 (ADR)

详细决策记录见 **[architecture/adr/](./adr/README.md)**（不可改写历史，仅可 Superseded）。

| ADR | 决策 | 状态 |
|-----|------|------|
| [0001](./adr/0001-shared-db-multi-tenancy.md) | 共享数据库 + tenant_id 多租户 | Accepted |
| [0002](./adr/0002-casbin-rbac.md) | Casbin RBAC + resource:action | Accepted |
| [0003](./adr/0003-goose-migrations.md) | Goose 数据库迁移 | Accepted |

*技术栈（Gin、Nuxt 3、Pinia）见 `.cursorrules`，后续若需独立 ADR 按 `0000-template.md` 补充。*

## 9. 潜在风险与缓解

- **风险**：多租户数据泄露 → **缓解**：中间件强制过滤 + RLS + 单元测试
- **风险**：RBAC 权限爆炸 → **缓解**：设计权限资源树 + 管理后台
- **风险**：i18n 文案不一致 → **缓解**：建立文案审核流程 + 自动化检查
- **风险**：性能（高并发租户）→ **缓解**：缓存（Redis）、读写分离、索引优化

---
