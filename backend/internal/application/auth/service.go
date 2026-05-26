package auth

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"crm-backend/internal/application/audit"
	"crm-backend/internal/domain"
	"crm-backend/internal/infrastructure/persistence"
	"crm-backend/internal/pkg/jwtutil"
	"crm-backend/internal/pkg/password"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTenantForbidden    = errors.New("tenant forbidden")
	ErrRoleForbidden      = errors.New("role forbidden")
	ErrEmailExists        = errors.New("email already exists")
	ErrDomainExists       = errors.New("domain already exists")
	ErrInvalidDomain      = errors.New("invalid domain")
)

var domainPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,48}[a-z0-9])?$`)

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

type RoleDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LoginResult struct {
	AccessToken   string      `json:"access_token"`
	RefreshToken  string      `json:"refresh_token"`
	ExpiresIn     int64       `json:"expires_in"`
	User          UserDTO     `json:"user"`
	Tenants       []TenantDTO `json:"tenants"`
	CurrentTenant *TenantDTO  `json:"current_tenant,omitempty"`
	Roles         []RoleDTO   `json:"roles,omitempty"`
	CurrentRole   *RoleDTO    `json:"current_role,omitempty"`
	Department    string      `json:"department,omitempty"`
}

type ServiceDeps struct {
	DB       *gorm.DB
	Enforcer *casbin.Enforcer
	Audit    *audit.Recorder
	RBAC     repository.RBACRepository
}

type Service struct {
	users           repository.UserRepository
	rbac            repository.RBACRepository
	jwtSecret       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	db              *gorm.DB
	enforcer        *casbin.Enforcer
	audit           *audit.Recorder
}

func NewService(users repository.UserRepository, jwtSecret string, accessTTL, refreshTTL time.Duration, deps *ServiceDeps) *Service {
	s := &Service{
		users:           users,
		jwtSecret:       jwtSecret,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
	if deps != nil {
		s.db = deps.DB
		s.enforcer = deps.Enforcer
		s.audit = deps.Audit
		s.rbac = deps.RBAC
	}
	return s
}

type RegisterInput struct {
	Email       string
	Password    string
	Name        string
	CompanyName string
	Domain      string
}

func (s *Service) Register(ctx context.Context, in RegisterInput, clientIP string) (*LoginResult, error) {
	domainSlug := normalizeDomain(in.Domain, in.CompanyName)
	if !domainPattern.MatchString(domainSlug) {
		return nil, ErrInvalidDomain
	}
	hash, err := password.Hash(in.Password)
	if err != nil {
		return nil, err
	}
	user, tenantID, err := s.users.RegisterWithTenant(ctx, repository.RegisterInput{
		Email:        strings.TrimSpace(strings.ToLower(in.Email)),
		PasswordHash: hash,
		Name:         strings.TrimSpace(in.Name),
		CompanyName:  strings.TrimSpace(in.CompanyName),
		Domain:       domainSlug,
	})
	if err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			return nil, ErrEmailExists
		}
		if errors.Is(err, repository.ErrDomainExists) {
			return nil, ErrDomainExists
		}
		return nil, err
	}
	if s.db != nil && s.enforcer != nil {
		if err := persistence.SyncCasbinPolicies(s.db, s.enforcer); err != nil {
			return nil, err
		}
	}
	s.recordAuth(ctx, tenantID, &user.ID, "auth.register", map[string]string{
		"email": user.Email, "domain": domainSlug,
	}, clientIP)
	return s.issueTokensForTenant(ctx, user, &tenantID, nil)
}

func (s *Service) Login(ctx context.Context, email, plainPassword, clientIP string) (*LoginResult, error) {
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
	result, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, err
	}
	// 单租户用户：登录响应直接带上 tenant/role/department，避免 FE 仅写 cookie 导致 Casbin 租户不一致
	if len(result.Tenants) == 1 {
		if tid, parseErr := uuid.Parse(result.Tenants[0].ID); parseErr == nil {
			if bound, bindErr := s.issueTokensForTenant(ctx, user, &tid, nil); bindErr == nil {
				result = bound
			}
		}
	}
	s.auditLoginWithIP(ctx, user.ID, result, clientIP)
	return result, nil
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

func (s *Service) SwitchTenant(ctx context.Context, userID uuid.UUID, email string, isSuperAdmin bool, tenantID uuid.UUID, clientIP string) (*LoginResult, error) {
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
	result, err := s.issueTokensForTenant(ctx, user, &tenantID, nil)
	if err != nil {
		return nil, err
	}
	s.recordAuth(ctx, tenantID, &userID, "auth.switch_tenant", map[string]string{
		"tenant_id": tenantID.String(),
	}, clientIP)
	return result, nil
}

func (s *Service) ListMyRoles(ctx context.Context, userID, tenantID uuid.UUID) ([]RoleDTO, error) {
	if s.rbac == nil {
		return []RoleDTO{}, nil
	}
	roles, err := s.rbac.ListUserRoles(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	return mapAuthRoles(roles), nil
}

func (s *Service) SwitchRole(ctx context.Context, userID uuid.UUID, email string, isSuperAdmin bool, tenantID, roleID uuid.UUID, clientIP string) (*LoginResult, error) {
	if s.rbac == nil {
		return nil, ErrRoleForbidden
	}
	ok, err := s.rbac.UserHasRole(ctx, tenantID, userID, roleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrRoleForbidden
	}
	user := &domain.User{
		ID:           userID,
		Email:        email,
		IsSuperAdmin: isSuperAdmin,
	}
	result, err := s.issueTokensForTenant(ctx, user, &tenantID, &roleID)
	if err != nil {
		return nil, err
	}
	s.recordAuth(ctx, tenantID, &userID, "auth.switch_role", map[string]string{
		"role_id": roleID.String(),
	}, clientIP)
	return result, nil
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
	var preferredRole *uuid.UUID
	if claims.ActiveRoleID != "" {
		if rid, err := uuid.Parse(claims.ActiveRoleID); err == nil {
			preferredRole = &rid
		}
	}
	return s.issueTokensForTenant(ctx, user, tenantID, preferredRole)
}

func (s *Service) issueTokens(ctx context.Context, user *domain.User) (*LoginResult, error) {
	return s.issueTokensForTenant(ctx, user, nil, nil)
}

func (s *Service) issueTokensForTenant(ctx context.Context, user *domain.User, tenantID, preferredRoleID *uuid.UUID) (*LoginResult, error) {
	var activeRoleID *uuid.UUID
	var roles []RoleDTO
	if tenantID != nil && s.rbac != nil {
		var err error
		activeRoleID, roles, err = s.pickActiveRole(ctx, *tenantID, user.ID, preferredRoleID)
		if err != nil {
			return nil, err
		}
	}
	access, expiresIn, err := jwtutil.GenerateAccess(
		s.jwtSecret, user.ID, user.Email, user.IsSuperAdmin, tenantID, activeRoleID, s.accessTokenTTL,
	)
	if err != nil {
		return nil, err
	}
	refresh, _, err := jwtutil.GenerateRefresh(
		s.jwtSecret, user.ID, user.Email, user.IsSuperAdmin, tenantID, activeRoleID, s.refreshTokenTTL,
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
	if len(roles) > 0 {
		result.Roles = roles
		if activeRoleID != nil {
			for _, r := range roles {
				if r.ID == activeRoleID.String() {
					copy := r
					result.CurrentRole = &copy
					break
				}
			}
		}
	}
	if tenantID != nil {
		if dept, err := s.users.UserDepartment(ctx, *tenantID, user.ID); err == nil && dept != "" {
			result.Department = dept
		}
	}
	return result, nil
}

func (s *Service) pickActiveRole(ctx context.Context, tenantID, userID uuid.UUID, preferred *uuid.UUID) (*uuid.UUID, []RoleDTO, error) {
	roleRows, err := s.rbac.ListUserRoles(ctx, tenantID, userID)
	if err != nil {
		return nil, nil, err
	}
	roles := mapAuthRoles(roleRows)
	if len(roles) == 0 {
		return nil, roles, nil
	}
	if preferred != nil {
		for _, r := range roleRows {
			if r.ID == *preferred {
				return preferred, roles, nil
			}
		}
	}
	firstID, err := uuid.Parse(roles[0].ID)
	if err != nil {
		return nil, roles, err
	}
	return &firstID, roles, nil
}

func mapAuthRoles(rows []domain.Role) []RoleDTO {
	out := make([]RoleDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, RoleDTO{
			ID:          r.ID.String(),
			Name:        r.Name,
			Description: r.Description,
		})
	}
	return out
}

func mapUser(user *domain.User) *UserDTO {
	return &UserDTO{
		ID:           user.ID.String(),
		Email:        user.Email,
		Name:         user.Name,
		IsSuperAdmin: user.IsSuperAdmin,
	}
}

func normalizeDomain(domain, companyName string) string {
	d := strings.TrimSpace(strings.ToLower(domain))
	if d != "" {
		return d
	}
	d = strings.ToLower(companyName)
	d = strings.ReplaceAll(d, " ", "-")
	var b strings.Builder
	for _, r := range d {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "org"
	}
	if len(out) > 50 {
		out = out[:50]
	}
	return out
}

func (s *Service) auditLoginWithIP(ctx context.Context, userID uuid.UUID, result *LoginResult, ip string) {
	if result == nil {
		return
	}
	tenantID := firstTenantID(result)
	if tenantID == uuid.Nil {
		return
	}
	s.recordAuth(ctx, tenantID, &userID, "auth.login", map[string]string{"email": result.User.Email}, ip)
}

func firstTenantID(result *LoginResult) uuid.UUID {
	if result.CurrentTenant != nil {
		if id, err := uuid.Parse(result.CurrentTenant.ID); err == nil {
			return id
		}
	}
	if len(result.Tenants) > 0 {
		if id, err := uuid.Parse(result.Tenants[0].ID); err == nil {
			return id
		}
	}
	return uuid.Nil
}

func (s *Service) recordAuth(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, action string, newValue any, ip string) {
	if s.audit == nil {
		return
	}
	s.audit.Record(ctx, audit.Entry{
		TenantID:     tenantID,
		UserID:       userID,
		Action:       action,
		ResourceType: "auth",
		NewValue:     newValue,
		IPAddress:    ip,
	})
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
