package crm

// ValidLeadStatuses per docs/api/phase-2-crm-ai.md §2.3
var ValidLeadStatuses = map[string]bool{
	"new":         true,
	"contacted":   true,
	"qualified":   true,
	"unqualified": true,
	"converted":   true,
}

func ValidLeadStatus(status string) bool {
	return ValidLeadStatuses[status]
}
