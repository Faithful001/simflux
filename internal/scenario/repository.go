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
	CreateScenario(ctx context.Context, scenario Scenario) error
	GetAllScenarios(ctx context.Context) ([]Scenario, error)
	GetScenarioByID(ctx context.Context, id int) (Scenario, error)
	UpdateScenario(ctx context.Context, scenario Scenario) error
	DeleteScenario(ctx context.Context, id int) error
}

func NewScenarioRepository(db *gorm.DB) (IScenarioRepository, error){
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &ScenarioRepository{db:db}, nil
}

func (r *ScenarioRepository) CreateScenario(ctx context.Context, scenario Scenario) error {
	return r.db.WithContext(ctx).Create(scenario).Error
}

func (r *ScenarioRepository) GetAllScenarios(ctx context.Context) ([]Scenario, error) {
	var scenarios []Scenario
	if err := r.db.WithContext(ctx).Find(&scenarios).Error; err != nil {
		return nil, err
	}

	return scenarios, nil
}
func (r *ScenarioRepository) GetScenarioByID(ctx context.Context, id int) (Scenario, error) {
	var scenario Scenario
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&scenario).Error; err != nil {
		return Scenario{}, err
	}

	return scenario, nil
}
func (r *ScenarioRepository) UpdateScenario(ctx context.Context, scenario Scenario) error {
	if err := r.db.WithContext(ctx).Find(scenario).Error; err != nil {
		return err
	}

	return nil
}
func (r *ScenarioRepository) DeleteScenario(ctx context.Context, id int) error {
	var scenarios []*Scenario
	if err := r.db.WithContext(ctx).Where("id = ?", id).Find(&scenarios).Error; err != nil {
		return err
	}

	return nil
}

