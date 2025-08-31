package repository

import (
	"context"
	"pos-api/app/model"

	"gorm.io/gorm"
)

type TableRepository interface {
	Create(ctx context.Context, t *model.Table) error
	GetByID(ctx context.Context, id uint64) (*model.Table, error)
	List(ctx context.Context) ([]model.Table, error)
	UpdateStatus(ctx context.Context, id uint64, status string) error
}

type tableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) TableRepository {
	return &tableRepository{db: db}
}

func (r *tableRepository) Create(ctx context.Context, t *model.Table) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *tableRepository) GetByID(ctx context.Context, id uint64) (*model.Table, error) {
	var tb model.Table
	if err := r.db.WithContext(ctx).First(&tb, id).Error; err != nil {
		return nil, err
	}
	return &tb, nil
}

func (r *tableRepository) List(ctx context.Context) ([]model.Table, error) {
	var tbs []model.Table
	if err := r.db.WithContext(ctx).Order("id asc").Find(&tbs).Error; err != nil {
		return nil, err
	}
	return tbs, nil
}

func (r *tableRepository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	return r.db.WithContext(ctx).Model(&model.Table{}).
		Where("id = ?", id).
		Update("status", status).Error
}
