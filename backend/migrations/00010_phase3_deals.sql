-- +goose Up
-- Phase 3: deals table (phase-3-deals-dashboard-schema.md)

CREATE TABLE deals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    stage VARCHAR(32) NOT NULL DEFAULT 'qualification',
    amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(8) NOT NULL DEFAULT 'CNY',
    probability SMALLINT NOT NULL DEFAULT 0,
    expected_close_date DATE,
    account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    lead_id UUID REFERENCES leads(id) ON DELETE SET NULL,
    contact_id UUID REFERENCES contacts(id) ON DELETE SET NULL,
    description TEXT,
    lost_reason VARCHAR(500),
    closed_at TIMESTAMPTZ,
    engagement_score SMALLINT NOT NULL DEFAULT 0,
    last_activity_at TIMESTAMPTZ,
    tags TEXT[] NOT NULL DEFAULT '{}',
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_deals_stage CHECK (
        stage IN ('qualification', 'proposal', 'negotiation', 'won', 'lost')
    ),
    CONSTRAINT chk_deals_probability CHECK (probability >= 0 AND probability <= 100),
    CONSTRAINT chk_deals_amount CHECK (amount >= 0),
    CONSTRAINT chk_deals_currency CHECK (currency IN ('CNY', 'USD'))
);

CREATE INDEX idx_deals_tenant_owner ON deals(tenant_id, owner_id);
CREATE INDEX idx_deals_tenant_stage ON deals(tenant_id, stage);
CREATE INDEX idx_deals_tenant_account ON deals(tenant_id, account_id);
CREATE INDEX idx_deals_tenant_expected_close ON deals(tenant_id, expected_close_date);
CREATE INDEX idx_deals_pipeline_sort ON deals(tenant_id, deleted_at, updated_at DESC);
CREATE INDEX idx_deals_tenant_lead ON deals(tenant_id, lead_id);

-- +goose Down
DROP TABLE IF EXISTS deals;
