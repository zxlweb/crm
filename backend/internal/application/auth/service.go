package auth

import (
	"context"
	"errors"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/pkg/jwtutil"
	"crm-backend/internal/pkg/password"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTenantForbidden    = errors.New("tenant forbidden")
)

type TenantDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type UserDTO struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	IsSuperAdmin bool   `json:"is_super_admin"`
}

type LoginResult struct {
	AccessToken   string      `json:"access_token"`
	RefreshToken  string      `json:"refresh_token"`
	ExpiresIn     int64       `json:"expires_in"`
	User          UserDTO     `json:"user"`
	Tenants       []TenantDTO `json:"tenants"`
	CurrentTenant *TenantDTO  `json:"current_tenant,omitempty"`
}

type Service struct {
	users           repository.UserRepository
	jwtSecret       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewService(users repository.UserRepository, jwtSecret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		users:           users,
		jwtSecret:       jwtSecret,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *Service) Login(ctx context.Context, email, plainPassword string) (*LoginResult, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !password.Verify(user.PasswordHash, plainPassword) {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokens(ctx, user)
}

func (s *Service) Profile(ctx context.Context, userID uuid.UUID) (*UserDTO, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapUser(user), nil
}

func (s *Service) ListTenants(ctx context.Context, userID uuid.UUID) ([]TenantDTO, error) {
	rows, err := s.users.ListActiveTenantsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapTenants(rows), nil
}

func (s *Service) SwitchTenant(ctx context.Context, userID uuid.UUID, email string, isSuperAdmin bool, tenantID uuid.UUID) (*LoginResult, error) {
	ok, err := s.users.UserBelongsToTenant(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrTenantForbidden
	}
	user := &domain.User{
		ID:           userID,
		Email:        email,
		IsSuperAdmin: isSuperAdmin,
	}
	return s.issueTokensForTenant(ctx, user, &tenantID)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	claims, err := jwtutil.Parse(s.jwtSecret, refreshToken, jwtutil.TokenTypeRefresh)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	user, err := s.users.FindByEmail(ctx, claims.Email)
	if err != nil || user.ID != userID {
		return nil, ErrInvalidCredentials
	}
	var tenantID *uuid.UUID
	if claims.TenantID != "" {
		if tid, err := uuid.Parse(claims.TenantID); err == nil {
			tenantID = &tid
		}
	}
	return s.issueTokensForTenant(ctx, user, tenantID)
}

func (s *Service) issueTokens(ctx context.Context, user *domain.User) (*LoginResult, error) {
	return s.issueTokensForTenant(ctx, user, nil)
}

func (s *Service) issueTokensForTenant(ctx context.Context, user *domain.User, tenantID *uuid.UUID) (*LoginResult, error) {
	access, expiresIn, err := jwtutil.GenerateAccess(
		s.jwtSecret, user.ID, user.Email, user.IsSuperAdmin, tenantID, s.accessTokenTTL,
	)
	if err != nil {
		return nil, err
	}
	refresh, _, err := jwtutil.GenerateRefresh(
		s.jwtSecret, user.ID, user.Email, user.IsSuperAdmin, tenantID, s.refreshTokenTTL,
	)
	if err != nil {
		return nil, err
	}
	tenants, err := s.users.ListActiveTenantsForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	result := &LoginResult{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		User:         *mapUser(user),
		Tenants:      mapTenants(tenants),
	}
	if tenantID != nil {
		for _, t := range result.Tenants {
			if t.ID == tenantID.String() {
				copy := t
				result.CurrentTenant = &copy
				break
			}
		}
	}
	return result, nil
}

func mapUser(user *domain.User) *UserDTO {
	return &UserDTO{
		ID:           user.ID.String(),
		Email:        user.Email,
		Name:         user.Name,
		IsSuperAdmin: user.IsSuperAdmin,
	}
}

func mapTenants(rows []repository.TenantBrief) []TenantDTO {
	out := make([]TenantDTO, 0, len(rows))
	for _, t := range rows {
		out = append(out, TenantDTO{
			ID:     t.ID.String(),
			Name:   t.Name,
			Domain: t.Domain,
		})
	}
	return out
}
