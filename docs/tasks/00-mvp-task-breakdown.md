# MVP 任务拆分清单

> **Card* / Chart* 与业务同迭代交付**：每个 Phase 同时交付后端 API + 业务页面 + ui-kit 组件 + 页面嵌入，避免「模块做完再补图」。组件场景见 [05-component-scenarios.md](../frontend-arch/05-component-scenarios.md)；双包与 ui-kit 见 [02-ui-kit-modules.md](../frontend-arch/02-ui-kit-modules.md)。

> **三轨并行（3 个 AI 对话框）**：Implementation 阶段可按 `[BE]` / `[FE]` / `[QA]` 分轨领任务，见 [parallel-implementation.md](./parallel-implementation.md)。开工前需 API 契约 + 前端切面；各轨只改约定目录，用 `⬜` / `🟡` / `✅` 标记状态。

---

## Chart / Card 同步原则

| 原则 | 说明 |
|------|------|
| **同迭代交付** | 模块 MVP 验收必须包含该 Phase 承诺的图表（见各 Phase「图表 Done」） |
| **先 API 后图** | 图表数据来自模块统计 API；开发期可短期 mock，禁止长期占坑 |
| **组件先行 ≤ 3 天** | 缺 ui-kit 组件时，在模块开工当周完成组件 + Story，再嵌页面 |
| **权限一致** | 图表查询与列表页相同：`tenant_id`、数据范围（本人/部门/全部） |
| **主题一致** | `useUiKitTheme()` / `ChartShell`，V1/V2 随设计系统 |

### 单模块图表 Done（验收勾选）

- [ ] 统计 API 有文档（路径、维度、筛选、权限）
- [ ] `@crm/ui-kit` 组件已实现（或复用）+ Vitest/Story 最小覆盖
- [ ] `apps/web` 对应 feature 页嵌入图表，空态/加载/错误态齐全
- [ ] 经理/个人视图数据范围正确（手动或 E2E 抽一条）

### 推荐双周节奏（示例）

```
Week 1  Mon–Tue  API 契约 + 列表/表单页
        Wed–Thu  补齐 ui-kit 组件 + Story
        Fri      统计 API + 联调
Week 2  Mon–Tue  图表嵌入 + 空态/权限
        Wed      单测 + 本清单勾选
        Thu–Fri  QA + 修复
```

**禁止**：Phase 3 做完 Dashboard 表格后，把图表整体推到后续「纯 UI Phase」。

### 文件落点

| 层级 | 路径 |
|------|------|
| 图表组件 | `packages/ui-kit/src/components/ui/chart/chart-*.vue` |
| 图表工具 | `packages/ui-kit/src/chart/utils/` |
| 统计 API 封装 | `apps/web/composables/use-*-stats.ts` |
| 业务嵌入 | `apps/web/components/feature/<module>/`、对应 `pages/` |
| 展示/回归 | `/charts` 新组件区块；Storybook（Phase 4 起） |

交付后：在本文件勾选业务与图表项，并在 [05-component-scenarios.md](../frontend-arch/05-component-scenarios.md) 总览表将 ui-kit 状态改为 ✅。

---

## Phase 0：基础架构（已完成 · v2.1 收尾，无 CI）

**业务**

- [x] 完善数据库核心 Schema + 迁移脚本（00001–00003，含 dev seed）
- [x] Backend 项目初始化（Gin + GORM + Casbin）
- [x] Frontend 项目初始化（Nuxt 3 + TypeScript + Tailwind + Pinia + i18n）
- [x] 多租户中间件 + Tenant Context（含单元/集成测试）
- [x] RBAC 基础框架（Casbin 集成 + 路由权限映射测试）
- [x] 自动化测试：Backend `make test`、Frontend Vitest、E2E 冒烟（`e2e/`）
- [x] QA + Code Review 留档 → [qa/phase-0-qa.md](../qa/phase-0-qa.md) · [reviews/phase-0-review.md](../reviews/phase-0-review.md) · [meeting-notes/phase-0-notes.md](../meeting-notes/phase-0-notes.md)

**图表 / 设计系统（基线，已完成）**

- [x] 双包组件模型（`02-ui-kit-modules`：ui-kit + web）
- [x] `ChartLine` / `ChartBar` / `ChartFunnel` / `ChartShell` + `/charts` 展示页
- [x] pnpm monorepo：`packages/ui-kit` + `apps/web`
- [x] 主题 bridge（`00-ui-kit-theme.ts` + `useUiKitTheme`）+ ui-kit `vite build`

| 模块 | 图表交付 | 状态 |
|------|----------|------|
| 设计系统 | 上述 Chart* + `CardMetric` / `CardShell` | ✅ |
| 主题 | `bridgeUiKitThemeFromApp`、`/charts` | ✅ |
| 业务嵌入 | 无要求 | — |

---

## Phase 1：认证与权限（有条件关闭 · 见 [phase-1-qa.md](../qa/phase-1-qa.md)）

**业务**

- [x] 用户登录 / Refresh Token
- [x] 用户注册（API + `/login` UI；E2E：`e2e/tests/phase1-auth.spec.ts`）
- [x] 多租户登录与切换
- [x] Super Admin 管理后台入口
- [x] RBAC 完整实现（角色、权限、分配；Casbin 新租户 Enforce ⬜ 见 QA Bug #3）
- [x] 前端权限控制组件（usePermission、PermissionGuard）
- [x] 审计日志基础记录

**QA / E2E**

- [x] [phase-1-qa.md](../qa/phase-1-qa.md) 留档（2026-05-22，注册链路补测）
- [ ] Casbin 新租户管理员全路由 Enforce 集成测
- [ ] Code Review → [reviews/phase-1-review.md](../reviews/phase-1-review.md)（可选）

**图表（与 Phase 1 同步）**

- [x] Admin 概览：租户活跃趋势 `ChartLine`（`/admin`，`GET /api/super-admin/stats/tenant-activity`）
- [x] 删除遗留目录 `frontend/`（已移除，仅用 `apps/web`）

| 模块 | 业务 | 同步图表 | 页面 |
|------|------|----------|------|
| 登录/租户 | 注册、审计 | 无统计图 | `/login`（表单内切换注册） |
| Super Admin | 后台框架 | `ChartLine` 租户趋势 | `/admin` |
| RBAC | 权限组件 | 无 | — |

**Phase 1 图表 Done**：Admin 首页 1 张 `ChartLine`（租户趋势）接真实或约定 API 契约。

---

## Phase 2：客户与线索 + 关系经营 / AI Preview

**PRD**：[phase-2-relationship-crm-prd.md](../prd/phase-2-relationship-crm-prd.md)（含 §15 老板演示、§4.6 情绪旅程）· **MVP 总览**：[00-crm-overview.md](../prd/00-crm-overview.md)

### Phase 2 并行任务（Implementation）

> 三轨同时推进时：各对话框首条消息声明 `【BE】` / `【FE】` / `【QA】`，只领本轨 `⬜` 行，完成后改 `✅`。细则见 [parallel-implementation.md](./parallel-implementation.md)。  
> **2.13**：`[BE]` 由架构师领（扩写 `docs/api/`，输入 PRD §15.3）；`[FE]` 可用 fixtures 与 2.3 并行，勿等真 LLM。

| ID | [BE] | [FE] | [QA] | 依赖 |
|----|------|------|------|------|
| 2.1 | Accounts API + 迁移 `00006` + `lifecycle_stage` — ✅ | `/accounts` 列表/表单 + `use-accounts` — ⬜ | Accounts HTTP 集成测（CRUD/筛选/租户/审计/只读字段）— ✅ | 契约 |
| 2.2 | Contacts API + 关联 — ⬜ | `/contacts` + `use-contacts` — ⬜ | Contacts 集成测 — ⬜ | 2.1 契约可参考 |
| 2.3 | Leads CRUD + 状态机 + `POST .../convert` — ✅ | Leads API 接入 + Nuxt UI + `layout: app` — ✅ | Leads HTTP 集成测（租户隔离 + 状态机 + convert 审计）— ✅ | 契约 |
| 2.4 | Leads 统计 API（来源） — ⬜ | `ChartDonut` + 嵌入 — ⬜ | 统计 API + 权限集成测 — ⬜ | 2.3 |
| 2.5 | Leads 趋势 API — ⬜ | `ChartLine` 嵌入 — ⬜ | 同上 — ⬜ | 2.4 |
| 2.6 | Leads 转化 API — ⬜ | `ChartFunnel` 嵌入 — ⬜ | 同上 — ⬜ | 2.5 |
| 2.7 | 导入、分配 API — ⬜ | 导入/分配 UI — ⬜ | 导入 E2E（可选）— ⬜ | 2.3 |
| 2.8 | Activity API + `event_type`/`direction` — ⬜ | ActivityTimeline + 摘要 + 可选 `ChartBar` — ⬜ | Activity 集成测 — ⬜ | 2.3 |
| 2.9 | `insights/evaluate` + engagement 字段 — ⬜ | 洞察侧栏 + `LifecycleBadge` — ⬜ | 规则引擎单测 + 洞察集成测 — ⬜ | 2.8 |
| 2.10 | segments 模板 API + 列表筛选 — ⬜ | 分群下拉 + URL 筛选 — ⬜ | 分群 count + 权限 — ⬜ | 2.3 |
| 2.11 | Activity `sentiment` 迁移 + 关键词规则 — ⬜ | Activity 情绪表单项 — ⬜ | 情绪字段 + 时间线 E2E — ⬜ | 2.8 |
| 2.12 | `emotion-journey` 聚合 API — ⬜ | `EmotionJourneyMap` Tab — ⬜ | 地图空态/触点联动 E2E — ⬜ | 2.11 |
| 2.13 | **架构师**：`docs/api/phase-2-crm-ai.md` + ADR-0004 — ✅ | `AiRelationPanel` + Copilot Mock + Preview 角标 + fixtures — ⬜ | §15.2 演示路径 E2E — ⬜ | 2.3 路由；契约已评审 |
| 2.E2E | — | Leads/详情 `data-testid` + Preview 路径 — ⬜ | CRUD + 权限 + Activity/洞察/Preview 冒烟 — ⬜ | 2.3 FE；2.13 演示 |

**业务**（勾选可与上表同步，或作汇总验收）

- [ ] 公司 (Accounts) CRUD + 搜索 + 生命周期阶段（BE ✅ / QA ✅；FE ⬜）
- [ ] 联系人 (Contacts) CRUD + 关联公司
- [x] 线索 (Leads) CRUD + 状态流转 + 转化（`POST /api/leads/:id/convert`）
- [ ] 线索导入（Excel）
- [ ] 线索分配功能
- [ ] 跟进记录（Activity）基础 + 时间线 + 情绪标注
- [ ] 规则洞察（≥2 条 INS）+ 预置分群（≥5）
- [ ] 情绪旅程地图 Tab
- [ ] **AI Preview**：详情 AI 侧栏 + Copilot Mock + 演示角标（§15）

**图表（缺组件则本迭代内补齐）**

- [ ] Leads 统计 API（来源、趋势、状态、转化）
- [ ] **`ChartDonut`** 组件 + Story（线索来源占比）
- [ ] Leads 报表区：`ChartLine` + `ChartDonut` + `ChartFunnel` + `ChartBar`（状态）
- [ ] Activity 摘要：`ChartBar` 跟进类型 TOP（可选，与 Activity 同迭代）

| 顺序 | 业务 | ui-kit | 嵌入页面 | 指标 |
|------|------|--------|----------|------|
| 2.1 | Accounts CRUD | — | `/accounts` | — |
| 2.2 | Contacts CRUD | — | `/contacts` | — |
| 2.3 | Leads CRUD + 状态 | `ChartBar` | Leads 报表 Tab | 状态分布 |
| 2.4 | Leads 统计 API | **`ChartDonut`**（新增） | Leads 报表 | 来源占比 |
| 2.5 | Leads 趋势 API | `ChartLine` | 同上 | 日/周新增 |
| 2.6 | Leads 转化 API | `ChartFunnel` | 同上 | 线索→合格→商机 |
| 2.7 | 导入、分配 | — | — | — |
| 2.8 | Activity 基础 | `ChartBar`（可选） | Activity 摘要 / 时间线 | 跟进 TOP |
| 2.9 | 规则洞察 | — | 详情洞察侧栏 | INS 命中 |
| 2.10 | 分群模板 | — | 列表分群筛选 | 高潜/流失等 |
| 2.11 | 情绪字段 | — | Activity 表单 | sentiment |
| 2.12 | 情绪旅程 API | `ChartLine`（曲线） | 情绪旅程 Tab | 情绪×时间 |
| 2.13 | AI 契约（架构师） | — | `AiRelationPanel` 等 | Preview 演示 |

**本阶段新增组件**：`ChartDonut`（2.4）；feature：`EmotionJourneyMap`、`AiRelationPanel`、`ActivityTimeline`（2.8–2.13）。  
**Phase 2 图表 Done**：Leads 报表区至少 **折线 + 环形 + 漏斗** 三张接 API。  
**Phase 2 演示 Done**：§15.2 三分钟路径 E2E 或彩排清单通过；Preview 数据不进生产表。

---

## Phase 3：商机与仪表盘

**业务**

- [ ] 商机 (Deals) CRUD + Pipeline 看板
- [ ] 销售阶段管理
- [ ] 仪表盘页面（个人 / 经理视图）

**图表（Dashboard 为必验收）**

- [ ] Deals 管道 `ChartFunnel` + 阶段金额 `ChartBar`（`/deals`）
- [ ] 赢单率 `ChartLine`（Deals 或 Dashboard）
- [ ] **`ChartSparkline`** 组件（KPI 卡内趋势）
- [ ] **`ChartGauge`** 组件（配额完成率）
- [ ] Dashboard：Sparkline×KPI + Line + Funnel + Gauge + Bar（经理排行）
- [ ] Dashboard 汇总 API + 数据权限（本人 / 部门 / 全部）

| 顺序 | 业务 | ui-kit | 嵌入页面 | 指标 |
|------|------|--------|----------|------|
| 3.1 | Deals + Pipeline | `ChartFunnel` | `/deals` | 阶段管道 |
| 3.2 | 阶段/金额 API | `ChartBar` | Deals 分析 | 金额、排行 |
| 3.3 | 赢单率 API | `ChartLine` | Deals / Dashboard | win rate |
| 3.4 | Dashboard 页面 | **`CardMetric`** | `/` 或 `/dashboard` | KPI 行 |
| 3.5 | 配额 API | **`ChartGauge`**（新增） | Dashboard | 目标 % |
| 3.6 | Dashboard 汇总 API | Line+Funnel+Bar+Gauge | Dashboard | 见 [05](../frontend-arch/05-component-scenarios.md) Part B |
| 3.7 | 经理视图 | `ChartBar` | Dashboard | 业绩排行 |

**本阶段新增组件**：`ChartSparkline`、`ChartGauge`（与 3.4、3.5 同迭代）。  
**Phase 3 图表 Done**：Dashboard **不得**仅有数字 KPI；至少 4 类图接 API。

---

## Phase 4：系统设置与收尾

**业务**

- [ ] 租户配置管理
- [ ] 自定义字段支持
- [ ] i18n 完整接入（中英）
- [x] 基础错误处理与统一响应格式
- [ ] Swagger API 文档
- [x] 权限与多租户集成测试（Phase 0 `scope_test` + QA）
- [ ] 部署文档

**图表**

- [ ] 审计：操作类型 `ChartDonut` / `ChartBar`（设置或审计页）
- [ ] **`ChartRadar`** 组件（租户健康度）
- [ ] Admin：`ChartLine` + `ChartDonut`（套餐）+ `ChartRadar` + `ChartBar`（租户 TOP）

| 顺序 | 业务 | ui-kit | 嵌入页面 | 指标 |
|------|------|--------|----------|------|
| 4.1 | 租户配置、自定义字段 | — | `/settings` | — |
| 4.2 | 审计日志 API | `ChartBar` / `ChartDonut` | 设置/审计 | 操作占比 |
| 4.3 | Admin 租户统计 | **`ChartRadar`**（新增） | `/admin` | 健康度 |
| 4.4 | Admin 套餐分布 | `ChartDonut` | `/admin` | 套餐占比 |
| 4.5 | i18n、Swagger、部署 | — | — | — |

**本阶段新增组件**：`ChartRadar`（与 4.3 同周）。

**图表工程（可与 Phase 4 并行）**

- [ ] Storybook：全 Chart 组件 V1/V2 Story
- [ ] VTU：本阶段新增 Chart 关键 props
- [ ] UI 视觉回归 CI（Chromatic 或等价）
- [ ] Playwright：Dashboard / Leads 报表视觉基线

---

## Phase 5+（MVP 之后）

**业务 / 组件（按报表迭代）**

- [ ] `ChartHeatmap` / `ChartScatter` / `ChartSankey`（报表中心）
- [ ] `ChartTreemap`（客户 ARR 分层）
- [ ] `ChartWaterfall`（Pipeline 金额变动）

组件选型详见 [05-component-scenarios.md](../frontend-arch/05-component-scenarios.md) Part B。

---

## 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-22 | 并入原 `frontend-arch/06-chart-module-dev-sync`（同步原则、Phase 矩阵、双周节奏） |
| 2026-05-22 | Phase 1 收尾：注册、审计日志、Admin 租户趋势 API + ChartLine |
| 2026-05-22 | Phase 2：关系经营任务 2.9–2.12、AI Preview 2.13（对齐 phase-2-relationship-crm-prd v0.3） |
| 2026-05-22 | Phase 2 Architect：phase-2-crm-ai.md、03-phase-2-crm-schema、ADR-0004、phase-2-notes |
| 2026-05-22 | Phase 2 BE：2.1 Accounts API + `00006`；2.3 `convert` + Leads HTTP 集成测 |
| 2026-05-22 | Phase 2 QA：2.1 `accounts_integration_test.go`；2.3 Leads HTTP 集成测 |
