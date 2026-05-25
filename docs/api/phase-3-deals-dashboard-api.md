# Phase 3 API 契约：商机 (Deals) + 仪表盘 (Dashboard)

**版本**：v1.0  
**日期**：2026-05-24  
**状态**：Accepted（Implementation 开工前评审基线）  
**输入**：[phase-3-deals-dashboard-prd.md](../prd/phase-3-deals-dashboard-prd.md) · [00-crm-overview.md](../prd/00-crm-overview.md) §Phase 3 · [00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md) §Phase 3  
**数据模型**：[04-phase-3-deals-dashboard-schema.md](../architecture/04-phase-3-deals-dashboard-schema.md)  
**索引**：[00-api-design.md](./00-api-design.md)

> HTML：`phase-3-deals-dashboard-api.html`（`make docs-html` 生成）

---

## 1. 范围与原则

| 类别 | Phase 3 交付 |
|------|----------------|
| **Deals (L1)** | 商机 CRUD、Pipeline 看板、阶段变更、统计 API |
| **Dashboard (L1)** | 汇总 KPI、漏斗、待办、经理排行、配额完成率 |
| **复用 Phase 2** | 多租户、`data_scope`、关系字段（`engagement_score` 等可选只读） |
| **不在 MVP** | 自定义 Pipeline 阶段编排 UI、审批流、多币种汇率 |

**通用约定**（继承 [00-api-design.md](./00-api-design.md)）：

- Header：`Authorization`、`X-Tenant-ID`
- 响应：`{ code, message, data, pagination? }`
- 列表默认 `page=1`、`page_size=20`（最大 100）
- 写操作写 `audit_logs`（见 §10）
- **禁止** body 传 `tenant_id`

**与 Phase 2 关系**：

- 线索 `POST /api/leads/:id/convert` 可创建 Account/Contact；**可选**同时创建 Deal（`create_deal` body，见 §4.3 扩展）
- Dashboard「今日优先」可继续聚合 Leads/Accounts；本契约提供 **统一汇总端点** 减少 FE N+1

---

## 2. 枚举与共享类型

### 2.1 商机阶段 `stage`（Pipeline）

| code | 中文 | 说明 |
|------|------|------|
| `qualification` | 需求确认 | 默认阶段；评估需求与预算 |
| `proposal` | 方案报价 | 已提交方案/报价 |
| `negotiation` | 商务谈判 | 价格与条款博弈 |
| `won` | 赢单 | **终态**；须写 `closed_at` |
| `lost` | 输单 | **终态**；须写 `closed_at`，建议 `lost_reason` |

**状态机**（非法迁移 → `400` + `invalid_stage_transition`）：

```
qualification → proposal → negotiation → won | lost
```

- `won` / `lost` 不可再变更 `stage`（除 Admin 纠错 PATCH，Phase 3 可不实现）
- 回退（如 `negotiation` → `proposal`）：**允许**，写审计 `deal.stage_change`

### 2.2 币种 `currency`

| code | 中文 |
|------|------|
| `CNY` | 人民币（默认） |
| `USD` | 美元 |

Phase 3 仅存储 code，不做汇率换算。

### 2.3 Dashboard `data_scope`

与 Phase 2 一致，由 `GET /api/rbac/my-permissions` 返回（可选字段）：

| code | 中文 | 过滤 |
|------|------|------|
| `self` | 本人 | `owner_id = 当前用户` |
| `department` | 部门 | `department_id` 匹配（种子可等同 `all`） |
| `all` | 全部 | 租户内可见范围 |

---

## 3. Deals（商机）资源

### 3.1 字段

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | uuid | 出 | |
| `title` | string | 是 | 商机名称 |
| `stage` | enum | 否 | 默认 `qualification` |
| `amount` | decimal(18,2) | 否 | 预计金额，≥ 0 |
| `currency` | enum | 否 | 默认 `CNY` |
| `probability` | int 0–100 | 否 | 成交概率；`won`=100，`lost`=0 建议 |
| `expected_close_date` | date | 否 | 预计成交日 |
| `account_id` | uuid | 否 | 关联公司 |
| `lead_id` | uuid | 否 | 来源线索 |
| `contact_id` | uuid | 否 | 主联系人 |
| `owner_id` | uuid | 否 | 负责人；默认当前用户 |
| `description` | text | 否 | |
| `tags` | string[] | 否 | |
| `lost_reason` | string | 否 | `stage=lost` 时建议必填 |
| `closed_at` | timestamptz | 出 | `won`/`lost` 时写入 |
| `engagement_score` | int | 出 | 只读；规则同 Lead（可选 Phase 3） |
| `last_activity_at` | timestamptz | 出 | 只读 |
| `created_at` / `updated_at` | timestamptz | 出 | |

### 3.2 CRUD

**权限**：`deals:view|create|update|delete`；数据范围见 §9。

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/deals` | view | Query 见下表 |
| POST | `/api/deals` | create | 创建 |
| GET | `/api/deals/:id` | view | 详情 |
| PUT | `/api/deals/:id` | update | 全量更新（含 `stage`，走状态机） |
| PATCH | `/api/deals/:id` | update | 部分更新 |
| DELETE | `/api/deals/:id` | delete | 软删 |

**列表 Query**

| 参数 | 说明 |
|------|------|
| `page`, `page_size` | 分页 |
| `search` | 标题模糊 |
| `stage` | 单阶段 |
| `stages` | 逗号分隔多阶段 |
| `owner_id` | 负责人 |
| `account_id` | 公司 |
| `lead_id` | 线索 |
| `expected_close_from`, `expected_close_to` | 预计成交日区间 |
| `min_amount`, `max_amount` | 金额筛选 |

**列表响应 `data`**：`{ items: Deal[], pagination }`

### 3.3 Pipeline 看板

**GET** `/api/deals/pipeline`

**权限**：`deals:view`

**Query**：`owner_id`（可选）、`account_id`（可选）；数据范围与列表一致。

**响应 data**：

```json
{
  "stages": [
    {
      "stage": "qualification",
      "count": 4,
      "amount_total": 580000,
      "items": [
        {
          "id": "uuid",
          "title": "云帆智造年度订阅",
          "amount": 280000,
          "currency": "CNY",
          "probability": 40,
          "expected_close_date": "2026-07-15",
          "account_id": "uuid",
          "owner_id": "uuid"
        }
      ]
    }
  ],
  "summary": {
    "open_count": 12,
    "open_amount": 3200000,
    "won_count_mtd": 2,
    "won_amount_mtd": 450000
  }
}
```

- `items`：每阶段默认最多返回 **20** 条（按 `updated_at DESC`）；完整列表走 `GET /api/deals?stage=`
- FE `ChartFunnel` 使用 `stages[].count` 或 `amount_total`（Query `metric=count|amount`，默认 `count`）

### 3.4 阶段变更（推荐）

**PUT** `/api/deals/:id/stage`

**权限**：`deals:update`

**请求**：

```json
{
  "stage": "negotiation",
  "lost_reason": null,
  "note": "客户确认采购清单"
}
```

**响应 data**：更新后的 `Deal` 对象。

> 亦可通过 `PATCH /api/deals/:id` 更新 `stage`；专用路径便于审计与 E2E。

### 3.5 统计 API（Deals 报表 / 任务 3.2–3.3）

**权限**：`deals:view`；数据范围与列表一致。

| 方法 | 路径 | 图表 | Query |
|------|------|------|-------|
| GET | `/api/deals/stats/by-stage` | `ChartBar` | `from`, `to`, `metric=count\|amount` |
| GET | `/api/deals/stats/win-rate` | `ChartLine` | `from`, `to`, `granularity=week\|month` |
| GET | `/api/deals/stats/velocity` | — | 平均阶段停留天数（Could） |

**`by-stage` 响应示例**：

```json
{
  "items": [
    { "label": "qualification", "value": 8, "amount": 1200000 },
    { "label": "proposal", "value": 5, "amount": 900000 }
  ],
  "total": 13
}
```

**`win-rate` 响应示例**：

```json
{
  "items": [
    { "period": "2026-W20", "won": 2, "lost": 1, "rate": 0.67 }
  ]
}
```

---

## 4. Dashboard（仪表盘）

**权限**：`leads:view` 或 `deals:view`（至少其一）；各 KPI 仅统计有权限的资源。

### 4.1 汇总 `summary`（任务 3.4）

**GET** `/api/dashboard/summary`

**Query**：`preview=1`（可选；演示租户可合并 fixture 元数据，**不写入库**）

**响应 data**：

```json
{
  "data_scope": "self",
  "kpis": {
    "leads_total": 5,
    "accounts_total": 2,
    "deals_total": 8,
    "deals_open_count": 6,
    "deals_open_amount": 2800000,
    "at_risk_total": 3,
    "avg_engagement": 62,
    "weekly_follow_ups": 12
  },
  "kpi_trends": {
    "leads_weekly_touch": 3,
    "accounts_weekly_touch": 2,
    "deals_weekly_new": 1,
    "engagement_delta": 15,
    "engagement_direction": "up"
  },
  "sparklines": {
    "leads": [2, 3, 1, 4, 3, 5, 4],
    "deals": [0, 1, 0, 2, 1, 1, 2]
  },
  "priorities": [
    {
      "entity_type": "lead",
      "entity_id": "uuid",
      "title": "华创科技",
      "reasons": ["7 天未跟进"],
      "suggestion": "今日电话确认方案",
      "score": 38,
      "engagement_score": 62,
      "is_preview": false
    }
  ]
}
```

- `priorities`：最多 **5** 条，规则复用 Phase 2（沉默、低健康度）；`is_preview` 仅演示数据为 true
- `sparklines`：近 **7** 个自然日点数，供 `ChartSparkline`
- FE 逐步用本端点替代 `useDashboard().loadSnapshot()` 内多次 `leads/accounts` 列表聚合

### 4.2 漏斗 `funnel`（任务 3.6）

**GET** `/api/dashboard/funnel`

**Query**：`scope=deals|leads`（默认 `deals`）

**响应 data**（`scope=deals`）：

```json
{
  "stages": [
    { "name": "qualification", "count": 8 },
    { "name": "proposal", "count": 5 },
    { "name": "negotiation", "count": 3 },
    { "name": "won", "count": 2 }
  ]
}
```

### 4.3 待办 `todo`（Could · 任务 3.6）

**GET** `/api/dashboard/todo`

**Query**：`date=YYYY-MM-DD`（默认今天）

**响应 data**：

```json
{
  "items": [
    {
      "id": "uuid",
      "time": "09:30",
      "title": "华创科技",
      "subtitle": "跟进 · 线索",
      "href_entity_type": "lead",
      "href_entity_id": "uuid"
    }
  ]
}
```

可与 FE 现有 `dashboard-calendar` fixture 对齐；BE 可从 Leads/Activities 聚合。

### 4.4 配额 `quota`（任务 3.5）

**GET** `/api/dashboard/quota`

**权限**：`deals:view`；经理 `data_scope=all` 可看团队配额。

**响应 data**：

```json
{
  "target_amount": 5000000,
  "won_amount_mtd": 3400000,
  "completion_rate": 0.68,
  "period": "2026-05"
}
```

- 配额存 `tenants.config.sales_quota` JSONB（见 Schema §3.2）
- FE `ChartGauge` 绑定 `completion_rate`

### 4.5 经理排行 `team-ranking`（任务 3.7）

**GET** `/api/dashboard/team-ranking`

**权限**：`deals:view` + `data_scope` 为 `department` 或 `all`

**Query**：`metric=won_amount|win_count|avg_engagement`（默认 `won_amount`）、`limit=10`

**响应 data**：

```json
{
  "items": [
    {
      "user_id": "uuid",
      "name": "张三",
      "value": 1200000,
      "rank": 1
    }
  ]
}
```

FE `ChartBar` 横向排行 + 现有 `dashboard-team-heatmap` 可共用用户维度。

---

## 5. 与 Leads Convert 扩展（可选 · 3.1）

**POST** `/api/leads/:id/convert` 增加可选字段：

```json
{
  "create_account": { "name": "华创科技" },
  "create_deal": {
    "title": "华创科技 — 年度合作",
    "amount": 280000,
    "stage": "qualification"
  }
}
```

响应 `data` 增加 `deal_id`。未传 `create_deal` 时行为与 Phase 2 一致。

---

## 6. RBAC

已有种子（`00002`）：`deals:view|create|update|delete`。

| 资源 | action | 说明 |
|------|--------|------|
| `deals` | `view` | 查看 |
| `deals` | `create` | 创建 |
| `deals` | `update` | 更新 / 改阶段 |
| `deals` | `delete` | 删除 |
| `dashboard` | `view` | 访问汇总端点（**新增** migration） |

**数据范围**（Service 层）：

| 角色模板 | deals |
|----------|-------|
| Sales | `owner_id = me` |
| Manager | `all`（种子） |
| Viewer | view only |

路由映射（`route.go`）：`/api/deals/*` → `deals`；`/api/dashboard/*` → `dashboard` + `view`。

---

## 7. 审计动作

| action | 触发 |
|--------|------|
| `deal.create` | POST deals |
| `deal.update` | PUT/PATCH |
| `deal.delete` | DELETE |
| `deal.stage_change` | PUT stage 或 PATCH stage |
| `deal.convert_from_lead` | convert 时 create_deal |

---

## 8. 错误码补充

| message | HTTP | 场景 |
|---------|------|------|
| `invalid_stage_transition` | 400 | 非法阶段迁移 |
| `deal_closed_readonly` | 400 | 终态商机禁止改 stage |
| `dashboard_scope_denied` | 403 | 无 `dashboard:view` |

---

## 9. 实现顺序（供 BE 领任务）

| 顺序 | 任务 ID | API |
|------|---------|-----|
| 1 | 3.1 | 迁移 `deals` + CRUD + pipeline + PUT stage |
| 2 | 3.4 | `GET /api/dashboard/summary` |
| 3 | 3.2–3.3 | `deals/stats/*` |
| 4 | 3.5 | `GET /api/dashboard/quota` |
| 5 | 3.6 | `funnel` + `todo` |
| 6 | 3.7 | `team-ranking` |
| 7 | — | 权限 `dashboard:view` 种子 + 集成测 |

---

## 10. 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-24 | v1.0：架构师 3.0 首版，对齐 Phase 3 并行表 3.1–3.7 |
| 2026-05-25 | 重命名为 `phase-3-deals-dashboard-api.md`（与 PRD / `phase-2-crm-ai.md` 规范一致） |
