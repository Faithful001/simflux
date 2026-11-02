package provider_definition

import "context"

type IProviderDefinitionService interface{
	Create(ctx context.Context, definition *ProviderDefinition) error
}

type ProviderDefinitionService struct {
	repo IProviderDefinitionRepository
}

func NewProviderDefinitionService(repo IProviderDefinitionRepository) IProviderDefinitionService {
	return &ProviderDefinitionService{repo: repo}
}

func (p *ProviderDefinitionService) Create(ctx context.Context, definition *ProviderDefinition) error {
	return p.repo.Create(ctx, definition)
}
