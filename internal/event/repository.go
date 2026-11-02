package event

import (
	"context"

	"gorm.io/gorm"
)

type IEventRepository interface {
	Create(ctx context.Context, event *Event) error
	Get(ctx context.Context) ([]Event, error)
	GetByID(ctx context.Context, id int64) (Event, error)
	GetByName(ctx context.Context, name string) (Event, error)
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, id int64) error
}

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) IEventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) Create(ctx context.Context, event *Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *EventRepository) Get(ctx context.Context) ([]Event, error) {
	var events []Event

	err := r.db.WithContext(ctx).Find(&events).Error
	return events, err
}

func (r *EventRepository) GetByID(ctx context.Context, id int64) (Event, error) {
	var event Event

	err := r.db.WithContext(ctx).First(&event, id).Error
	return event, err
}

func (r *EventRepository) GetByName(ctx context.Context, name string) (Event, error) {
	var event Event

	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&event).Error

	return event, err
}

func (r *EventRepository) Update(ctx context.Context, event *Event) error {
	result := r.db.WithContext(ctx).
		Model(&Event{}).
		Where("id = ?", event.ID).
		Updates(map[string]any{
			"name": event.Name,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *EventRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&Event{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

