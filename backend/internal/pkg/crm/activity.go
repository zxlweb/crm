package crm

var ValidActivityEventTypes = map[string]bool{
	"note": true, "call": true, "email": true, "meeting": true,
	"wechat": true, "visit": true, "system": true,
}

var ValidActivityDirections = map[string]bool{
	"inbound": true, "outbound": true,
}

var ValidActivitySubjectTypes = map[string]bool{
	"lead": true, "contact": true, "account": true,
}

var ValidActivitySentiments = map[string]bool{
	"positive": true, "neutral": true, "hesitant": true, "negative": true, "unknown": true,
}

var ValidSentimentSources = map[string]bool{
	"manual": true, "rule": true,
}

func ValidActivityEventType(v string) bool   { return ValidActivityEventTypes[v] }
func ValidActivityDirection(v string) bool   { return v == "" || ValidActivityDirections[v] }
func ValidActivitySubjectType(v string) bool { return ValidActivitySubjectTypes[v] }
func ValidActivitySentiment(v string) bool   { return ValidActivitySentiments[v] }
func ValidSentimentSource(v string) bool {
	if v == "" {
		return true
	}
	return ValidSentimentSources[v]
}
