-- +goose Up
-- Demo tenant deals for Pipeline / Dashboard (aligns with phase-3 PRD fixtures)

INSERT INTO deals (
    id, tenant_id, owner_id, title, stage, amount, currency, probability,
    expected_close_date, account_id, lead_id, description, tags,
    created_by, created_at, updated_at
) VALUES
(
    'd1000000-0000-4000-8000-000000000001',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '云帆智造年度订阅',
    'qualification',
    280000.00,
    'CNY',
    40,
    (CURRENT_DATE + INTERVAL '52 days')::date,
    'c1000000-0000-4000-8000-000000000002',
    'a1000000-0000-4000-8000-000000000002',
    'SaaS 年度合作意向',
    ARRAY['SaaS', '重点'],
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '12 days',
    NOW() - INTERVAL '1 day'
),
(
    'd1000000-0000-4000-8000-000000000002',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '明德医疗集团扩容',
    'proposal',
    450000.00,
    'CNY',
    55,
    (CURRENT_DATE + INTERVAL '30 days')::date,
    'c1000000-0000-4000-8000-000000000001',
    NULL,
    '医疗信息化扩容方案',
    ARRAY['医疗'],
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '25 days',
    NOW() - INTERVAL '2 days'
),
(
    'd1000000-0000-4000-8000-000000000003',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '华创科技 POC 转正式',
    'negotiation',
    180000.00,
    'CNY',
    70,
    (CURRENT_DATE + INTERVAL '14 days')::date,
    NULL,
    'a1000000-0000-4000-8000-000000000001',
    'POC 通过后商务谈判',
    ARRAY[]::text[],
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '40 days',
    NOW() - INTERVAL '3 hours'
),
(
    'd1000000-0000-4000-8000-000000000004',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '北辰物流 Q2 签约',
    'won',
    320000.00,
    'CNY',
    100,
    CURRENT_DATE,
    NULL,
    NULL,
    NULL,
    ARRAY['赢单'],
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    NOW() - INTERVAL '60 days',
    NOW() - INTERVAL '5 days'
)
ON CONFLICT (id) DO NOTHING;

UPDATE deals
SET closed_at = updated_at
WHERE id = 'd1000000-0000-4000-8000-000000000004'
  AND closed_at IS NULL;

-- +goose Down
DELETE FROM deals
WHERE tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
  AND id IN (
    'd1000000-0000-4000-8000-000000000001',
    'd1000000-0000-4000-8000-000000000002',
    'd1000000-0000-4000-8000-000000000003',
    'd1000000-0000-4000-8000-000000000004'
  );
