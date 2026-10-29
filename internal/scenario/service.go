package scenario

import "context"

type ScenarioService interface {
    CreateScenario(ctx context.Context, scenario Scenario) error
    GetAllScenarios(ctx context.Context) ([]Scenario, error)
    GetScenarioByID(ctx context.Context, id int) (Scenario, error)
    UpdateScenario(ctx context.Context, scenario Scenario) error
    DeleteScenario(ctx context.Context, id int) error
}

type scenarioService struct {
    repo IScenarioRepository
}

func NewScenarioService(repo IScenarioRepository) ScenarioService {
    return &scenarioService{repo: repo}
}

func (s *scenarioService) CreateScenario(ctx context.Context, scenario Scenario) error {
    return s.repo.CreateScenario(ctx, scenario)
}

func (s *scenarioService) GetAllScenarios(ctx context.Context) ([]Scenario, error) {
    return s.repo.GetAllScenarios(ctx)
}

func (s *scenarioService) GetScenarioByID(ctx context.Context, id int) (Scenario, error) {
    return s.repo.GetScenarioByID(ctx, id)
}

func (s *scenarioService) UpdateScenario(ctx context.Context, scenario Scenario) error {
    return s.repo.UpdateScenario(ctx, scenario)
}

func (s *scenarioService) DeleteScenario(ctx context.Context, id int) error {
    return s.repo.DeleteScenario(ctx, id)
}