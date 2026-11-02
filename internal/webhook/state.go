package webhook

import "time"

type WebhookState struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	EventName string    `gorm:"type:varchar(255);not null" json:"event_name"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (WebhookState) TableName() string {
	return "webhook_states"
}

