# CRM Backend

Go + Gin + GORM + PostgreSQL + Casbin 多租户 CRM 后端。

## 技术栈

- Go 1.25+
- Gin
- GORM + PostgreSQL
- Casbin (RBAC)
- Goose (数据库迁移)

## 快速开始

### Windows：Go 与 PATH

安装 [Go](https://go.dev/dl/) 后，将 `C:\Program Files\Go\bin` 加入**用户**或**系统**环境变量 `Path`（设置 → 系统 → 关于 → 高级系统设置 → 环境变量）。修改 PATH 后需**完全退出并重新打开 Cursor**（仅新开终端有时仍读不到旧环境）。也可在 `backend` 目录执行 `.\run.ps1` 启动，无需依赖全局 `go` 命令。若 `8080` 被占用，用 `netstat -ano | findstr :8080` 查 PID 后 `taskkill /PID <pid> /F` 再启动。

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

Windows（终端找不到 `go` 时）：

```powershell
.\run.ps1
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
