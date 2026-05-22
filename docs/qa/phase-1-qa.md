# QA 测试计划 - Phase 1 认证与权限

**测试日期**：2026-05-22  
**测试人**：AI + 开发者  
**关联笔记**：[phase-1-notes.md](../meeting-notes/phase-1-notes.md)

## 1. 测试范围

- 登录 / Refresh Token
- **自助注册**（新租户 + Tenant Admin）→ 进入应用
- 多租户列表与切换
- Super Admin `/admin` 入口
- RBAC API 与前端 `my-permissions` / `PermissionGuard`
- 审计日志（注册/登录写 `audit_logs`）
- 中英文登录页文案（抽测）

**不在本次范围**：忘记密码、邮件验证、Casbin 全路由 Enforce 根因修复（见 Bug #3 技术债）

## 2. 与 Phase 0 / 任务清单的差异说明

| 来源 | 写了什么 | 实际验收了什么 |
|------|----------|----------------|
| [phase-0-qa.md](./phase-0-qa.md) | Phase 0 **通过**，可进入 Phase 1 | 未测注册；Casbin 端到端**留 Phase 1** |
| [00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md) | Phase 1 多项 `[x]` | 多为实现完成勾选，**此前无本文档** |
| [phase-1-notes.md](../meeting-notes/phase-1-notes.md) | 1.1–1.5 完成 | 本地验证以 **demo Super Admin 登录** 为主 |

**结论**：此前「Phase 1 过了」易误解为注册 E2E 已过；本次补测聚焦 **新用户注册链路**。

## 3. 自动化测试

| 类型 | 命令 | 结果 |
|------|------|------|
| Backend | `cd backend && make test` | 沿用 Phase 0（无 Register 成功路径集成测） |
| Frontend Vitest | `pnpm --filter @crm/web test` | 未单独回归 Phase 1 |
| E2E Playwright | `cd e2e && npm test -- tests/phase1-auth.spec.ts` | 见 §4（需本地 DB + 前后端） |

## 4. 核心功能测试（手动 + E2E）

| # | 场景 | 步骤 | 预期 | 结果 |
|---|------|------|------|------|
| 1 | Demo 登录 | `/login`，`admin@demo.com` / `password123` | 跳转 `/admin` | ✅ |
| 2 | Super Admin 概览 | `/admin` | 卡片 + 租户列表 + 图表 | ✅（笔记已验） |
| 3 | **新用户注册** | `/login?mode=register`，填名称/公司/邮箱/密码，域名留空 | `POST /register` 201，跳转 `/` | ✅（修复后） |
| 4 | 注册后权限 | 注册完成瞬间 | `GET /my-permissions` 不阻断进首页 | ✅（中间件白名单 + 前端容错） |
| 5 | 注册入口 UI | 登录页底部「注册」 | 可切换注册表单 | ✅（`LoginAuthModeSwitch` 已移除；底部链接） |
| 6 | Refresh | `POST /api/auth/refresh` | 新 access_token | ⬜ 未回归 |
| 7 | 租户切换 | `POST /switch-tenant` | `current_tenant` 更新 | ⬜ 未回归 |
| 8 | 角色管理页 | `/settings/roles`（Tenant Admin） | 列表可访问 | ⬜ 依赖 Casbin Enforce，可能 403 |

**E2E**：`e2e/tests/phase1-auth.spec.ts` — 用例「register new tenant admin lands on home」覆盖 #3。

## 5. 多租户 / RBAC

| # | 场景 | 结果 | 说明 |
|---|------|------|------|
| 1 | 新注册租户数据隔离 | ✅ | 每注册创建独立 `tenants` + `user_tenants` |
| 2 | Tenant Admin 拥有 19 条 `role_permissions` | ✅ | DB 可查 |
| 3 | 新用户 `GET /rbac/my-permissions`（修复前） | ❌ → ✅ | 修复前 Casbin 路由 Enforce 403 |
| 4 | 新用户 `GET /rbac/roles` 等写操作 | ⬜ | 仍走 Casbin，**可能 403**（技术债） |
| 5 | Super Admin 绕过 RBAC | ✅ | `is_super_admin` 跳过中间件 |

## 6. Bug 记录

| 序号 | 描述 | 严重程度 | 状态 |
|------|------|----------|------|
| 1 | 注册时 `domain` 留空，前端 Zod 将 `undefined` 当非法，`safeParse` 失败 | 高 | ✅ 已修 `apps/web/schemas/register.ts`（`.optional()`） |
| 2 | 注册 API 201 后 `my-permissions` 403，页面显示「无权限访问该资源」且不跳转 | 高 | ✅ 已修：`rbac.go` 对 `GET .../my-permissions` 放行；`login.vue` 加载失败不阻断 |
| 3 | 新租户管理员 Casbin `Enforce` 对一般 RBAC 路由仍失败（DB 有策略） | 中 | ⬜ 技术债；需集成测 + 核对 `SyncCasbinPolicies` / matcher |
| 4 | 登录页误用 `AuthModeSwitch`，Nuxt 实际组件名为 `LoginAuthModeSwitch`，注册入口不显示 | 中 | ✅ 已改为底部链接切换 |
| 5 | 无 `phase-1-qa.md`、无注册 E2E，任务清单误标「注册已完成」 | 低 | ✅ 本文档 + `phase1-auth.spec.ts` |

## 7. 测试结论

**有条件通过** — Phase 1 可进入 Phase 2，但须知晓：

1. **注册 → 首页** 主路径已通（含本次修复）。  
2. **Casbin 全链路** 对新租户管理员仍未完全验证；`my-permissions` 为**临时白名单**，非根因修复。  
3. 建议在 Phase 2 并行轨 `[QA]` 补：Register 集成测、Casbin Enforce 集成测、`/settings/roles` 租户 Admin 冒烟。

**Phase 1 正式关闭条件（建议）**：

- [x] 本文档落档  
- [x] E2E `phase1-auth.spec.ts` 纳入 CI / 本地 `make e2e-test`  
- [ ] Bug #3 根因修复或 ADR 记录「白名单策略」  
- [ ] `phase-1-review.md`（Code Review 留档，可选）
