package emotion

import (
	"context"
	"fmt"
	"strings"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

type Service struct {
	activities repository.ActivityRepository
}

func NewService(activities repository.ActivityRepository) *Service {
	return &Service{activities: activities}
}

// SubjectInput carries entity context for journey aggregation.
type SubjectInput struct {
	SubjectType      string
	SubjectID        uuid.UUID
	LifecycleCurrent string
	CreatedAt        time.Time
	ConvertedAt      *time.Time
}

func (s *Service) BuildJourney(ctx context.Context, tenantID uuid.UUID, in SubjectInput, rangeQuery string) (JourneyResponse, error) {
	since, err := ParseRange(rangeQuery)
	if err != nil {
		return JourneyResponse{}, err
	}
	lifecycle := in.LifecycleCurrent
	if lifecycle == "" {
		lifecycle = "acquire"
	}
	if s.activities == nil {
		return EmptyJourney(in.SubjectType, in.SubjectID, lifecycle), nil
	}

	acts, err := s.activities.ListForJourney(ctx, tenantID, in.SubjectType, in.SubjectID, since)
	if err != nil {
		return JourneyResponse{}, err
	}

	points := make([]JourneyPoint, 0, len(acts))
	for _, a := range acts {
		points = append(points, activityToPoint(a, lifecycle))
	}

	milestones := buildMilestones(in)
	bands := currentLifecycleBand(lifecycle, in.CreatedAt)
	summary := buildSummary(points)

	return JourneyResponse{
		SubjectType:      in.SubjectType,
		SubjectID:        in.SubjectID,
		LifecycleCurrent: lifecycle,
		LifecycleBands:   bands,
		Points:           points,
		Milestones:       milestones,
		Summary:          summary,
	}, nil
}

func activityToPoint(a domain.Activity, lifecycle string) JourneyPoint {
	pt := JourneyPoint{
		ActivityID:           a.ID,
		At:                   a.OccurredAt.UTC(),
		EventType:            a.EventType,
		Direction:            a.Direction,
		LifecycleStageAtTime: lifecycle,
		Label:                activityLabel(a),
	}
	if a.Sentiment != nil {
		s := *a.Sentiment
		pt.Sentiment = &s
	}
	pt.SentimentScore = SentimentScore(ptrStr(a.Sentiment))
	if a.SentimentSource != nil {
		pt.SentimentSource = *a.SentimentSource
	}
	return pt
}

func ptrStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func activityLabel(a domain.Activity) string {
	prefix := eventTypeLabel(a.EventType)
	body := strings.TrimSpace(a.Body)
	if body == "" {
		return prefix
	}
	if len(body) > 80 {
		body = body[:80] + "…"
	}
	return prefix + "：" + body
}

func eventTypeLabel(eventType string) string {
	switch eventType {
	case "call":
		return "电话"
	case "email":
		return "邮件"
	case "meeting":
		return "会议"
	case "note":
		return "备注"
	case "task":
		return "任务"
	default:
		return eventType
	}
}

func currentLifecycleBand(stage string, createdAt time.Time) []LifecycleBand {
	from := createdAt.UTC()
	if from.IsZero() {
		from = time.Now().UTC().AddDate(0, -3, 0)
	}
	return []LifecycleBand{{
		Stage: stage,
		From:  from,
		To:    time.Now().UTC(),
	}}
}

func buildMilestones(in SubjectInput) []Milestone {
	if in.ConvertedAt == nil {
		return []Milestone{}
	}
	return []Milestone{{
		Type:  "converted",
		At:    in.ConvertedAt.UTC(),
		Label: "转为商机",
	}}
}

func buildSummary(points []JourneyPoint) JourneySummary {
	if len(points) == 0 {
		return JourneySummary{Trend: "flat"}
	}

	var scored []int
	var lastSentiment *string
	for i := len(points) - 1; i >= 0; i-- {
		p := points[i]
		if p.Sentiment != nil {
			lastSentiment = p.Sentiment
			break
		}
	}
	for _, p := range points {
		if p.SentimentScore != nil {
			scored = append(scored, *p.SentimentScore)
		}
	}

	summary := JourneySummary{
		CurrentSentiment: lastSentiment,
		Trend:            computeTrend(scored),
	}
	if days := daysSincePositive(points); days != nil {
		summary.DaysSincePositive = days
	}
	return summary
}

func computeTrend(scores []int) string {
	if len(scores) < 2 {
		return "flat"
	}
	last := scores[len(scores)-1]
	prev := scores[len(scores)-2]
	switch {
	case last > prev:
		return "up"
	case last < prev:
		return "down"
	default:
		return "flat"
	}
}

func daysSincePositive(points []JourneyPoint) *int {
	now := time.Now().UTC()
	for i := len(points) - 1; i >= 0; i-- {
		p := points[i]
		if p.Sentiment != nil && isPositiveSentiment(*p.Sentiment) {
			d := int(now.Sub(p.At).Hours() / 24)
			if d < 0 {
				d = 0
			}
			return &d
		}
	}
	return nil
}

// FormatRangeError returns API-facing range error text.
func FormatRangeError() string {
	return fmt.Sprintf("invalid range, use 30d|90d|all")
}
