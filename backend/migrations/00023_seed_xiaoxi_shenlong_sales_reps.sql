-- +goose Up
-- 神龙云计算 · 补充一线销售代表（销售经理成员列表 / 团队排行演示）
INSERT INTO users (id, email, password_hash, name, is_super_admin)
VALUES
    (
        'f1000001-000b-4000-8000-00000000000b',
        'zhaoyun@xiaoxi.com',
        '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse',
        '赵云',
        false
    ),
    (
        'f1000001-000c-4000-8000-00000000000c',
        'sunli@xiaoxi.com',
        '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse',
        '孙丽',
        false
    )
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_tenants (user_id, tenant_id, department)
VALUES
    ('f1000001-000b-4000-8000-00000000000b', 'f1111111-1111-4111-8111-111111111111', '神龙云计算'),
    ('f1000001-000c-4000-8000-00000000000c', 'f1111111-1111-4111-8111-111111111111', '神龙云计算')
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, tenant_id)
VALUES
    ('f1000001-000b-4000-8000-00000000000b', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-000c-4000-8000-00000000000c', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111')
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM user_roles
WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111'
  AND user_id IN (
    'f1000001-000b-4000-8000-00000000000b',
    'f1000001-000c-4000-8000-00000000000c'
  );
DELETE FROM user_tenants
WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111'
  AND user_id IN (
    'f1000001-000b-4000-8000-00000000000b',
    'f1000001-000c-4000-8000-00000000000c'
  );
DELETE FROM users
WHERE id IN (
    'f1000001-000b-4000-8000-00000000000b',
    'f1000001-000c-4000-8000-00000000000c'
);
