package rbac

import (
	"context"
	"errors"

	"crm-backend/internal/domain"
	"crm-backend/internal/infrastructure/persistence"
	"crm-backend/internal/repository"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrInvalidPermissions = errors.New("invalid permission ids")
	ErrInvalidRoles       = errors.New("invalid role ids")
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
}

func NewService(repo repository.RBACRepository, db *gorm.DB, enforcer *casbin.Enforcer) *Service {
	return &Service{repo: repo, db: db, enforcer: enforcer}
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

func (s *Service) MyPermissions(ctx context.Context, tenantID, userID uuid.UUID) ([]PermissionGroup, error) {
	rows, err := s.repo.ListUserPermissions(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	return groupPermissions(rows), nil
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
		return nil, err
	}
	if err := persistence.SyncCasbinPolicies(s.db, s.enforcer); err != nil {
		return nil, err
	}
	return s.ListUserRoles(ctx, tenantID, userID)
}

func (s *Service) Check(ctx context.Context, tenantID uuid.UUID, userID string, isSuperAdmin bool, resource, action string) (bool, error) {
	if isSuperAdmin {
		return true, nil
	}
	return s.enforcer.Enforce(userID, tenantID.String(), resource, action)
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
