-- +goose Up
-- Demo user with multiple roles for role-switch testing (password123)

INSERT INTO users (id, email, password_hash, name, is_super_admin)
VALUES (
    'c4c4c4c4-c4c4-c4c4-c4c4-c4c4c4c4c4c4',
    'multi-role@demo.com',
    '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse',
    'Demo Multi Role',
    false
) ON CONFLICT (email) DO NOTHING;

INSERT INTO user_tenants (user_id, tenant_id)
VALUES (
    'c4c4c4c4-c4c4-c4c4-c4c4-c4c4c4c4c4c4',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
) ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, tenant_id)
VALUES
    (
        'c4c4c4c4-c4c4-c4c4-c4c4-c4c4c4c4c4c4',
        'e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
    ),
    (
        'c4c4c4c4-c4c4-c4c4-c4c4-c4c4c4c4c4c4',
        'e2e2e2e2-e2e2-e2e2-e2e2-e2e2e2e2e2e2',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
    )
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM user_roles
WHERE user_id = 'c4c4c4c4-c4c4-c4c4-c4c4-c4c4c4c4c4c4';
DELETE FROM user_tenants WHERE user_id = 'c4c4c4c4-c4c4-c4c4-c4c4-c4c4c4c4c4c4';
DELETE FROM users WHERE email = 'multi-role@demo.com';
