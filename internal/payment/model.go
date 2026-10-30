package payment

import (
	"encoding/json"
	"time"
)

type Payment struct {
    ID        int64          	`gorm:"primaryKey;autoIncrement" json:"id"`
    Provider  string         	`gorm:"not null;index" json:"provider"`
    Reference string         	`gorm:"not null;uniqueIndex" json:"reference"`
    Status    string         	`gorm:"not null;index" json:"status"`

    Data      json.RawMessage 	`gorm:"type:jsonb;not null" json:"data"`

    CreatedAt time.Time 		`gorm:"not null;autoCreateTime" json:"created_at"`
    UpdatedAt time.Time 		`gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (Payment) TableName() string {
    return "payments"
}