# Phase 1 笔记

**时间**：2026-05-22  
**当前**：1.1–1.3 已完成

## 1.1 登录 / Refresh（已完成）

- `POST /api/auth/login`、`POST /api/auth/refresh`
- 响应含 `user` 对象（含 `is_super_admin`）

## 1.2 多租户列表与切换（已完成）

- `GET /api/auth/tenants`、`POST /api/auth/switch-tenant`
- JWT `tenant_id` + `TenantMiddleware` 回退

## 1.3 Super Admin 管理后台入口（已完成）

### Backend

- 路由组 `/api/super-admin`：`Auth` + `SuperAdminMiddleware`（无 Tenant/RBAC）
- `GET /overview`、`GET /tenants`、`GET /tenants/:id`、`PATCH /tenants/:id`
- `GET /api/auth/profile`（authOnly，无需租户）
- 迁移 `00004`：第二演示租户 `acme`

### Frontend

- `/login` 登录页（Super Admin 登录后跳转 `/admin`）
- `/admin` 管理后台（概览卡片 + 租户列表 + 启停）
- 顶栏「管理后台」入口（仅 `is_super_admin` 可见）
- `middleware/super-admin.ts` 路由守卫

### 本地验证

```bash
# 登录（Super Admin）
curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@demo.com","password":"password123"}' | jq

# 平台概览（用返回的 access_token）
curl -s http://localhost:8080/api/super-admin/overview \
  -H "Authorization: Bearer <token>" | jq

# 租户列表
curl -s "http://localhost:8080/api/super-admin/tenants?page=1&page_size=20" \
  -H "Authorization: Bearer <token>" | jq
```

前端：访问 `http://localhost:3000/login`，用 demo 账号登录后进入 `/admin`。

## 1.4 RBAC 完整 API（已完成）

- `GET/POST/PUT /api/rbac/roles`、`POST .../permissions`
- `GET/POST /api/rbac/users/:id/roles`
- `GET /api/rbac/my-permissions`、`POST /api/rbac/check`
- `GET /api/rbac/permission-items`（含 UUID）
- 迁移 `00005`：Demo 租户 `Tenant Admin` 角色 + 全权限

## 1.5 前端权限组件（已完成）

- `useRbac`、`PermissionGuard`、Super Admin 绕过
- `/settings/roles` 角色管理页
- 登录后加载 `my-permissions`

### 下一步（1.6）

- 审计日志基础记录
