-- +goose Up
-- Phase 2.10: 5 preset segment templates per tenant (PRD §4.3.3)

INSERT INTO segment_templates (id, tenant_id, code, name_i18n_key, filter_json, is_system)
SELECT gen_random_uuid(), t.id, s.code, s.name_key, s.filter_json::jsonb, true
FROM tenants t
CROSS JOIN (
    VALUES
        (
            'high_value',
            'segments.high_value.name',
            '{"entity":"lead","rules":[{"field":"amount","op":"gt","value_ref":"high_value_amount"}]}'
        ),
        (
            'churn_risk',
            'segments.churn_risk.name',
            '{"entity":"both","rules":[{"field":"days_since_last_activity","op":"gt","value_ref":"days_silent"}]}'
        ),
        (
            'new_potential',
            'segments.new_potential.name',
            '{"entity":"lead","rules":[{"field":"days_since_created","op":"lte","value":7},{"field":"status","op":"neq","value":"qualified"}]}'
        ),
        (
            'needs_activation',
            'segments.needs_activation.name',
            '{"entity":"both","rules":[{"field":"lifecycle_stage","op":"eq","value":"acquire"},{"field":"has_outbound_activity","op":"eq","value":false}]}'
        ),
        (
            'revive_pool',
            'segments.revive_pool.name',
            '{"entity":"both","rules":[{"field":"lifecycle_stage","op":"eq","value":"revive"}]}'
        )
) AS s(code, name_key, filter_json)
WHERE NOT EXISTS (
    SELECT 1 FROM segment_templates st
    WHERE st.tenant_id = t.id AND st.code = s.code
);

-- +goose Down
DELETE FROM segment_templates WHERE is_system = true;
