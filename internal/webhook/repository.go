package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Faithful001/simflux.git/internal/event"
	"gorm.io/gorm"
)

type IWebhookRepository interface {
	Create(ctx context.Context, sigHeader string, eventName string, fields map[string]any) error
	CreateWebhook(ctx context.Context, webhook *Webhook) error

	Get(ctx context.Context) ([]Webhook, error)
	GetWithEvent(ctx context.Context) ([]Webhook, error)

	GetByID(ctx context.Context, id int64) (Webhook, error)
	GetByIDWithEvent(ctx context.Context, id int64) (Webhook, error)

	GetByEventID(ctx context.Context, eventID int64) ([]Webhook, error)
	GetByEventIDWithEvent(ctx context.Context, eventID int64) ([]Webhook, error)

	GetByEventName(ctx context.Context, eventName string) ([]Webhook, error)

	Update(ctx context.Context, webhook *Webhook) error
	Delete(ctx context.Context, id int64) error
	Clear(ctx context.Context) error
}

type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) IWebhookRepository {
	return &WebhookRepository{
		db: db,
	}
}

func (w *WebhookRepository) Create(
	ctx context.Context,
	sigHeader string,
	eventName string,
	fields map[string]any,
) error {
	var fieldsJSON []byte
	var err error

	if fields != nil {
		fieldsJSON, err = json.Marshal(fields)
		if err != nil {
			return fmt.Errorf("failed to encode webhook fields: %w", err)
		}
	} else {
		fieldsJSON = []byte("{}")
	}

	var ev event.Event
	if err := w.db.WithContext(ctx).
		Where("name = ?", eventName).
		FirstOrCreate(&ev, event.Event{Name: eventName}).Error; err != nil {
		return fmt.Errorf("failed to resolve event: %w", err)
	}

	webhook := Webhook{
		SigHeader: sigHeader,
		Fields:    fieldsJSON,
		EventID:   ev.ID,
	}

	return w.db.WithContext(ctx).Create(&webhook).Error
}

func (w *WebhookRepository) CreateWebhook(ctx context.Context, webhook *Webhook) error {
	return w.db.WithContext(ctx).Create(webhook).Error
}

func (w *WebhookRepository) Get(ctx context.Context) ([]Webhook, error) {
	var webhooks []Webhook

	err := w.db.WithContext(ctx).Find(&webhooks).Error
	return webhooks, err
}

func (w *WebhookRepository) GetWithEvent(ctx context.Context) ([]Webhook, error) {
	var webhooks []Webhook

	err := w.db.WithContext(ctx).
		Preload("Event").
		Find(&webhooks).Error

	return webhooks, err
}

func (w *WebhookRepository) GetByID(ctx context.Context, id int64) (Webhook, error) {
	var webhook Webhook

	err := w.db.WithContext(ctx).First(&webhook, id).Error
	return webhook, err
}

func (w *WebhookRepository) GetByIDWithEvent(ctx context.Context, id int64) (Webhook, error) {
	var webhook Webhook

	err := w.db.WithContext(ctx).
		Preload("Event").
		First(&webhook, id).Error

	return webhook, err
}

func (w *WebhookRepository) GetByEventID(ctx context.Context, eventID int64) ([]Webhook, error) {
	var webhooks []Webhook

	err := w.db.WithContext(ctx).
		Where("event_id = ?", eventID).
		Find(&webhooks).Error

	return webhooks, err
}

func (w *WebhookRepository) GetByEventIDWithEvent(ctx context.Context, eventID int64) ([]Webhook, error) {
	var webhooks []Webhook

	err := w.db.WithContext(ctx).
		Preload("Event").
		Where("event_id = ?", eventID).
		Find(&webhooks).Error

	return webhooks, err
}

func (w *WebhookRepository) GetByEventName(ctx context.Context, eventName string) ([]Webhook, error) {
	var webhooks []Webhook

	err := w.db.WithContext(ctx).
		Preload("Event").
		Joins("JOIN events ON events.id = webhooks.event_id").
		Where("events.name = ?", eventName).
		Find(&webhooks).Error

	return webhooks, err
}

func (w *WebhookRepository) Update(ctx context.Context, webhook *Webhook) error {
	result := w.db.WithContext(ctx).
		Model(&Webhook{}).
		Where("id = ?", webhook.ID).
		Updates(map[string]any{
			"sig_header": webhook.SigHeader,
			"fields":     webhook.Fields,
			"event_id":   webhook.EventID,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (w *WebhookRepository) Delete(ctx context.Context, id int64) error {
	result := w.db.WithContext(ctx).Delete(&Webhook{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (w *WebhookRepository) Clear(ctx context.Context) error {
	return w.db.WithContext(ctx).
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&Webhook{}).Error
}

