# Phase 3 PRD：商机 Pipeline 与成交仪表盘

**产品名称**：EnterpriseFlow CRM  
**版本**：Phase 3 v0.2  
**日期**：2026-05-26  
**状态**：已交付（含部门数据范围 v0.2 增补）  
**MVP 总览**：[00-crm-overview.md](./00-crm-overview.md)  
**前置依赖**：Phase 0 架构、Phase 1 认证/RBAC、Phase 2 客户/线索/Activity/洞察主链（可有条件关闭）  
**关联任务**：[00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md) § Phase 3  
**API 契约（Accepted）**：[phase-3-deals-dashboard-api.md](../api/phase-3-deals-dashboard-api.md) v1.0  
**数据模型**：[04-phase-3-deals-dashboard-schema.md](../architecture/04-phase-3-deals-dashboard-schema.md)  
**前端切面**：[phase-3-notes.md](../meeting-notes/phase-3-notes.md) §2b

---

## 0. 执行摘要（给决策者）

Phase 2 解决了 **「谁是我的客户、关系健康吗、今天该跟进谁」**。Phase 3 回答销售团队的下一问：**「钱在哪、能不能赢、离目标还差多少」**。

| 范式问题 | Phase 2 答案 | Phase 3 目标答案 |
|----------|--------------|------------------|
| 关系经营 vs 成交经营 | 健康度、情绪、洞察、今日优先 | **Pipeline 可视化 + 阶段推进 + 金额预测** |
| 个人作战 vs 团队管理 | 销售工作台 KPI（线索/公司） | **个人/经理 Dashboard**：赢单率、配额、排行 |
| Preview vs 生产 | Zone E 迷你漏斗为 fixtures | **Zone E 接真实 Deals API**；Preview 角标仅保留 AI 预测块 |

**结论**：Phase 3 是 MVP **Must Have** 的最后一块业务拼图（MoSCoW：`Deals + Dashboard`）。交付标准不是「能录入商机」，而是 **Pipeline 看板 + 至少 4 类图表接 API + 数据范围正确**。

---

## 1. 背景与业务目标

### 1.1 业务背景

B2B 销售在 Phase 2 已能管理线索与客户关系，但缺少：

1. **成交管道**：线索「合格」之后去哪？金额、阶段、预计关单日如何跟踪？  
2. **管理视角**：经理无法一眼看到团队 Pipeline 总额、阶段分布、赢单率与配额完成度。  
3. **工作台闭环**：首页 Zone E 的迷你漏斗仍是 Preview fixtures，无法支撑「从关系到成交」的完整叙事。

### 1.2 与 Phase 2 的衔接

| Phase 2 产出 | Phase 3 消费方式 |
|--------------|------------------|
| `POST /api/leads/:id/convert` | 转化时可选 `create_deal`（写入 `lead_id` 关联） |
| `engagement_score`、情绪、洞察 | Dashboard「需关注」与 Deals 详情侧栏只读展示（不重复算） |
| 销售工作台 `/`（§11.5） | KPI 行扩展 **商机指标**；Zone E 漏斗/热力 **生产化**（漏斗接 API，热力保留规则/Mock 至 L2 真算） |
| Leads 报表 `ChartFunnel` | Deals Pipeline 复用组件，**数据源切换**为 Deals 统计 API |

### 1.3 Phase 3 业务目标（可度量）

| 目标 ID | 描述 | 指标（Phase 3 末） |
|---------|------|-------------------|
| G1 | 销售能在 **Pipeline 看板** 管理商机阶段 | `/deals` 看板可拖拽/变更阶段；CRUD 全链路 |
| G2 | 经理能在 **Dashboard** 看团队成交健康度 | 经理视图：Pipeline 漏斗 + 赢单率 + 业绩排行接 API |
| G3 | 个人销售能看 **配额完成率** | `ChartGauge` 接配额 API；本人数据范围正确 |
| G4 | 首页 KPI **含商机维度** | `GET /api/dashboard/summary` 替代客户端 N 次聚合 |
| G5 | 线索转化到商机 **可追溯** | 从 Lead 详情可跳转关联 Deal；审计 `deal.create` / `lead.convert` |
| G6 | 为 L2 预测铺路 | 「关系降温」名单规则版上线（Should）；契约预留 `scores.predict` |

---

## 2. 用户角色 (Persona) 与场景

| 角色 | Phase 3 核心场景 | 诉求 | 数据范围（`data_scope`） |
|------|------------------|------|-------------------------|
| **Sales** | 每日看 Pipeline、推进阶段、录入金额 | 我的商机在哪一列、本月离配额还差多少 | `self`：仅 `owner_id = 本人` |
| **Sales Manager** | 晨会看 **本部门** Pipeline 与 **部门内成员** 排行 | 谁卡在某个阶段、本团队赢单率 | `department`：同 `user_tenants.department` 成员负责的数据 |
| **Tenant Admin** | 看全集团 Pipeline、**按部门** 业绩排行 | 各事业群贡献、全租户配额 | `all`：全租户；排行 **按部门汇总** |
| **Viewer** | 只读 Pipeline 与 Dashboard | 无创建/阶段变更 | 同 Sales（`self`）或按角色 seed |

### 2.1 场景故事

**场景 A — 线索转商机**  
销售在线索详情点击「转化」→ 创建 Account/Contact → **勾选「同时创建商机」** → 系统生成 Deal（阶段 `qualification`，金额可后填）→ 跳转 `/deals/:id`。

**场景 B — Pipeline 推进**  
销售打开 `/deals` 看板 → 将「北辰物流」从「方案报价」拖到「商务谈判」→ 系统写 `deal.stage_change` 审计 → 看板与漏斗图同步更新。

**场景 C — 租户管理员晨会（全集团）**  
潘总（租户管理员）打开 `/` Dashboard → `data_scope=all` → 看到全租户 KPI、Pipeline 漏斗、配额 Gauge、**部门业绩排行**（灵狐数据、神龙云计算等各事业群本月赢单金额）。

**场景 D — 销售经理晨会（本部门）**  
李婷（灵狐数据 · 销售经理）登录 → `data_scope=department` → 线索/客户/商机/仪表盘 KPI **仅本部门**；团队排行为 **本部门成员**（非全集团）；不可见莫邪互娱等其它部门数据。

**场景 E — 销售代表**  
王磊（销售代表）登录 → `data_scope=self` → 仅本人 Pipeline；**不展示**团队排行区块（API 403 或 FE 按 scope 隐藏）。

---

## 3. 功能需求（User Story + Acceptance Criteria）

### 3.1 模块映射：任务清单 ↔ 本 PRD

| 任务 ID | PRD 章节 | 交付摘要 |
|---------|----------|----------|
| 3.0 | §10 架构师输入 | `phase-3-deals-dashboard-api.md`、Schema 纲要、2b 评审 |
| 3.UI | §4.5、§5.2 | `ChartSparkline`、`ChartGauge` + `/charts` 案例 |
| 3.1 | §4.1–4.2 | Deals CRUD + Pipeline + `/deals` |
| 3.2–3.3 | §4.3 | 阶段金额、赢单率统计 |
| 3.4–3.6 | §4.4–4.5 | Dashboard summary + 全图嵌入 |
| 3.5 | §4.4.3 | 配额 API + Gauge |
| 3.7 | §4.4.4、§5.3.1 | 团队/部门业绩排行 + 部门数据范围 |
| 3.E2E | §8 | E2E 冒烟路径 |

### 4.1 商机 (Deals) — 数据模型

> **字段与枚举以 [phase-3-deals-dashboard-api.md](../api/phase-3-deals-dashboard-api.md) §2–§3 为准**；下表为 PM 语义摘要。

| 字段 | 说明 |
|------|------|
| `title`, `amount`, `currency` | 商机名称与预计金额（默认 CNY，可选 USD） |
| `stage` | Pipeline 阶段，见下表 |
| `probability` | 0–100；`won`=100、`lost`=0 建议 |
| `expected_close_date` | 预计关单日 |
| `owner_id` | 负责人，默认当前用户 |
| `account_id` / `contact_id` / `lead_id` | 关联公司、联系人、来源线索 |
| `lost_reason`, `closed_at` | 终态关单信息 |
| `engagement_score`, `last_activity_at` | 只读，规则同 Phase 2（可选） |

#### 4.1.1 Pipeline 阶段枚举（Accepted 契约）

| Code | 中文默认 | 说明 |
|------|----------|------|
| `qualification` | 需求确认 | 默认阶段 |
| `proposal` | 方案报价 | |
| `negotiation` | 商务谈判 | |
| `won` | 赢单 | **终态**；写 `closed_at` |
| `lost` | 输单 | **终态**；建议 `lost_reason` |

**状态机**：`qualification → proposal → negotiation → won | lost`；允许回退；非法迁移 → `invalid_stage_transition`（详见 phase-3-deals-dashboard-api §2.1）。

#### 4.1.2 User Story — Deals CRUD

| ID | 作为 | 我要 | 以便 | AC |
|----|------|------|------|-----|
| DL-01 | 销售 | 新建商机并关联公司/联系人 | 跟踪具体机会 | 表单 Zod 校验；`deals:create`；写审计 |
| DL-02 | 销售 | 在列表/看板查看我的商机 | 管理 Pipeline | 默认 `owner_id=me`；分页/筛选 stage、account |
| DL-03 | 销售 | 编辑金额与预计关单日 | 更新预测 | `deals:update`；仅 owner 或经理（view_all） |
| DL-04 | 销售 | 删除误建商机 | 保持管道干净 | `deals:delete`；软删或硬删由架构师定（须审计） |
| DL-05 | 销售 | 从线索转化时可选创建商机 | convert Body 增 `create_deal`；写入 `lead_id`；响应含 `deal_id` |
| DL-06 | 销售 | 在商机详情看到来源线索 | 有 `lead_id` 时展示链接 `/leads/:id` |

#### 4.1.3 User Story — Pipeline 看板

| ID | 作为 | 我要 | AC |
|----|------|------|-----|
| PL-01 | 销售 | 在 `/deals` 按阶段分列查看看板 | 每列展示该阶段商机卡片（标题、金额、owner、预计关单） |
| PL-02 | 销售 | 拖拽或菜单变更阶段 | `PUT /api/deals/:id/stage`；乐观 UI 可失败回滚 |
| PL-03 | 经理/管理员 | 查看团队 Pipeline | `department`：本部门各阶段商机；`all`：全租户 |
| PL-04 | 系统 | 看板顶部展示 Pipeline 漏斗 | `GET /api/deals/pipeline` → `ChartFunnel` |

**UI 要求**（FE）：

- 看板 `data-testid="deals-pipeline"`；新建 `deal-create-btn`  
- 空列显示引导文案 + CTA  
- 移动端：看板横向滚动或切换为列表+阶段筛选（Could）

### 4.2 Deals 详情页

| ID | 作为 | 我要 | AC |
|----|------|------|-----|
| DD-01 | 销售 | 查看商机详情 Tab：概览 / 时间线 | 时间线复用 Activity（`subject_type=deal` 或关联 account） |
| DD-02 | 销售 | 快捷关单 | won/lost 按钮 + 确认 Modal |

### 4.3 Deals 分析区（`/deals` 报表 Tab 或分析区块）

| ID | 作为 | 我要 | AC |
|----|------|------|-----|
| DS-01 | 销售 | 看各阶段 **金额合计** 柱图 | `ChartBar` 接 `GET /api/deals/stats/by-stage` |
| DS-02 | 经理 | 看 **赢单率** 趋势 | `ChartLine` 接 `GET /api/deals/stats/win-rate`；按周/月 |
| DS-03 | 系统 | 统计 API 遵守数据范围 | 与列表相同 `ScopeParams`（`self` / `department` / `all`）；集成测覆盖 tenant + scope |

**赢单率定义（MVP）**：

```
win_rate = closed_won_count / (closed_won_count + closed_lost_count)
```

分母为 0 时返回 `null`，前端显示「—」空态。

### 4.4 Dashboard（`/` 工作台升级）

Phase 2 工作台 IA（§11.5）保留；Phase 3 **生产化**以下能力：

#### 4.4.1 汇总 API 替代客户端聚合

| ID | 作为 | 我要 | AC |
|----|------|------|-----|
| DB-01 | 销售 | 首页 KPI 一次加载 | `GET /api/dashboard/summary` 返回 KPI + sparkline 序列 |
| DB-02 | 系统 | 避免 N+1 | 禁止首页对 leads/accounts 全量 list 再算（现状 `use-dashboard` 应重构） |

**`summary` 响应语义**：以 [phase-3-deals-dashboard-api.md §4.1](../api/phase-3-deals-dashboard-api.md) 为准，含 `kpis`、`kpi_trends`、`sparklines`、`priorities`。

#### 4.4.2 Dashboard 图表区（Must — 不得仅数字）

| 图表 | 组件 | 数据 API |
|------|------|----------|
| KPI sparklines | `ChartSparkline` | `summary.sparklines` |
| 赢单率趋势 | `ChartLine` | `GET /api/deals/stats/win-rate` |
| Pipeline 漏斗 | `ChartFunnel` | `GET /api/dashboard/funnel` |
| 配额完成 | `ChartGauge` | `GET /api/dashboard/quota` |
| 团队/部门排行 | `ChartBar` | `GET /api/dashboard/team-ranking`（见 §4.4.4） |
| 阶段金额 | `ChartBar` | `GET /api/deals/stats/by-stage`（Deals 页） |

#### 4.4.3 配额（Quota）

| ID | 作为 | 我要 | AC |
|----|------|------|-----|
| QT-01 | 销售 | 看本周期配额完成 % | Gauge 0–100%；超额可 >100 显示 |
| QT-02 | 经理/管理员 | 看团队汇总配额 | `department` 聚合本部门；`all` 聚合全租户 |

- 粒度：用户或租户（`tenants.config.sales_quota` JSONB）  
- `target_amount` / `won_amount_mtd` / `completion_rate` 见 phase-3-deals-dashboard-api §4.4

#### 4.4.4 数据范围与团队排行（Accepted v0.2）

**解析规则**（BE `datascope.Resolver`，FE 读 `summary.data_scope`）：

| 条件 | 如何判定 | CRM 列表/统计（Leads、Accounts、Contacts、Deals、Dashboard summary/funnel/quota） |
|------|----------|-------------------------------------------------------------------------------------|
| `all` | 用户具备 `rbac:manage`（租户管理员） | 全租户 `owner_id` |
| `department` | 角色名为 **销售经理** / `Sales Manager`，且 `user_tenants.department` 非空 | 仅 `owner_id` 属于 **同部门** 的成员（含 `owner_id IS NULL` 的未分配记录，与 MVP 一致） |
| `self` | 其它角色（销售代表、只读等） | 仅 `owner_id = 当前用户` |

部门归属存 **`user_tenants.department`**（一级部门名称，可与钉钉通讯录对齐）；**设置 → 成员** 展示部门列（Phase 4 交付，见 [phase-4-system-settings-close-prd.md](./phase-4-system-settings-close-prd.md) §3.6）。

**`GET /api/dashboard/team-ranking`**

| `data_scope` | 谁可见排行 | `group_by` | 展示语义 |
|--------------|------------|------------|----------|
| `self` | 否（403 或 FE 隐藏 `dashboard-team-ranking`） | — | — |
| `department` | 销售经理 | `user` | 本部门 **成员** 本月赢单金额 Top N |
| `all` | 租户管理员 | `department` | 全集团 **各部门** 本月赢单金额 Top N |

响应示例字段：

```json
{
  "group_by": "department",
  "items": [
    { "department": "灵狐数据", "name": "灵狐数据", "value": 1200000, "rank": 1 }
  ]
}
```

```json
{
  "group_by": "user",
  "items": [
    { "user_id": "…", "name": "李婷", "value": 800000, "rank": 1 }
  ]
}
```

| 其它 | 要求 |
|------|------|
| Viewer | 只读；无新建 Deal 入口；`self` 时无排行 |
| FE 标题 | `group_by=department` →「部门业绩排行」；`user` →「团队业绩排行（本部门）」 |
| 一致性 | Deals / Leads / Accounts / Contacts / Activity 触达统计 **共用** 同一 `ScopeParams`，禁止经理 Leads 全量、Deals 仅本人 |

**演示**：小西科技集团租户见 [xiaoxi-boss-demo.md](../demo/xiaoxi-boss-demo.md)（`ceo@xiaoxi.com` vs `linghu@xiaoxi.com`）。

#### 4.4.5 Zone E 生产化策略

| 区块 | Phase 2 | Phase 3 |
|------|---------|---------|
| 迷你商机漏斗 | fixtures + Preview 角标 | **真实** `dashboard/funnel`；移除 Preview 角标 |
| 团队关系热力 | fixtures Mock | **Should**：规则版（成员 × 平均 engagement）；**Could** 仍 Mock + 角标至 L2 真算 |

### 4.5 ui-kit 组件（与业务同迭代）

| 组件 | 场景 | DoD |
|------|------|-----|
| **`ChartSparkline`** | KPI 卡内 7 点折线 | `/charts` 案例 + i18n + Vitest 最小覆盖 |
| **`ChartGauge`** | 配额完成率 | 同上 |

复用已有：`ChartFunnel`、`ChartLine`、`ChartBar`、`CardMetric`、`ChartShell`。

### 4.6 AI — 经理「关系降温」名单（Should Have）

接 Phase 2 情绪与 L2 规则，**不做真 ML**。

| ID | 作为 | 我要 | AC |
|----|------|------|-----|
| AI-01 | 经理 | 在 Dashboard 看到「关系降温」列表 | 规则：`engagement_score` 周环比下降 ≥ 阈值 或 连续 negative sentiment |
| AI-02 | 系统 | Preview 与规则分区 | 真规则无角标；Mock 分数带 `AiPreviewBadge` |
| AI-03 | 经理 | 点击跳转客户/商机详情 | 链接 `/leads/:id` 或 `/accounts/:id` |

**Out of Scope（Phase 3）**：真 LLM 预测、自动阶段推荐、RAG。

---

## 5. 多租户与 RBAC 要求

### 5.1 多租户

- 所有 `deals` 及 Dashboard 聚合查询 **强制** `tenant_id` 来自 JWT Context，禁止 body 传 tenant。  
- 跨租户 Deal ID 访问返回 `404`（非 403，防枚举）。  
- Dashboard 汇总 **不得** 泄漏其他租户 KPI。

### 5.2 权限资源

| 资源 | Actions | 说明 |
|------|---------|------|
| `deals` | `view`, `create`, `update`, `delete` | 已 seed（00002） |
| `dashboard` | `view` | 访问 `/api/dashboard/*`（迁移新增） |

### 5.3 数据范围（ABAC）

| 策略 | 适用 | 判定方式（MVP） |
|------|------|-----------------|
| `self` | 销售代表、只读等 | 默认 |
| `department` | 销售经理 | 角色名 + `user_tenants.department` |
| `all` | 租户管理员 | `rbac:manage` |

### 5.3.1 部门模型（MVP）

| 项 | 说明 |
|----|------|
| 存储 | `user_tenants.department`（`VARCHAR`，可空） |
| 与业务标签 | 线索/商机上的 `tags`（事业群标签）**独立**；数据范围按 **负责人所属部门**，不按 tag 过滤 |
| 组织树 | **无** 多级部门表；一级部门字符串即可支撑 MVP |
| 后续 | Phase 5+ 可对接钉钉部门 API 同步；RBAC UI 配置 `data_scope`（Phase 4 RB-08 Should） |

**集成测试必须覆盖**：

- 租户 A 用户不可见租户 B 的 Deal  
- `self` 不可见他人 Pipeline  
- 销售经理 A 部门不可见 B 部门 `owner_id` 的 Deal  
- 租户管理员 `team-ranking` 为 `group_by=department`；销售经理为 `group_by=user` 且仅本部门成员  
- `self` 调用 `team-ranking` → 403

### 5.4 审计

| action | 触发 |
|--------|------|
| `deal.create` / `deal.update` / `deal.delete` | CRUD |
| `deal.stage_change` | 阶段变更（含 from/to） |
| `deal.close_won` / `deal.close_lost` | 关单（stage=won/lost） |
| `deal.convert_from_lead` | convert 时 create_deal |

---

## 6. i18n 与本地化需求

- 所有用户可见文案走 `nuxt/i18n`，**key 英文**。  
- Pipeline 阶段：API 返回 `stage` code + 可选 `stage_label`；前端用 `deals.stage.{code}` 兜底。  
- 金额：默认 `zh-CN` 显示 `¥1,234,567.00`；`en-US` 显示 `CNY 1,234,567.00`（MVP 单币种）。  
- Dashboard KPI label、Gauge 中心文案、空态、错误态双语齐全。

---

## 7. 非功能需求

| 类别 | 要求 |
|------|------|
| **性能** | `GET /api/dashboard/summary` P95 < 300ms（单租户 200 用户、≤5k deals）；看板页首屏 < 1.5s |
| **安全** | 阶段变更、关单须鉴权 + 审计；禁止 IDOR |
| **可用性** | 看板加载/空态/错误态；图表 API 失败单图降级不影响整页 |
| **兼容** | 现有 `/` 工作台布局不推翻；增量接 API |
| **测试** | BE 集成测（CRUD/scope/stats）；FE Vitest；QA E2E `phase3-deals` + Dashboard 冒烟 |

---

## 8. 验收与测试要点（QA 输入）

### 8.1 核心 E2E 路径

1. 登录 → 新建商机 → 看板可见 → 变更阶段 → 漏斗更新  
2. 线索 convert + `create_deal` → Deal 详情有 `lead_id`  
3. Dashboard summary 一次加载（替代多列表聚合）→ KPI 下钻  
4. **租户管理员**：`team-ranking` 可见且 `group_by=department`；**销售经理**（同部门）：`group_by=user` 且仅本部门成员；**销售代表**：无排行、不可见他人 Deal  
5. 配额 Gauge 与本月 won 金额一致（seed 数据下）；经理配额为本部门汇总、管理员为全租户  
6. （小西演示）`ceo@xiaoxi.com` vs `linghu@xiaoxi.com` vs `moye@xiaoxi.com` 数据范围对比

### 8.2 图表 Done（与 task-breakdown 一致）

- [ ] Deals：`ChartFunnel` + `ChartBar` + `ChartLine` 接 API  
- [ ] Dashboard：`ChartSparkline`×KPI + `ChartLine` + `ChartFunnel` + `ChartGauge` + `ChartBar`（经理）  
- [ ] `/charts` 含 Sparkline、Gauge 案例区块

---

## 9. 优先级 (MoSCoW)

| 级别 | Phase 3 范围 |
|------|--------------|
| **Must** | Deals CRUD + Pipeline + stage API + `/deals` 页 |
| **Must** | Deals 统计（by-stage、win-rate、pipeline） |
| **Must** | `GET /api/dashboard/summary` + KPI Sparkline |
| **Must** | Dashboard 至少 4 类图接 API（Sparkline/Line/Funnel/Gauge/Bar 取 4+） |
| **Must** | 多租户 + RBAC + 数据范围集成测 |
| **Must** | 部门数据范围（`self` / `department` / `all`）+ team-ranking（`user` / `department` 分组） |
| **Should** | convert 时 `create_deal`；配额 API |
| **Should** | 关系降温名单（规则版）；Zone E 热力规则化 |
| **Could** | 看板移动端优化；独立 `/today` |
| **Won't** | 多币种、自定义 Pipeline 阶段（Phase 4）、真 ML 赢单预测、审批流 |

---

## 10. 成功指标

| 指标 | 目标 |
|------|------|
| Pipeline 可用性 | 销售可在 3 分钟内完成「建商机 → 改阶段 → 看漏斗」 |
| Dashboard 加载 | summary 一次请求；较 Phase 2 减少 ≥50% 首页 API 调用次数 |
| 数据隔离 | 租户/scope 集成测 100% 通过 |
| 图表验收 | task-breakdown「Phase 3 图表 Done」全勾选 |
| 演示升级 | 演示租户 Zone E 漏斗 **无 Preview 角标**（真实数据） |

---

## 11. 风险与依赖

| 风险 | 影响 | 缓解 |
|------|------|------|
| Phase 2 convert 未带 `create_deal` | 断链 | Phase 3 BE 扩展 convert；FE 转化 Modal 增勾选项 |
| 首页 `use-dashboard` 重构量大 | 延期 | 3.4 优先 summary API；FE 分 PR 替换 |
| 配额无设置 UI | 只能 seed | MVP 文档说明；Phase 4 设置页 |
| ChartSparkline/Gauge 未就绪 | Dashboard 阻塞 | 3.UI 与 3.1 并行，Week 1 必交付 |
| Phase 2.7 导入/分配未完成 | 数据少 | 不挡 Phase 3；QA 用 seed |

**依赖**：

- 架构师：`docs/api/phase-3-deals-dashboard-api.md`（3.0 门禁）  
- Phase 2：`leads`、`accounts`、`activities` 表与权限  
- ui-kit：主题 bridge 已有

---

## 12. 架构师输入清单

> **状态：已完成**。Implementation 以 [phase-3-deals-dashboard-api.md](../api/phase-3-deals-dashboard-api.md) v1.0 Accepted + [04-phase-3-deals-dashboard-schema.md](../architecture/04-phase-3-deals-dashboard-schema.md) 为准，本节保留 PM→架构师原始输入追溯。

| # | 能力 | 契约章节 |
|---|------|----------|
| 1 | Deals CRUD | phase-3-deals-dashboard-api §3.2 |
| 2 | Pipeline | §3.3 |
| 3 | Stage | §3.4 |
| 4 | Stats | §3.5 |
| 5 | Dashboard | §4.1–4.5 |
| 6 | Convert 扩展 | §5 |
| 7 | RBAC / 审计 | §6–§7 |

迁移：`00008_deals.sql`（见 Schema §5）。

---

## 13. 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-25 | PM v0.1 初稿：Deals Pipeline、Dashboard 生产化、MoSCoW；对齐 phase-3-deals-dashboard-api v1.0 Accepted |
| 2026-05-26 | PM v0.2：部门数据范围（租户管理员 `all` + 部门排行；销售经理 `department` + 成员排行）；`user_tenants.department`；对齐实现与 xiaoxi 演示 |
