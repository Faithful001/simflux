package delivery_log

import (
	"context"

	"gorm.io/gorm"
)

type IDeliveryLogRepository interface {
	Append(ctx context.Context, log DeliveryLog) error
	Recent(ctx context.Context, limit int) ([]DeliveryLog, error)
	Stats(ctx context.Context) (DeliveryStats, error)
}

type DeliveryStats struct {
	Total      int64
	Successful int64
	Failed     int64
	AvgMs      float64
}

type DeliveryLogRepository struct {
	db *gorm.DB
}

func NewDeliveryLogRepository(db *gorm.DB) IDeliveryLogRepository {
	return &DeliveryLogRepository{db: db}
}

func (r *DeliveryLogRepository) Append(ctx context.Context, log DeliveryLog) error {
	return r.db.WithContext(ctx).Create(&log).Error
}

func (r *DeliveryLogRepository) Recent(ctx context.Context, limit int) ([]DeliveryLog, error) {
	var logs []DeliveryLog
	err := r.db.WithContext(ctx).Order("created_at desc").Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *DeliveryLogRepository) Stats(ctx context.Context) (DeliveryStats, error) {
	var stats DeliveryStats

	r.db.WithContext(ctx).Model(&DeliveryLog{}).Count(&stats.Total)
	r.db.WithContext(ctx).Model(&DeliveryLog{}).Where("success = ?", true).Count(&stats.Successful)
	r.db.WithContext(ctx).Model(&DeliveryLog{}).Where("success = ?", false).Count(&stats.Failed)

	var avgMs *float64
	r.db.WithContext(ctx).Model(&DeliveryLog{}).Select("AVG(duration_ms)").Scan(&avgMs)
	if avgMs != nil {
		stats.AvgMs = *avgMs
	}

	return stats, nil
}


