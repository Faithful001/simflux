package webhook

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type WebhookConfigResult struct {
	URL       string
	HeaderKey string
	HeaderVal string
}

type IWebhookDispatcher interface {
	GetWebhookConfig() (WebhookConfigResult, error)

	// SetActiveEvent persists an event name override. Any subsequent webhook
	// delivery will use this event instead of the provider definition default.
	SetActiveEvent(eventName string) error

	// GetActiveEvent returns the persisted event override, or ("", nil) when
	// no override is set.
	GetActiveEvent() (string, error)

	// ResetActiveEvent clears the event override so the default is used.
	ResetActiveEvent() error
}

type WebhookDispatcher struct {
	db *gorm.DB
}

func NewWebhookDispatcher(db *gorm.DB) IWebhookDispatcher {
	return &WebhookDispatcher{db: db}
}

func (d *WebhookDispatcher) GetWebhookConfig() (WebhookConfigResult, error) {
	var result WebhookConfig

	if err := d.db.First(&result).Error; err != nil {
		return WebhookConfigResult{}, err
	}

	return WebhookConfigResult{
		URL:       result.URL,
		HeaderKey: result.SigHeaderKey,
		HeaderVal: result.SigHeaderValue,
	}, nil
}

func (d *WebhookDispatcher) SetActiveEvent(eventName string) error {
	// Delete any existing state row, then insert the new one (upsert pattern).
	if err := d.db.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&WebhookState{}).Error; err != nil {
		return fmt.Errorf("failed to clear webhook state: %w", err)
	}

	state := WebhookState{EventName: eventName}
	if err := d.db.Create(&state).Error; err != nil {
		return fmt.Errorf("failed to set active event: %w", err)
	}

	return nil
}

func (d *WebhookDispatcher) GetActiveEvent() (string, error) {
	var state WebhookState

	if err := d.db.First(&state).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read webhook state: %w", err)
	}

	return state.EventName, nil
}

func (d *WebhookDispatcher) ResetActiveEvent() error {
	if err := d.db.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&WebhookState{}).Error; err != nil {
		return fmt.Errorf("failed to reset webhook state: %w", err)
	}

	return nil
}
