package crm

import "testing"

func TestResolveQuotaTarget_Department(t *testing.T) {
	cfg := TenantCRMConfig{
		SalesQuota: SalesQuota{Amount: 28_000_000, Period: "2026-05"},
		DepartmentQuotas: map[string]SalesQuota{
			"神龙云计算": {Amount: 8_000_000, Period: "2026-05"},
		},
	}
	target, period, scope := cfg.ResolveQuotaTarget("department", "神龙云计算")
	if scope != "department" || target != 8_000_000 || period != "2026-05" {
		t.Fatalf("got target=%v period=%q scope=%q", target, period, scope)
	}
	target, _, scope = cfg.ResolveQuotaTarget("department", "未知部门")
	if scope != "department" || target != 0 {
		t.Fatalf("unknown dept: target=%v scope=%q", target, scope)
	}
	target, _, scope = cfg.ResolveQuotaTarget("all", "")
	if scope != "tenant" || target != 28_000_000 {
		t.Fatalf("tenant: target=%v scope=%q", target, scope)
	}
}
