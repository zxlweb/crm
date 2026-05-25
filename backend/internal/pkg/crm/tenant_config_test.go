package crm

import "testing"

func TestInferSentimentFromBody(t *testing.T) {
	rules := DefaultSentimentKeywordRules()
	s, ok := InferSentimentFromBody("客户觉得太贵了，要再考虑一下", rules)
	if !ok || s != "hesitant" {
		t.Fatalf("got %q %v", s, ok)
	}
	s, ok = InferSentimentFromBody("非常失望，要投诉", rules)
	if !ok || s != "negative" {
		t.Fatalf("got %q %v", s, ok)
	}
}
