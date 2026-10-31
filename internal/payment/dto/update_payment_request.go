package dto

type UpdatePaymentRequest struct {
	ProviderDefinitionID int64  `json:"provider_definition_id" binding:"required"`
	Reference            string `json:"reference" binding:"omitempty"`
	Status               string `json:"status" binding:"omitempty"`
}