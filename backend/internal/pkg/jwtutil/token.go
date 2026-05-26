package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrWrongType    = errors.New("wrong token type")
)

type Claims struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	TokenType    string `json:"token_type"`
	TenantID     string `json:"tenant_id,omitempty"`
	ActiveRoleID string `json:"active_role_id,omitempty"`
	jwt.RegisteredClaims
}

func GenerateAccess(secret string, userID uuid.UUID, email string, isSuperAdmin bool, tenantID, activeRoleID *uuid.UUID, ttl time.Duration) (string, int64, error) {
	return sign(secret, userID, email, isSuperAdmin, tenantID, activeRoleID, TokenTypeAccess, ttl)
}

func GenerateRefresh(secret string, userID uuid.UUID, email string, isSuperAdmin bool, tenantID, activeRoleID *uuid.UUID, ttl time.Duration) (string, int64, error) {
	return sign(secret, userID, email, isSuperAdmin, tenantID, activeRoleID, TokenTypeRefresh, ttl)
}

func sign(secret string, userID uuid.UUID, email string, isSuperAdmin bool, tenantID, activeRoleID *uuid.UUID, tokenType string, ttl time.Duration) (string, int64, error) {
	now := time.Now()
	exp := now.Add(ttl)
	claims := Claims{
		UserID:       userID.String(),
		Email:        email,
		IsSuperAdmin: isSuperAdmin,
		TokenType:    tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
	}
	if tenantID != nil && *tenantID != uuid.Nil {
		claims.TenantID = tenantID.String()
	}
	if activeRoleID != nil && *activeRoleID != uuid.Nil {
		claims.ActiveRoleID = activeRoleID.String()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return signed, int64(ttl.Seconds()), nil
}

func Parse(secret, tokenStr string, expectedType string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !t.Valid {
		return nil, ErrInvalidToken
	}
	if expectedType != "" && claims.TokenType != expectedType {
		return nil, ErrWrongType
	}
	return claims, nil
}
