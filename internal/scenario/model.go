package scenario

import (
	"time"
)

type Scenario struct {
	ID     		int  		`gorm:"type:bigint;primaryKey;autoIncrement" json:"id"`
	Name   		string 		`gorm:"type:varchar(255);not null" json:"name"`
	Type   		string		`gorm:"type:varchar(100);default:'success'" json:"type"`
	Status 		string 		`gorm:"type:varchar(100);default:'ready'" json:"status"`
	CreatedAt 	time.Time 	`gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (Scenario) TableName() string {
	return "scenarios"
}

