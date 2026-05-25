package emotion

import "strings"

// SentimentScore maps API §7 sentiment codes to numeric scores.
func SentimentScore(sentiment string) *int {
	switch strings.ToLower(strings.TrimSpace(sentiment)) {
	case "positive":
		v := 2
		return &v
	case "neutral":
		v := 0
		return &v
	case "hesitant":
		v := -1
		return &v
	case "negative":
		v := -2
		return &v
	default:
		return nil
	}
}

func isPositiveSentiment(sentiment string) bool {
	return strings.EqualFold(sentiment, "positive")
}
