# API 设计文档（MVP）

**版本**：v0.1  
**日期**：2026-05-21  
**Base URL**：`http://localhost:8080`

> HTML：[00-api-design.html](./00-api-design.html)

---

## 统一响应格式

所有接口返回：

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "pagination": { "page": 1, "page_size": 20, "total": 100 }
}
```

| HTTP | code | 场景 |
|------|------|------|
| 200 | 200 | 成功 |
| 201 | 201 | 创建成功 |
| 400 | 400 | 参数错误 |
| 401 | 401 | 未认证 |
| 403 | 403 | 无权限 |
| 404 | 404 | 资源不存在 |
| 500 | 500 | 服务器错误 |

---

## 通用请求头

| Header | 必填 | 说明 |
|--------|------|------|
| `Authorization` | 受保护接口 | `Bearer <access_token>` |
| `X-Tenant-ID` | 受保护业务接口 | 当前租户 UUID |
| `Content-Type` | POST/PUT/PATCH | `application/json` |

---

## 认证与租户

| 方法 | 接口 | 描述 | 认证 |
|------|------|------|------|
| GET | `/health` | 健康检查 | 否 |
| POST | `/api/auth/login` | 登录 | 否 |
| POST | `/api/auth/register` | 自助注册（创建租户 + 管理员） | 否 |
| POST | `/api/auth/refresh` | 刷新 Token | 否 |
| GET | `/api/auth/tenants` | 可访问租户列表 | 是 |
| POST | `/api/auth/switch-tenant` | 切换租户 | 是 |
| GET | `/api/auth/profile` | 当前用户（无需租户） | 是 |
| GET | `/api/auth/me` | 当前用户 + 租户上下文 | 是 |

### POST /api/auth/login

**请求**
```json
{ "email": "user@example.com", "password": "secret123" }
```

**响应 data**
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "expires_in": 3600,
  "user": { "id": "uuid", "email": "user@example.com", "name": "Admin", "is_super_admin": false },
  "tenants": [{ "id": "uuid", "name": "Acme", "domain": "acme" }],
  "current_tenant": null
}
```

登录后 `current_tenant` 为空；调用切换租户后返回当前租户，且 access/refresh Token 的 JWT 载荷含 `tenant_id`。

**Demo 账号**（迁移 `00003_seed_dev` 后）：`admin@demo.com` / `password123`

### POST /api/auth/register

**请求**
```json
{
  "email": "owner@acme.com",
  "password": "secret123",
  "name": "Jane Owner",
  "company_name": "Acme Corp",
  "domain": "acme"
}
```

`domain` 可选；留空则根据 `company_name` 生成小写 slug。冲突时返回 409（邮箱或域名已存在）。

**响应 data**：与 login 相同（HTTP 201），`current_tenant` 为新创建的租户。

### GET /api/auth/tenants

仅需 `Authorization: Bearer <access_token>`，**不需要** `X-Tenant-ID`。

**响应 data**：`[{ "id": "uuid", "name": "Acme", "domain": "acme" }]`

### POST /api/auth/switch-tenant

仅需 Bearer Token，**不需要** `X-Tenant-ID`。

**请求**
```json
{ "tenant_id": "uuid" }
```

**响应 data**：与 login 相同；`current_tenant` 为切换后的租户。

### POST /api/auth/refresh

**请求**
```json
{ "refresh_token": "<refresh_token>" }
```

**响应 data**：与 login 相同结构（含新 access_token、refresh_token、tenants）

### GET /api/auth/profile

仅需 Bearer Token，**不需要** `X-Tenant-ID`。

**响应 data**：`{ "id", "email", "name", "is_super_admin" }`

---

## Super Admin（跨租户）

独立路由组 `/api/super-admin`，仅需 Bearer Token + `is_super_admin=true`，**不需要** `X-Tenant-ID`，不经过 Tenant / RBAC 中间件。

| 方法 | 接口 | 描述 |
|------|------|------|
| GET | `/api/super-admin/overview` | 平台概览（租户数、用户数） |
| GET | `/api/super-admin/stats/tenant-activity` | 租户活跃趋势（近 N 日登录/新增租户） |
| GET | `/api/super-admin/tenants` | 租户分页列表 |
| GET | `/api/super-admin/tenants/:id` | 租户详情 |
| PATCH | `/api/super-admin/tenants/:id` | 启用/停用租户 `{ "is_active": true }` |

列表 Query：`page`、`page_size`、`search`、`is_active`。

---

## RBAC 权限

| 方法 | 接口 | 描述 |
|------|------|------|
| GET | `/api/rbac/permissions` | 权限字典（按 resource 分组） |
| GET | `/api/rbac/permission-items` | 权限平铺列表（含 id，用于勾选） |
| GET | `/api/rbac/my-permissions` | 当前用户在本租户的权限 |
| GET | `/api/rbac/roles` | 角色列表 |
| POST | `/api/rbac/roles` | 创建角色 |
| PUT | `/api/rbac/roles/:id` | 更新角色 |
| POST | `/api/rbac/roles/:id/permissions` | 分配权限 `{ "permission_ids": ["uuid"] }` |
| GET | `/api/rbac/users/:id/roles` | 用户角色 |
| POST | `/api/rbac/users/:id/roles` | 分配角色 `{ "role_ids": ["uuid"] }` |
| POST | `/api/rbac/check` | 权限检查 `{ "resource", "action" }` → `{ "allowed": true }` |

权限标识：`resource:action`（如 `leads:view`）。变更角色/用户角色后自动同步 Casbin。

---

## 客户与线索（Phase 2）

> **详细契约**（字段、洞察、情绪旅程、统计、AI Preview）：[phase-2-crm-ai.md](./phase-2-crm-ai.md)

### Accounts（公司）

- `GET|POST /api/accounts`、`GET|PUT|PATCH|DELETE /api/accounts/:id`
- `GET /api/accounts/:id/emotion-journey`
- `POST /api/accounts/:id/insights/evaluate`

### Contacts（联系人）

- `GET|POST /api/contacts`、`GET|PUT|PATCH|DELETE /api/contacts/:id`
- `GET /api/accounts/:id/contacts`
- `GET /api/contacts/:id/emotion-journey`

### Leads（线索）

- `GET|POST /api/leads`、`GET|PUT|PATCH|DELETE /api/leads/:id`
- `POST /api/leads/:id/convert`、`POST /api/leads/:id/assign`
- `POST /api/leads/import`
- `GET /api/leads/stats/*`（来源、趋势、漏斗、状态、健康度）
- `GET /api/leads/:id/emotion-journey`

### Activities / Insights / Segments / Copilot

见 [phase-2-crm-ai.md](./phase-2-crm-ai.md) §6–§11

### Deals + Dashboard（商机与仪表盘，Phase 3）

**详约**：[phase-3-deals-dashboard-api.md](./phase-3-deals-dashboard-api.md)（v1.0 Accepted）

| 模块 | 代表路径 |
|------|----------|
| Deals CRUD | `GET/POST /api/deals`、`GET/PUT/PATCH/DELETE /api/deals/:id` |
| Pipeline | `GET /api/deals/pipeline`、`PUT /api/deals/:id/stage` |
| Deals 统计 | `GET /api/deals/stats/by-stage`、`GET /api/deals/stats/win-rate` |
| Dashboard | `GET /api/dashboard/summary`、`/funnel`、`/todo`、`/quota`、`/team-ranking` |

---

## 设计原则

- 业务接口由中间件注入租户，**前端不传 `tenant_id`  body 字段**
- 分页：`page`、`page_size` Query 参数
- 过滤：Query 参数，命名与资源字段一致
- 写操作记录 `audit_logs`：登录/注册/切换租户、RBAC 变更、Super Admin 租户启停（`action` 如 `auth.login`、`rbac.role.create`）
