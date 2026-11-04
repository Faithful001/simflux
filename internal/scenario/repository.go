package scenario

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type ScenarioRepository struct {
	db *gorm.DB
}

type IScenarioRepository interface {
	CreateScenario(ctx context.Context, scenario *Scenario) error
	GetAllScenarios(ctx context.Context) ([]Scenario, error)
	GetScenarioByID(ctx context.Context, id int64) (Scenario, error)
	UpdateScenario(ctx context.Context, scenario *Scenario) error
	DeleteScenario(ctx context.Context, id int64) error
}

func NewScenarioRepository(db *gorm.DB) (IScenarioRepository, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &ScenarioRepository{db: db}, nil
}

func (r *ScenarioRepository) CreateScenario(ctx context.Context, scenario *Scenario) error {
	return r.db.WithContext(ctx).Create(scenario).Error
}

func (r *ScenarioRepository) GetAllScenarios(ctx context.Context) ([]Scenario, error) {
	var scenarios []Scenario
	if err := r.db.WithContext(ctx).Find(&scenarios).Error; err != nil {
		return nil, err
	}
	return scenarios, nil
}

func (r *ScenarioRepository) GetScenarioByID(ctx context.Context, id int64) (Scenario, error) {
	var scenario Scenario
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&scenario).Error; err != nil {
		return Scenario{}, err
	}
	return scenario, nil
}

func (r *ScenarioRepository) UpdateScenario(ctx context.Context, scenario *Scenario) error {
	result := r.db.WithContext(ctx).
		Model(&Scenario{}).
		Where("id = ?", scenario.ID).
		Updates(map[string]any{
			"name":   scenario.Name,
			"type":   scenario.Type,
			"status": scenario.Status,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ScenarioRepository) DeleteScenario(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&Scenario{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
