package crm

// Lead status transitions per phase-2-crm-ai / QA matrix (PATCH only; converted via /convert).
var leadStatusTransitions = map[string]map[string]bool{
	"new":         {"contacted": true},
	"contacted":   {"qualified": true, "unqualified": true},
	"qualified":   {"unqualified": true},
	"unqualified": {},
	"converted":   {},
}

// CanTransitionLeadStatus reports whether PATCH may change status from → to.
// Setting converted is only allowed through POST /leads/:id/convert.
func CanTransitionLeadStatus(from, to string) bool {
	if from == to {
		return true
	}
	if to == "converted" {
		return false
	}
	next, ok := leadStatusTransitions[from]
	if !ok {
		return false
	}
	return next[to]
}

// CanConvertLead reports whether the lead may be converted from its current status.
func CanConvertLead(status string) bool {
	return status == "qualified"
}
