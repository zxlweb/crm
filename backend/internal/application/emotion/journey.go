package emotion

import (
	"time"

	"github.com/google/uuid"
)

// JourneyResponse matches phase-2-crm-ai.md §7.
type JourneyResponse struct {
	SubjectType      string           `json:"subject_type"`
	SubjectID        uuid.UUID        `json:"subject_id"`
	LifecycleCurrent string           `json:"lifecycle_current"`
	LifecycleBands   []LifecycleBand  `json:"lifecycle_bands"`
	Points           []JourneyPoint   `json:"points"`
	Milestones       []Milestone      `json:"milestones"`
	Summary          JourneySummary   `json:"summary"`
}

type LifecycleBand struct {
	Stage string    `json:"stage"`
	From  time.Time `json:"from"`
	To    time.Time `json:"to"`
}

type JourneyPoint struct {
	ActivityID          uuid.UUID `json:"activity_id"`
	At                  time.Time `json:"at"`
	EventType           string    `json:"event_type"`
	Direction           string    `json:"direction,omitempty"`
	Sentiment           *string   `json:"sentiment"`
	SentimentScore      *int      `json:"sentiment_score"`
	SentimentSource     string    `json:"sentiment_source,omitempty"`
	Label               string    `json:"label,omitempty"`
	LifecycleStageAtTime string   `json:"lifecycle_stage_at_time,omitempty"`
}

type Milestone struct {
	Type  string    `json:"type"`
	At    time.Time `json:"at"`
	Label string    `json:"label"`
}

type JourneySummary struct {
	CurrentSentiment   *string `json:"current_sentiment"`
	Trend              string  `json:"trend"`
	DaysSincePositive  *int    `json:"days_since_positive,omitempty"`
}

// EmptyJourney returns a valid empty journey when activities are not recorded yet.
func EmptyJourney(subjectType string, subjectID uuid.UUID, lifecycle string) JourneyResponse {
	if lifecycle == "" {
		lifecycle = "acquire"
	}
	return JourneyResponse{
		SubjectType:      subjectType,
		SubjectID:        subjectID,
		LifecycleCurrent: lifecycle,
		LifecycleBands:   []LifecycleBand{},
		Points:           []JourneyPoint{},
		Milestones:       []Milestone{},
		Summary: JourneySummary{
			Trend: "flat",
		},
	}
}

// EmptyAccountJourney returns a valid empty journey for accounts without activities yet.
func EmptyAccountJourney(subjectID uuid.UUID, lifecycle string) JourneyResponse {
	return EmptyJourney("account", subjectID, lifecycle)
}

// EmptyLeadJourney returns a valid empty journey for leads without activities yet.
func EmptyLeadJourney(subjectID uuid.UUID, lifecycle string) JourneyResponse {
	return EmptyJourney("lead", subjectID, lifecycle)
}
