-- +goose Up
-- 小西 · 神龙云计算事业部补充销售数据（本月赢单 + 在途，陈洋 shenlong@xiaoxi.com）

INSERT INTO accounts (
    id, tenant_id, owner_id, name, industry, website,
    lifecycle_stage, engagement_score, last_activity_at, tags,
    created_by, created_at, updated_at
) VALUES (
    'f3100001-0009-4000-8000-000000000009',
    'f1111111-1111-4111-8111-111111111111',
    'f1000001-0005-4000-8000-000000000005',
    '东湖高新区政务云运营中心',
    'government',
    NULL,
    'retain',
    86,
    NOW() - INTERVAL '2 days',
    ARRAY['神龙云计算', '政企', '赢单客户'],
    'f1000001-0005-4000-8000-000000000005',
    NOW() - INTERVAL '45 days',
    NOW() - INTERVAL '2 days'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO contacts (
    id, tenant_id, account_id, owner_id, first_name, last_name, email, phone,
    is_primary, lifecycle_stage, engagement_score, last_activity_at, tags,
    created_by, created_at, updated_at
) VALUES (
    'f6100001-0005-4000-8000-000000000005',
    'f1111111-1111-4111-8111-111111111111',
    'f3100001-0009-4000-8000-000000000009',
    'f1000001-0005-4000-8000-000000000005',
    '张',
    '主任',
    'zhang.zr@donghu-gov.example.com',
    '027-88880009',
    true,
    'retain',
    84,
    NOW() - INTERVAL '2 days',
    ARRAY['神龙云计算', '决策人'],
    'f1000001-0005-4000-8000-000000000005',
    NOW() - INTERVAL '40 days',
    NOW() - INTERVAL '2 days'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO leads (
    id, tenant_id, owner_id, title, status, source, amount, expected_close_date,
    lifecycle_stage, engagement_score, last_activity_at, tags, relationship_health,
    converted_account_id, created_by, created_at, updated_at
) VALUES (
    'f4100001-000a-4000-8000-00000000000a',
    'f1111111-1111-4111-8111-111111111111',
    'f1000001-0005-4000-8000-000000000005',
    '东湖高新 · 光谷政务云一期',
    'converted',
    'referral',
    1680000,
    (CURRENT_DATE - INTERVAL '5 days')::date,
    'retain',
    86,
    NOW() - INTERVAL '2 days',
    ARRAY['神龙云计算', '政企', '赢单'],
    'high',
    'f3100001-0009-4000-8000-000000000009',
    'f1000001-0005-4000-8000-000000000005',
    NOW() - INTERVAL '50 days',
    NOW() - INTERVAL '2 days'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO deals (
    id, tenant_id, owner_id, title, stage, amount, currency, probability,
    expected_close_date, account_id, lead_id, description, tags,
    created_by, created_at, updated_at
) VALUES
(
    'f5100001-000b-4000-8000-00000000000b',
    'f1111111-1111-4111-8111-111111111111',
    'f1000001-0005-4000-8000-000000000005',
    '神龙云计算 · 光谷政务云一期',
    'won',
    1680000,
    'CNY',
    100,
    (CURRENT_DATE - INTERVAL '5 days')::date,
    'f3100001-0009-4000-8000-000000000009',
    'f4100001-000a-4000-8000-00000000000a',
    '专有云一期建设已验收回款',
    ARRAY['神龙云计算', '政企', '赢单'],
    'f1000001-0005-4000-8000-000000000005',
    NOW() - INTERVAL '48 days',
    NOW() - INTERVAL '5 days'
),
(
    'f5100001-000c-4000-8000-00000000000c',
    'f1111111-1111-4111-8111-111111111111',
    'f1000001-0005-4000-8000-000000000005',
    '神龙云计算 · 湖北交投云资源续费',
    'proposal',
    960000,
    'CNY',
    55,
    (CURRENT_DATE + INTERVAL '32 days')::date,
    'f3100001-0003-4000-8000-000000000003',
    NULL,
    '混合云资源包年度续费，方案评审中',
    ARRAY['神龙云计算', '政企'],
    'f1000001-0005-4000-8000-000000000005',
    NOW() - INTERVAL '12 days',
    NOW() - INTERVAL '2 days'
)
ON CONFLICT (id) DO NOTHING;

UPDATE deals
SET closed_at = updated_at
WHERE id = 'f5100001-000b-4000-8000-00000000000b'
  AND closed_at IS NULL;

INSERT INTO activities (
    id, tenant_id, subject_type, subject_id, event_type, direction, body,
    metadata, sentiment, sentiment_source, occurred_at, created_by
) VALUES
(
    'f7100001-0006-4000-8000-000000000006',
    'f1111111-1111-4111-8111-111111111111',
    'deal',
    'f5100001-000b-4000-8000-00000000000b',
    'meeting',
    'outbound',
    '光谷政务云一期终验通过，财务确认本月回款。',
    '{"label":"会议：终验回款"}'::jsonb,
    'positive',
    'manual',
    NOW() - INTERVAL '5 days',
    'f1000001-0005-4000-8000-000000000005'
),
(
    'f7100001-0007-4000-8000-000000000007',
    'f1111111-1111-4111-8111-111111111111',
    'deal',
    'f5100001-000c-4000-8000-00000000000c',
    'email',
    'outbound',
    '向交投信息中心发送续费方案与 SLA 对比表。',
    '{"label":"邮件：续费方案"}'::jsonb,
    'positive',
    'manual',
    NOW() - INTERVAL '2 days',
    'f1000001-0005-4000-8000-000000000005'
)
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM activities WHERE id IN (
    'f7100001-0006-4000-8000-000000000006',
    'f7100001-0007-4000-8000-000000000007'
);
DELETE FROM deals WHERE id IN (
    'f5100001-000b-4000-8000-00000000000b',
    'f5100001-000c-4000-8000-00000000000c'
);
DELETE FROM leads WHERE id = 'f4100001-000a-4000-8000-00000000000a';
DELETE FROM contacts WHERE id = 'f6100001-0005-4000-8000-000000000005';
DELETE FROM accounts WHERE id = 'f3100001-0009-4000-8000-000000000009';
