-- +goose Up
-- Grant Phase 4 permissions to Demo Tenant Admin role

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'dddddddd-dddd-dddd-dddd-dddddddddddd', p.id
FROM permissions p
WHERE (p.resource, p.action) IN (
    ('settings', 'view'),
    ('settings', 'update'),
    ('custom_fields', 'view'),
    ('custom_fields', 'update'),
    ('audit', 'view'),
    ('audit', 'export')
)
ON CONFLICT DO NOTHING;

UPDATE tenants
SET plan = 'professional'
WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

UPDATE tenants
SET config = COALESCE(config, '{}'::jsonb) || jsonb_build_object(
    'default_locale', 'zh-CN',
    'timezone', 'Asia/Shanghai',
    'business_switches', jsonb_build_object(
        'ai_preview_enabled', false,
        'lead_import_mode', 'manual_review'
    )
)
WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

-- +goose Down
DELETE FROM role_permissions
WHERE role_id = 'dddddddd-dddd-dddd-dddd-dddddddddddd'
  AND permission_id IN (
    SELECT id FROM permissions
    WHERE (resource, action) IN (
        ('settings', 'view'),
        ('settings', 'update'),
        ('custom_fields', 'view'),
        ('custom_fields', 'update'),
        ('audit', 'view'),
        ('audit', 'export')
    )
  );

UPDATE tenants SET plan = 'standard' WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
