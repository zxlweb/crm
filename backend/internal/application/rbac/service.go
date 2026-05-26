package rbac

import (
	"context"
	"errors"
	"time"

	"crm-backend/internal/application/appscope"
	"crm-backend/internal/domain"
	"crm-backend/internal/infrastructure/persistence"
	"crm-backend/internal/pkg/activerole"
	"crm-backend/internal/pkg/datascope"
	"crm-backend/internal/pkg/rbacutil"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrInvalidPermissions = errors.New("invalid permission ids")
	ErrInvalidRoles       = errors.New("invalid role ids")
	ErrRoleForbidden      = errors.New("role forbidden")
	ErrMemberNotInTenant  = errors.New("user not in tenant")
)

type PermissionGroup struct {
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

type PermissionItemDTO struct {
	ID          string `json:"id"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type RoleDTO struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	IsSystem      bool     `json:"is_system"`
	PermissionIDs []string `json:"permission_ids"`
	UserCount     int64    `json:"user_count"`
}

type MemberRoleDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsSystem bool   `json:"is_system"`
}

type MemberDTO struct {
	ID         string          `json:"id"`
	Email      string          `json:"email"`
	Name       string          `json:"name"`
	AvatarURL  string          `json:"avatar_url"`
	Department string          `json:"department,omitempty"`
	Roles      []MemberRoleDTO `json:"roles"`
	JoinedAt   string          `json:"joined_at"`
}

type CheckRequest struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type CheckResult struct {
	Allowed bool `json:"allowed"`
}

type Service struct {
	repo      repository.RBACRepository
	db        *gorm.DB
	enforcer  *casbin.Enforcer
	scope     appscope.Provider
}

func NewService(repo repository.RBACRepository, db *gorm.DB, enforcer *casbin.Enforcer, scope appscope.Provider) *Service {
	return &Service{repo: repo, db: db, enforcer: enforcer, scope: scope}
}

func (s *Service) ListPermissionDictionary(ctx context.Context) ([]PermissionGroup, error) {
	rows, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return groupPermissions(rows), nil
}

func (s *Service) ListPermissionItems(ctx context.Context) ([]PermissionItemDTO, error) {
	rows, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]PermissionItemDTO, 0, len(rows))
	for _, p := range rows {
		items = append(items, PermissionItemDTO{
			ID: p.ID.String(), Resource: p.Resource, Action: p.Action, Description: p.Description,
		})
	}
	return items, nil
}

func (s *Service) MyPermissions(ctx context.Context, tenantID, userID uuid.UUID, activeRoleID *uuid.UUID) ([]PermissionGroup, error) {
	var rows []domain.Permission
	var err error
	if activeRoleID != nil {
		ok, err := s.repo.UserHasRole(ctx, tenantID, userID, *activeRoleID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrRoleForbidden
		}
		rows, err = s.repo.ListRolePermissions(ctx, *activeRoleID)
	} else {
		rows, err = s.repo.ListUserPermissions(ctx, tenantID, userID)
	}
	if err != nil {
		return nil, err
	}
	return groupPermissions(rows), nil
}

func (s *Service) ListMembers(ctx context.Context, tenantID, userID uuid.UUID) ([]MemberDTO, error) {
	rows, err := s.repo.ListTenantMembers(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	rows = filterMembersByScope(rows, s.scope.Params(ctx, tenantID, userID))
	byUser := make(map[uuid.UUID]*MemberDTO)
	order := make([]uuid.UUID, 0)
	for _, row := range rows {
		m, ok := byUser[row.UserID]
		if !ok {
			m = &MemberDTO{
				ID:         row.UserID.String(),
				Email:      row.Email,
				Name:       row.Name,
				AvatarURL:  row.AvatarURL,
				Department: row.Department,
				Roles:      []MemberRoleDTO{},
				JoinedAt:   row.CreatedAt.UTC().Format(time.RFC3339),
			}
			byUser[row.UserID] = m
			order = append(order, row.UserID)
		}
		if row.RoleID != nil {
			m.Roles = append(m.Roles, MemberRoleDTO{
				ID: row.RoleID.String(), Name: row.RoleName, IsSystem: row.IsSystem,
			})
		}
	}
	out := make([]MemberDTO, 0, len(order))
	for _, id := range order {
		out = append(out, *byUser[id])
	}
	return out, nil
}

func filterMembersByScope(rows []repository.TenantMemberRow, scope datascope.ScopeParams) []repository.TenantMemberRow {
	switch scope.Level {
	case datascope.LevelAll:
		return rows
	case datascope.LevelDepartment:
		if scope.Department == "" {
			return filterMembersByUser(rows, scope.UserID)
		}
		out := make([]repository.TenantMemberRow, 0, len(rows))
		for _, row := range rows {
			if row.Department == scope.Department {
				out = append(out, row)
			}
		}
		return out
	default:
		return filterMembersByUser(rows, scope.UserID)
	}
}

func filterMembersByUser(rows []repository.TenantMemberRow, userID uuid.UUID) []repository.TenantMemberRow {
	out := make([]repository.TenantMemberRow, 0, 1)
	for _, row := range rows {
		if row.UserID == userID {
			out = append(out, row)
		}
	}
	return out
}

func (s *Service) ListRoles(ctx context.Context, tenantID uuid.UUID) ([]RoleDTO, error) {
	rows, err := s.repo.ListRoles(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return mapRoles(rows), nil
}

func (s *Service) CreateRole(ctx context.Context, tenantID uuid.UUID, name, description string) (*RoleDTO, error) {
	role := &domain.Role{
		TenantID:    tenantID,
		Name:        name,
		Description: description,
	}
	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	row, err := s.repo.FindRole(ctx, tenantID, role.ID)
	if err != nil {
		return nil, err
	}
	dto := mapRole(*row)
	return &dto, nil
}

func (s *Service) UpdateRole(ctx context.Context, tenantID, roleID uuid.UUID, name, description string) (*RoleDTO, error) {
	if err := s.repo.UpdateRole(ctx, tenantID, roleID, name, description); err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	row, err := s.repo.FindRole(ctx, tenantID, roleID)
	if err != nil {
		return nil, err
	}
	dto := mapRole(*row)
	return &dto, nil
}

func (s *Service) SetRolePermissions(ctx context.Context, tenantID, roleID uuid.UUID, permissionIDs []uuid.UUID) (*RoleDTO, error) {
	if _, err := s.repo.FindRole(ctx, tenantID, roleID); err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	ok, err := s.repo.PermissionsExist(ctx, permissionIDs)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidPermissions
	}
	if err := s.repo.SetRolePermissions(ctx, roleID, permissionIDs); err != nil {
		return nil, err
	}
	if err := persistence.SyncCasbinPolicies(s.db, s.enforcer); err != nil {
		return nil, err
	}
	row, err := s.repo.FindRole(ctx, tenantID, roleID)
	if err != nil {
		return nil, err
	}
	dto := mapRole(*row)
	return &dto, nil
}

func (s *Service) ListUserRoles(ctx context.Context, tenantID, userID uuid.UUID) ([]RoleDTO, error) {
	roles, err := s.repo.ListUserRoles(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	out := make([]RoleDTO, 0, len(roles))
	for _, role := range roles {
		row, err := s.repo.FindRole(ctx, tenantID, role.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, mapRole(*row))
	}
	return out, nil
}

func (s *Service) SetUserRoles(ctx context.Context, tenantID, userID uuid.UUID, roleIDs []uuid.UUID) ([]RoleDTO, error) {
	ok, err := s.repo.RolesBelongToTenant(ctx, tenantID, roleIDs)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidRoles
	}
	if err := s.repo.SetUserRoles(ctx, tenantID, userID, roleIDs); err != nil {
		if errors.Is(err, repository.ErrMemberNotInTenant) {
			return nil, ErrMemberNotInTenant
		}
		return nil, err
	}
	if err := persistence.SyncCasbinPolicies(s.db, s.enforcer); err != nil {
		return nil, err
	}
	return s.ListUserRoles(ctx, tenantID, userID)
}

func (s *Service) Check(ctx context.Context, tenantID uuid.UUID, userID string, isSuperAdmin bool, activeRoleID, resource, action string) (bool, error) {
	if isSuperAdmin {
		return true, nil
	}
	if activeRoleID != "" {
		ctx = activerole.WithID(ctx, activeRoleID)
	}
	return rbacutil.Enforce(ctx, s.enforcer, userID, tenantID.String(), resource, action)
}

func groupPermissions(rows []domain.Permission) []PermissionGroup {
	grouped := make(map[string][]string)
	for _, p := range rows {
		grouped[p.Resource] = append(grouped[p.Resource], p.Action)
	}
	result := make([]PermissionGroup, 0, len(grouped))
	for resource, actions := range grouped {
		result = append(result, PermissionGroup{Resource: resource, Actions: actions})
	}
	return result
}

func mapRoles(rows []repository.RoleWithPermissions) []RoleDTO {
	out := make([]RoleDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapRole(row))
	}
	return out
}

func mapRole(row repository.RoleWithPermissions) RoleDTO {
	permIDs := make([]string, 0, len(row.PermissionIDs))
	for _, id := range row.PermissionIDs {
		permIDs = append(permIDs, id.String())
	}
	return RoleDTO{
		ID:            row.Role.ID.String(),
		Name:          row.Role.Name,
		Description:   row.Role.Description,
		IsSystem:      row.Role.IsSystem,
		PermissionIDs: permIDs,
		UserCount:     row.UserCount,
	}
}
