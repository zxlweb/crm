package main

import (
	"crm-backend/internal/application/auth"
	rbacapp "crm-backend/internal/application/rbac"
	"crm-backend/internal/application/superadmin"
	"crm-backend/internal/config"
	"crm-backend/internal/infrastructure/persistence"
	httphandler "crm-backend/internal/interfaces/http"
	"crm-backend/internal/interfaces/middleware"
	"crm-backend/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db := persistence.InitDB(cfg)
	enforcer := persistence.InitCasbin(db)

	userRepo := repository.NewUserRepository(db)
	tenantRepo := repository.NewTenantRepository(db)

	authSvc := auth.NewService(
		userRepo,
		cfg.JWTSecret,
		cfg.AccessTokenTTL(),
		cfg.RefreshTokenTTL(),
	)
	authHTTP := httphandler.NewAuthHandlers(authSvc)

	superAdminSvc := superadmin.NewService(tenantRepo)
	superAdminHTTP := httphandler.NewSuperAdminHandlers(superAdminSvc)

	rbacRepo := repository.NewRBACRepository(db)
	rbacSvc := rbacapp.NewService(rbacRepo, db, enforcer)
	rbacHTTP := httphandler.NewRBACHandlers(rbacSvc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/health", httphandler.HealthHandler(db))

	public := r.Group("/api")
	{
		authGroup := public.Group("/auth")
		authGroup.POST("/login", authHTTP.Login)
		authGroup.POST("/refresh", authHTTP.Refresh)
	}

	authOnly := r.Group("/api")
	authOnly.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		authOnly.GET("/auth/profile", authHTTP.Profile)
		authOnly.GET("/auth/tenants", authHTTP.ListTenants)
		authOnly.POST("/auth/switch-tenant", authHTTP.SwitchTenant)
	}

	superAdmin := r.Group("/api/super-admin")
	superAdmin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	superAdmin.Use(middleware.SuperAdminMiddleware())
	{
		superAdmin.GET("/overview", superAdminHTTP.Overview)
		superAdmin.GET("/tenants", superAdminHTTP.ListTenants)
		superAdmin.GET("/tenants/:id", superAdminHTTP.GetTenant)
		superAdmin.PATCH("/tenants/:id", superAdminHTTP.PatchTenant)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	protected.Use(middleware.TenantMiddleware())
	protected.Use(middleware.RBACMiddleware(enforcer))
	{
		protected.GET("/auth/me", httphandler.CurrentUserHandler())
		protected.GET("/rbac/permissions", rbacHTTP.ListPermissions)
		protected.GET("/rbac/permission-items", rbacHTTP.ListPermissionItems)
		protected.GET("/rbac/my-permissions", rbacHTTP.MyPermissions)
		protected.GET("/rbac/roles", rbacHTTP.ListRoles)
		protected.POST("/rbac/roles", rbacHTTP.CreateRole)
		protected.PUT("/rbac/roles/:id", rbacHTTP.UpdateRole)
		protected.POST("/rbac/roles/:id/permissions", rbacHTTP.AssignRolePermissions)
		protected.GET("/rbac/users/:id/roles", rbacHTTP.ListUserRoles)
		protected.POST("/rbac/users/:id/roles", rbacHTTP.AssignUserRoles)
		protected.POST("/rbac/check", rbacHTTP.Check)
	}

	log.Printf("🚀 CRM Backend started on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
