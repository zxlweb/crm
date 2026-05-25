-- +goose Up
-- Demo 租户月度销售配额（Dashboard quota Gauge 联调）
UPDATE tenants
SET config = COALESCE(config, '{}'::jsonb) || jsonb_build_object(
    'sales_quota', jsonb_build_object(
        'amount', 5000000,
        'currency', 'CNY',
        'period', to_char(CURRENT_DATE, 'YYYY-MM')
    )
)
WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

-- +goose Down
UPDATE tenants
SET config = config - 'sales_quota'
WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
