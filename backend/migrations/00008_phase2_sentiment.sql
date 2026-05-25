-- +goose Up
-- Phase 2.11: Activity sentiment field constraints + default tenant keyword rules

ALTER TABLE activities
    DROP CONSTRAINT IF EXISTS chk_activities_sentiment;

ALTER TABLE activities
    ADD CONSTRAINT chk_activities_sentiment
    CHECK (sentiment IS NULL OR sentiment IN ('positive', 'neutral', 'hesitant', 'negative', 'unknown'));

ALTER TABLE activities
    DROP CONSTRAINT IF EXISTS chk_activities_sentiment_source;

ALTER TABLE activities
    ADD CONSTRAINT chk_activities_sentiment_source
    CHECK (sentiment_source IS NULL OR sentiment_source IN ('manual', 'rule'));

-- Demo tenant: default keyword rules (PRD §4.6.2)
UPDATE tenants
SET config = config || '{
  "sentiment_keyword_rules": [
    {"keywords": ["太贵", "考虑一下", "犹豫", "再想想"], "sentiment": "hesitant"},
    {"keywords": ["投诉", "失望", "不满", "生气"], "sentiment": "negative"},
    {"keywords": ["满意", "感谢", "不错"], "sentiment": "positive"}
  ]
}'::jsonb
WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

-- +goose Down
ALTER TABLE activities DROP CONSTRAINT IF EXISTS chk_activities_sentiment;
ALTER TABLE activities DROP CONSTRAINT IF EXISTS chk_activities_sentiment_source;

UPDATE tenants
SET config = config - 'sentiment_keyword_rules'
WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
