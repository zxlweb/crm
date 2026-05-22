-- +goose Up
-- Demo Corp 演示线索、活动、公司（与 apps/web/fixtures 对齐，仅本地/测试）

INSERT INTO accounts (
    id, tenant_id, owner_id, name, industry, website,
    lifecycle_stage, engagement_score, last_activity_at, tags,
    created_by, created_at, updated_at
) VALUES (
    'c1000000-0000-4000-8000-000000000001',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '明德医疗集团',
    'healthcare',
    'https://mingde-health.example.com',
    'retain',
    88,
    NOW() - INTERVAL '2 days',
    ARRAY['已转化', '大客户'],
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '90 days',
    NOW() - INTERVAL '2 days'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO accounts (
    id, tenant_id, owner_id, name, industry, website,
    lifecycle_stage, engagement_score, last_activity_at, tags,
    created_by, created_at, updated_at
) VALUES (
    'c1000000-0000-4000-8000-000000000002',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '云帆智造',
    'manufacturing',
    'https://yunfan.example.com',
    'activate',
    78,
    NOW() - INTERVAL '1 day',
    ARRAY['SaaS'],
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '20 days',
    NOW() - INTERVAL '1 day'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO leads (
    id, tenant_id, owner_id, title, status, source, amount, expected_close_date,
    lifecycle_stage, engagement_score, last_activity_at, tags, relationship_health,
    converted_account_id, created_by, created_at, updated_at
) VALUES
(
    'a1000000-0000-4000-8000-000000000001',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '⭐ 华创科技',
    'qualified',
    'referral',
    280000,
    '2026-06-30',
    'grow',
    62,
    NOW() - INTERVAL '3 days',
    ARRAY['重点', '华东'],
    'medium',
    NULL,
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '45 days',
    NOW() - INTERVAL '3 days'
),
(
    'a1000000-0000-4000-8000-000000000002',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '云帆智造',
    'contacted',
    'website',
    120000,
    '2026-07-15',
    'activate',
    78,
    NOW() - INTERVAL '1 day',
    ARRAY['SaaS'],
    'high',
    NULL,
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '20 days',
    NOW() - INTERVAL '1 day'
),
(
    'a1000000-0000-4000-8000-000000000003',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '北辰物流',
    'new',
    'exhibition',
    95000,
    NULL,
    'acquire',
    28,
    NOW() - INTERVAL '12 days',
    ARRAY[]::text[],
    'low',
    NULL,
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '8 days',
    NOW() - INTERVAL '8 days'
),
(
    'a1000000-0000-4000-8000-000000000004',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '星海教育',
    'unqualified',
    'cold_call',
    0,
    NULL,
    'acquire',
    15,
    NOW() - INTERVAL '30 days',
    ARRAY['低意向'],
    'low',
    NULL,
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '60 days',
    NOW() - INTERVAL '25 days'
),
(
    'a1000000-0000-4000-8000-000000000005',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '明德医疗',
    'converted',
    'partner',
    450000,
    '2026-04-01',
    'retain',
    88,
    NOW() - INTERVAL '2 days',
    ARRAY['已转化'],
    'high',
    'c1000000-0000-4000-8000-000000000001',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '90 days',
    NOW() - INTERVAL '2 days'
)
ON CONFLICT (id) DO NOTHING;

-- 华创科技：情绪旅程演示活动
INSERT INTO activities (
    id, tenant_id, subject_type, subject_id, event_type, direction, body,
    metadata, sentiment, sentiment_source, occurred_at, created_by
) VALUES
(
    'f1000000-0000-4000-8000-000000000001',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'lead',
    'a1000000-0000-4000-8000-000000000001',
    'email',
    'outbound',
    '产品方案已发送，附 ROI 测算表。',
    '{"label":"邮件：产品方案发送"}'::jsonb,
    'positive',
    'manual',
    NOW() - INTERVAL '28 days',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
),
(
    'f1000000-0000-4000-8000-000000000002',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'lead',
    'a1000000-0000-4000-8000-000000000001',
    'call',
    'inbound',
    '客户确认预算区间，待内部审批。',
    '{"label":"电话：预算范围确认"}'::jsonb,
    'neutral',
    'manual',
    NOW() - INTERVAL '21 days',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
),
(
    'f1000000-0000-4000-8000-000000000003',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'lead',
    'a1000000-0000-4000-8000-000000000001',
    'meeting',
    'outbound',
    '讨论价格与交付周期，客户对交付时间有顾虑。',
    '{"label":"会议：价格与交付周期"}'::jsonb,
    'hesitant',
    'rule',
    NOW() - INTERVAL '14 days',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
),
(
    'f1000000-0000-4000-8000-000000000004',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'lead',
    'a1000000-0000-4000-8000-000000000001',
    'wechat',
    'inbound',
    '提及竞品对比，需要补充案例材料。',
    '{"label":"微信：竞品对比顾虑"}'::jsonb,
    'negative',
    'manual',
    NOW() - INTERVAL '7 days',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
),
(
    'f1000000-0000-4000-8000-000000000005',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'lead',
    'a1000000-0000-4000-8000-000000000001',
    'call',
    'outbound',
    '跟进 ROI 材料阅读情况。',
    '{"label":"电话：ROI 材料跟进"}'::jsonb,
    'hesitant',
    'manual',
    NOW() - INTERVAL '3 days',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO activities (
    id, tenant_id, subject_type, subject_id, event_type, direction, body,
    metadata, sentiment, sentiment_source, occurred_at, created_by
) VALUES
(
    'f1000000-0000-4000-8000-000000000010',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'lead',
    'a1000000-0000-4000-8000-000000000002',
    'email',
    'outbound',
    '发送产品白皮书。',
    '{"label":"邮件：白皮书"}'::jsonb,
    'positive',
    'manual',
    NOW() - INTERVAL '5 days',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
),
(
    'f1000000-0000-4000-8000-000000000011',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'lead',
    'a1000000-0000-4000-8000-000000000002',
    'call',
    'inbound',
    '预约下周演示。',
    '{"label":"电话：预约演示"}'::jsonb,
    'positive',
    'manual',
    NOW() - INTERVAL '1 day',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
)
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM activities WHERE tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
  AND subject_id IN (
    'a1000000-0000-4000-8000-000000000001',
    'a1000000-0000-4000-8000-000000000002'
  );
DELETE FROM leads WHERE tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
  AND id IN (
    'a1000000-0000-4000-8000-000000000001',
    'a1000000-0000-4000-8000-000000000002',
    'a1000000-0000-4000-8000-000000000003',
    'a1000000-0000-4000-8000-000000000004',
    'a1000000-0000-4000-8000-000000000005'
  );
DELETE FROM accounts WHERE tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
  AND id IN (
    'c1000000-0000-4000-8000-000000000001',
    'c1000000-0000-4000-8000-000000000002'
  );
