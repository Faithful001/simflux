package dto

import "encoding/json"

type UpdatePaymentRequest struct {
	Provider  string `json:"provider" binding:"required"`
	Reference string `json:"reference" binding:"required"`
	Status    string `json:"status" binding:"required"`

	Data json.RawMessage `json:"data" binding:"omitempty"`
}