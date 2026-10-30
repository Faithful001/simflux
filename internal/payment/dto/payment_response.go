package dto

import (
	"encoding/json"
	"time"
)

type PaymentResponse struct {
    ID        int64           `json:"id"`
    Provider  string          `json:"provider"`
    Reference string          `json:"reference"`
    Status    string          `json:"status"`
    Data      json.RawMessage `json:"-"`
    CreatedAt time.Time       `json:"created_at"`
    UpdatedAt time.Time       `json:"updated_at"`
}

func (p PaymentResponse) MarshalJSON() ([]byte, error) {
    var data map[string]any

    if err := json.Unmarshal(p.Data, &data); err != nil {
        return nil, err
    }

    response := map[string]any{
        "id":         p.ID,
        "provider":   p.Provider,
        "reference":  p.Reference,
        "status":     p.Status,
        "created_at": p.CreatedAt,
        "updated_at": p.UpdatedAt,
    }

    for key, value := range data {
        response[key] = value
    }

    return json.Marshal(response)
}

