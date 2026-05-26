-- +goose Up
-- 小西 · 各事业部月度配额（销售经理工作台，与 user_tenants.department 键一致）
UPDATE tenants
SET config = config || jsonb_build_object(
    'department_quotas', jsonb_build_object(
        '神龙云计算', jsonb_build_object(
            'amount', 8000000,
            'currency', 'CNY',
            'period', to_char(CURRENT_DATE, 'YYYY-MM')
        ),
        '灵狐数据', jsonb_build_object(
            'amount', 6000000,
            'currency', 'CNY',
            'period', to_char(CURRENT_DATE, 'YYYY-MM')
        ),
        '麒麟事业群', jsonb_build_object(
            'amount', 5000000,
            'currency', 'CNY',
            'period', to_char(CURRENT_DATE, 'YYYY-MM')
        ),
        '莫邪互娱', jsonb_build_object(
            'amount', 4000000,
            'currency', 'CNY',
            'period', to_char(CURRENT_DATE, 'YYYY-MM')
        ),
        '南京小东', jsonb_build_object(
            'amount', 3000000,
            'currency', 'CNY',
            'period', to_char(CURRENT_DATE, 'YYYY-MM')
        )
    )
)
WHERE domain = 'xiaoxi';

-- +goose Down
UPDATE tenants
SET config = config - 'department_quotas'
WHERE domain = 'xiaoxi';
