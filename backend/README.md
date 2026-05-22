# CRM Backend

Go + Gin + GORM + PostgreSQL + Casbin 多租户 CRM 后端。

## 技术栈

- Go 1.25+
- Gin
- GORM + PostgreSQL
- Casbin (RBAC)
- Goose (数据库迁移)

## 快速开始

### 1. 环境变量

```bash
cp .env.example .env
# 编辑 DB_* 与 JWT_SECRET
```

### 2. 数据库迁移（推荐）

```bash
# 在 backend 目录下执行
make migrate-up
make migrate-status
```

### 3. 启动 API

```bash
make run
# 或
go run ./cmd/api/
```

### 4. 健康检查

```bash
curl http://localhost:8080/health
```

### 5. 运行测试

```bash
make test
```

## 开发说明

| 命令 | 说明 |
|------|------|
| `make migrate-up` | 执行 Goose 迁移 |
| `make migrate-down` | 回滚一个版本 |
| `make migrate-status` | 查看迁移状态 |
| `AUTO_MIGRATE=true make run` | 开发时用 GORM AutoMigrate（可选） |

## 目录结构

```
cmd/api/          # HTTP 服务入口
cmd/migrate/      # Goose 迁移 CLI
migrations/       # SQL 迁移脚本
internal/domain/  # 领域模型
internal/infrastructure/persistence/  # DB + Casbin
internal/interfaces/http/             # HTTP Handlers
internal/interfaces/middleware/       # 中间件（Auth / Tenant / RBAC）
internal/pkg/tenant/                  # 租户 Context
internal/pkg/rbac/                    # 路由 → 权限映射
```

## 多租户

- 请求头：`X-Tenant-ID: <uuid>`
- `TenantMiddleware` 写入 `context.Context`
- 查询业务数据使用 `persistence.DBFromContext(db, ctx)`

## RBAC

- 权限格式：`resource:action`（如 `leads:view`）
- 策略存储：`roles` / `permissions` / `role_permissions` / `user_roles`
- 启动时 `SyncCasbinPolicies` 从 DB 同步到 Casbin
