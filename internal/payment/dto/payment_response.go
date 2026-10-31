package dto

import (
	"encoding/json"
	"time"
)

type PaymentResponse struct {
	ID                   int64          `json:"id"`
	ProviderDefinitionID int64          `json:"provider_definition_id"`
	Reference            string         `json:"reference"`
	Status               string         `json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	Fields               map[string]any `json:"-"`
}

func (p PaymentResponse) MarshalJSON() ([]byte, error) {
	response := make(map[string]any)

	for key, value := range p.Fields {
		response[key] = value
	}

	response["id"] = p.ID
	response["provider_definition_id"] = p.ProviderDefinitionID
	response["reference"] = p.Reference
	response["status"] = p.Status
	response["created_at"] = p.CreatedAt
	response["updated_at"] = p.UpdatedAt

	return json.Marshal(response)
}