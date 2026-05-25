package emotion

import (
	"errors"
	"time"
)

var ErrInvalidRange = errors.New("invalid_range")

// ParseRange maps query range=30d|90d|all to a lower bound (nil = all time).
func ParseRange(raw string) (since *time.Time, err error) {
	switch raw {
	case "", "90d":
		t := time.Now().UTC().AddDate(0, 0, -90)
		return &t, nil
	case "30d":
		t := time.Now().UTC().AddDate(0, 0, -30)
		return &t, nil
	case "all":
		return nil, nil
	default:
		return nil, ErrInvalidRange
	}
}
