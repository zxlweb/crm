package persistence

import (
	"log"
	"time"

	"crm-backend/internal/config"
	"crm-backend/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if cfg.AutoMigrate {
		if err := autoMigrate(db); err != nil {
			log.Fatalf("GORM AutoMigrate 失败: %v", err)
		}
		log.Println("✅ GORM AutoMigrate 完成（开发模式）")
	}

	log.Println("✅ 数据库连接成功")
	return db
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Tenant{},
		&domain.User{},
		&domain.UserTenant{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RolePermission{},
		&domain.UserRole{},
		&domain.AuditLog{},
		&domain.Lead{},
	)
}

func PingDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
