package provider_definition

import (
	"context"

	"gorm.io/gorm"
)

type IProviderDefinitionRepository interface {
	Create(ctx context.Context, definition *ProviderDefinition) error
	Get(ctx context.Context) ([]ProviderDefinition, error)
	GetByID(ctx context.Context, id int64) (ProviderDefinition, error)
	GetByName(ctx context.Context, name string) (ProviderDefinition, error)
	Update(ctx context.Context, definition *ProviderDefinition) error
	Delete(ctx context.Context, id int64) error
}

type ProviderDefinitionRepository struct {
	db *gorm.DB
}

func NewProviderDefinitionRepository(db *gorm.DB) IProviderDefinitionRepository {
	return &ProviderDefinitionRepository{
		db: db,
	}
}

func (r *ProviderDefinitionRepository) Create(ctx context.Context, definition *ProviderDefinition) error {
	return r.db.WithContext(ctx).Create(definition).Error
}

func (r *ProviderDefinitionRepository) Get(ctx context.Context) ([]ProviderDefinition, error) {
	var definitions []ProviderDefinition

	err := r.db.WithContext(ctx).Find(&definitions).Error
	return definitions, err
}

func (r *ProviderDefinitionRepository) GetByID(ctx context.Context, id int64) (ProviderDefinition, error) {
	var definition ProviderDefinition

	err := r.db.WithContext(ctx).First(&definition, id).Error
	return definition, err
}

func (r *ProviderDefinitionRepository) GetByName(ctx context.Context, name string) (ProviderDefinition, error) {
	var definition ProviderDefinition

	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&definition).Error

	return definition, err
}

func (r *ProviderDefinitionRepository) Update(ctx context.Context, definition *ProviderDefinition) error {
	result := r.db.WithContext(ctx).
		Model(&ProviderDefinition{}).
		Where("id = ?", definition.ID).
		Updates(map[string]any{
			"name":          definition.Name,
			"provider_name": definition.ProviderName,
			"version":       definition.Version,
			"data":          definition.Data,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *ProviderDefinitionRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).
		Delete(&ProviderDefinition{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
