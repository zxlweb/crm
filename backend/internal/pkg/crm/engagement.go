package crm

import "time"

// ComputeEngagementScore calculates 0–100 rule-based score (PRD §4.3.4 baseline).
func ComputeEngagementScore(status, lifecycle string, activitiesLast7d int) int16 {
	score := activitiesLast7d * 4
	if score > 40 {
		score = 40
	}
	switch status {
	case "contacted":
		score += 10
	case "qualified":
		score += 20
	case "converted":
		score += 30
	case "unqualified":
		score += 5
	}
	switch lifecycle {
	case "activate":
		score += 15
	case "grow":
		score += 25
	case "retain":
		score += 30
	case "revive":
		score += 10
	}
	if score > 100 {
		score = 100
	}
	return int16(score)
}

// DaysSince returns whole days since t; if t is nil returns a large value.
func DaysSince(t *time.Time) int {
	if t == nil {
		return 9999
	}
	d := time.Since(t.UTC())
	if d < 0 {
		return 0
	}
	return int(d.Hours() / 24)
}
