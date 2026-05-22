# QA 测试计划 - Phase 0 基础架构

**测试日期**：2026-05-22  
**测试人**：AI + 开发者  
**结论**：通过（自动化测试 + 冒烟验证）

## 1. 测试范围

- 基础架构脚手架
- 数据库迁移
- 健康检查 API
- 多租户 Context / GORM Scope
- RBAC 路由权限映射
- 前端首页与 composables 逻辑
- E2E 冒烟（health + 首页）

## 2. 自动化测试

| 类型 | 命令 | 结果 |
|------|------|------|
| Backend Unit/Integration | `cd backend && make test` | ✅ 通过 |
| Frontend Vitest | `cd frontend && npm run test` | ✅ 通过（2 tests） |
| E2E Playwright | `make e2e-test` | 需本地 DB + 前后端 |

### 用例覆盖

**Backend**
- [x] `tenant` Context 读写与 panic
- [x] `rbac.RouteToPermission` 路径映射
- [x] `TenantScope` 租户隔离（SQLite 内存）
- [x] `GET /health` 成功与 DB 不可用

**Frontend**
- [x] `toPermissionMap` / `canAccess`

**E2E**
- [x] `GET /health` 返回 `status: ok`
- [x] 首页展示 CRM 标题

## 3. 手动测试

| # | 场景 | 结果 |
|---|------|------|
| 1 | `docker compose up -d` 启动 Postgres | 通过 |
| 2 | `make migrate-up` 迁移到 v3 | 通过 |
| 3 | `curl localhost:8080/health` | 通过 |
| 4 | 浏览器打开 `localhost:3000` | 通过 |

## 4. 多租户 / RBAC（Phase 0 范围）

| # | 场景 | 结果 |
|---|------|------|
| 1 | `X-Tenant-ID` 中间件写入 Context | 代码审查通过，有单元测试 |
| 2 | GORM `DBFromContext` 过滤 `tenant_id` | 集成测试通过 |
| 3 | Casbin 从 DB 同步策略 | 启动日志确认，端到端留 Phase 1 |

## 5. Bug 记录

| 序号 | 描述 | 严重程度 | 状态 |
|------|------|----------|------|
| — | 无阻塞缺陷 | — | — |

## 6. 测试结论

**通过** — Phase 0 可关闭，进入 Phase 1（认证与权限）。
