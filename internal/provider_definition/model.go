package provider_definition

import (
	"encoding/json"
	"time"
)

type ProviderDefinition struct {
	ID           	int64 				`gorm:"primaryKey;autoIncrement" json:"id"`
	Name         	string 				`gorm:"type:varchar(255);not null" json:"name"`
	ProviderName	string				`gorm:"type:varchar(255);not null" json:"provider_name"`
	Version      	string 				`gorm:"type:varchar(100);not null" json:"version"`
	Fields 			json.RawMessage		`gorm:"type:jsonb;not null" json:"fields"`

	Payment			json.RawMessage		`gorm:"type:jsonb;not null" json:"payment"`
	Webhook			json.RawMessage		`gorm:"type:jsonb;not null" json:"webhook"`

	CreatedAt 		time.Time 			`gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt 		time.Time 			`gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (ProviderDefinition) TableName() string {
	return "provider_definitions"
}