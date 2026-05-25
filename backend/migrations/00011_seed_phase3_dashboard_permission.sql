-- +goose Up
-- dashboard:view for Phase 3 summary endpoints (phase-3-deals-dashboard-api.md §6)

INSERT INTO permissions (id, resource, action, description) VALUES
    (gen_random_uuid(), 'dashboard', 'view', '查看仪表盘汇总')
ON CONFLICT (resource, action) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'dddddddd-dddd-dddd-dddd-dddddddddddd', p.id
FROM permissions p
WHERE p.resource = 'dashboard' AND p.action = 'view'
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE resource = 'dashboard' AND action = 'view'
);
DELETE FROM permissions WHERE resource = 'dashboard' AND action = 'view';
