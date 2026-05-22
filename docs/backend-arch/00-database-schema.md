#### **核心设计原则**
- 所有业务表都包含 `tenant_id`
- 软删除（`deleted_at`）
- 审计字段（`created_by`, `updated_by`）
- 主键使用 `UUID`
- 时间字段使用 `timestamptz`

#### **SQL 表结构（核心表）**

```sql
-- 1. 租户表
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(100) UNIQUE,
    logo_url TEXT,
    config JSONB DEFAULT '{}',           -- 租户个性化配置
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. 用户表（全局）
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name VARCHAR(100),
    avatar_url TEXT,
    is_super_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- 3. 用户-租户关联表（支持用户属于多个租户）
CREATE TABLE user_tenants (
    user_id UUID REFERENCES users(id),
    tenant_id UUID REFERENCES tenants(id),
    PRIMARY KEY (user_id, tenant_id)
);

-- 4. 角色表
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT false,     -- 系统内置角色不可删除
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. 权限资源表
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL,      -- 如 "leads", "deals"
    action VARCHAR(50) NOT NULL,         -- 如 "create", "view", "delete"
    description TEXT,
    UNIQUE(resource, action)
);

-- 6. 角色-权限关联
CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id),
    permission_id UUID REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

-- 7. 用户-角色关联（租户内）
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id),
    role_id UUID REFERENCES roles(id),
    tenant_id UUID REFERENCES tenants(id),
    PRIMARY KEY (user_id, role_id, tenant_id)
);

-- 8. 审计日志表（重要）
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID,
    action VARCHAR(100),
    resource_type VARCHAR(100),
    resource_id UUID,
    old_value JSONB,
    new_value JSONB,
    ip_address INET,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 示例业务表（Leads）
CREATE TABLE leads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    owner_id UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    status VARCHAR(50),
    source VARCHAR(50),
    amount DECIMAL(15,2),
    expected_close_date DATE,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- 索引建议
CREATE INDEX idx_leads_tenant ON leads(tenant_id);
CREATE INDEX idx_leads_owner ON leads(owner_id);
```

#### **GORM Model 示例**（Go）
我后续可以给你完整 Model 文件，这里先给一个示例：

```go
type Tenant struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name      string
    Domain    string
    Config    datatypes.JSON
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Lead struct {
    ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
    TenantID        uuid.UUID `gorm:"index"`
    OwnerID         *uuid.UUID
    Title           string
    // ...
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt `gorm:"index"`
}
```



