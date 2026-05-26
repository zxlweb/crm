package rbac

import (
	"context"
	"testing"
	"time"

	"crm-backend/internal/domain"
	"crm-backend/internal/repository"

	"github.com/google/uuid"
)

type mockRBACRepo struct {
	perms []domain.Permission
	roles []repository.RoleWithPermissions
}

func (m *mockRBACRepo) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	return m.perms, nil
}

func (m *mockRBACRepo) ListRoles(ctx context.Context, tenantID uuid.UUID) ([]repository.RoleWithPermissions, error) {
	return m.roles, nil
}

func (m *mockRBACRepo) FindRole(ctx context.Context, tenantID, roleID uuid.UUID) (*repository.RoleWithPermissions, error) {
	for i := range m.roles {
		if m.roles[i].Role.ID == roleID {
			return &m.roles[i], nil
		}
	}
	return nil, repository.ErrRoleNotFound
}

func (m *mockRBACRepo) CreateRole(ctx context.Context, role *domain.Role) error {
	m.roles = append(m.roles, repository.RoleWithPermissions{Role: *role})
	return nil
}

func (m *mockRBACRepo) UpdateRole(ctx context.Context, tenantID, roleID uuid.UUID, name, description string) error {
	return nil
}

func (m *mockRBACRepo) SetRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return nil
}

func (m *mockRBACRepo) ListUserRoles(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Role, error) {
	return nil, nil
}

func (m *mockRBACRepo) SetUserRoles(ctx context.Context, tenantID, userID uuid.UUID, roleIDs []uuid.UUID) error {
	return nil
}

func (m *mockRBACRepo) ListUserPermissions(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Permission, error) {
	return m.perms, nil
}

func (m *mockRBACRepo) ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]domain.Permission, error) {
	return m.perms, nil
}

func (m *mockRBACRepo) UserHasRole(ctx context.Context, tenantID, userID, roleID uuid.UUID) (bool, error) {
	return true, nil
}

func (m *mockRBACRepo) PermissionsExist(ctx context.Context, ids []uuid.UUID) (bool, error) {
	return true, nil
}

func (m *mockRBACRepo) RolesBelongToTenant(ctx context.Context, tenantID uuid.UUID, roleIDs []uuid.UUID) (bool, error) {
	return true, nil
}

func TestService_ListPermissionDictionary(t *testing.T) {
	repo := &mockRBACRepo{
		perms: []domain.Permission{
			{Resource: "leads", Action: "view"},
			{Resource: "leads", Action: "create"},
		},
	}
	svc := NewService(repo, nil, nil)

	groups, err := svc.ListPermissionDictionary(context.Background())
	if err != nil || len(groups) != 1 || len(groups[0].Actions) != 2 {
		t.Fatalf("groups: %+v err=%v", groups, err)
	}
}

func TestService_ListRoles(t *testing.T) {
	tid := uuid.New()
	rid := uuid.New()
	repo := &mockRBACRepo{
		roles: []repository.RoleWithPermissions{{
			Role: domain.Role{
				ID: rid, TenantID: tid, Name: "Admin", IsSystem: true, CreatedAt: time.Now(),
			},
			PermissionIDs: []uuid.UUID{uuid.New()},
			UserCount:     1,
		}},
	}
	svc := NewService(repo, nil, nil)

	roles, err := svc.ListRoles(context.Background(), tid)
	if err != nil || len(roles) != 1 || roles[0].Name != "Admin" {
		t.Fatalf("roles: %+v err=%v", roles, err)
	}
}
