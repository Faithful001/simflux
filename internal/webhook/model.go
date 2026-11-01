package webhook

import "time"

type Webhook struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SigHeader string    `gorm:"type:varchar(255);not null" json:"sig_header"` // webhook signature header
	Body      int64     `gorm:"type:jsonb" json:"body"`                       // what the webhook should send in the request body (payload)
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (Webhook) TableName() string {
	return "webhooks"
}