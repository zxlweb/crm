package crm

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrInvalidSegmentCode = errors.New("invalid_segment_code")
	ErrSegmentNotLead     = errors.New("segment_not_applicable_to_leads")
	ErrSegmentNotAccount  = errors.New("segment_not_applicable_to_accounts")
)

// ValidSegmentCodes are Phase 2 preset templates (PRD §4.3.3 / API §2.5).
var ValidSegmentCodes = map[string]bool{
	"high_value":        true,
	"churn_risk":        true,
	"new_potential":     true,
	"needs_activation":  true,
	"revive_pool":       true,
}

func ValidSegmentCode(code string) bool { return ValidSegmentCodes[code] }

// SegmentApplyOpts carries tenant thresholds for segment filters.
type SegmentApplyOpts struct {
	DaysSilent      int
	HighValueAmount float64
}

func (o SegmentApplyOpts) daysSilent() int {
	if o.DaysSilent < 1 {
		return 7
	}
	return o.DaysSilent
}

func (o SegmentApplyOpts) highValueAmount() float64 {
	if o.HighValueAmount <= 0 {
		return 100000
	}
	return o.HighValueAmount
}

// ApplyLeadSegmentFilter applies preset segment logic to a leads query (alias: leads).
func ApplyLeadSegmentFilter(q *gorm.DB, code string, opts SegmentApplyOpts) error {
	if !ValidSegmentCode(code) {
		return ErrInvalidSegmentCode
	}
	switch code {
	case "high_value":
		return q.Where("leads.amount > ?", opts.highValueAmount()).Error
	case "churn_risk":
		cutoff := time.Now().UTC().AddDate(0, 0, -opts.daysSilent())
		return q.Where("leads.last_activity_at IS NULL OR leads.last_activity_at < ?", cutoff).Error
	case "new_potential":
		since := time.Now().UTC().AddDate(0, 0, -7)
		return q.Where("leads.created_at >= ? AND leads.status <> ?", since, "qualified").Error
	case "needs_activation":
		return q.Where("leads.lifecycle_stage = ?", "acquire").
			Where(`NOT EXISTS (
				SELECT 1 FROM activities a
				WHERE a.tenant_id = leads.tenant_id
				  AND a.subject_type = 'lead'
				  AND a.subject_id = leads.id
				  AND a.direction = 'outbound'
				  AND a.deleted_at IS NULL
			)`).Error
	case "revive_pool":
		return q.Where("leads.lifecycle_stage = ?", "revive").Error
	default:
		return ErrInvalidSegmentCode
	}
}

// ApplyAccountSegmentFilter applies preset segment logic to an accounts query (alias: accounts).
func ApplyAccountSegmentFilter(q *gorm.DB, code string, opts SegmentApplyOpts) error {
	if !ValidSegmentCode(code) {
		return ErrInvalidSegmentCode
	}
	switch code {
	case "high_value":
		// Accounts: high engagement as proxy for high value
		return q.Where("accounts.engagement_score >= ?", int16(70)).Error
	case "churn_risk":
		cutoff := time.Now().UTC().AddDate(0, 0, -opts.daysSilent())
		return q.Where("accounts.last_activity_at IS NULL OR accounts.last_activity_at < ?", cutoff).Error
	case "new_potential":
		since := time.Now().UTC().AddDate(0, 0, -7)
		return q.Where("accounts.created_at >= ? AND accounts.lifecycle_stage = ?", since, "acquire").Error
	case "needs_activation":
		return q.Where("accounts.lifecycle_stage = ?", "acquire").
			Where(`NOT EXISTS (
				SELECT 1 FROM activities a
				WHERE a.tenant_id = accounts.tenant_id
				  AND a.subject_type = 'account'
				  AND a.subject_id = accounts.id
				  AND a.direction = 'outbound'
				  AND a.deleted_at IS NULL
			)`).Error
	case "revive_pool":
		return q.Where("accounts.lifecycle_stage = ?", "revive").Error
	default:
		return ErrInvalidSegmentCode
	}
}

// SegmentEntityForCode returns which list types support the segment (all presets apply to both).
func SegmentEntityForCode(code string) (supportsLead, supportsAccount bool, err error) {
	if !ValidSegmentCode(code) {
		return false, false, ErrInvalidSegmentCode
	}
	return true, true, nil
}
