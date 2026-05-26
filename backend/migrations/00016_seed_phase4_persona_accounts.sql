-- +goose Up
-- Phase 4 PRD §2 Persona demo accounts (Demo Corp tenant)
-- Password for all: password123 (bcrypt, same as admin@demo.com)

-- Roles
INSERT INTO roles (id, tenant_id, name, description, is_system)
VALUES
    (
        'e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        'Sales Manager',
        '销售经理（演示）：审计/设置只读，无导出与配置变更',
        true
    ),
    (
        'e2e2e2e2-e2e2-e2e2-e2e2-e2e2e2e2e2e2',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        'Viewer',
        '只读（演示）：查看设置与审计，不可修改',
        true
    )
ON CONFLICT (id) DO NOTHING;

-- Sales Manager: team-oriented read + CRM view (no settings/custom_fields update, no audit export)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1', p.id
FROM permissions p
WHERE (p.resource, p.action) IN (
    ('settings', 'view'),
    ('custom_fields', 'view'),
    ('audit', 'view'),
    ('dashboard', 'view'),
    ('leads', 'view'),
    ('deals', 'view'),
    ('accounts', 'view'),
    ('contacts', 'view')
)
ON CONFLICT DO NOTHING;

-- Viewer: read-only across settings + audit + CRM lists
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'e2e2e2e2-e2e2-e2e2-e2e2-e2e2e2e2e2e2', p.id
FROM permissions p
WHERE (p.resource, p.action) IN (
    ('settings', 'view'),
    ('custom_fields', 'view'),
    ('audit', 'view'),
    ('dashboard', 'view'),
    ('leads', 'view'),
    ('deals', 'view'),
    ('accounts', 'view'),
    ('contacts', 'view')
)
ON CONFLICT DO NOTHING;

-- Users (password123)
INSERT INTO users (id, email, password_hash, name, is_super_admin)
VALUES
    (
        'c1c1c1c1-c1c1-c1c1-c1c1-c1c1c1c1c1c1',
        'tenant-admin@demo.com',
        '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse',
        'Demo Tenant Admin',
        false
    ),
    (
        'c2c2c2c2-c2c2-c2c2-c2c2-c2c2c2c2c2c2',
        'manager@demo.com',
        '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse',
        'Demo Sales Manager',
        false
    ),
    (
        'c3c3c3c3-c3c3-c3c3-c3c3-c3c3c3c3c3c3',
        'viewer@demo.com',
        '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse',
        'Demo Viewer',
        false
    )
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_tenants (user_id, tenant_id)
VALUES
    ('c1c1c1c1-c1c1-c1c1-c1c1-c1c1c1c1c1c1', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('c2c2c2c2-c2c2-c2c2-c2c2-c2c2c2c2c2c2', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('c3c3c3c3-c3c3-c3c3-c3c3-c3c3c3c3c3c3', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa')
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, tenant_id)
VALUES
    (
        'c1c1c1c1-c1c1-c1c1-c1c1-c1c1c1c1c1c1',
        'dddddddd-dddd-dddd-dddd-dddddddddddd',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
    ),
    (
        'c2c2c2c2-c2c2-c2c2-c2c2-c2c2c2c2c2c2',
        'e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
    ),
    (
        'c3c3c3c3-c3c3-c3c3-c3c3-c3c3c3c3c3c3',
        'e2e2e2e2-e2e2-e2e2-e2e2-e2e2e2e2e2e2',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
    )
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM user_roles
WHERE user_id IN (
    'c1c1c1c1-c1c1-c1c1-c1c1-c1c1c1c1c1c1',
    'c2c2c2c2-c2c2-c2c2-c2c2-c2c2c2c2c2c2',
    'c3c3c3c3-c3c3-c3c3-c3c3-c3c3c3c3c3c3'
);
DELETE FROM user_tenants
WHERE user_id IN (
    'c1c1c1c1-c1c1-c1c1-c1c1-c1c1c1c1c1c1',
    'c2c2c2c2-c2c2-c2c2-c2c2-c2c2c2c2c2c2',
    'c3c3c3c3-c3c3-c3c3-c3c3-c3c3c3c3c3c3'
);
DELETE FROM users
WHERE email IN (
    'tenant-admin@demo.com',
    'manager@demo.com',
    'viewer@demo.com'
);
DELETE FROM role_permissions
WHERE role_id IN (
    'e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1',
    'e2e2e2e2-e2e2-e2e2-e2e2-e2e2e2e2e2e2'
);
DELETE FROM roles
WHERE id IN (
    'e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1',
    'e2e2e2e2-e2e2-e2e2-e2e2-e2e2e2e2e2e2'
);
