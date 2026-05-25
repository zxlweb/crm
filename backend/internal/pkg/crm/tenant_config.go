package crm

import (
	"encoding/json"
	"strings"

	"gorm.io/datatypes"
)

// SentimentKeywordRule maps body keywords to a sentiment (tenant config).
type SentimentKeywordRule struct {
	Keywords  []string `json:"keywords"`
	Sentiment string   `json:"sentiment"`
}

// InsightThresholds from tenants.config.insight_thresholds.
type InsightThresholds struct {
	DaysSilent      int     `json:"days_silent"`
	HighValueAmount float64 `json:"high_value_amount"`
}

// SalesQuota from tenants.config.sales_quota (Phase 3 Dashboard).
type SalesQuota struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Period   string  `json:"period"`
}

// TenantCRMConfig parses Phase 2 keys from tenants.config JSONB.
type TenantCRMConfig struct {
	InsightThresholds     InsightThresholds      `json:"insight_thresholds"`
	SentimentKeywordRules []SentimentKeywordRule `json:"sentiment_keyword_rules"`
	SalesQuota            SalesQuota             `json:"sales_quota"`
}

func ParseTenantCRMConfig(raw datatypes.JSON) TenantCRMConfig {
	cfg := TenantCRMConfig{
		InsightThresholds: InsightThresholds{DaysSilent: 7, HighValueAmount: 100000},
		SentimentKeywordRules: DefaultSentimentKeywordRules(),
	}
	if len(raw) == 0 {
		return cfg
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return cfg
	}
	if v, ok := m["insight_thresholds"]; ok {
		_ = json.Unmarshal(v, &cfg.InsightThresholds)
	}
	if cfg.InsightThresholds.DaysSilent < 1 {
		cfg.InsightThresholds.DaysSilent = 7
	}
	if v, ok := m["sentiment_keyword_rules"]; ok {
		var rules []SentimentKeywordRule
		if err := json.Unmarshal(v, &rules); err == nil && len(rules) > 0 {
			cfg.SentimentKeywordRules = rules
		}
	}
	if v, ok := m["sales_quota"]; ok {
		_ = json.Unmarshal(v, &cfg.SalesQuota)
	}
	return cfg
}

// DefaultSentimentKeywordRules matches PRD §4.6.2 examples.
func DefaultSentimentKeywordRules() []SentimentKeywordRule {
	return []SentimentKeywordRule{
		{Keywords: []string{"太贵", "考虑一下", "犹豫", "再想想"}, Sentiment: "hesitant"},
		{Keywords: []string{"投诉", "失望", "不满", "生气"}, Sentiment: "negative"},
		{Keywords: []string{"满意", "感谢", "不错"}, Sentiment: "positive"},
	}
}

// InferSentimentFromBody returns sentiment + true when a keyword rule matches.
func InferSentimentFromBody(body string, rules []SentimentKeywordRule) (string, bool) {
	if body == "" {
		return "", false
	}
	for _, rule := range rules {
		if !ValidActivitySentiment(rule.Sentiment) {
			continue
		}
		for _, kw := range rule.Keywords {
			if kw != "" && containsKeyword(body, kw) {
				return rule.Sentiment, true
			}
		}
	}
	return "", false
}

func containsKeyword(body, keyword string) bool {
	return keyword != "" && strings.Contains(body, keyword)
}
