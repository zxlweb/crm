# Phase 2 阶段笔记 — 关系经营 / AI Preview

**日期**：2026-05-22  
**状态**：Architect 阶段完成 → **可开启三轨 Implementation**

---

## 1. 架构产出（2a）

| 文档 | 路径 |
|------|------|
| API 契约 v1.0 | [docs/api/phase-2-crm-ai.md](../api/phase-2-crm-ai.md) |
| 数据模型 / 迁移纲要 | [docs/architecture/03-phase-2-crm-schema.md](../architecture/03-phase-2-crm-schema.md) |
| ADR-0004 规则洞察 + AI Preview | [docs/architecture/adr/0004-rule-insights-and-ai-preview.md](../architecture/adr/0004-rule-insights-and-ai-preview.md) |

**评审结论**：契约已覆盖 PRD §15.3 能力清单；BE 按 §15 实现顺序领 `2.1`；FE/QA 可对契约开发。

---

## 2. 前端切面（2b）

### 2.1 路由（Nuxt `apps/web/pages/`）

| 路由 | 页面 | layout |
|------|------|--------|
| `/accounts` | 列表 + `segment` Query | `admin` |
| `/accounts/[id]` | Tab：概览 / 时间线 / 情绪旅程 / 洞察 | `admin` |
| `/contacts` | 列表 | `admin` |
| `/contacts/[id]` | 同 accounts 详情结构 | `admin` |
| `/leads` | 列表 + **报表 Tab** | `admin` |
| `/leads/[id]` | 详情 + **右侧 AI 侧栏**（宽屏） | `admin` |

权限：`usePermission()` — `accounts:view` 等；AI 侧栏另判 `tenant.config.ai_enabled`（`use-tenant` 扩展读 config）。

### 2.2 Composables

| 名称 | 职责 |
|------|------|
| `use-accounts` | CRUD + 列表筛选 |
| `use-contacts` | CRUD + `account_id` |
| `use-leads` | CRUD + convert + assign + import |
| `use-activities` | 时间线 CRUD + summary |
| `use-insights` | `POST .../insights/evaluate` + adopt 预填 |
| `use-emotion-journey` | `GET .../emotion-journey` |
| `use-lead-stats` | 四个 stats 端点 → Chart 数据 |
| `use-segments` | 分群列表 + count |
| `use-ai-preview` | `ai_enabled` / fixtures / `?preview=1` |
| `use-ai-copilot` | Copilot 适配层（501 → fixture） |

### 2.3 组件落点

| 组件 | 路径 |
|------|------|
| `ActivityTimeline` | `components/feature/crm/activity-timeline.vue` |
| `InsightCard`, `LifecycleBadge` | `components/feature/crm/` |
| `EmotionJourneyMap` | `components/feature/crm/emotion-journey-map.vue` |
| `AiRelationPanel`, `AiCopilotDrawer`, `AiPreviewBadge` | `components/feature/ai/` |
| `ChartDonut`（新增） | `packages/ui-kit/.../chart-donut.vue` |

Fixtures：`apps/web/fixtures/ai-preview/*.json`（与契约 payload 同形）。

### 2.4 数据流

```
API → use-* (auth + tenant + permission)
        → feature/crm/* + feature/ai/*
        → @crm/ui-kit Chart* / Card*
insights/evaluate (rule) 与 insights/preview (preview) 分区渲染，禁止混排无角标 Mock
```

### 2.5 `data-testid`（QA 约定）

| 元素 | testid |
|------|--------|
| Leads 新建 | `lead-create-btn` |
| 详情情绪 Tab | `tab-emotion-journey` |
| AI 侧栏 | `ai-relation-panel` |
| Preview 角标 | `ai-preview-badge` |
| 采纳建议 | `insight-adopt-btn` |

---

## 3. 三轨领任务指引

复制到各对话框首条消息：

- **BE**：`【BE】Phase 2 — 2.1 Accounts API`，读 `docs/api/phase-2-crm-ai.md` §4 + `03-phase-2-crm-schema.md`，只改 `backend/`
- **FE**：`【FE】Phase 2 — 2.3 Leads 页骨架`，契约/mock，只改 `apps/web/`、`packages/ui-kit/`
- **QA**：`【QA】Phase 2 — 2.3 Leads 集成测计划`，读 PRD 验收标准，只改测试目录

**同步检查点**：2.1 契约 API 可用 → FE 换真实 API；2.3 路由 testid 稳定 → QA 补 E2E。

**QA 计划**：[phase-2-2.3-leads-qa-plan.md](../qa/phase-2-2.3-leads-qa-plan.md)（集成测 + E2E 用例追溯 PRD）

---

## 4. 阻塞 / 风险

| 项 | 说明 |
|----|------|
| Casbin 新租户 Enforce | Phase 1 QA 遗留，不阻塞 2.3，但 2.7 assign 前需验证 |
| `ChartDonut` | ui-kit 未实现，2.4 当周补齐 |

---

## 5. 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-22 | Architect 2a/2b 完成，开放 Implementation |
