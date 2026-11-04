package webhook

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type IWebhookDelivery interface {
	Deliver(
		ctx context.Context,
		url string,
		body []byte,
		headerKey string,
		headerVal string,
	) error
}

type WebhookDelivery struct {
	httpClient *http.Client
}

func NewWebhookDelivery(httpClient *http.Client) IWebhookDelivery {
	return &WebhookDelivery{
		httpClient: httpClient,
	}
}

func (d *WebhookDelivery) Deliver(
	ctx context.Context,
	url string,
	body []byte,
	headerKey string,
	headerVal string,
) error {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if headerKey != "" {
		httpReq.Header.Set(headerKey, headerVal)
	}

	resp, err := d.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook endpoint responded with status %d", resp.StatusCode)
	}

	return nil
}
