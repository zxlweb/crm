# ADR-0004: Phase 2 规则洞察引擎与 AI Preview 分离

**状态**：Accepted  
**日期**：2026-05-22  
**决策者**：系统架构师

## 背景（Context）

Phase 2 PRD 要求同时交付 **L1 生产能力**（时间线、规则洞察、情绪旅程）与 **P0 AI Preview**（老板演示、Mock Copilot）。若将 Preview 与生产洞察共用同一端点且无 `source` 区分，会导致：

1. 销售误把 Mock 概率当作真实预测；
2. FE/BE 在 Phase 3 接真模型时需推倒响应结构；
3. 租户未开 `ai_enabled` 时仍可能暴露未鉴权 Copilot。

## 决策（Decision）

1. **洞察双通道**  
   - **生产**：`POST /api/{entity}/:id/insights/evaluate` — 仅 L1 规则（`source=rule`），不返回 ML 分数。  
   - **Preview**：`POST /api/insights/preview` — 返回 `AiCapabilityResult`，`source=preview`，需 `ai_enabled` + 角标/i18n disclaimer。

2. **Copilot / 预测 / 动态分群**  
   - 路径独立：`/api/copilot/*`、`/api/scores/predict`、`/api/segments/dynamic/refresh`。  
   - Phase 2 默认 **HTTP 501** + `code: AI_NOT_READY`，或 `X-CRM-Preview: 1` / `ai_preview_mode=fixtures` 返回 fixture。  
   - **禁止** Preview 响应写入 `leads` / `activities` 表。

3. **统一包装**  
   - 所有非 CRUD AI 能力使用 `AiCapabilityResult`（`source`, `capability`, `payload`, `disclaimer`, `request_id`）。  
   - 业务错误码：`AI_DISABLED` (403)、`AI_PREVIEW_ONLY` (200)、`AI_NOT_READY` (501)。

4. **租户开关**  
   - `tenants.config.ai_enabled` 默认 `false`；`ai_preview_mode` 仅 `off|fixtures|live`，Phase 2 实现 `fixtures`。  
   - 演示租户种子单独迁移，不与生产租户混用。

5. **规则引擎位置**  
   - 仅在 **Go Service 层**求值；规则 ID 常量化（INS-001～006）；阈值来自 `tenant.config.insight_thresholds`。  
   - 前端只展示结果与 `explanation_key`，不传规则 DSL。

6. **情绪数据来源**  
   - Phase 2 允许 `manual` + `rule`（关键词）；`sentiment_source=ai` 字段预留但 **禁止写入**。  
   - 情绪旅程与 Activity 列表 **同一聚合服务**，禁止两套序列。

## 备选方案（Options Considered）

| 方案 | 优点 | 缺点 |
|------|------|------|
| A. 单一 `/api/ai/invoke` 网关 | 扩展统一 | Phase 2 过度设计；RBAC 难映射 |
| B. FE 纯 Mock 无后端契约 | 最快 | 演示与生产环境不一致；违背 workflow |
| **C. 分路径 + AiCapabilityResult（采纳）** | 清晰分层、可并行三轨 | 端点数量略多 |

## 后果（Consequences）

### 正面

- BE/FE/QA 可依据 [phase-2-crm-ai.md](../../api/phase-2-crm-ai.md) 并行开工。  
- Preview 与生产洞察 UI 可分区展示（PRD §15.4）。  
- Phase 3+ 接模型时仅替换 Copilot/Preview 实现，不改 CRUD。

### 负面 / 风险

- 端点与权限表需同步维护（`copilot:use`, `insights:view`）。  
- 演示租户数据需单独 seed，避免污染集成测。

### 后续行动

- [ ] 迁移 `00006` + 权限种子  
- [ ] `route.go` 对 `assign`、`insights` 子路径如需特殊 action 再扩展（当前首段 resource 即可）  
- [ ] Phase 3 ADR：真实 `scores.predict` 模型服务边界

## 关联文档

- PRD: [phase-2-relationship-crm-prd.md](../../prd/phase-2-relationship-crm-prd.md) §15  
- API: [phase-2-crm-ai.md](../../api/phase-2-crm-ai.md)  
- Schema: [03-phase-2-crm-schema.md](../03-phase-2-crm-schema.md)  
- Tasks: [00-mvp-task-breakdown.md](../../tasks/00-mvp-task-breakdown.md) Phase 2
