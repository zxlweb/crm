# Phase 3 阶段笔记 — 商机与仪表盘

**日期**：2026-05-24  
**状态**：PRD v0.1 + Architect 3.0 完成 → **可开启 Implementation 三轨**  
**PRD**：[phase-3-deals-dashboard-prd.md](../prd/phase-3-deals-dashboard-prd.md)

---

## 1. 架构产出（3.0）

| 文档 | 路径 | 状态 |
|------|------|------|
| Deals + Dashboard API 契约 v1.0 | [docs/api/phase-3-deals-dashboard-api.md](../api/phase-3-deals-dashboard-api.md) | ✅ Accepted |
| 数据模型 / 迁移纲要 | [docs/architecture/04-phase-3-deals-dashboard-schema.md](../architecture/04-phase-3-deals-dashboard-schema.md) | ✅ |

**评审结论**：契约覆盖并行表 3.1–3.7 路径；`deals` 五阶段 Pipeline + Dashboard `summary/funnel/quota/team-ranking`；BE 按 phase-3-deals-dashboard-api §9 实现顺序领任务。

---

## 2. 前端切面（2b）

### 2.1 路由（`apps/web/pages/`）

| 路由 | 页面 | layout |
|------|------|--------|
| `/deals` | Pipeline 看板 + 列表入口 | `app` |
| `/deals/[id]` | 商机详情 / 编辑 | `app` |
| `/` | 工作台（`GET /api/dashboard/summary` 替代多列表聚合） | `app` |

权限：`deals:view|create|update|delete`；`dashboard:view`（汇总端点）。

### 2.2 Composables

| 名称 | 职责 |
|------|------|
| `use-deals` | CRUD + `fetchPipeline` + `updateStage` |
| `use-deal-stats` | `by-stage`、`win-rate`（3.2–3.3） |
| `use-dashboard-stats` | `summary`、`funnel`、`quota`、`teamRanking`（3.4–3.7） |

### 2.3 组件落点

| 组件 | 路径 |
|------|------|
| Deals 看板 | `components/feature/deals/deals-pipeline-board.vue`（建议） |
| Pipeline 图 | `ChartFunnel` @ `/deals` |
| KPI 行 | `dashboard-kpi-row.vue` + **`ChartSparkline`**（`summary.sparklines`） |
| 配额 | **`ChartGauge`** ← `GET /api/dashboard/quota` |
| 经理排行 | `ChartBar` ← `team-ranking` |

### 2.4 `data-testid`（QA 约定）

| 元素 | testid |
|------|--------|
| 新建商机 | `deal-create-btn` |
| Pipeline 看板 | `deals-pipeline` |
| 阶段列 | `deals-pipeline-stage-{stage}` |
| Dashboard KPI 行 | `dashboard-kpi-row` |
| 配额仪表 | `dashboard-quota-gauge` |

---

## 3. BE 3.1–3.7 交付（2026-05-25 冲刺）

- 迁移：`00010`–`00012`（`deals` + `dashboard:view` + demo seed）
- **Deals**：CRUD、pipeline、stage、`stats/by-stage`、`stats/win-rate`
- **Dashboard**：`summary`、`funnel`、`quota`、`team-ranking`；`todo` 返回 `items: []`（P2 stub）
- 测试：`deals_integration_test.go`、`dashboard_integration_test.go`；`make backend-test` 通过
- **FE**：`use-deals` / `use-dashboard-stats` 可接真 API；`make migrate-up` 后联调

## 4. FE 联调（2026-05-25）

- **Mock 策略**：默认走真实 API；仅 `NUXT_PUBLIC_USE_DEALS_MOCK=true` / `NUXT_PUBLIC_USE_DASHBOARD_MOCK=true` 时强制 fixture
- **已验证端点**（Demo 租户 `admin@demo.com` / `password123`）：
  - `GET /api/dashboard/summary` → KPI + sparklines + priorities
  - `GET /api/dashboard/quota|funnel` → Gauge + Funnel
  - `GET /api/deals/stats/by-stage` → Deals 分析 Tab
- **经理排行**：`data_scope=self` 时 BE 返回 403；FE 隐藏 `dashboard-team-ranking` 列，不阻断 quota/funnel
- **本地**：BE `:8080` + FE `pnpm --filter @crm/web dev`；确认未设置 mock env

## 5. 三轨领任务指引

- **BE**：3.1–3.7 ✅
- **FE**：`【FE】Phase 3 — 3.UI` 与 `3.1` 并行；`3.4` 接 `dashboard/summary`
- **QA**：`【QA】Phase 3 — [phase-3-deals-qa-plan.md](../qa/phase-3-deals-qa-plan.md)`（已就绪，按子任务 3.1→3.E2E 领测）

**Week 1**：3.1 BE + 3.UI FE + 3.4 summary；**Week 2**：3.2–3.7 + 3.E2E。

---

## 6. 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-24 | 初稿：Phase 3 开工切面 |
| 2026-05-24 | 3.0 完成：phase-3-deals-dashboard-api v1.0、04-phase-3 schema，开放 Implementation |
| 2026-05-25 | 链 Phase 3 详细 PRD v0.1 |
| 2026-05-25 | BE 3.1：Deals CRUD + Pipeline + stage + convert create_deal |
| 2026-05-25 | BE 3.1–3.7 一晚冲刺：Dashboard + Deals 统计全 API |
| 2026-05-25 | FE 联调：切真实 API，经理排行 403 降级，移除 silent mock fallback |
| 2026-05-25 | 联调 polish：Casbin 域模型修复、Demo 销售配额、赢单率单点、文案清理 |
