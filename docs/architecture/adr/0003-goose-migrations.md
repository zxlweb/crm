# ADR-0003: Goose 作为数据库 Schema 迁移工具

**状态**：Accepted  
**日期**：2026-05-21  
**决策者**：Backend Developer

## 背景（Context）

Schema 需版本化、可回滚、可在 CI/CD 中执行。GORM AutoMigrate 适合本地快速试验，但不利于审查 SQL 与生产一致性。

## 决策（Decision）

**生产与团队协作用 [Goose](https://github.com/pressly/goose)** 管理迁移：

- SQL 文件位于 `backend/migrations/`
- CLI 入口 `backend/cmd/migrate/`（`make migrate-up`）
- GORM `AutoMigrate` 仅通过 `AUTO_MIGRATE=true` 在本地可选启用

## 备选方案（Options Considered）

| 方案 | 优点 | 缺点 |
|------|------|------|
| Goose ✅ | SQL 可审、Up/Down、生态成熟 | 需维护 SQL 与 Model 一致 |
| golang-migrate | 同样成熟 | 团队未统一前不引入双工具 |
| 仅 GORM AutoMigrate | 开发快 | 难 Code Review、回滚弱 |

## 后果（Consequences）

### 正面
- `00001_init_core_schema` / `00002_seed_permissions` 可重复部署
- 与 `.cursorrules` 后端规范一致

### 负面 / 风险
- Domain Model 与 SQL 需人工保持同步（变更时同时改两处）

### 后续行动
- [ ] CI 中加入 `make migrate-up` 烟雾步骤
- [ ] 新表继续 `0000N_*.sql` 递增编号

## 关联文档

- [backend-arch/00-database-schema.md](../../backend-arch/00-database-schema.md)
- [backend/README.md](../../../backend/README.md)
