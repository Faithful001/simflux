package webhook

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Faithful001/simflux.git/internal/delivery_log"
)

type IWebhookDelivery interface {
	Deliver(
		ctx context.Context,
		url string,
		body []byte,
		headerKey string,
		headerVal string,
		reference string,
		eventName string,
	) error
}

type WebhookDelivery struct {
	httpClient  *http.Client
	deliveryLog delivery_log.IDeliveryLogRepository
}

func NewWebhookDelivery(httpClient *http.Client, deliveryLog delivery_log.IDeliveryLogRepository) IWebhookDelivery {
	return &WebhookDelivery{
		httpClient:  httpClient,
		deliveryLog: deliveryLog,
	}
}

func (d *WebhookDelivery) Deliver(
	ctx context.Context,
	url string,
	body []byte,
	headerKey string,
	headerVal string,
	reference string,
	eventName string,
) error {
	start := time.Now()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if headerKey != "" {
		httpReq.Header.Set(headerKey, headerVal)
	}

	resp, err := d.httpClient.Do(httpReq)
	durationMs := time.Since(start).Milliseconds()

	var (
		httpStatus int
		success    bool
		errMsg     string
	)

	if err != nil {
		errMsg = err.Error()
		d.appendLog(ctx, reference, eventName, url, 0, durationMs, false, errMsg)
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	httpStatus = resp.StatusCode
	success = resp.StatusCode >= 200 && resp.StatusCode < 300

	if !success {
		errMsg = fmt.Sprintf("endpoint responded with status %d", httpStatus)
	}

	d.appendLog(ctx, reference, eventName, url, httpStatus, durationMs, success, errMsg)

	if !success {
		return fmt.Errorf("webhook endpoint responded with status %d", httpStatus)
	}

	return nil
}

func (d *WebhookDelivery) appendLog(
	ctx context.Context,
	reference, eventName, url string,
	httpStatus int,
	durationMs int64,
	success bool,
	errMsg string,
) {
	if d.deliveryLog == nil {
		return
	}
	_ = d.deliveryLog.Append(ctx, delivery_log.DeliveryLog{
		PaymentReference: reference,
		Event:            eventName,
		TargetURL:        url,
		HTTPStatus:       httpStatus,
		DurationMs:       durationMs,
		Success:          success,
		ErrorMessage:     errMsg,
	})
}
