package repository

import (
	"context"

	"crm-backend/internal/domain"

	"gorm.io/gorm"
)

type AuditRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
}

type GormAuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &GormAuditRepository{db: db}
}

func (r *GormAuditRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
