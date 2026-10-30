package dto

import "encoding/json"

type CreatePayment struct {
	Provider  string `json:"provider" binding:"required"`
	Reference string `json:"reference" binding:"required"`
	Status    string `json:"status" binding:"required"`

	Data json.RawMessage `json:"data" binding:"omitempty"`
}