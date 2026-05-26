-- +goose Up
ALTER TABLE user_tenants ADD COLUMN IF NOT EXISTS department VARCHAR(100);

CREATE INDEX IF NOT EXISTS idx_user_tenants_tenant_department
    ON user_tenants (tenant_id, department);

-- 小西科技集团 · 部门归属（与 docs/demo/xiaoxi-boss-demo.md 一致）
UPDATE user_tenants SET department = '总裁办' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-0001-4000-8000-000000000001';
UPDATE user_tenants SET department = '技术中心' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-0002-4000-8000-000000000002';
UPDATE user_tenants SET department = '灵狐数据' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-0003-4000-8000-000000000003';
UPDATE user_tenants SET department = '莫邪互娱' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-0004-4000-8000-000000000004';
UPDATE user_tenants SET department = '神龙云计算' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-0005-4000-8000-000000000005';
UPDATE user_tenants SET department = '南京小东' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-0006-4000-8000-000000000006';
UPDATE user_tenants SET department = '麒麟事业群' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-0007-4000-8000-000000000007';
UPDATE user_tenants SET department = '马来西亚' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-0008-4000-8000-000000000008';
UPDATE user_tenants SET department = '数据中心实施部' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-0009-4000-8000-000000000009';
UPDATE user_tenants SET department = '产品交付部' WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111' AND user_id = 'f1000001-000a-4000-8000-00000000000a';

-- +goose Down
UPDATE user_tenants SET department = NULL WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111';
DROP INDEX IF EXISTS idx_user_tenants_tenant_department;
ALTER TABLE user_tenants DROP COLUMN IF EXISTS department;
