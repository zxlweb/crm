package main

import (
	"crm-backend/internal/application/account"
	"crm-backend/internal/application/audit"
	"crm-backend/internal/application/auth"
	actapp "crm-backend/internal/application/activity"
	contactapp "crm-backend/internal/application/contact"
	"crm-backend/internal/application/emotion"
	dashapp "crm-backend/internal/application/dashboard"
	dealapp "crm-backend/internal/application/deal"
	leadapp "crm-backend/internal/application/lead"
	segmentapp "crm-backend/internal/application/segment"
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
	auditRepo := repository.NewAuditRepository(db)
	auditRec := audit.NewRecorder(auditRepo)

	authSvc := auth.NewService(
		userRepo,
		cfg.JWTSecret,
		cfg.AccessTokenTTL(),
		cfg.RefreshTokenTTL(),
		&auth.ServiceDeps{DB: db, Enforcer: enforcer, Audit: auditRec},
	)
	authHTTP := httphandler.NewAuthHandlers(authSvc)

	superAdminSvc := superadmin.NewService(tenantRepo)
	superAdminHTTP := httphandler.NewSuperAdminHandlers(superAdminSvc, auditRec)

	rbacRepo := repository.NewRBACRepository(db)
	rbacSvc := rbacapp.NewService(rbacRepo, db, enforcer)
	rbacHTTP := httphandler.NewRBACHandlers(rbacSvc, auditRec)

	accountRepo := repository.NewAccountRepository(db)
	accountSvc := account.NewService(accountRepo, tenantRepo, enforcer)
	leadRepo := repository.NewLeadRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	emotionSvc := emotion.NewService(activityRepo)

	accountHTTP := httphandler.NewAccountHandlers(accountSvc, auditRec, emotionSvc)

	contactRepo := repository.NewContactRepository(db)
	contactSvc := contactapp.NewService(contactRepo, accountRepo, enforcer)
	contactHTTP := httphandler.NewContactHandlers(contactSvc, auditRec, emotionSvc)

	dealRepo := repository.NewDealRepository(db)
	dealSvc := dealapp.NewService(dealRepo, accountRepo, enforcer)
	dealHTTP := httphandler.NewDealHandlers(dealSvc, auditRec)

	dashboardSvc := dashapp.NewService(leadRepo, accountRepo, dealRepo, activityRepo, tenantRepo, userRepo, enforcer)
	dashboardHTTP := httphandler.NewDashboardHandlers(dashboardSvc, enforcer)

	leadSvc := leadapp.NewService(leadRepo, accountRepo, activityRepo, tenantRepo, enforcer, dealSvc)
	leadHTTP := httphandler.NewLeadHandlers(leadSvc, auditRec, emotionSvc)

	activitySvc := actapp.NewService(activityRepo, leadRepo, accountRepo, contactRepo, tenantRepo, enforcer)
	activityHTTP := httphandler.NewActivityHandlers(activitySvc, auditRec)

	segmentRepo := repository.NewSegmentRepository(db)
	segmentSvc := segmentapp.NewService(segmentRepo, leadRepo, accountRepo, tenantRepo, enforcer)
	segmentHTTP := httphandler.NewSegmentHandlers(segmentSvc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/health", httphandler.HealthHandler(db))

	public := r.Group("/api")
	{
		authGroup := public.Group("/auth")
		authGroup.POST("/login", authHTTP.Login)
		authGroup.POST("/register", authHTTP.Register)
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
		superAdmin.GET("/stats/tenant-activity", superAdminHTTP.TenantActivityTrend)
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

		protected.GET("/accounts", accountHTTP.List)
		protected.POST("/accounts", accountHTTP.Create)
		protected.GET("/accounts/:id", accountHTTP.Get)
		protected.PUT("/accounts/:id", accountHTTP.Put)
		protected.PATCH("/accounts/:id", accountHTTP.Patch)
		protected.DELETE("/accounts/:id", accountHTTP.Delete)
		protected.GET("/accounts/:id/contacts", contactHTTP.ListByAccount)
		protected.GET("/accounts/:id/emotion-journey", accountHTTP.EmotionJourney)
		protected.POST("/accounts/:id/insights/evaluate", accountHTTP.EvaluateInsights)

		protected.GET("/contacts", contactHTTP.List)
		protected.POST("/contacts", contactHTTP.Create)
		protected.GET("/contacts/:id", contactHTTP.Get)
		protected.PUT("/contacts/:id", contactHTTP.Put)
		protected.PATCH("/contacts/:id", contactHTTP.Patch)
		protected.DELETE("/contacts/:id", contactHTTP.Delete)
		protected.GET("/contacts/:id/emotion-journey", contactHTTP.EmotionJourney)
		protected.POST("/contacts/:id/insights/evaluate", contactHTTP.EvaluateInsights)

		protected.GET("/leads/stats/by-source", leadHTTP.StatsBySource)
		protected.GET("/leads/stats/by-status", leadHTTP.StatsByStatus)
		protected.GET("/leads/stats/trend", leadHTTP.StatsTrend)
		protected.GET("/leads/stats/funnel", leadHTTP.StatsFunnel)

		protected.GET("/leads", leadHTTP.List)
		protected.POST("/leads", leadHTTP.Create)
		protected.GET("/leads/:id", leadHTTP.Get)
		protected.PUT("/leads/:id", leadHTTP.Put)
		protected.PATCH("/leads/:id", leadHTTP.Patch)
		protected.DELETE("/leads/:id", leadHTTP.Delete)
		protected.GET("/dashboard/summary", dashboardHTTP.Summary)
		protected.GET("/dashboard/funnel", dashboardHTTP.Funnel)
		protected.GET("/dashboard/quota", dashboardHTTP.Quota)
		protected.GET("/dashboard/team-ranking", dashboardHTTP.TeamRanking)
		protected.GET("/dashboard/todo", dashboardHTTP.Todo)

		protected.GET("/deals/stats/by-stage", dealHTTP.StatsByStage)
		protected.GET("/deals/stats/win-rate", dealHTTP.StatsWinRate)
		protected.GET("/deals/pipeline", dealHTTP.Pipeline)
		protected.GET("/deals", dealHTTP.List)
		protected.POST("/deals", dealHTTP.Create)
		protected.GET("/deals/:id", dealHTTP.Get)
		protected.PUT("/deals/:id", dealHTTP.Put)
		protected.PATCH("/deals/:id", dealHTTP.Patch)
		protected.DELETE("/deals/:id", dealHTTP.Delete)
		protected.PUT("/deals/:id/stage", dealHTTP.PutStage)

		protected.POST("/leads/:id/convert", leadHTTP.Convert)
		protected.GET("/leads/:id/emotion-journey", leadHTTP.EmotionJourney)
		protected.POST("/leads/:id/insights/evaluate", leadHTTP.EvaluateInsights)

		protected.GET("/segments", segmentHTTP.List)
		protected.GET("/segments/:code/count", segmentHTTP.Count)

		protected.GET("/activities/summary", activityHTTP.Summary)
		protected.GET("/activities", activityHTTP.List)
		protected.POST("/activities", activityHTTP.Create)
		protected.GET("/activities/:id", activityHTTP.Get)
		protected.PATCH("/activities/:id", activityHTTP.Patch)
		protected.DELETE("/activities/:id", activityHTTP.Delete)
	}

	log.Printf("🚀 CRM Backend started on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
