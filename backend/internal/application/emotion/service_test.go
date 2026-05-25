package emotion

import (
	"context"
	"testing"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

type stubActivityRepo struct {
	items []domain.Activity
}

func (s *stubActivityRepo) List(context.Context, uuid.UUID, repository.ActivityListFilter) ([]domain.Activity, int64, error) {
	return nil, 0, nil
}

func (s *stubActivityRepo) ListForJourney(_ context.Context, _ uuid.UUID, subjectType string, subjectID uuid.UUID, since *time.Time) ([]domain.Activity, error) {
	out := make([]domain.Activity, 0, len(s.items))
	for _, a := range s.items {
		if a.SubjectType != subjectType || a.SubjectID != subjectID {
			continue
		}
		if since != nil && a.OccurredAt.Before(*since) {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

func (s *stubActivityRepo) GetByID(context.Context, uuid.UUID, uuid.UUID) (*domain.Activity, error) {
	return nil, repository.ErrActivityNotFound
}

func (s *stubActivityRepo) Create(context.Context, *domain.Activity) error { return nil }
func (s *stubActivityRepo) Update(context.Context, *domain.Activity) error { return nil }
func (s *stubActivityRepo) SoftDelete(context.Context, uuid.UUID, uuid.UUID) error { return nil }
func (s *stubActivityRepo) SummaryByEventType(context.Context, uuid.UUID, string, uuid.UUID) ([]repository.LabelCount, int64, error) {
	return nil, 0, nil
}
func (s *stubActivityRepo) LatestOccurredAt(context.Context, uuid.UUID, string, uuid.UUID) (*time.Time, error) {
	return nil, nil
}
func (s *stubActivityRepo) CountBySubject(context.Context, uuid.UUID, string, uuid.UUID) (int64, error) {
	return 0, nil
}

func (s *stubActivityRepo) CountSince(context.Context, uuid.UUID, time.Time) (int64, error) {
	return 0, nil
}

func (s *stubActivityRepo) CountLeadTouchesSince(context.Context, uuid.UUID, time.Time, bool, uuid.UUID) (int64, error) {
	return 0, nil
}

func (s *stubActivityRepo) CountAccountTouchesSince(context.Context, uuid.UUID, time.Time, bool, uuid.UUID) (int64, error) {
	return 0, nil
}

func TestBuildJourney_AggregatesPointsAndSummary(t *testing.T) {
	subjectID := uuid.New()
	pos := "positive"
	hes := "hesitant"
	repo := &stubActivityRepo{items: []domain.Activity{
		{
			ID: uuid.New(), SubjectType: "lead", SubjectID: subjectID,
			EventType: "call", OccurredAt: time.Now().UTC().Add(-48 * time.Hour), Sentiment: &pos,
		},
		{
			ID: uuid.New(), SubjectType: "lead", SubjectID: subjectID,
			EventType: "email", OccurredAt: time.Now().UTC().Add(-24 * time.Hour), Sentiment: &hes,
		},
	}}
	svc := NewService(repo)
	j, err := svc.BuildJourney(context.Background(), uuid.New(), SubjectInput{
		SubjectType: "lead", SubjectID: subjectID, LifecycleCurrent: "grow",
		CreatedAt: time.Now().UTC().Add(-90 * 24 * time.Hour),
	}, "90d")
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Points) != 2 {
		t.Fatalf("points: %d", len(j.Points))
	}
	if j.Summary.Trend != "down" {
		t.Fatalf("trend: %s", j.Summary.Trend)
	}
	if j.Summary.CurrentSentiment == nil || *j.Summary.CurrentSentiment != "hesitant" {
		t.Fatalf("current sentiment: %v", j.Summary.CurrentSentiment)
	}
}

func TestBuildJourney_ConvertedMilestone(t *testing.T) {
	subjectID := uuid.New()
	convertedAt := time.Now().UTC().Add(-time.Hour)
	svc := NewService(&stubActivityRepo{})
	j, err := svc.BuildJourney(context.Background(), uuid.New(), SubjectInput{
		SubjectType: "lead", SubjectID: subjectID, LifecycleCurrent: "grow",
		CreatedAt: time.Now().UTC().Add(-30 * 24 * time.Hour), ConvertedAt: &convertedAt,
	}, "all")
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Milestones) != 1 || j.Milestones[0].Type != "converted" {
		t.Fatalf("milestones: %+v", j.Milestones)
	}
}

func TestParseRange_Invalid(t *testing.T) {
	if _, err := ParseRange("1y"); err != ErrInvalidRange {
		t.Fatalf("err: %v", err)
	}
}

func TestSentimentScore_UnknownNil(t *testing.T) {
	if SentimentScore("unknown") != nil {
		t.Fatal("unknown should map to nil score")
	}
}
