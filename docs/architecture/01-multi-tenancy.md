# 01 多租户架构设计

**版本**：v0.1  
**日期**：2026-05-21

## 1. 多租户需求概述

- 支持多个企业（Tenant）共用同一套系统
- 数据完全隔离（用户只能看到自己租户的数据）
- 租户级配置（主题、功能开关、i18n 默认语言等）
- 支持 Super Admin 跨租户管理
- 用户可属于多个租户并切换

## 2. 选型对比

| 方案 | 隔离性 | 开发难度 | 运维成本 | 扩展性 | 推荐指数 |
|------|--------|----------|----------|--------|----------|
| 共享数据库 + tenant_id | 中 | 低 | 低 | 高 | ★★★★★ |
| Schema 隔离 | 高 | 中 | 中 | 中 | ★★★ |
| 独立数据库 | 最高 | 高 | 高 | 低 | ★ |

**最终决策**：**共享数据库 + tenant_id + PostgreSQL RLS**（Row Level Security）

## 3. 实现方案

### 数据库层面
- 所有业务表必须包含 `tenant_id UUID NOT NULL`
- 核心共享表（tenants、users、roles 等）不强制 tenant_id 或设为特殊值
- 启用 PostgreSQL RLS 策略（生产环境强烈推荐）

### Backend (Go) 实现
- **Tenant Middleware**（最优先执行）：
  - 从 JWT 中解析 `tenant_id`
  - 将 `tenant_id` 注入 Context / Request Scope
- **Tenant-aware Query**：
  - GORM：使用 `Scopes` 自动添加 `WHERE tenant_id = ?`
  - sqlc：封装 Query 函数，强制传入 tenant_id
- **数据隔离关键点**：
  - 所有 CRUD 操作必须经过 Tenant Filter
  - 禁止使用全局查询（提供 `WithTenant()` 方法）

### Frontend (Nuxt 3) 实现
- `composables/useTenant.ts`：管理当前租户信息
- API 请求自动携带 `tenant-id` Header
- 路由守卫检查租户权限

## 4. 核心表结构（示例）

```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    domain VARCHAR(100) UNIQUE,
    config JSONB,                    -- 租户个性化配置
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 业务表示例
CREATE TABLE leads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    owner_id UUID,                   -- 数据归属人
    name VARCHAR(200),
    -- ... 其他字段
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- RLS 策略示例
CREATE POLICY tenant_isolation ON leads
    USING (tenant_id = current_setting('app.current_tenant')::uuid);