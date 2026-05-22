package jwtutil

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateAndParseAccess(t *testing.T) {
	secret := "test-secret"
	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	token, exp, err := GenerateAccess(secret, uid, "a@b.com", true, nil, time.Hour)
	if err != nil || exp != 3600 {
		t.Fatalf("generate: err=%v exp=%d", err, exp)
	}

	claims, err := Parse(secret, token, TokenTypeAccess)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != uid.String() || claims.Email != "a@b.com" || !claims.IsSuperAdmin {
		t.Fatalf("claims: %+v", claims)
	}
}

func TestAccessTokenWithTenantID(t *testing.T) {
	secret := "test-secret"
	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	tid := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	token, _, err := GenerateAccess(secret, uid, "a@b.com", false, &tid, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := Parse(secret, token, TokenTypeAccess)
	if err != nil || claims.TenantID != tid.String() {
		t.Fatalf("tenant in claims: %+v err=%v", claims, err)
	}
}

func TestRefreshRejectedAsAccess(t *testing.T) {
	secret := "test-secret"
	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	refresh, _, err := GenerateRefresh(secret, uid, "a@b.com", false, nil, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(secret, refresh, TokenTypeAccess); err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}
