# Phase 0 笔记

**时间**：2026-05-22  
**状态**：收尾完成（无 CI，含本地测试）

## 验收标准

- [x] Goose 迁移可执行（00001–00003）
- [x] Backend `go test ./...` 通过
- [x] Frontend `npm run test`（Vitest）通过
- [x] E2E `e2e/` Playwright 冒烟用例可运行（需本地 DB + 前后端）
- [x] ADR 0001–0003 已记录
- [x] 多租户 Context + GORM Scope 有测试覆盖
- [x] RBAC 路由映射有单元测试

## 交付物清单

| 类别 | 路径 |
|------|------|
| 迁移 | `backend/migrations/` |
| 后端测试 | `backend/internal/**/*_test.go` |
| 前端测试 | `frontend/utils/permissions.test.ts` |
| E2E | `e2e/tests/phase0.spec.ts` |
| QA 报告 | [qa/phase-0-qa.md](../qa/phase-0-qa.md) |
| Code Review | [reviews/phase-0-review.md](../reviews/phase-0-review.md) |

## 本地启动

```bash
docker compose up -d
cd backend && cp .env.example .env && make migrate-up && make run
cd frontend && npm run dev
```

## 运行测试

```bash
make backend-test
make frontend-test
# E2E（需服务已启动或 Playwright 自动拉起）
make e2e-test
```

## 已知限制

- 登录/API 业务仍为 Phase 1 占位（501）
- 未配置 CI（按项目要求跳过）
- Casbin 集成测试依赖完整 user_roles 数据，留待 Phase 1
