# Phase 2 API 契约：客户关系 + 情绪旅程 + AI Preview

**版本**：v1.0  
**日期**：2026-05-22  
**状态**：Accepted（Implementation 开工前评审基线）  
**输入**：[phase-2-relationship-crm-prd.md](../prd/phase-2-relationship-crm-prd.md) §4–§6、§15.3  
**索引**：[00-api-design.md](./00-api-design.md)

> HTML：`phase-2-crm-ai.html`（`make docs-html` 生成）

---

## 1. 范围与原则

| 类别 | Phase 2 交付 |
|------|----------------|
| **生产 (L1)** | Accounts / Contacts / Leads CRUD、Activity 时间线、规则洞察、分群、统计 API、情绪旅程聚合 |
| **Preview (P0)** | Copilot / 预测分数 / 增强洞察 — **501 或 fixture**，`source: preview` |
| **预留 (L2/L3)** | 路径与 `AiCapabilityResult` 定稿；真算 Phase 3+ |

**通用约定**（继承 00-api-design）：

- Header：`Authorization`、`X-Tenant-ID`（Super Admin / auth 例外）
- 响应：`{ code, message, data, pagination? }`
- 列表默认 `page=1`、`page_size=20`（最大 100）
- 写操作写 `audit_logs`（见 §8）
- **禁止** body 传 `tenant_id`；由中间件注入

**新增 Header（AI Preview）**：

| Header | 说明 |
|--------|------|
| `X-CRM-Preview: 1` | 仅 `demo` 角色或 `tenant.config.ai_preview_mode=fixtures` 时生效；强制 fixture |
| `X-CRM-AI-Capability` | 可选，Copilot 场景见 §2.7 |

---

## 2. 枚举与共享类型

> 下列 **code** 为 API / DB 存贮值；**中文** 为默认展示名（租户可配置显示名，i18n key 见 `apps/web/locales`）。

### 2.1 生命周期 `lifecycle_stage`

| code | 中文 | 说明 |
|------|------|------|
| `acquire` | 获客 | 新线索未有效触达 |
| `activate` | 激活 | 首次有效互动 |
| `grow` | 成长 | 需求明确、推进中 |
| `retain` | 留存 | 已成交 / 稳定合作 |
| `revive` | 唤醒 | 曾活跃现沉默 |

### 2.2 关系健康 `relationship_health`

规则计算，**只读**（客户端不可 POST/PATCH 写入）。

| code | 中文 | 说明 |
|------|------|------|
| `high` | 健康 | 互动与阶段信号良好 |
| `medium` | 一般 | 需关注，未达风险阈值 |
| `low` | 风险 | 沉默或负面信号偏多 |

### 2.3 Lead 状态 `status`

| code | 中文 | 说明 |
|------|------|------|
| `new` | 新建 | 初始状态 |
| `contacted` | 已联系 | 已有触达 |
| `qualified` | 合格 | 意向明确 |
| `unqualified` | 不合格 | 归档 / 培育池 |
| `converted` | 已转化 | 须关联 `converted_account_id`（及可选 Contact） |

合法迁移：`new` → `contacted` → `qualified` → `unqualified` \| `converted`（非法返回 `invalid_status_transition`）。

### 2.4 Activity

| 字段 | 类型 | 说明 |
|------|------|------|
| `event_type` | enum | 见下表 |
| `direction` | enum | 见下表（`note` / `system` 可为空） |
| `subject_type` | enum | 见下表 |
| `subject_id` | uuid | 主体 ID |
| `sentiment` | enum nullable | 见 §2.4.1 |
| `sentiment_source` | enum | 见 §2.4.2（Phase 2 **禁止**写 `ai`） |
| `metadata` | jsonb | 时长、渠道、人工行为等 |
| `occurred_at` | timestamptz | 业务发生时间（默认 now） |

**`event_type`**

| code | 中文 |
|------|------|
| `note` | 备注 |
| `call` | 电话 |
| `email` | 邮件 |
| `meeting` | 会议 |
| `wechat` | 微信 |
| `visit` | 拜访 |
| `system` | 系统 |

**`direction`**

| code | 中文 |
|------|------|
| `inbound` | 入站（客户 / 对方发起） |
| `outbound` | 出站（我方发起） |

**`subject_type`**

| code | 中文 |
|------|------|
| `lead` | 线索 |
| `contact` | 联系人 |
| `account` | 公司 |

#### 2.4.1 情绪 `sentiment`

| code | 中文 | `sentiment_score`（情绪旅程聚合） |
|------|------|----------------------------------|
| `positive` | 积极 | 2 |
| `neutral` | 中性 | 0 |
| `hesitant` | 犹豫 | -1 |
| `negative` | 消极 | -2 |
| `unknown` | 未知 | `null` |

#### 2.4.2 情绪来源 `sentiment_source`

| code | 中文 | Phase 2 |
|------|------|---------|
| `manual` | 人工标注 | 允许 |
| `rule` | 规则推断 | 允许 |
| `ai` | AI 推断 | **禁止写入**（Phase 5+） |

### 2.5 分群 `segment_code`（预置）

| code | 中文 | 筛选逻辑（示例） |
|------|------|------------------|
| `high_value` | 高价值 | 预计金额 > 租户阈值 |
| `churn_risk` | 流失预警 | N 天无 Activity（默认 7） |
| `new_potential` | 潜在新客 | 近 7 天创建且未 `qualified` |
| `needs_activation` | 待激活 | `lifecycle=acquire` 且无 `outbound` |
| `revive_pool` | 唤醒池 | `lifecycle=revive` |

### 2.6 租户 AI 配置 `ai_preview_mode`

| code | 中文 | Phase 2 |
|------|------|---------|
| `off` | 关闭 | 不展示 Preview |
| `fixtures` | 演示夹具 | **默认演示**；`X-CRM-Preview: 1` 生效 |
| `live` | 真机推理 | 预留（Phase 3+） |

### 2.7 Copilot 场景 `scene` / `X-CRM-AI-Capability`

| code | 中文 |
|------|------|
| `followup_script` | 跟进话术 |
| `email_draft` | 邮件草稿 |
| `summarize` | 摘要 |

### 2.8 情绪旅程 Query `range`

| code | 中文 |
|------|------|
| `30d` | 近 30 天 |
| `90d` | 近 90 天（默认） |
| `all` | 全部 |

### 2.9 统计 Query `granularity`（`/api/leads/stats/trend`）

| code | 中文 |
|------|------|
| `day` | 按日 |
| `week` | 按周 |

### 2.10 数据范围 `data_scope`（`GET /api/rbac/my-permissions` 可选字段）

| code | 中文 |
|------|------|
| `self` | 本人（`owner_id = me`） |
| `department` | 部门 |
| `all` | 全部（经理种子） |

### 2.11 情绪旅程里程碑 `milestones[].type`（示例）

| code | 中文 |
|------|------|
| `converted` | 已转化 |
| `owner_changed` | 负责人变更 |
| `insight_triggered` | 洞察触发 |

### 2.12 `AiCapabilityResult`（Preview / 未来真 AI 统一包）

用于 Copilot、预测分、增强洞察等 **非 CRUD** 能力响应（可嵌在 `data` 内或作为 `data` 本体）：

```json
{
  "source": "rule",
  "capability": "insights.evaluate",
  "payload": {},
  "disclaimer": "ai.preview.disclaimer",
  "request_id": "uuid"
}
```

**`source`**

| code | 中文 | 说明 |
|------|------|------|
| `rule` | 规则引擎 | L1 生产逻辑 |
| `preview` | 演示数据 | Mock / fixture |
| `model` | 模型推理 | Phase 3+ 真 LLM |

**业务错误码（`code` 字段，HTTP 见下表）**：

| code | HTTP | 说明 |
|------|------|------|
| `AI_DISABLED` | 403 | `tenant.config.ai_enabled=false` |
| `AI_PREVIEW_ONLY` | 200 | 仅 Preview；body 含 `AiCapabilityResult` 且 `source=preview` |
| `AI_NOT_READY` | 501 | 能力未实现；FE 可降级 fixture |

---

## 3. 租户配置（AI 开关）

存储于 `tenants.config` JSONB（读写经 Admin API，Phase 2 可先 **种子 + PATCH 内部接口**）：

| 键 | 类型 | 默认 | 说明 |
|----|------|------|------|
| `ai_enabled` | bool | `false` | 是否展示 AI 模块 |
| `ai_preview_mode` | string | `off` | 见 §2.6（Phase 2 演示租户为 `fixtures`） |
| `insight_thresholds` | object | 内置默认 | 如 `days_silent: 7` |
| `sentiment_keyword_rules` | array | 可选 | 关键词 → sentiment（租户可配） |

**GET** `/api/tenant/config` — 返回当前租户 `config`（需 `settings:tenant_config` 或 Admin）  
**PATCH** `/api/tenant/config` — 合并更新上述键

演示租户种子（迁移 `00006`）：`ai_enabled=true`, `ai_preview_mode=fixtures`

---

## 4. Accounts（公司）

**权限**：`accounts:view|create|update|delete`；数据范围：本人 `owner_id` / 部门（Phase 2 同 leads：`view_all` 扩展为经理看全部，见 §8）

### 4.1 资源字段

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | uuid | 出 | |
| `name` | string | 是 | 公司名称 |
| `industry` | string | 否 | |
| `website` | string | 否 | |
| `owner_id` | uuid | 否 | 负责人 |
| `lifecycle_stage` | enum | 否 | 默认 `acquire` |
| `relationship_health` | enum | 出 | 只读 |
| `engagement_score` | int 0–100 | 出 | 只读 |
| `last_activity_at` | timestamptz | 出 | 只读 |
| `tags` | string[] | 否 | |
| `created_at` / `updated_at` | timestamptz | 出 | |

### 4.2 接口

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/accounts` | view | Query：`page`, `page_size`, `search`, `lifecycle_stage`, `relationship_health`, `segment`, `owner_id` |
| POST | `/api/accounts` | create | 创建 |
| GET | `/api/accounts/:id` | view | 详情 |
| PUT | `/api/accounts/:id` | update | 全量更新（含 `lifecycle_stage`，写审计） |
| PATCH | `/api/accounts/:id` | update | 部分更新 |
| DELETE | `/api/accounts/:id` | delete | 软删 |
| GET | `/api/accounts/:id/emotion-journey` | view | §7 |
| POST | `/api/accounts/:id/insights/evaluate` | view | §6（主体 account） |

**列表响应 `data`**：`{ items: Account[], pagination }`

---

## 5. Contacts（联系人）

| 字段 | 类型 | 说明 |
|------|------|------|
| `account_id` | uuid nullable | 关联公司 |
| `first_name` / `last_name` | string | |
| `email` / `phone` | string | |
| `is_primary` | bool | 主联系人 |
| 关系字段 | 同 Account | `lifecycle_stage`, `engagement_score`, … |

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/contacts` | 列表；Query 同 accounts + `account_id` |
| POST | `/api/contacts` | 创建 |
| GET | `/api/contacts/:id` | 详情 |
| PUT/PATCH | `/api/contacts/:id` | 更新 |
| DELETE | `/api/contacts/:id` | 软删 |
| GET | `/api/accounts/:id/contacts` | 公司下联系人 |
| GET | `/api/contacts/:id/emotion-journey` | §7 |
| POST | `/api/contacts/:id/insights/evaluate` | §6 |

---

## 6. Leads（线索）

### 6.1 资源字段

在现有 `leads` 表基础上扩展（见 [03-phase-2-crm-schema.md](../architecture/03-phase-2-crm-schema.md)）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `title` | string | 线索标题 |
| `status` | enum | §2.3 |
| `source` | string | 来源 |
| `amount` | decimal | 预计金额 |
| `expected_close_date` | date | |
| `owner_id` | uuid | |
| `lifecycle_stage` | enum | |
| `relationship_health` | enum | 只读 |
| `engagement_score` | int | 只读 |
| `last_activity_at` | timestamptz | 只读 |
| `tags` | string[] | |
| `converted_account_id` | uuid nullable | `status=converted` 时必填 |
| `converted_contact_id` | uuid nullable | 可选 |

### 6.2 CRUD

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/leads` | view | Query：`status`, `source`, `owner_id`, `lifecycle_stage`, `relationship_health`, `segment`, `search` |
| POST | `/api/leads` | create | |
| GET | `/api/leads/:id` | view | |
| PUT/PATCH | `/api/leads/:id` | update | 含状态变更 |
| DELETE | `/api/leads/:id` | delete | 软删 |
| POST | `/api/leads/:id/convert` | update | Body：`{ account_id?, contact_id?, create_account?: { name } }`；写审计 `lead.convert` |
| POST | `/api/leads/:id/assign` | **assign** | Body：`{ owner_id, note? }`；创建 system Activity；需 `leads:assign` |
| POST | `/api/leads/import` | create | multipart Excel；返回 `{ imported, failed, errors_url? }` |
| GET | `/api/leads/import/template` | view | 下载模板（Could） |
| GET | `/api/leads/:id/emotion-journey` | view | §7 |
| POST | `/api/leads/:id/insights/evaluate` | view | §6.3 |

**状态机**：非法迁移返回 `400` + `{ code: 400, message: "invalid_status_transition" }`

### 6.3 规则洞察 `insights/evaluate`

**POST** `/api/{leads|contacts|accounts}/:id/insights/evaluate`

**权限**：与主体 `view` 一致（resource 为第一段路径）

**响应 data**：

```json
{
  "items": [
    {
      "id": "INS-001",
      "priority": 1,
      "title_key": "insight.silent_risk.title",
      "body_key": "insight.silent_risk.body",
      "suggested_action": {
        "activity_event_type": "call",
        "activity_direction": "outbound",
        "title": "今日回访"
      },
      "rule_id": "INS-001",
      "explanation_key": "insight.silent_risk.explanation"
    }
  ],
  "engagement_score": 62,
  "engagement_explanation_key": "engagement.last_7_days",
  "lifecycle_stage": "grow",
  "relationship_health": "medium"
}
```

- 最多返回 **3** 条（按 `priority`）
- 规则集：INS-001～006（PRD §4.3.2、§4.6.6）
- **服务端求值**；禁止客户端传规则表达式

**POST** `/api/insights/preview`（Preview 增强包）

**权限**：`insights:view` + `ai_enabled`

**请求**：`{ subject_type, subject_id }`

**响应**：`AiCapabilityResult`，`payload` 含 Mock 流失概率、增强文案；`source=preview`

---

## 7. 情绪旅程 `emotion-journey`

**GET** `/api/{leads|contacts|accounts}/:id/emotion-journey`

**Query**：`range=30d|90d|all`（默认 `90d`）

**响应 data**（与 PRD §4.6.3 一致）：

```json
{
  "subject_type": "lead",
  "subject_id": "uuid",
  "lifecycle_current": "grow",
  "lifecycle_bands": [
    { "stage": "acquire", "from": "2026-01-01T00:00:00Z", "to": "2026-02-01T00:00:00Z" }
  ],
  "points": [
    {
      "activity_id": "uuid",
      "at": "2026-03-15T10:00:00Z",
      "event_type": "call",
      "direction": "outbound",
      "sentiment": "hesitant",
      "sentiment_score": -1,
      "sentiment_source": "manual",
      "label": "电话：价格顾虑",
      "lifecycle_stage_at_time": "grow"
    }
  ],
  "milestones": [
    { "type": "converted", "at": "2026-04-01T12:00:00Z", "label": "转为商机" }
  ],
  "summary": {
    "current_sentiment": "hesitant",
    "trend": "down",
    "days_since_positive": 14
  }
}
```

`sentiment_score` 映射：`positive=2`, `neutral=0`, `hesitant=-1`, `negative=-2`, `unknown=null`

Preview：`X-CRM-Preview: 1` 可返回固定 fixture（与生产结构相同）

---

## 8. Activities（跟进 / 时间线）

**权限**：`activities:view|create|update|delete`

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/activities` | Query：`subject_type`, `subject_id`, `page`, `page_size`；按 `occurred_at` DESC |
| POST | `/api/activities` | 创建 |
| GET | `/api/activities/:id` | 详情 |
| PATCH | `/api/activities/:id` | 更新 body、sentiment 等 |
| DELETE | `/api/activities/:id` | 软删 |
| GET | `/api/activities/summary` | ChartBar 用：按 `event_type` TOP N（Query：`subject_type`, `subject_id` 或租户级 + 数据范围） |

**POST body 示例**：

```json
{
  "subject_type": "lead",
  "subject_id": "uuid",
  "event_type": "call",
  "direction": "outbound",
  "body": "客户对价格有顾虑",
  "sentiment": "hesitant",
  "occurred_at": "2026-05-22T10:00:00Z",
  "metadata": { "duration_minutes": 15 }
}
```

创建/更新后：异步或同步刷新主体 `last_activity_at`、`engagement_score`（同事务亦可）。

---

## 9. 分群 Segments

**权限**：`segments:view`（使用模板）、`segments:manage`（Admin 改阈值）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/segments` | 租户预置模板列表 `{ code, name_key, description_key, filter, count? }` |
| GET | `/api/segments/:code/count` | 当前用户数据范围内的匹配数量 |
| PATCH | `/api/segments/:code` | Admin 更新 `filter` / 阈值（manage） |

列表接口 Query `segment=<code>` 与分群 filter 语义一致。

---

## 10. 统计 API（Leads 报表）

**权限**：`leads:view`；数据范围与列表一致（`scope=self|department|all`，由角色配置 + Service 层过滤）

| 方法 | 路径 | 图表 | Query |
|------|------|------|-------|
| GET | `/api/leads/stats/by-source` | ChartDonut | `from`, `to` |
| GET | `/api/leads/stats/trend` | ChartLine | `from`, `to`, `granularity=day\|week`, `group_by=lifecycle?` |
| GET | `/api/leads/stats/funnel` | ChartFunnel | `from`, `to` |
| GET | `/api/leads/stats/by-status` | ChartBar | `from`, `to` |
| GET | `/api/leads/stats/by-health` | ChartBar | `relationship_health` 分布 |

**响应示例**（`by-source`）：

```json
{
  "items": [
    { "label": "web", "value": 120, "percentage": 0.4 }
  ],
  "total": 300
}
```

**funnel**：

```json
{
  "stages": [
    { "name": "new", "count": 100 },
    { "name": "contacted", "count": 60 },
    { "name": "qualified", "count": 25 }
  ],
  "conversion_rates": [
    { "from": "new", "to": "contacted", "rate": 0.6 }
  ]
}
```

---

## 11. Copilot & AI Preview（P0）

**权限**：`copilot:use`（演示角色默认有）；且 `tenant.config.ai_enabled=true`

| 方法 | 路径 | capability | Phase 2 |
|------|------|------------|---------|
| POST | `/api/copilot/summarize` | `copilot.summarize` | 501 或 Preview fixture |
| POST | `/api/copilot/generate` | `copilot.generate` | Preview fixture |
| POST | `/api/copilot/chat` | `copilot.chat` | 501 |
| POST | `/api/scores/predict` | `scores.predict` | Preview fixture |
| POST | `/api/segments/dynamic/refresh` | `segments.dynamic` | Preview fixture |

**POST** `/api/copilot/generate`

**请求**：

```json
{
  "scene": "followup_script",
  "subject_type": "lead",
  "subject_id": "uuid",
  "context": { "last_activity_id": "uuid" }
}
```

**响应 data**：`AiCapabilityResult`，`payload`：

```json
{
  "text": "您好张总，上次您提到对部署周期有顾虑……",
  "scene": "followup_script"
}
```

**POST** `/api/copilot/summarize` — Body：`{ activity_id }` 或 `{ text }`；Phase 2 返回 501 + `AI_NOT_READY` 或 Preview 摘要。

---

## 12. 生命周期建议（LC-03）

**POST** `/api/{leads|contacts|accounts}/:id/lifecycle/suggest`

**响应 data**：

```json
{
  "suggested_stage": "grow",
  "reason_key": "lifecycle.suggest.positive_qualified",
  "requires_confirmation": true
}
```

**POST** `.../lifecycle/confirm` — Body：`{ stage }`；写入主体并审计 `lifecycle.change`。

---

## 13. RBAC 增量

迁移 `00006` 新增 permissions：

| resource | action | 说明 |
|----------|--------|------|
| `leads` | `assign` | 分配线索 |
| `activities` | `view`, `create`, `update`, `delete` | 跟进 |
| `insights` | `view` | 洞察只读 |
| `segments` | `view`, `manage` | 分群 |
| `copilot` | `use` | Copilot / Preview |

**数据范围**（Service 层，非 Casbin 路径）：

| 角色模板 | leads/accounts/contacts | activities |
|----------|-------------------------|------------|
| Sales | `owner_id = me` | 本人相关主体 |
| Manager | `department` 或 `all`（配置） | 部门 |
| Viewer | view only | view only |

`GET /api/rbac/my-permissions` 增加返回 `data_scope: self|department|all`（可选字段，Phase 2 默认 `self`，经理种子 `all`）。

---

## 14. 审计动作

| action | 触发 |
|--------|------|
| `lead.convert` | 转化 |
| `lead.assign` | 分配 |
| `lead.import` | 导入 |
| `lifecycle.change` | 阶段变更 |
| `activity.sentiment_change` | 情绪修改 |
| `insight.adopt` | 采纳建议（FE 上报或创建 Activity 时带 `metadata.insight_id`） |

---

## 15. 实现顺序（供 BE 领任务）

1. 迁移 `00006` + 域模型  
2. Accounts → Contacts → Leads 扩展  
3. Activities + `last_activity_at` 触发  
4. `insights/evaluate` + engagement 计算  
5. Segments 模板 + 列表 `segment` 筛选  
6. Stats 四个端点  
7. `emotion-journey` 聚合  
8. Copilot 501/Preview + tenant config  
9. 权限种子 + 路由注册  

---

## 16. 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-22 | v1.0：架构师首版，对齐 PRD v0.3 §15.3 |
| 2026-05-24 | §2 枚举补充中文注释表（对齐 PRD / i18n） |
