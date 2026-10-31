package dto

type CreatePaymentRequest struct {
	ProviderDefinitionID int64          `json:"provider_definition_id" binding:"required"`
	Reference            string         `json:"reference" binding:"required"`
	Status               string         `json:"status" binding:"required"`
	Fields               map[string]any `json:"fields"`
}