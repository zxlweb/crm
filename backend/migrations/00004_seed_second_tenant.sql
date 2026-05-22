-- +goose Up
INSERT INTO tenants (id, name, domain, is_active)
VALUES (
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    'Acme Inc',
    'acme',
    true
) ON CONFLICT (domain) DO NOTHING;

-- +goose Down
DELETE FROM tenants WHERE domain = 'acme';
