package webhook

import "time"

type WebhookConfig struct {
	ID        		int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	URL       		string    `gorm:"type:varchar(255);not null" json:"url"`
	SigHeaderKey 	string    `gorm:"type:varchar(255);not null" json:"sig_header_key"`
	SigHeaderValue 	string    `gorm:"type:varchar(255);not null" json:"sig_header_value"`
	CreatedAt 		time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt 		time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (WebhookConfig) TableName() string {
	return "webhook_configs"
}