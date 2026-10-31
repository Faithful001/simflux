package payment

import (
	"encoding/json"
	"time"

	provider "github.com/Faithful001/simflux.git/internal/provider_definition"
)

type Payment struct {
    ID                      int64          	            `gorm:"primaryKey;autoIncrement" json:"id"`
    Reference               string                      `gorm:"not null;uniqueIndex:idx_provider_reference" json:"reference"`
    Status                  string                      `gorm:"not null;index" json:"status"`

    Attributes              json.RawMessage             `gorm:"type:jsonb" json:"-"`

    ProviderDefinitionID    int64                       `gorm:"not null;uniqueIndex:idx_provider_reference" json:"provider_definition_id"`
	ProviderDefinition		provider.ProviderDefinition `gorm:"foreignKey:ProviderDefinitionID;references:ID" json:"payment"`  

    CreatedAt               time.Time                   `gorm:"not null;autoCreateTime" json:"created_at"`
    UpdatedAt               time.Time 		            `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (Payment) TableName() string {
    return "payments"
}