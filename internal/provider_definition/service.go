package provider_definition

import "context"

type IProviderDefinitionService interface {
	SaveActive(ctx context.Context, definition *ProviderDefinition) error
	GetActive(ctx context.Context) (ProviderDefinition, error)
	Create(ctx context.Context, definition *ProviderDefinition) error
}

type ProviderDefinitionService struct {
	repo IProviderDefinitionRepository
}

func NewProviderDefinitionService(repo IProviderDefinitionRepository) IProviderDefinitionService {
	return &ProviderDefinitionService{repo: repo}
}

func (p *ProviderDefinitionService) SaveActive(ctx context.Context, definition *ProviderDefinition) error {
	return p.repo.SaveActive(ctx, definition)
}

func (p *ProviderDefinitionService) GetActive(ctx context.Context) (ProviderDefinition, error) {
	return p.repo.GetActive(ctx)
}

func (p *ProviderDefinitionService) Create(ctx context.Context, definition *ProviderDefinition) error {
	return p.repo.Create(ctx, definition)
}
