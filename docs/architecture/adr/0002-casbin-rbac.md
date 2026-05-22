# ADR-0002: Casbin RBAC + resource:action 权限模型

**状态**：Accepted  
**日期**：2026-05-21  
**决策者**：Software Architect

## 背景（Context）

系统需要租户内角色权限、支持用户多角色，且权限需可配置、可审计。MVP 要先打通鉴权链路，再逐步引入 ABAC（数据范围）。

## 决策（Decision）

采用 **Casbin** 作为权限引擎，权限标识统一为 **`resource:action`**（如 `leads:view`）。

策略来源：
- 全局表 `permissions` 定义资源与操作
- 租户表 `roles` + `role_permissions` + `user_roles` 关联
- 服务启动时 `SyncCasbinPolicies` 从 DB 同步到 Casbin 内存

HTTP 中间件顺序：`Auth` → `Tenant` → `RBAC`；`is_super_admin` 跳过 RBAC。

## 备选方案（Options Considered）

| 方案 | 优点 | 缺点 |
|------|------|------|
| Casbin ✅ | 成熟、支持 dom（租户）、可扩展 ABAC | 需同步策略、学习模型语法 |
| 自研权限表 + 硬编码 | 简单直观 | 策略复杂后难维护 |
| OPA | 策略能力强 | MVP 过重 |

## 后果（Consequences）

### 正面
- 与文档 [02-rbac-design.md](../02-rbac-design.md) 一致
- 前端 `can(resource, action)` 可直接对齐

### 负面 / 风险
- 策略变更需重启或实现热加载（当前为启动同步）
- 路由到 `resource:action` 的映射需维护（见 `pkg/rbac/route.go`）

### 后续行动
- [ ] Phase 1 实现角色/权限分配 API
- [ ] 评估 Casbin GORM Adapter 持久化策略
- [ ] ABAC（view_own / view_all）单独 ADR

## 关联文档

- [02-rbac-design.md](../02-rbac-design.md)
- [api/00-api-design.md](../../api/00-api-design.md) RBAC 章节
