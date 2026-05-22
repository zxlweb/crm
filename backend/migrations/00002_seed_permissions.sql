-- +goose Up
INSERT INTO permissions (id, resource, action, description) VALUES
    (gen_random_uuid(), 'leads', 'view', '查看线索'),
    (gen_random_uuid(), 'leads', 'create', '创建线索'),
    (gen_random_uuid(), 'leads', 'update', '更新线索'),
    (gen_random_uuid(), 'leads', 'delete', '删除线索'),
    (gen_random_uuid(), 'deals', 'view', '查看商机'),
    (gen_random_uuid(), 'deals', 'create', '创建商机'),
    (gen_random_uuid(), 'deals', 'update', '更新商机'),
    (gen_random_uuid(), 'deals', 'delete', '删除商机'),
    (gen_random_uuid(), 'accounts', 'view', '查看客户'),
    (gen_random_uuid(), 'accounts', 'create', '创建客户'),
    (gen_random_uuid(), 'accounts', 'update', '更新客户'),
    (gen_random_uuid(), 'accounts', 'delete', '删除客户'),
    (gen_random_uuid(), 'contacts', 'view', '查看联系人'),
    (gen_random_uuid(), 'contacts', 'create', '创建联系人'),
    (gen_random_uuid(), 'contacts', 'update', '更新联系人'),
    (gen_random_uuid(), 'contacts', 'delete', '删除联系人'),
    (gen_random_uuid(), 'rbac', 'view', '查看权限配置'),
    (gen_random_uuid(), 'rbac', 'manage', '管理角色与权限'),
    (gen_random_uuid(), 'settings', 'tenant_config', '租户配置管理')
ON CONFLICT (resource, action) DO NOTHING;

-- +goose Down
DELETE FROM permissions WHERE resource IN ('leads', 'deals', 'accounts', 'contacts', 'rbac', 'settings');
