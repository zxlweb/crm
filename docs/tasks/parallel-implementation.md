# Implementation 三轨并行领任务指南

> 与 [.cursor/rules/development-workflow.mdc](../../.cursor/rules/development-workflow.mdc) §3 配套。用于 **3 个 AI 对话框**（或 3 位开发者）同时推进同一 Phase。

---

## 1. 何时可以并行

| 条件 | 说明 |
|------|------|
| API 契约 | `docs/api/` 或 PRD 附录已写明路径、请求/响应、错误码、租户与 RBAC |
| 前端切面 | 路由、`use-*` 名、feature/ui-kit 落点（§2b） |
| 任务已拆轨 | 本 Phase 在 [00-mvp-task-breakdown.md](./00-mvp-task-breakdown.md) 有 `[BE]` / `[FE]` / `[QA]` 列或子表 |

缺任一项 → 先完成 Architect 阶段，**不要**三轨硬并行。

---

## 2. 轨道与规则文件

| 轨道 | 首条消息前缀 | Cursor 规则（建议 @） |
|------|--------------|------------------------|
| Backend | `【BE】` | `backend-developer.mdc` + `development-workflow.mdc` |
| Frontend | `【FE】` | `frontend-developer.mdc` + `development-workflow.mdc`（架构/design 用 `frontend-architect.mdc`） |
| Test | `【QA】` | `qa-tester.mdc` + `development-workflow.mdc` |

---

## 3. 领任务与状态标记

在 `docs/tasks/` 对应勾选行使用 **状态 + 轨道**（可多行拆 BE/FE/QA）：

```markdown
- [x] **[BE] 2.3** Leads CRUD API + 迁移 — ✅
- [x] **[FE] 2.3** Leads API 接入 + Nuxt UI + `layout: app` — ✅
- [ ] **[QA] 2.3** Leads 集成测 + E2E 创建线索 — ⬜ 待领
```

| 标记 | 含义 |
|------|------|
| `⬜ 待领` | 无人做，可领 |
| `🟡 进行中 (...)` | 已被某轨占用，**其他对话框勿领同一 ID** |
| `✅ 完成` | 本轨交付完成（含该轨要求的自动化测试） |
| `🔴 BLOCKED (...)` | 依赖他轨；在 `docs/meeting-notes/phase-N-notes.md` 写一条阻塞说明 |

**子任务 ID**：与 Phase 表一致（如 `2.3`、`2.4`），便于三轨对话对齐。

### 3.1 任务表更新边界（防越轨）

`00-mvp-task-breakdown.md` 中 Phase 并行表列为 **`[BE] | [FE] | [QA]`** 时：

| 规则 | 说明 |
|------|------|
| 只改本列 | `【BE】` 对话仅更新 `[BE]` 列的 `⬜/🟡/✅`；FE、QA 同理 |
| 不改他轨 | 禁止替 FE 标「页面完成」、禁止替 QA 标「E2E 完成」 |
| 业务汇总 | 底部 `- [ ]` / `- [x]` 业务勾选由 **PM 或用户** 验收后改；Implementation 轨默认不改 |
| 全局约束 | 根目录 `.cursorrules`（alwaysApply）+ `.cursor/rules/parallel-track-guard.mdc`（编辑 `docs/tasks/` 时） |

---

## 4. 文件边界（减少 Git 冲突）

| 轨道 | 允许修改 | 禁止修改（除非显式 BLOCKED 协调） |
|------|----------|-----------------------------------|
| **BE** | `internal/`、`cmd/`、`migrations/`、`docs/api/`（实现备注）、`**/*_test.go`（包内单测） | `apps/web/`、`packages/ui-kit/`、`e2e/` |
| **FE** | `apps/web/`、`packages/ui-kit/`、`apps/web/**/*.test.ts`、`packages/ui-kit/**/*.test.ts` | `internal/`、`migrations/`、`e2e/`（E2E 归 QA 轨） |
| **QA** | `e2e/`、`**/*_test.go`、`**/*.test.ts`、`tests/`、测试 fixture/seed 脚本 | 业务实现文件（可 **新增** 测试，勿改 Handler/Vue 业务逻辑；必要 mock 放测试目录） |

跨轨改同一文件 → 先停手，在 phase notes 记 BLOCKED，由人工或主对话合并。

---

## 5. 依赖与 mock 约定

```
契约 (docs/api)
    ├─► [BE] 实现 + Go 集成测
    ├─► [FE] composable 对契约；开发期可用 MSW/静态 mock，路径与字段与契约一致
    └─► [QA] 用例按 PRD 验收标准写；BE 未就绪 → 测 mock 或 `test.skip` + BLOCKED 说明

[FE] 页面需 E2E 时：关键按钮/表格加 `data-testid`（与 QA 轨约定命名）
[BE] API 行为变更：当天更新 docs/api + phase notes，避免 FE/QA 静默失败
```

---

## 6. 同步检查点（建议每月/每子任务 ID）

1. **BE → FE**：子任务 ID 的 API 可调用（或官方 mock 文件路径）→ FE 去掉临时 mock  
2. **FE → QA**：路由、`data-testid`、登录态说明 → QA 提交 E2E  
3. **QA → ALL**：CI 失败用例指回对应轨（在任务清单备注失败套件名）

---

## 7. 对话框开场白模板（复制即用）

**Backend**

```
【BE】Phase 2 — 领 [BE] 且 ⬜ 的子任务。
@backend-developer @development-workflow
先读 docs/tasks/00-mvp-task-breakdown.md Phase 2 并行表、docs/api/ 契约、docs/tasks/parallel-implementation.md。
只改 backend/；同步 Go 测试；完成则在 tasks 标 ✅。
```

**Frontend**

```
【FE】Phase 2 — 领 [FE] 且 ⬜ 的子任务。
@frontend-developer @development-workflow
先读 2b 前端切面、docs/api/ 契约；API 未就绪用契约 mock。
只改 apps/web、packages/ui-kit 与对应 Vitest；完成标 ✅。
```

**Test**

```
【QA】Phase 2 — 领 [QA] 且 ⬜ 的子任务。
先读 PRD 验收标准与 docs/api/；编写/补充 Go 集成、Vitest、Playwright（e2e/）。
BE/FE 未就绪处用 mock 或 skip 并记 BLOCKED 到 meeting-notes。
```

---

## 8. Phase 任务表模板（复制到各 Phase 下）

```markdown
### Phase X 并行任务（Implementation）

| ID | [BE] | [FE] | [QA] | 依赖 |
|----|------|------|------|------|
| X.1 | … API | … 页面 | … 集成测 | 契约 |
| X.2 | … | … | … E2E | X.1 FE 路由 |

状态列 shorthand：`⬜` `🟡` `✅` `🔴`
```

Phase 2 已按此格式拆表示例（含 2.9–2.13 关系经营 / AI Preview），见 [00-mvp-task-breakdown.md § Phase 2](./00-mvp-task-breakdown.md#phase-2客户与线索--关系经营--ai-preview)。
