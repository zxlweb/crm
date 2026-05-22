-- +goose Up
-- Phase 2 CRM: accounts, contacts, activities, segment_templates, leads extensions, permissions

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    industry VARCHAR(100),
    website TEXT,
    lifecycle_stage VARCHAR(32) NOT NULL DEFAULT 'acquire',
    engagement_score SMALLINT NOT NULL DEFAULT 0,
    last_activity_at TIMESTAMPTZ,
    tags TEXT[] NOT NULL DEFAULT '{}',
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_accounts_tenant ON accounts(tenant_id);
CREATE INDEX idx_accounts_tenant_owner ON accounts(tenant_id, owner_id);
CREATE INDEX idx_accounts_tenant_lifecycle ON accounts(tenant_id, lifecycle_stage);
CREATE INDEX idx_accounts_list ON accounts(tenant_id, deleted_at, updated_at DESC);

CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    is_primary BOOLEAN NOT NULL DEFAULT false,
    lifecycle_stage VARCHAR(32) NOT NULL DEFAULT 'acquire',
    engagement_score SMALLINT NOT NULL DEFAULT 0,
    last_activity_at TIMESTAMPTZ,
    tags TEXT[] NOT NULL DEFAULT '{}',
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_contacts_tenant ON contacts(tenant_id);
CREATE INDEX idx_contacts_tenant_account ON contacts(tenant_id, account_id);
CREATE INDEX idx_contacts_tenant_email ON contacts(tenant_id, email);

CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    subject_type VARCHAR(20) NOT NULL,
    subject_id UUID NOT NULL,
    event_type VARCHAR(32) NOT NULL,
    direction VARCHAR(16),
    body TEXT,
    metadata JSONB NOT NULL DEFAULT '{}',
    sentiment VARCHAR(20),
    sentiment_source VARCHAR(16),
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_activities_subject ON activities(tenant_id, subject_type, subject_id, occurred_at DESC);

CREATE TABLE segment_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name_i18n_key VARCHAR(100) NOT NULL,
    filter_json JSONB NOT NULL DEFAULT '{}',
    is_system BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_segment_templates_tenant_code ON segment_templates(tenant_id, code);

CREATE TABLE lifecycle_stage_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    subject_type VARCHAR(20) NOT NULL,
    subject_id UUID NOT NULL,
    from_stage VARCHAR(32),
    to_stage VARCHAR(32) NOT NULL,
    changed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_lifecycle_history_subject ON lifecycle_stage_history(tenant_id, subject_type, subject_id, created_at DESC);

ALTER TABLE leads
    ADD COLUMN IF NOT EXISTS lifecycle_stage VARCHAR(32) DEFAULT 'acquire',
    ADD COLUMN IF NOT EXISTS engagement_score SMALLINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS last_activity_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS relationship_health VARCHAR(16),
    ADD COLUMN IF NOT EXISTS converted_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS converted_contact_id UUID REFERENCES contacts(id) ON DELETE SET NULL;

-- Phase 2 RBAC permissions
INSERT INTO permissions (id, resource, action, description) VALUES
    (gen_random_uuid(), 'leads', 'assign', '分配线索'),
    (gen_random_uuid(), 'activities', 'view', '查看跟进'),
    (gen_random_uuid(), 'activities', 'create', '创建跟进'),
    (gen_random_uuid(), 'activities', 'update', '更新跟进'),
    (gen_random_uuid(), 'activities', 'delete', '删除跟进'),
    (gen_random_uuid(), 'insights', 'view', '查看洞察'),
    (gen_random_uuid(), 'segments', 'view', '查看分群'),
    (gen_random_uuid(), 'segments', 'manage', '管理分群'),
    (gen_random_uuid(), 'copilot', 'use', '使用 Copilot / AI Preview')
ON CONFLICT (resource, action) DO NOTHING;

-- Grant new permissions to demo Tenant Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'dddddddd-dddd-dddd-dddd-dddddddddddd', p.id
FROM permissions p
WHERE NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = 'dddddddd-dddd-dddd-dddd-dddddddddddd' AND rp.permission_id = p.id
);

-- Default tenant config keys for all tenants
UPDATE tenants SET config = config || '{
  "ai_enabled": false,
  "ai_preview_mode": "off",
  "insight_thresholds": { "days_silent": 7, "high_value_amount": 100000 },
  "sentiment_keyword_rules": []
}'::jsonb
WHERE config = '{}'::jsonb OR NOT (config ? 'ai_enabled');

-- Demo tenant: AI preview fixtures
UPDATE tenants SET config = config || '{
  "ai_enabled": true,
  "ai_preview_mode": "fixtures"
}'::jsonb
WHERE domain = 'demo';

-- +goose Down
ALTER TABLE leads
    DROP COLUMN IF EXISTS converted_contact_id,
    DROP COLUMN IF EXISTS converted_account_id,
    DROP COLUMN IF EXISTS relationship_health,
    DROP COLUMN IF EXISTS tags,
    DROP COLUMN IF EXISTS last_activity_at,
    DROP COLUMN IF EXISTS engagement_score,
    DROP COLUMN IF EXISTS lifecycle_stage;

DROP TABLE IF EXISTS lifecycle_stage_history;
DROP TABLE IF EXISTS segment_templates;
DROP TABLE IF EXISTS activities;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS accounts;

DELETE FROM permissions WHERE resource IN ('activities', 'insights', 'segments', 'copilot')
    OR (resource = 'leads' AND action = 'assign');
