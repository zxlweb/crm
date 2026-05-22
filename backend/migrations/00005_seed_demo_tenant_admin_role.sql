-- +goose Up
INSERT INTO roles (id, tenant_id, name, description, is_system)
VALUES (
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'Tenant Admin',
    '租户管理员（演示）',
    true
) ON CONFLICT (id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'dddddddd-dddd-dddd-dddd-dddddddddddd', p.id
FROM permissions p
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, tenant_id)
VALUES (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
) ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM user_roles WHERE role_id = 'dddddddd-dddd-dddd-dddd-dddddddddddd';
DELETE FROM role_permissions WHERE role_id = 'dddddddd-dddd-dddd-dddd-dddddddddddd';
DELETE FROM roles WHERE id = 'dddddddd-dddd-dddd-dddd-dddddddddddd';
