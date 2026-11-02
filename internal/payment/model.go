package payment

import (
	"encoding/json"
	"time"
)

type Payment struct {
	ID        int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Reference string          `gorm:"type:varchar(255);not null;uniqueIndex:idx_payments_reference" json:"reference"`
	Status    string          `gorm:"type:varchar(100);not null;index" json:"status"`
	Fields    json.RawMessage `gorm:"type:jsonb" json:"-"`
	CreatedAt time.Time       `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (Payment) TableName() string {
	return "payments"
}