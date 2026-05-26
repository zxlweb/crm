-- +goose Up
-- 小西科技集团 · 按钉钉组织架构的老板演示数据（本地/演示环境）
-- 登录：ceo@xiaoxi.com / password123 · 租户域名 xiaoxi

-- ── 租户 ─────────────────────────────────────────────────────────────
INSERT INTO tenants (id, name, domain, is_active, config)
VALUES (
    'f1111111-1111-4111-8111-111111111111',
    '小西科技集团',
    'xiaoxi',
    true,
    jsonb_build_object(
        'default_locale', 'zh-CN',
        'sales_quota', jsonb_build_object(
            'amount', 28000000,
            'currency', 'CNY',
            'period', to_char(CURRENT_DATE, 'YYYY-MM')
        )
    )
) ON CONFLICT (domain) DO UPDATE SET
    name = EXCLUDED.name,
    config = tenants.config || EXCLUDED.config;

-- ── 角色 ─────────────────────────────────────────────────────────────
INSERT INTO roles (id, tenant_id, name, description, is_system)
VALUES
    (
        'f2111111-1111-4111-8111-111111111111',
        'f1111111-1111-4111-8111-111111111111',
        'Tenant Admin',
        '集团管理员（潘总 / 总裁办）',
        true
    ),
    (
        'f2121111-1111-4111-8111-111111111111',
        'f1111111-1111-4111-8111-111111111111',
        'Sales Rep',
        '事业群销售（各业务部门负责人）',
        true
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'f2111111-1111-4111-8111-111111111111', p.id FROM permissions p
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'f2121111-1111-4111-8111-111111111111', p.id
FROM permissions p
WHERE (p.resource, p.action) IN (
    ('dashboard', 'view'),
    ('leads', 'view'), ('leads', 'create'), ('leads', 'update'), ('leads', 'assign'),
    ('deals', 'view'), ('deals', 'create'), ('deals', 'update'),
    ('accounts', 'view'), ('accounts', 'create'), ('accounts', 'update'),
    ('contacts', 'view'), ('contacts', 'create'), ('contacts', 'update'),
    ('activities', 'view'), ('activities', 'create'),
    ('insights', 'view'), ('segments', 'view'), ('copilot', 'use')
)
ON CONFLICT DO NOTHING;

-- password123
INSERT INTO users (id, email, password_hash, name, is_super_admin)
VALUES
    ('f1000001-0001-4000-8000-000000000001', 'ceo@xiaoxi.com',       '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '潘卫国', false),
    ('f1000001-0002-4000-8000-000000000002', 'tech@xiaoxi.com',      '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '张伟',   false),
    ('f1000001-0003-4000-8000-000000000003', 'linghu@xiaoxi.com',    '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '李婷',   false),
    ('f1000001-0004-4000-8000-000000000004', 'moye@xiaoxi.com',      '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '王磊',   false),
    ('f1000001-0005-4000-8000-000000000005', 'shenlong@xiaoxi.com',  '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '陈洋',   false),
    ('f1000001-0006-4000-8000-000000000006', 'nanjing@xiaoxi.com',   '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '赵敏',   false),
    ('f1000001-0007-4000-8000-000000000007', 'kylin@xiaoxi.com',     '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '刘洋',   false),
    ('f1000001-0008-4000-8000-000000000008', 'my@xiaoxi.com',        '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '林薇',   false),
    ('f1000001-0009-4000-8000-000000000009', 'dc@xiaoxi.com',        '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '周强',   false),
    ('f1000001-000a-4000-8000-00000000000a', 'delivery@xiaoxi.com',  '$2a$10$nJa4XI2mJcQ1X7arOl5Ak.ykJkz32Lt2/vBvld0ly6G..rdxv/bse', '孙丽',   false)
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_tenants (user_id, tenant_id)
SELECT u.id, 'f1111111-1111-4111-8111-111111111111'
FROM users u
WHERE u.email LIKE '%@xiaoxi.com'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id, tenant_id)
VALUES
    ('f1000001-0001-4000-8000-000000000001', 'f2111111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-0002-4000-8000-000000000002', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-0003-4000-8000-000000000003', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-0004-4000-8000-000000000004', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-0005-4000-8000-000000000005', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-0006-4000-8000-000000000006', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-0007-4000-8000-000000000007', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-0008-4000-8000-000000000008', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-0009-4000-8000-000000000009', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111'),
    ('f1000001-000a-4000-8000-00000000000a', 'f2121111-1111-4111-8111-111111111111', 'f1111111-1111-4111-8111-111111111111')
ON CONFLICT DO NOTHING;

-- ── 客户（accounts）──────────────────────────────────────────────────
INSERT INTO accounts (
    id, tenant_id, owner_id, name, industry, website,
    lifecycle_stage, engagement_score, last_activity_at, tags,
    created_by, created_at, updated_at
) VALUES
('f3100001-0001-4000-8000-000000000001', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0003-4000-8000-000000000003', '华东数智科技', 'technology', 'https://huadong-data.example.com', 'grow', 82, NOW() - INTERVAL '1 day', ARRAY['灵狐数据', '重点'], 'f1000001-0003-4000-8000-000000000003', NOW() - INTERVAL '60 days', NOW() - INTERVAL '1 day'),
('f3100001-0002-4000-8000-000000000002', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0004-4000-8000-000000000004', '星耀互动娱乐', 'gaming', 'https://xingyao-play.example.com', 'activate', 71, NOW() - INTERVAL '2 days', ARRAY['莫邪互娱'], 'f1000001-0004-4000-8000-000000000004', NOW() - INTERVAL '35 days', NOW() - INTERVAL '2 days'),
('f3100001-0003-4000-8000-000000000003', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0005-4000-8000-000000000005', '江城政务云', 'government', NULL, 'retain', 90, NOW() - INTERVAL '3 hours', ARRAY['神龙云计算', '政企'], 'f1000001-0005-4000-8000-000000000005', NOW() - INTERVAL '120 days', NOW() - INTERVAL '3 hours'),
('f3100001-0004-4000-8000-000000000004', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0007-4000-8000-000000000007', '麒麟智造股份', 'manufacturing', 'https://kylin-mfg.example.com', 'grow', 76, NOW() - INTERVAL '4 days', ARRAY['麒麟事业群'], 'f1000001-0007-4000-8000-000000000007', NOW() - INTERVAL '45 days', NOW() - INTERVAL '4 days'),
('f3100001-0005-4000-8000-000000000005', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0008-4000-8000-000000000008', '吉隆坡零售集团', 'retail', NULL, 'activate', 58, NOW() - INTERVAL '6 days', ARRAY['马来西亚', '出海'], 'f1000001-0008-4000-8000-000000000008', NOW() - INTERVAL '22 days', NOW() - INTERVAL '6 days'),
('f3100001-0006-4000-8000-000000000006', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0009-4000-8000-000000000009', '浦发金融科技', 'finance', 'https://pf-fintech.example.com', 'retain', 85, NOW() - INTERVAL '1 day', ARRAY['数据中心实施部'], 'f1000001-0009-4000-8000-000000000009', NOW() - INTERVAL '200 days', NOW() - INTERVAL '1 day'),
('f3100001-0007-4000-8000-000000000007', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0002-4000-8000-000000000002', '内部研发平台', 'technology', NULL, 'retain', 95, NOW() - INTERVAL '12 hours', ARRAY['技术中心'], 'f1000001-0002-4000-8000-000000000002', NOW() - INTERVAL '300 days', NOW() - INTERVAL '12 hours'),
('f3100001-0008-4000-8000-000000000008', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0006-4000-8000-000000000006', '宁东供应链', 'logistics', NULL, 'acquire', 42, NOW() - INTERVAL '10 days', ARRAY['南京小东'], 'f1000001-0006-4000-8000-000000000006', NOW() - INTERVAL '14 days', NOW() - INTERVAL '10 days')
ON CONFLICT (id) DO NOTHING;

-- ── 联系人 ───────────────────────────────────────────────────────────
INSERT INTO contacts (
    id, tenant_id, account_id, owner_id, first_name, last_name, email, phone,
    is_primary, lifecycle_stage, engagement_score, last_activity_at, tags,
    created_by, created_at, updated_at
) VALUES
('f6100001-0001-4000-8000-000000000001', 'f1111111-1111-4111-8111-111111111111', 'f3100001-0001-4000-8000-000000000001', 'f1000001-0003-4000-8000-000000000003', '陈', '总', 'chen.zong@huadong-data.example.com', '13800001001', true, 'grow', 80, NOW() - INTERVAL '1 day', ARRAY['灵狐数据', '决策人'], 'f1000001-0003-4000-8000-000000000003', NOW() - INTERVAL '50 days', NOW() - INTERVAL '1 day'),
('f6100001-0002-4000-8000-000000000002', 'f1111111-1111-4111-8111-111111111111', 'f3100001-0002-4000-8000-000000000002', 'f1000001-0004-4000-8000-000000000004', '周', '策划', 'zhou@xingyao-play.example.com', '13800002002', true, 'activate', 68, NOW() - INTERVAL '2 days', ARRAY['莫邪互娱'], 'f1000001-0004-4000-8000-000000000004', NOW() - INTERVAL '30 days', NOW() - INTERVAL '2 days'),
('f6100001-0003-4000-8000-000000000003', 'f1111111-1111-4111-8111-111111111111', 'f3100001-0003-4000-8000-000000000003', 'f1000001-0005-4000-8000-000000000005', '李', '处长', 'li.chu@jiangcheng.gov.cn', '027-88880001', true, 'retain', 88, NOW() - INTERVAL '3 hours', ARRAY['神龙云计算'], 'f1000001-0005-4000-8000-000000000005', NOW() - INTERVAL '100 days', NOW() - INTERVAL '3 hours'),
('f6100001-0004-4000-8000-000000000004', 'f1111111-1111-4111-8111-111111111111', 'f3100001-0004-4000-8000-000000000004', 'f1000001-0007-4000-8000-000000000007', '王', '采购', 'wang.cg@kylin-mfg.example.com', '13800004004', true, 'grow', 74, NOW() - INTERVAL '4 days', ARRAY['麒麟事业群'], 'f1000001-0007-4000-8000-000000000007', NOW() - INTERVAL '40 days', NOW() - INTERVAL '4 days')
ON CONFLICT (id) DO NOTHING;

-- ── 线索 ─────────────────────────────────────────────────────────────
INSERT INTO leads (
    id, tenant_id, owner_id, title, status, source, amount, expected_close_date,
    lifecycle_stage, engagement_score, last_activity_at, tags, relationship_health,
    converted_account_id, created_by, created_at, updated_at
) VALUES
('f4100001-0001-4000-8000-000000000001', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0003-4000-8000-000000000003', '⭐ 华东数智 · 数据中台年度', 'qualified', 'referral', 1280000, (CURRENT_DATE + INTERVAL '38 days')::date, 'grow', 82, NOW() - INTERVAL '1 day', ARRAY['灵狐数据', '重点', '华东'], 'high', NULL, 'f1000001-0003-4000-8000-000000000003', NOW() - INTERVAL '55 days', NOW() - INTERVAL '1 day'),
('f4100001-0002-4000-8000-000000000002', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0004-4000-8000-000000000004', '星耀互动 · 联运分成方案', 'contacted', 'partner', 680000, (CURRENT_DATE + INTERVAL '52 days')::date, 'activate', 71, NOW() - INTERVAL '2 days', ARRAY['莫邪互娱'], 'medium', NULL, 'f1000001-0004-4000-8000-000000000004', NOW() - INTERVAL '28 days', NOW() - INTERVAL '2 days'),
('f4100001-0003-4000-8000-000000000003', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0005-4000-8000-000000000005', '江城政务云 · 二期扩容', 'qualified', 'inbound', 2200000, (CURRENT_DATE + INTERVAL '21 days')::date, 'retain', 90, NOW() - INTERVAL '3 hours', ARRAY['神龙云计算', '政企'], 'high', 'f3100001-0003-4000-8000-000000000003', 'f1000001-0005-4000-8000-000000000005', NOW() - INTERVAL '90 days', NOW() - INTERVAL '3 hours'),
('f4100001-0004-4000-8000-000000000004', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0007-4000-8000-000000000007', '麒麟智造 · MES 集成', 'new', 'exhibition', 950000, (CURRENT_DATE + INTERVAL '70 days')::date, 'activate', 58, NOW() - INTERVAL '5 days', ARRAY['麒麟事业群'], 'medium', NULL, 'f1000001-0007-4000-8000-000000000007', NOW() - INTERVAL '12 days', NOW() - INTERVAL '5 days'),
('f4100001-0005-4000-8000-000000000005', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0008-4000-8000-000000000008', '吉隆坡零售 · POS 出海', 'contacted', 'website', 420000, (CURRENT_DATE + INTERVAL '45 days')::date, 'activate', 55, NOW() - INTERVAL '6 days', ARRAY['马来西亚', '出海'], 'medium', NULL, 'f1000001-0008-4000-8000-000000000008', NOW() - INTERVAL '18 days', NOW() - INTERVAL '6 days'),
('f4100001-0006-4000-8000-000000000006', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0009-4000-8000-000000000009', '浦发金科 · IDC 运维续签', 'converted', 'partner', 1560000, (CURRENT_DATE + INTERVAL '10 days')::date, 'retain', 85, NOW() - INTERVAL '1 day', ARRAY['数据中心实施部'], 'high', 'f3100001-0006-4000-8000-000000000006', 'f1000001-0009-4000-8000-000000000009', NOW() - INTERVAL '180 days', NOW() - INTERVAL '1 day'),
('f4100001-0007-4000-8000-000000000007', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0006-4000-8000-000000000006', '宁东供应链 · WMS 试点', 'new', 'cold_call', 280000, NULL, 'acquire', 35, NOW() - INTERVAL '10 days', ARRAY['南京小东'], 'low', NULL, 'f1000001-0006-4000-8000-000000000006', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days'),
('f4100001-0008-4000-8000-000000000008', 'f1111111-1111-4111-8111-111111111111', 'f1000001-000a-4000-8000-00000000000a', '集团客户 · 交付 SLA 升级包', 'qualified', 'referral', 380000, (CURRENT_DATE + INTERVAL '28 days')::date, 'grow', 72, NOW() - INTERVAL '2 days', ARRAY['产品交付部'], 'high', NULL, 'f1000001-000a-4000-8000-00000000000a', NOW() - INTERVAL '25 days', NOW() - INTERVAL '2 days'),
('f4100001-0009-4000-8000-000000000009', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0002-4000-8000-000000000002', '技术中心 · AI Copilot 内测扩容', 'contacted', 'inbound', 0, NULL, 'retain', 92, NOW() - INTERVAL '12 hours', ARRAY['技术中心'], 'high', 'f3100001-0007-4000-8000-000000000007', 'f1000001-0002-4000-8000-000000000002', NOW() - INTERVAL '5 days', NOW() - INTERVAL '12 hours')
ON CONFLICT (id) DO NOTHING;

-- ── 商机 ─────────────────────────────────────────────────────────────
INSERT INTO deals (
    id, tenant_id, owner_id, title, stage, amount, currency, probability,
    expected_close_date, account_id, lead_id, description, tags,
    created_by, created_at, updated_at
) VALUES
('f5100001-0001-4000-8000-000000000001', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0003-4000-8000-000000000003', '灵狐数据 · 华东数智中台签约', 'negotiation', 1280000, 'CNY', 72, (CURRENT_DATE + INTERVAL '18 days')::date, 'f3100001-0001-4000-8000-000000000001', 'f4100001-0001-4000-8000-000000000001', '数据中台 + BI 模块', ARRAY['灵狐数据', '重点'], 'f1000001-0003-4000-8000-000000000003', NOW() - INTERVAL '40 days', NOW() - INTERVAL '1 day'),
('f5100001-0002-4000-8000-000000000002', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0004-4000-8000-000000000004', '莫邪互娱 · 星耀 Q3 联运', 'proposal', 680000, 'CNY', 48, (CURRENT_DATE + INTERVAL '42 days')::date, 'f3100001-0002-4000-8000-000000000002', 'f4100001-0002-4000-8000-000000000002', '手游联运分成', ARRAY['莫邪互娱'], 'f1000001-0004-4000-8000-000000000004', NOW() - INTERVAL '20 days', NOW() - INTERVAL '2 days'),
('f5100001-0003-4000-8000-000000000003', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0005-4000-8000-000000000005', '神龙云计算 · 政务云二期', 'negotiation', 2200000, 'CNY', 78, (CURRENT_DATE + INTERVAL '14 days')::date, 'f3100001-0003-4000-8000-000000000003', 'f4100001-0003-4000-8000-000000000003', '政务专有云扩容', ARRAY['神龙云计算', '政企'], 'f1000001-0005-4000-8000-000000000005', NOW() - INTERVAL '75 days', NOW() - INTERVAL '3 hours'),
('f5100001-0004-4000-8000-000000000004', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0007-4000-8000-000000000007', '麒麟事业群 · MES 首单', 'qualification', 950000, 'CNY', 35, (CURRENT_DATE + INTERVAL '60 days')::date, 'f3100001-0004-4000-8000-000000000004', 'f4100001-0004-4000-8000-000000000004', '制造执行系统 POC', ARRAY['麒麟事业群'], 'f1000001-0007-4000-8000-000000000007', NOW() - INTERVAL '10 days', NOW() - INTERVAL '5 days'),
('f5100001-0005-4000-8000-000000000005', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0008-4000-8000-000000000008', '马来西亚 · 零售 POS 出海', 'proposal', 420000, 'CNY', 52, (CURRENT_DATE + INTERVAL '35 days')::date, 'f3100001-0005-4000-8000-000000000005', 'f4100001-0005-4000-8000-000000000005', '东南亚本地化部署', ARRAY['马来西亚', '出海'], 'f1000001-0008-4000-8000-000000000008', NOW() - INTERVAL '15 days', NOW() - INTERVAL '6 days'),
('f5100001-0006-4000-8000-000000000006', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0009-4000-8000-000000000009', '数据中心 · 浦发 IDC 续签', 'won', 1560000, 'CNY', 100, CURRENT_DATE, 'f3100001-0006-4000-8000-000000000006', 'f4100001-0006-4000-8000-000000000006', '年度运维与实施', ARRAY['数据中心实施部', '赢单'], 'f1000001-0009-4000-8000-000000000009', NOW() - INTERVAL '90 days', NOW() - INTERVAL '4 days'),
('f5100001-0007-4000-8000-000000000007', 'f1111111-1111-4111-8111-111111111111', 'f1000001-000a-4000-8000-00000000000a', '产品交付 · SLA 升级包', 'qualification', 380000, 'CNY', 40, (CURRENT_DATE + INTERVAL '25 days')::date, NULL, 'f4100001-0008-4000-8000-000000000008', '交付团队扩容', ARRAY['产品交付部'], 'f1000001-000a-4000-8000-00000000000a', NOW() - INTERVAL '18 days', NOW() - INTERVAL '2 days'),
('f5100001-0008-4000-8000-000000000008', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0003-4000-8000-000000000003', '灵狐数据 · 广告归因 POC', 'won', 320000, 'CNY', 100, (CURRENT_DATE - INTERVAL '8 days')::date, NULL, NULL, 'POC 转正式', ARRAY['灵狐数据', '赢单'], 'f1000001-0003-4000-8000-000000000003', NOW() - INTERVAL '50 days', NOW() - INTERVAL '8 days'),
('f5100001-0009-4000-8000-000000000009', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0005-4000-8000-000000000005', '神龙云计算 · 某省信创云', 'qualification', 3500000, 'CNY', 28, (CURRENT_DATE + INTERVAL '90 days')::date, NULL, NULL, '信创云招标跟进', ARRAY['神龙云计算', '战略'], 'f1000001-0005-4000-8000-000000000005', NOW() - INTERVAL '7 days', NOW() - INTERVAL '1 day'),
('f5100001-000a-4000-8000-00000000000a', 'f1111111-1111-4111-8111-111111111111', 'f1000001-0007-4000-8000-000000000007', '麒麟事业群 · 集团统采框架', 'won', 2100000, 'CNY', 100, (CURRENT_DATE - INTERVAL '15 days')::date, 'f3100001-0004-4000-8000-000000000004', NULL, '年度框架已签', ARRAY['麒麟事业群', '赢单'], 'f1000001-0007-4000-8000-000000000007', NOW() - INTERVAL '120 days', NOW() - INTERVAL '15 days')
ON CONFLICT (id) DO NOTHING;

UPDATE deals SET closed_at = updated_at
WHERE id IN (
    'f5100001-0006-4000-8000-000000000006',
    'f5100001-0008-4000-8000-000000000008',
    'f5100001-000a-4000-8000-00000000000a'
) AND closed_at IS NULL;

-- ── 跟进活动（旗舰线索情绪旅程）────────────────────────────────────
INSERT INTO activities (
    id, tenant_id, subject_type, subject_id, event_type, direction, body,
    metadata, sentiment, sentiment_source, occurred_at, created_by
) VALUES
('f7100001-0001-4000-8000-000000000001', 'f1111111-1111-4111-8111-111111111111', 'lead', 'f4100001-0001-4000-8000-000000000001', 'email', 'outbound', '发送数据中台方案与三年 ROI 测算。', '{"label":"邮件：方案与 ROI"}'::jsonb, 'positive', 'manual', NOW() - INTERVAL '25 days', 'f1000001-0003-4000-8000-000000000003'),
('f7100001-0002-4000-8000-000000000002', 'f1111111-1111-4111-8111-111111111111', 'lead', 'f4100001-0001-4000-8000-000000000001', 'meeting', 'outbound', '与陈总确认预算池，待集团 CFO 会签。', '{"label":"会议：预算确认"}'::jsonb, 'neutral', 'manual', NOW() - INTERVAL '18 days', 'f1000001-0003-4000-8000-000000000003'),
('f7100001-0003-4000-8000-000000000003', 'f1111111-1111-4111-8111-111111111111', 'lead', 'f4100001-0001-4000-8000-000000000001', 'call', 'inbound', '客户提及友商报价，需补充案例。', '{"label":"电话：竞品顾虑"}'::jsonb, 'hesitant', 'rule', NOW() - INTERVAL '10 days', 'f1000001-0003-4000-8000-000000000003'),
('f7100001-0004-4000-8000-000000000004', 'f1111111-1111-4111-8111-111111111111', 'lead', 'f4100001-0001-4000-8000-000000000001', 'wechat', 'inbound', '陈总反馈方案认可，推进合同评审。', '{"label":"微信：推进签约"}'::jsonb, 'positive', 'manual', NOW() - INTERVAL '1 day', 'f1000001-0003-4000-8000-000000000003'),
('f7100001-0005-4000-8000-000000000005', 'f1111111-1111-4111-8111-111111111111', 'deal', 'f5100001-0003-4000-8000-000000000003', 'meeting', 'outbound', '政务云二期技术评审通过。', '{"label":"会议：技术评审"}'::jsonb, 'positive', 'manual', NOW() - INTERVAL '3 hours', 'f1000001-0005-4000-8000-000000000005')
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM activities WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111';
DELETE FROM deals WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111';
DELETE FROM leads WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111';
DELETE FROM contacts WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111';
DELETE FROM accounts WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111';
DELETE FROM user_roles WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111';
DELETE FROM role_permissions WHERE role_id IN (
    'f2111111-1111-4111-8111-111111111111',
    'f2121111-1111-4111-8111-111111111111'
);
DELETE FROM roles WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111';
DELETE FROM user_tenants WHERE tenant_id = 'f1111111-1111-4111-8111-111111111111';
DELETE FROM users WHERE email LIKE '%@xiaoxi.com';
DELETE FROM tenants WHERE id = 'f1111111-1111-4111-8111-111111111111';
