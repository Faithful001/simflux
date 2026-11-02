package webhook

import (
	"encoding/json"
	"time"

	"github.com/Faithful001/simflux.git/internal/event"
)

type Webhook struct {
	ID        int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	SigHeader string          `gorm:"type:varchar(255);not null" json:"sig_header"` // webhook signature header
	Fields    json.RawMessage `gorm:"type:jsonb" json:"fields"`                       // what the webhook should send in the request body (payload)
	EventID   int64           `gorm:"not null;index:idx_webhooks_event_id" json:"event_id"`
	Event     event.Event     `gorm:"foreignKey:EventID;references:ID" json:"-"`
	CreatedAt time.Time       `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (Webhook) TableName() string {
	return "webhooks"
}