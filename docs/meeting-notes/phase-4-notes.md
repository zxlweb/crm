# Phase 4 阶段笔记 — 系统设置与收尾

**日期**：2026-05-26  
**状态**：BE 4.1–4.5 完成 → **BE 可联调**  
**PRD**：[phase-4-system-settings-close-prd.md](../prd/phase-4-system-settings-close-prd.md)  
**API**：[phase-4-system-settings-api.md](../api/phase-4-system-settings-api.md) v1.0 Accepted

---

## 1. BE 交付（4.1–4.5）

| ID | 内容 | 状态 |
|----|------|------|
| 4.1 | `GET/PATCH /api/settings/tenant`、`GET /api/settings/features`、自定义字段 CRUD | ✅ |
| 4.2 | `GET /api/audit/stats/*`、`GET /api/audit/export`（CSV + 限频） | ✅ |
| 4.3 | `GET /api/super-admin/stats/tenant-health` | ✅ |
| 4.4 | `plan-distribution`、`top-tenants` | ✅ |
| 4.5 | `backend/api/openapi-phase4.yaml`、`backend/README.md` Phase 4 节 | ✅ |

- 迁移：`00014_phase4_settings.sql`、`00015_seed_phase4_permissions.sql`
- 测试：`phase4_integration_test.go`、`tenant_settings_test.go`、RBAC 路由测
- 配置存储：`tenants.config` JSONB（locale/timezone/switches/quota）；`custom_fields` 表

---

## 2. BE 可联调

**FE / QA 可去掉 mock，对接真实 API。**

### 本地启动

1. PostgreSQL 已运行，`.env` 配置 `DB_*`、`JWT_SECRET`
2. 迁移：`cd backend && make migrate-up`（需包含 00014–00017）
3. 启动 API：`make run`（或 Windows `.\run.ps1`）→ `http://localhost:8080`
4. 健康检查：`curl http://localhost:8080/health`

### Demo 账号（PRD §2 Persona）

迁移：`00016_seed_phase4_persona_accounts.sql`（`make migrate-up`）。**统一密码**：`password123`  
**租户**：Demo Corp · `aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa`（登录后选该租户或 `switch-tenant`）

| Persona | 邮箱 | 角色 | 典型可验证能力 |
|---------|------|------|----------------|
| **Super Admin**（平台管理员） | `admin@demo.com` | 平台超管 + 租户管理员 | `/admin` 跨租户健康度/套餐/TOP；租户内 settings 全权限 |
| **Tenant Admin**（租户管理员） | `tenant-admin@demo.com` | Tenant Admin | `/settings` 改配置与自定义字段；`/settings/audit` 看图 + **导出** |
| **Sales Manager**（销售经理） | `manager@demo.com` | Sales Manager | `/settings` 只读；`/settings/audit` 看图，**无导出**；CRM 只读 |
| **Viewer**（只读） | `viewer@demo.com` | Viewer | 同经理侧只读；不可保存设置、不可管理字段 |
| **多角色切换** | `multi-role@demo.com` | Sales Manager + Viewer | 顶栏「角色」下拉切换；权限随当前角色变化（`POST /api/auth/switch-role`） |

迁移 `00017_seed_multi_role_demo_user.sql` 提供 `multi-role@demo.com`（密码同上）。

**场景对照（PRD）**

| 场景 | 建议账号 |
|------|----------|
| S1 租户配置保存 | `tenant-admin@demo.com` |
| S2 自定义字段 | `tenant-admin@demo.com` |
| S3 审计图（经理） | `manager@demo.com` |
| S4 Super Admin 运营图 | `admin@demo.com`（无需 `X-Tenant-ID` 调 super-admin API） |

**越权抽检**：`viewer@demo.com` 访问 `PATCH /api/settings/tenant` 或导出审计 → 应 `403`。

请求头：`Authorization: Bearer <access_token>`、`X-Tenant-ID`（租户接口）；Super Admin 统计接口仅需 Bearer。

### 建议验收顺序

1. `GET /api/settings/tenant` → `PATCH` 改 `default_locale` / `business_switches`
2. `POST /api/settings/custom-fields`（`lead` + `select`）→ 列表 → `DELETE` 逻辑删
3. `GET /api/audit/stats/by-action?from=2026-05-01&to=2026-05-26`
4. Super Admin：`GET /api/super-admin/stats/tenant-health`（无 `X-Tenant-ID`）

OpenAPI：`backend/api/openapi-phase4.yaml`。

---

## 3. 角色与权限（系统设置内）

- 入口：**系统设置** → Tab **「角色与权限」**（`/settings?tab=roles`）；旧路径 `/settings/roles` 自动跳转。
- 布局：左侧角色列表，右侧按模块（客户与商机 / 商业智能 / 系统设置 / …）树形勾选权限。
- 编辑：需 `rbac:manage`（如 `tenant-admin@demo.com`）；`manager@demo.com` / `viewer@demo.com` 仅查看或不可见该 Tab。
- **运行时角色切换**：顶栏「角色」下拉（仅当用户绑定 ≥2 角色时显示）；切换后 JWT `active_role_id` 更新，权限按当前角色重算。

---

## 4. 前端切面（2b，架构师）

见 PRD §10；FE 轨领 `/settings`、`/settings/audit`、`/admin` 图表嵌入。

---

## 4. FE 真实 API 联调（2026-05-26）

**环境**：`NUXT_PUBLIC_API_BASE=http://localhost:8080`；`NUXT_PUBLIC_USE_SETTINGS_MOCK=false`、`NUXT_PUBLIC_USE_ADMIN_INSIGHTS_MOCK=false`（默认，见 `apps/web/.env.example`）。

### 代码切换

| 模块 | 变更 |
|------|------|
| `use-settings` / `use-custom-fields` / `use-audit-stats` | 默认走真实 API；mock 仅 `useSettingsMock=true` |
| `use-admin-tenant-insights` | 默认走真实 API；mock 仅 `useAdminInsightsMock=true` |
| `utils/phase4-api-normalize.ts` | 兼容 `{ items }` 列表、审计 `Bucket`→`date`、Go 无 json tag 的 PascalCase 字段 |
| `nuxt.config` | 注册 `feature/settings` 组件目录；暴露 Phase 4 mock 开关 |
| 页面 | `/settings` → 审计页入口；Admin 健康度/套餐/TOP 接真实 super-admin 统计 |

### 自动化验收（角色切换）

| 类型 | 命令 | 说明 |
|------|------|------|
| API smoke | `cd apps/web && pnpm smoke:phase4-role` | 需 API `http://localhost:8080`、迁移 00017 |
| E2E | `cd e2e && pnpm test tests/phase4-role-switch.spec.ts` | Playwright 会复用已启动的 web/API（见 `playwright.config.ts`） |

账号：`multi-role@demo.com` / `password123`（Sales Manager + Viewer）。

### 联调检查项（本地）

1. 登录 `admin@demo.com` / `password123`，顶栏租户 `aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa`
2. `/settings`：加载租户配置 → PATCH 改语言/开关 → 自定义字段 CRUD
3. `/settings/audit`：三张图有数据；导出 CSV（需 `audit:export`）
4. `/admin`（Super Admin）：健康度雷达 + 套餐 Donut + TOP Bar

### 已知兼容

- BE 审计趋势字段为 `bucket`（非 `date`），FE 已映射。
- BE `ActionCountRow` 等 repository 行类型无 `json` tag 时可能返回 PascalCase，FE normalize 已兜底。
- 审计导出为 **纯 CSV 响应**（非 `{ code, data }` 包装），`exportCsv` 使用裸 `fetch` + `format=csv`。

---

## 5. 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-26 | BE 4.1–4.5 完成，标记 BE 可联调 |
| 2026-05-26 | 迁移 00016：PRD Persona 演示账号（tenant-admin / manager / viewer） |
| 2026-05-26 | FE：系统设置内「角色与权限」Tab，树形权限配置（`/settings?tab=roles`） |
| 2026-05-26 | FE 切换真实 API + normalize 层；联调说明记入本节 |
| 2026-05-26 | 多角色切换：`multi-role@demo.com`、顶栏角色下拉、`POST /api/auth/switch-role` |
