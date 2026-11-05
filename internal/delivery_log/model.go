package delivery_log

import "time"

type DeliveryLog struct {
	ID               uint      `gorm:"primarykey"`
	PaymentReference string    `gorm:"index"`
	Event            string
	TargetURL        string
	HTTPStatus       int
	DurationMs       int64
	Success          bool
	ErrorMessage     string
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}
