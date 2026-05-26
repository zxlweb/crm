# CRM 系统 MVP PRD

**产品名称**：EnterpriseFlow CRM（企业流 CRM）  
**版本**：MVP v0.3  
**日期**：2026-05-26（权限与数据范围修订）  
**目标用户**：中小企业销售团队、客服团队  
**Phase 2 详细 PRD**：[phase-2-relationship-crm-prd.md](./phase-2-relationship-crm-prd.md)（关系经营 + AI 演示预留）

---

## 1. 背景与业务目标

随着企业数字化转型，中小企业需要一套**轻量、易用、安全的多租户 CRM**。MVP 之后的产品方向是 **AI 客户关系管理**：从「管流程、记台账」升级为 **「经营关系、理解客户、对话式跟进」**。

**核心业务目标**（MVP 阶段）：

- 实现客户全生命周期管理（含关系阶段与跟进，见 Phase 2 PRD）
- 提升销售团队效率和成交率
- 支持多租户 SaaS 模式，快速服务多个企业客户
- 建立完善的权限与数据隔离体系
- 支持中英文双语

**AI 客户关系目标**（战略层，与 MVP 交付分层）：

| 范式 | 目标 |
|------|------|
| 管理 → **经营关系** | 生命周期 + 关系健康 + 下一步行动 |
| 推送 → **对话** | 跟进时间线 + 语境化建议（非群发） |
| 记录 → **理解** | 行为/情绪信号 → 洞察 → 预测与 Copilot（分期上线） |

> **MVP 现实**：真 LLM / RAG / ML 预测不在 MVP 内一次性交付；但 **产品形态、UI、接口预留、演示数据** 必须在 Phase 2 就绪，用于对内汇报与对外路演。细则见 Phase 2 PRD **§15**。

---

## 2. 用户角色 (Personas)

| 角色 | 描述 | 核心诉求 |
|------|------|----------|
| Tenant Admin（租户管理员） | 企业负责人或 IT 管理员 | 全集团数据、**按部门** 业绩排行；成员/部门/角色管理；租户配置、**AI 能力开关** |
| Sales Manager（销售经理） | 销售团队主管 | **本部门** 团队业绩与成员排行、分配线索、预测性优先名单（演示/后续） |
| Sales（销售人员） | 普通销售 | 跟进线索、**Copilot 建议、情绪旅程** |
| Viewer（只读用户） | 部分查看权限用户 | 查看相关数据 |

---

## 3. MVP 功能范围（核心模块）

### Phase 0–1（已完成 / 关闭中）

- 多租户、RBAC、认证、审计、Super Admin、前端权限组件

### Phase 2 — 客户与线索 + **AI 演示壳**（进行中）

**业务 Must**（传统 CRM 主链）：

- 公司（Accounts）、联系人（Contacts）、线索（Leads）CRUD
- 线索导入、分配、状态流转、Activity 跟进
- Leads 报表（来源、趋势、漏斗）

**关系经营 Must / Should**（见 [phase-2-relationship-crm-prd.md](./phase-2-relationship-crm-prd.md)）：

- 生命周期阶段、规则洞察、情绪标注、**情绪旅程地图**
- **AI Preview 壳**：Copilot 侧栏、预测洞察卡、智能摘要按钮、演示数据（**可 Mock**）

### Phase 3 — 商机与仪表盘

**详细 PRD**：[phase-3-deals-dashboard-prd.md](./phase-3-deals-dashboard-prd.md)

- 商机 (Deals) CRUD + Pipeline 看板 + 阶段推进
- 个人/经理/管理员 Dashboard（KPI + Sparkline/Line/Funnel/Gauge/Bar 接 API）
- **数据范围**：`self`（本人）/ `department`（同 `user_tenants.department` 成员数据，**不含** `owner_id` 空记录）/ `all`（租户管理员 `rbac:manage`）
- **分角色工作台**：销售经理精简 KPI + 事业部配额；管理员部门排行 + Zone E；销售代表全个人视图（见 Phase 3 §4.4.6）
- **事业部配额**：`tenants.config.department_quotas`（经理）与 `sales_quota`（管理员）
- 首页 Zone E：**仅** `data_scope=all` 展示；一线角色隐藏演示样例区
- AI（Should）：经理「关系降温」名单（规则版，接 Phase 2 情绪与 engagement）

### Phase 4 — 系统设置与收尾

**详细 PRD**：[phase-4-system-settings-close-prd.md](./phase-4-system-settings-close-prd.md)

- 租户配置、自定义字段、**系统设置内角色与权限**、**成员列表按数据范围过滤**、**登录/切换租户与会话**（AUTH-01–04）、i18n、Swagger、审计报表、部署文档闭环

### Phase 5+ — **AI 生产能力**

- LLM 自动情绪、NLP 字段抽取、RAG 知识应答、个性旅程编排、多渠道触达

### 不在 MVP 一次性交付、但 Phase 2 必须「看得见」的 AI 能力

| 能力 | MVP 交付形态 | 真能力阶段 |
|------|--------------|------------|
| 情绪旅程地图 | UI + Mock/规则数据 | L1 规则 → L3 LLM |
| 关系洞察 / Next Best Action | 规则 + Preview 文案 | L2 预测 |
| Copilot（摘要、话术、邮件草稿） | UI + Mock 或 501 | L3 |
| 流失/复购概率 | Preview 分数 + 说明条 | L2 |
| RAG 对话助手 | 详情页「智能问答」占位 Tab | L3 |

---

## 4. 非功能需求

- **多租户数据隔离**：严格隔离；演示数据仅 `preview` 模式，禁止写入生产表
- **性能**：单租户 ≤ 200 用户时，页面加载 < 1.5s；Preview 走静态 JSON 不阻塞主链
- **安全**：JWT、密码加密、操作审计；AI 调用需 `ai_enabled` + 审计（真 AI 阶段）
- **国际化**：完整中英文；演示文案双语
- **技术约束**：Nuxt 3 + Go + PostgreSQL
- **可扩展性**：`tenant.config.ai_*`、`X-CRM-Preview`、统一 `AiCapability` 响应包（契约由架构师落 `docs/api/`）

---

## 5. 用户故事示例（部分）

**线索管理**

- 作为销售，我可以新建线索，填写基本信息和跟进记录，并标注客户情绪
- 作为销售经理，我可以查看 **本部门** 的线索/商机与 **部门内成员** 业绩排行
- 作为租户管理员，我可以查看 **全集团** 数据与各 **部门** 业绩排行，并在设置中查看成员所属部门
- 作为销售，我可以在详情页看到 **AI 建议的下一步**（演示期为样例数据）

**AI 演示（Phase 2 Preview）**

- 作为 **老板**，我打开演示账号，在 3 分钟内看到：情绪旅程、洞察卡片、Copilot 生成跟进话术（Mock）
- 作为管理员，我可以关闭租户 AI 演示，仅显示已上线的规则能力

**权限与数据范围**

- 作为管理员，我可以创建角色并分配权限；查看本租户**全部**成员
- 作为销售经理，我只能看到**本事业部** CRM 数据、成员排行与成员列表（`data_scope=department`）
- 作为销售代表，我只能看到自己负责的线索和商机（`data_scope=self`）
- 作为任意用户，登录后 JWT 租户与当前租户一致，避免残留 cookie 导致 403（见 Phase 4 §3.7）

---

## 6. 优先级 (MoSCoW)

| 级别 | 范围 |
|------|------|
| **Must Have** | 认证 + 多租户 + RBAC + Contacts + Leads + Deals + Dashboard + Phase 2 CRUD 主链 |
| **Should Have** | 自定义字段、导入导出、Activity、情绪旅程、**AI Preview UI 全套（Mock）** |
| **Could Have** | 邮件通知、今日待跟进聚合页 |
| **Won't Have（MVP 真交付）** | 生产级 LLM、RAG 知识库、ML 模型、自动营销群发、移动 App |
| **Won't Have 但 Phase 2 必须有壳** | 见上表「看得见」— 在 [phase-2-relationship-crm-prd.md](./phase-2-relationship-crm-prd.md) §15 验收 |

> 原 v0.1「Won't Have：AI 智能分析」已废止，改为 **「Won't：生产 AI；Must/Should：AI 产品与接口预留 + 演示」**。

---

## 7. 成功指标 (OKRs)

**Objective 1**：验证产品市场可用性（多租户 CRM 主链）

- 完成 MVP 开发并成功部署
- ≥3 租户同时使用且数据隔离正常
- 核心页面响应 < 2s；权限无越权

**Objective 2**：验证 AI CRM 叙事可被决策者理解（Phase 2 Preview）

- 演示路径 ≤ 3 分钟：登录 → Lead 详情 → 情绪旅程 → Copilot 采纳建议（Mock）
- 决策者访谈：能复述「三转变」（经营关系 / 对话 / 理解）≥ 2/3 人
- Preview 与生产数据隔离 100%（无 Mock 写入生产表）

---

## 8. 文档索引（勿散落新 md）

| 文档 | 职责 |
|------|------|
| **本文件** | MVP 总览、AI 战略、MoSCoW |
| [phase-2-relationship-crm-prd.md](./phase-2-relationship-crm-prd.md) | Phase 2 全量需求 + AI Preview + 架构师输入清单 |
| [phase-3-deals-dashboard-prd.md](./phase-3-deals-dashboard-prd.md) | Phase 3 商机 Pipeline + Dashboard 生产化 |
| [phase-4-system-settings-close-prd.md](./phase-4-system-settings-close-prd.md) | Phase 4 系统设置、角色权限、自定义字段、审计报表、文档收尾 |
| `docs/tasks/00-mvp-task-breakdown.md` | 任务勾选 |
| `docs/api/*` | **架构师** API 契约（PM 不写） |
| `docs/meeting-notes/phase-N-notes.md` | 阶段验收与汇报纪要 |

---

## 9. 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-21 | MVP v0.1 初稿 |
| 2026-05-22 | v0.2：AI 客户关系战略；废止「MVP 不做 AI」；链 Phase 2 PRD §15 演示预留 |
| 2026-05-25 | 链 Phase 3 详细 PRD（Deals + Dashboard） |
| 2026-05-26 | 链 Phase 4 详细 PRD（系统设置与收尾） |
| 2026-05-26 | Phase 4 PRD v0.2：补充系统设置内角色与权限配置 |
| 2026-05-26 | MVP v0.3：同步 Phase 3 v0.3（分角色工作台、事业部配额、Zone E 可见性）与 Phase 4 v0.4（成员 scope、登录/租户会话） |
