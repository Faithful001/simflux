package scenario

import "context"

type IScenarioService interface {
    Create(ctx context.Context, scenario Scenario) error
    Get(ctx context.Context) ([]Scenario, error)
    GetByID(ctx context.Context, id int) (Scenario, error)
    Update(ctx context.Context, scenario Scenario) error
    Delete(ctx context.Context, id int) error
}

type ScenarioService struct {
    repo IScenarioRepository
}

func NewScenarioService(repo IScenarioRepository) IScenarioService {
    return &ScenarioService{repo: repo}
}

func (s *ScenarioService) Create(ctx context.Context, scenario Scenario) error {
    return s.repo.CreateScenario(ctx, scenario)
}

func (s *ScenarioService) Get(ctx context.Context) ([]Scenario, error) {
    return s.repo.GetAllScenarios(ctx)
}

func (s *ScenarioService) GetByID(ctx context.Context, id int) (Scenario, error) {
    return s.repo.GetScenarioByID(ctx, id)
}

func (s *ScenarioService) Update(ctx context.Context, scenario Scenario) error {
    return s.repo.UpdateScenario(ctx, scenario)
}

func (s *ScenarioService) Delete(ctx context.Context, id int) error {
    return s.repo.DeleteScenario(ctx, id)
}