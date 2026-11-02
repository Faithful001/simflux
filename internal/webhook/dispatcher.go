package webhook

import (
	"gorm.io/gorm"
)

type WebhookConfigResult struct {
	URL       string
	HeaderKey string
	HeaderVal string
}

type IWebhookDispatcher interface {
	GetWebhookConfig() (WebhookConfigResult, error)
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


