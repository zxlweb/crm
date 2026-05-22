#!/bin/bash

# ==================== CRM Backend 目录结构创建脚本 ====================
echo "🚀 开始创建 Go Backend 项目目录结构..."

# 创建主目录
mkdir -p backend
cd backend

# 创建所有目录
mkdir -p cmd/api
mkdir -p internal/config
mkdir -p internal/domain
mkdir -p internal/application/{auth,tenant,rbac,lead,contact,deal,dashboard,settings}
mkdir -p internal/infrastructure/{persistence,cache,casbin,logger}
mkdir -p internal/interfaces/http/handlers
mkdir -p internal/interfaces/middleware
mkdir -p internal/interfaces/api
mkdir -p internal/pkg/{utils,errors,response}
mkdir -p migrations
mkdir -p docs/swagger

# 创建核心占位文件
touch cmd/api/main.go
touch internal/config/config.go
touch go.mod
touch go.sum
touch README.md

# 创建部分重要空文件（便于后续复制代码）
touch internal/domain/tenant.go
touch internal/domain/user.go
touch internal/domain/lead.go
touch internal/infrastructure/persistence/gorm.go
tes/middleware/rbac.go

echo '
module crm-backend

go 1.23
' > go.mod

echo '
package main

import "fmt"

func main() {
    fmt.Println("🚀 CRM Backend started successfully!")
}
' > cmd/api/main.go

echo '# CRM Backend

基于 Clean Architecture 的企业级 CRM 后端系统。

## 技术栈
- Go 1.23+
- Gin
- GORM + PostgreSQL
- Casbin (RBAC)
- Goose (迁移)
' > README.md

echo "✅ 目录结构创建完成！"

# 显示目录树
echo "📁 生成的目录结构："
tree -L 3

