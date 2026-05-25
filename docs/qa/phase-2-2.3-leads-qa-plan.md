# QA 测试计划 — Phase 2.3 Leads CRUD + 状态

**子任务 ID**：`2.3`（`[QA]` 轨）  
**测试日期**：待执行  
**测试人**：待填  
**输入**：[phase-2-relationship-crm-prd.md](../prd/phase-2-relationship-crm-prd.md) §4.5、§5 · [phase-2-crm-ai.md](../api/phase-2-crm-ai.md) §6 · [phase-2-notes.md](../meeting-notes/phase-2-notes.md) §2.5  
**关联任务**：[00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md) 表行 `2.3`

---

## 1. 测试范围

### 1.1 本计划覆盖（2.3 Done 条件）

| 类别 | 内容 |
|------|------|
| **API 集成测** | `GET/POST /api/leads`、`GET/PUT/PATCH/DELETE /api/leads/:id`；Query 筛选；状态合法/非法迁移；`POST .../convert`（最小路径）；租户隔离；RBAC `leads:*` |
| **E2E（冒烟）** | 登录 → `/leads` 列表 → 新建线索 → 详情可见 → 改状态 → 列表/详情一致 |
| **契约** | 响应包 `{ code, message, data, pagination? }`；禁止 body 传 `tenant_id`；只读字段不可写 |
| **审计** | `lead.convert`（convert 用例）；状态/owner 变更建议抽 1 条查 `audit_logs` |

### 1.2 不在 2.3（由后续子任务领）

| 子任务 | 内容 | QA 文档 |
|--------|------|---------|
| 2.4–2.6 | `GET /api/leads/stats/*`、报表 Tab、Chart 数据 | 本计划 §8 仅列依赖用例 ID |
| 2.7 | `import`、`assign` | `L-INT-07x` / `L-E2E-07` 占位 |
| 2.8+ | Activity 时间线、洞察、情绪旅程、AI Preview | `2.E2E` / `2.13` 专项计划 |
| 2.10 | 分群下拉 + URL | 列表 `segment` Query 可在 2.3 做 API 单测，UI 归 2.10 |

### 1.3 前置与环境

| 项 | 要求 |
|----|------|
| BE | 迁移 `00006` + Leads 路由注册完成 |
| FE | `/leads`、`/leads/:id` 可访问；`data-testid` 见 §6 |
| 账号 | Demo 租户 `admin@demo.com` / `password123`（Super Admin，绕过 Casbin）；**另需** 种子 `Sales` + `Manager` 各一（见 §3） |
| 命令 | `make backend-test`；`cd e2e && npm test -- tests/phase2-leads.spec.ts`（待建） |
| 技术债 | [phase-1-qa.md](./phase-1-qa.md) Bug #3：新租户 Casbin — **不阻塞** 2.3 单接口 CRUD，但 **阻塞** Sales/Manager 越权矩阵全绿 |

---

## 2. PRD / 契约 → 验收追溯

| PRD / 契约要点 | 验收标准摘要 | 测试 ID |
|----------------|--------------|---------|
| §4.5 Leads CRUD | 创建/读/改/删/搜索；关系字段可读 | L-INT-01～05 |
| §4.5 状态机 | `new→contacted→qualified→unqualified\|converted`；非法迁移 400 | L-INT-10～12 |
| §4.5 `converted` | 须关联 Account/Contact；审计 | L-INT-13 |
| §2.3 只读字段 | `engagement_score`、`relationship_health`、`last_activity_at` 不可 POST/PATCH 写入 | L-INT-06 |
| §5 RBAC | Sales 本人；Manager `view_all`；`assign` 仅经理 | L-INT-20～23 |
| §5 租户 | 强制 `tenant_id`；跨租户 404/403 | L-INT-30～31 |
| §7 性能 | 列表默认 20 条；详情时间线首屏（2.8 后） | L-PERF-01 |
| Must E2E（§8） | Leads CRUD + 权限 + 1 条 Activity 洞察联动 | **2.3 仅 CRUD+权限**；洞察归 `2.E2E` |
| 任务 2.3 增强 | 洞察侧栏、状态变更可选关联 Activity | UI 冒烟占位 `L-E2E-04`（FE 就绪后启用） |
| 00-overview 用户故事 | 销售新建线索；经理分配/漏斗 | 分配归 2.7；漏斗归 2.6 |

---

## 3. 测试数据与角色

### 3.1 建议种子（BE 迁移或测试 helper）

| 角色 | 邮箱建议 | `data_scope` | `leads` 权限 |
|------|----------|--------------|--------------|
| `sales_a` | `sales-a@demo.com` | `self` | view, create, update（无 assign） |
| `sales_b` | `sales-b@demo.com` | `self` | 同上 |
| `manager` | `manager@demo.com` | `all` | view, create, update, assign |
| `viewer` | `viewer@demo.com` | `self` | view only |

每个用户至少 1 条 **本人** `owner_id` 线索；`sales_b` 线索对 `sales_a` 不可见。

### 3.2 标准 Lead 载荷（集成测）

```json
{
  "title": "QA Lead 华创科技",
  "status": "new",
  "source": "web",
  "amount": 50000,
  "lifecycle_stage": "acquire",
  "tags": ["qa", "phase2"]
}
```

### 3.3 状态迁移矩阵（集成测必跑）

| 从 \ 到 | contacted | qualified | unqualified | converted | new |
|---------|-----------|-----------|-------------|-----------|-----|
| **new** | ✅ 允许 | ⬜ 契约未写清则按实现 | ⬜ | ❌ | — |
| **contacted** | — | ✅ | ✅ | ❌ | ❌ 回退 |
| **qualified** | ❌ | — | ✅ | ✅（+ convert body） | ❌ |
| **unqualified** | ❌ | ❌ | — | ❌ | ❌ |
| **converted** | ❌ 终态 | ❌ | ❌ | — | ❌ |

非法迁移统一断言：`HTTP 400`，`message` 含 `invalid_status_transition`（[phase-2-crm-ai.md](../api/phase-2-crm-ai.md) §6.2）。

---

## 4. API 集成测试用例

**建议文件**：`backend/internal/interfaces/http/leads_test.go`（或 `internal/application/leads/service_test.go` + HTTP 薄层）

**运行**：`cd backend && go test ./internal/interfaces/http/... -run Leads -count=1`  
**模式**：httptest + 测试 DB（与 `scope_test` 同库）或 testcontainers；请求头带 `Authorization` + `X-Tenant-ID`。

| ID | 场景 | 步骤 | 预期 |
|----|------|------|------|
| L-INT-01 | 创建 | `POST /api/leads` 合法 body | `201/200`，`data.id` UUID，`tenant_id` 为当前租户，`owner_id` 默认当前用户 |
| L-INT-02 | 列表分页 | `GET /api/leads?page=1&page_size=20` | `pagination` 存在；仅本租户 |
| L-INT-03 | 搜索 | `GET /api/leads?search=华创` | 命中 title 子串 |
| L-INT-04 | 筛选 | `?status=new&source=web&lifecycle_stage=acquire` | 结果均匹配 Query |
| L-INT-05 | 详情/更新/软删 | `GET` → `PATCH` title → `DELETE` → `GET` | 更新生效；删除后列表不可见或 404 |
| L-INT-06 | 只读字段 | PATCH `engagement_score: 99` | 忽略或 400；DB 仍为系统值 |
| L-INT-10 | 合法迁移 | new→contacted→qualified | 每次 200，`status` 递增 |
| L-INT-11 | 非法迁移 | qualified→new | 400 + `invalid_status_transition` |
| L-INT-12 | 终态 | converted 后再 PATCH status | 400 |
| L-INT-13 | 转化 | qualified + `POST .../convert` 带 `create_account` | `status=converted`，`converted_account_id` 非空；`audit_logs.action=lead.convert` |
| L-INT-14 | convert 缺关联 | converted 无 account/contact | 400（契约） |
| L-INT-20 | Sales 越权读 | A 的 token 访问 B 的 `GET /api/leads/:id` | 403 或 404（不泄露存在性） |
| L-INT-21 | Sales 越权改 | A 改 B 的 lead | 403 |
| L-INT-22 | 无 create | viewer `POST /api/leads` | 403 |
| L-INT-23 | Manager 全量 | manager `GET /api/leads` | 含 A、B 线索条数 ≥ 2 |
| L-INT-30 | 跨租户 | Tenant-A token + Tenant-B 的 `X-Tenant-ID` | 403 或空列表 |
| L-INT-31 | 无租户头 | 缺 `X-Tenant-ID` | 401/400 |
| L-INT-40 | 参数校验 | `title` 空、`status=invalid_enum` | 400 |
| L-INT-41 | SQL 注入 | `search=' OR 1=1--` | 无异常泄漏；结果为空或安全转义 |
| L-INT-42 | XSS 存储 | title `<script>alert(1)</script>` | 存原文；列表 API 不执行（FE 转义归 2.E2E 可选） |

**列表 Query `segment`**（若 2.3 已接筛选 DSL）：增加 L-INT-04b，与 `GET /api/segments/:code/count` 结果一致（完整分群归 2.10）。

---

## 5. E2E 测试用例（Playwright）

**建议文件**：`e2e/tests/phase2-leads.spec.ts`  
**前置**：`playwright.config.ts` 已起 backend + web；使用 **Demo 或 seeds** 登录（优先 `manager@demo.com` 测全列表）。

| ID | 场景 | 步骤 | 预期 | testid / 路由 |
|----|------|------|------|----------------|
| L-E2E-01 | 进入列表 | 登录 → `/leads` | 表格或空态；无 500 | — |
| L-E2E-02 | 新建线索 | 点 `lead-create-btn` → 填标题/来源 → 保存 | 列表出现新行；URL 可进详情 | `lead-create-btn` |
| L-E2E-03 | 状态流转 | 详情改状态 contacted → qualified | UI 与刷新后 API 一致 | 状态控件 testid 待 FE：`lead-status-select` |
| L-E2E-04 | 侧栏占位 | 打开 `/leads/:id` | 宽屏可见洞察区或「暂无洞察」（2.9 前允许空） | `ai-relation-panel`（可选） |
| L-E2E-05 | 权限 UI | `sales_a` 登录，不见他人线索名 | 列表无 B 的 title | — |
| L-E2E-06 | 报表 Tab 不回归 | 点报表 Tab | 2.3 允许 skeleton/空态；**不断言图表数据** | `leads-tab-reports`（建议 FE 补充） |

**同步检查点**（[phase-2-notes.md](../meeting-notes/phase-2-notes.md)）：FE 提交 `lead-create-btn` 等 testid 后 QA 再合入 `L-E2E-02`。

---

## 6. `data-testid` 约定（与 FE 对齐）

| 元素 | testid | 状态 |
|------|--------|------|
| 新建线索 | `lead-create-btn` | 已约定 |
| 情绪 Tab | `tab-emotion-journey` | 2.12 |
| AI 侧栏 | `ai-relation-panel` | 2.13 |
| 采纳建议 | `insight-adopt-btn` | 2.9 |
| **建议 2.3 新增** | `lead-list-table`、`lead-form-title`、`lead-form-submit`、`lead-status-select`、`leads-tab-list`、`leads-tab-reports` | 待 FE |

---

## 7. 安全 / i18n / 非功能（抽测）

| ID | 类型 | 说明 |
|----|------|------|
| L-SEC-01 | 认证 | 无 token 访问 `/api/leads` → 401 |
| L-SEC-02 | 租户 | 见 L-INT-30/31 |
| L-I18N-01 | 文案 | `/leads` 切换 zh-CN / en-US，表头与空态键存在（无裸 key） |
| L-PERF-01 | 列表 | `GET /api/leads` P95 &lt; 500ms（本地 20 条种子） |

---

## 8. 统计 API 用例（2.4–2.6，已实现）

| ID | 子任务 | 简述 | 自动化 |
|----|--------|------|--------|
| L-STAT-01 | 2.4 | `GET /api/leads/stats/by-source` 与列表同源权限 / total 一致 | `TestLeadsHTTP_Stats_LSTAT01_*` |
| L-STAT-02 | 2.5 | `stats/trend` 返回 categories + series | `TestLeadsHTTP_Stats_LSTAT02_*` |
| L-STAT-03 | 2.6 | `stats/funnel` 各阶段 count 之和 = 列表 total | `TestLeadsHTTP_Stats_LSTAT03_*` |
| L-STAT-04 | 2.4 | `by-status` total 与列表一致 | `TestLeadsHTTP_Stats_ByStatusMatchesList` |
| L-STAT-05 | 2.4 | 跨租户 stats total=0 | `TestLeadsHTTP_Stats_TenantIsolation` |
| L-STAT-06 | 2.4 | 非法日期 `to < from` → 400 | `TestLeadsHTTP_Stats_InvalidDateRange` |
| L-STAT-07 | 2.4 | 无 `leads:view` → 403（stats 与 list 一致） | `TestLeadsHTTP_Stats_RequiresLeadsView` |
| L-STAT-08 | 2.4 | 仅 `view` 可访问 stats | `TestLeadsHTTP_Stats_ViewerCanAccess` |
| L-STAT-09 | 2.4 | owner scope：stats 与 list 同范围 | `TestLeadsHTTP_Stats_OwnerScopeMatchesList` |

## 8b. 下游依赖用例（非 2.3/2.4–2.6，仅登记）

| ID | 子任务 | 简述 |
|----|--------|------|
| L-INT-07a | 2.7 | `POST /api/leads/import` multipart |
| L-INT-07b | 2.7 | `POST .../assign` + system Activity |
| L-ACT-01 | 2.8 | `POST /api/activities` 后 `last_activity_at` 更新 |
| L-INS-01 | 2.9 | `POST .../insights/evaluate` INS-001 命中 |
| L-E2E-DEMO | 2.13 / 2.E2E | §15.2 三分钟路径 |

---

## 9. 执行顺序与入口准则

```mermaid
flowchart LR
  A[BE 2.3 API 可用] --> B[L-INT-01～13]
  B --> C[L-INT-20～31 RBAC]
  C --> D[FE testid 稳定]
  D --> E[L-E2E-01～03]
  E --> F[勾选任务表 2.3 QA]
```

| 入口 | 条件 |
|------|------|
| 开始集成测 | `POST /api/leads` 在 Swagger/本地可 200 |
| 开始 E2E | `lead-create-btn` 已合入 + 登录种子可用 |
| **2.3 QA 通过** | §4 全部 **Must** ID 绿；§5 L-E2E-01～03 绿；无 P0/P1 开放 Bug |

---

## 10. 自动化落点（实现时创建）

| 类型 | 路径 |
|------|------|
| Go 集成测 | `backend/internal/interfaces/http/leads_integration_test.go`、`leads_stats_integration_test.go` |
| 域/状态机单测 | `backend/internal/application/leads/status_test.go` |
| Playwright | `e2e/tests/phase2-leads.spec.ts` |
| 测试辅助登录 | `e2e/helpers/auth.ts`（复用 phase1 模式） |

**Makefile**：`cd backend && go test ./internal/interfaces/http/... -run 'LeadsHTTP'`（含 Stats）。

---

## 11. Bug 记录

| 序号 | 描述 | 严重程度 | 状态 |
|------|------|----------|------|
| — | （执行后填写） | | |

---

## 12. 测试结论

**状态**：⬜ 未执行（计划已评审）

**正式关闭 2.3 `[QA]` 行条件**：

- [ ] 本文档 §4 Must 用例已实现且 CI/本地通过  
- [ ] §5 L-E2E-01～03 通过  
- [ ] [00-mvp-task-breakdown.md](../tasks/00-mvp-task-breakdown.md) 表 `2.3` QA 列改 `✅`  
- [ ] 无 P0/P1；P2 记入 §11 并关联 BE/FE 子任务  

**说明**：统计图（2.4–2.6）、导入分配（2.7）、洞察/情绪/演示（2.8–2.13、2.E2E）在各自 QA 计划中验收，不重复计入 2.3 通过条件。

---

## 修订记录

| 日期 | 说明 |
|------|------|
| 2026-05-22 | 首版：2.3 Leads 集成测 + E2E 计划，对齐 PRD §4.5/§5 与 API 契约 §6 |
