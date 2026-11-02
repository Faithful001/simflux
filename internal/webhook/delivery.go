package webhook

import (
	"bytes"
	"context"
	"net/http"
)

type IWebhookDelivery interface{
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
		return err
	}

	httpReq.Header.Set(headerKey, headerVal)

	_, err = d.httpClient.Do(httpReq)

	if err != nil {
		return err
	}

	return nil
}