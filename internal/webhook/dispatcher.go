package webhook

import (
	"encoding/json"
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
	SaveWebhookConfig(url, headerKey, headerVal string) error
	SetActiveEvent(eventName string) error
	GetActiveEvent() (string, error)
	ResetActiveEvent() error
}

type WebhookDispatcher struct {
	db *gorm.DB
}

func NewWebhookDispatcher(db *gorm.DB) IWebhookDispatcher {
	return &WebhookDispatcher{db: db}
}

func (d *WebhookDispatcher) SaveWebhookConfig(url, headerKey, headerVal string) error {
	if err := d.db.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&WebhookConfig{}).Error; err != nil {
		return fmt.Errorf("failed to clear previous webhook config: %w", err)
	}

	config := WebhookConfig{
		URL:            url,
		SigHeaderKey:   headerKey,
		SigHeaderValue: headerVal,
	}

	if err := d.db.Create(&config).Error; err != nil {
		return fmt.Errorf("failed to save webhook config: %w", err)
	}

	return nil
}

func (d *WebhookDispatcher) GetWebhookConfig() (WebhookConfigResult, error) {
	var results []WebhookConfig

	if err := d.db.Limit(1).Find(&results).Error; err == nil && len(results) > 0 {
		return WebhookConfigResult{
			URL:       results[0].URL,
			HeaderKey: results[0].SigHeaderKey,
			HeaderVal: results[0].SigHeaderValue,
		}, nil
	}

	// Fallback: check provider_definitions table if webhook_configs was not populated
	var pds []struct {
		Webhook json.RawMessage
	}
	if err := d.db.Table("provider_definitions").Order("id desc").Limit(1).Find(&pds).Error; err == nil && len(pds) > 0 && len(pds[0].Webhook) > 0 {
		var fallback struct {
			URL       string `json:"url"`
			Signature struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"signature"`
		}
		if err := json.Unmarshal(pds[0].Webhook, &fallback); err == nil && fallback.URL != "" {
			return WebhookConfigResult{
				URL:       fallback.URL,
				HeaderKey: fallback.Signature.Key,
				HeaderVal: fallback.Signature.Value,
			}, nil
		}
	}

	return WebhookConfigResult{}, errors.New("no webhook configuration found; run `simflux load` first")
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
	var states []WebhookState

	if err := d.db.Limit(1).Find(&states).Error; err != nil {
		return "", fmt.Errorf("failed to read webhook state: %w", err)
	}

	if len(states) == 0 {
		return "", nil
	}

	return states[0].EventName, nil
}

func (d *WebhookDispatcher) ResetActiveEvent() error {
	if err := d.db.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&WebhookState{}).Error; err != nil {
		return fmt.Errorf("failed to reset webhook state: %w", err)
	}

	return nil
}
