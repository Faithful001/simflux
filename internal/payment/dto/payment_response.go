package dto

import "time"

type PaymentResponse struct {
	ID                   int64     `json:"id"`
	ProviderDefinitionID int64     `json:"provider_definition_id"`
	Reference            string    `json:"reference"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}