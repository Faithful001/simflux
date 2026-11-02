package dto

type CreatePaymentRequest struct {
	Reference string         `json:"reference" binding:"required"`
	Status    string         `json:"status" binding:"required"`
	Fields    map[string]any `json:"fields"`
}
