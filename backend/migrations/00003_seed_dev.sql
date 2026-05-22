-- +goose Up
-- 开发环境演示数据（仅本地/测试使用）
INSERT INTO tenants (id, name, domain, is_active)
VALUES (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'Demo Corp',
    'demo',
    true
) ON CONFLICT (domain) DO NOTHING;

INSERT INTO users (id, email, password_hash, name, is_super_admin)
VALUES (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'admin@demo.com',
    '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse',
    'Demo Admin',
    true
) ON CONFLICT (email) DO NOTHING;

INSERT INTO user_tenants (user_id, tenant_id)
VALUES (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
) ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM user_tenants WHERE tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
DELETE FROM users WHERE email = 'admin@demo.com';
DELETE FROM tenants WHERE domain = 'demo';
