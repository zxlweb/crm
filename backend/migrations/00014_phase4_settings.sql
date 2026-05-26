-- +goose Up
-- Phase 4: custom_fields, tenant plan, audit indexes, RBAC permissions

ALTER TABLE tenants
    ADD COLUMN IF NOT EXISTS plan VARCHAR(50) NOT NULL DEFAULT 'standard';

CREATE TABLE IF NOT EXISTS custom_fields (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    entity_type VARCHAR(20) NOT NULL,
    field_key VARCHAR(100) NOT NULL,
    field_label JSONB NOT NULL DEFAULT '{}',
    field_type VARCHAR(20) NOT NULL,
    required BOOLEAN NOT NULL DEFAULT false,
    options JSONB,
    default_value JSONB,
    display_order INT NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, entity_type, field_key)
);

CREATE INDEX IF NOT EXISTS idx_custom_fields_tenant_entity
    ON custom_fields(tenant_id, entity_type, is_active, display_order);

CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_action_created
    ON audit_logs(tenant_id, action, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_actor_created
    ON audit_logs(tenant_id, user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_tenants_plan_active
    ON tenants(plan, is_active);

INSERT INTO permissions (id, resource, action, description) VALUES
    (gen_random_uuid(), 'settings', 'view', '查看租户设置'),
    (gen_random_uuid(), 'settings', 'update', '更新租户设置'),
    (gen_random_uuid(), 'custom_fields', 'view', '查看自定义字段'),
    (gen_random_uuid(), 'custom_fields', 'update', '管理自定义字段'),
    (gen_random_uuid(), 'audit', 'view', '查看审计统计'),
    (gen_random_uuid(), 'audit', 'export', '导出审计日志'),
    (gen_random_uuid(), 'admin_tenant_insights', 'view', '查看跨租户运营统计')
ON CONFLICT (resource, action) DO NOTHING;

-- +goose Down
DROP INDEX IF EXISTS idx_tenants_plan_active;
DROP INDEX IF EXISTS idx_audit_logs_tenant_actor_created;
DROP INDEX IF EXISTS idx_audit_logs_tenant_action_created;
DROP INDEX IF EXISTS idx_custom_fields_tenant_entity;
DROP TABLE IF EXISTS custom_fields;
ALTER TABLE tenants DROP COLUMN IF EXISTS plan;
DELETE FROM permissions WHERE resource IN ('settings', 'custom_fields', 'audit', 'admin_tenant_insights');
