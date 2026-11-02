package dto

type UpdatePaymentRequest struct {
	Reference string `json:"reference" binding:"required"`
	Status    string `json:"status" binding:"required"`
}
