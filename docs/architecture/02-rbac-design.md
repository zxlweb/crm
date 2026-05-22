# 02 RBAC + ABAC 权限系统设计

**版本**：v0.1  
**日期**：2026-05-21

## 1. 权限模型

采用 **RBAC + ABAC 混合模型**：

- **RBAC**：用户 → 角色 → 权限（基础权限）
- **ABAC**：基于属性（tenant_id、owner_id、department 等）的动态权限

## 2. 核心实体关系

- **Tenant** ←→ **Role**（租户可自定义角色）
- **User** ←→ **Role**（支持多角色）
- **Role** ←→ **Permission**
- **Permission**：资源 + 操作（例如 `leads:create`, `reports:view_all`）

## 3. 权限资源设计（推荐）

权限格式统一为：`resource:action`

**示例**：
- `customers:view`
- `customers:create`
- `deals:delete`
- `reports:export`
- `settings:tenant_config`

**数据权限**（ABAC）：
- `leads:view_own` / `leads:view_department` / `leads:view_all`

## 4. Backend 实现方案

**推荐方案**：**Casbin + Go 适配器**

- 支持多租户（`tenant` 作为 domain）
- 支持 RBAC + ABAC + 自定义策略
- 策略缓存（Redis）

**备用方案**：自实现权限服务 + 缓存

**中间件执行顺序**：
1. Auth Middleware（JWT 解析）
2. Tenant Middleware
3. RBAC Middleware（Casbin Enforce）

**代码示例**（Gin）：
```go
if !enforcer.Enforce(userID, tenantID, "leads", "create") {
    return c.JSON(403, gin.H{"error": "无权限"})
}