package control

import (
	"context"
	"net/http"

	"github.com/Faithful001/simflux.git/internal/delivery_log"
	"github.com/Faithful001/simflux.git/internal/event"
	"github.com/Faithful001/simflux.git/internal/payment"
	"github.com/Faithful001/simflux.git/internal/webhook"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	dispatcher  webhook.IWebhookDispatcher
	eventRepo   event.IEventRepository
	paymentRepo payment.IPaymentRepository
	deliveryLog delivery_log.IDeliveryLogRepository
}

func NewHandler(
	dispatcher webhook.IWebhookDispatcher,
	eventRepo event.IEventRepository,
	paymentRepo payment.IPaymentRepository,
	deliveryLog delivery_log.IDeliveryLogRepository,
) *Handler {
	return &Handler{
		dispatcher:  dispatcher,
		eventRepo:   eventRepo,
		paymentRepo: paymentRepo,
		deliveryLog: deliveryLog,
	}
}

// GET /_simflux/state
func (h *Handler) GetState(c *gin.Context) {
	activeEvent, _ := h.dispatcher.GetActiveEvent()
	webhookCfg, _ := h.dispatcher.GetWebhookConfig()
	events, _ := h.eventRepo.Get(c.Request.Context())

	eventNames := make([]string, len(events))
	for i, e := range events {
		eventNames[i] = e.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"active_event":  activeEvent,
		"webhook_url":   webhookCfg.URL,
		"header_key":    webhookCfg.HeaderKey,
		"event_options": eventNames,
	})
}

// POST /_simflux/webhook/event  body: {"event":"payment.success"}
func (h *Handler) SetEvent(c *gin.Context) {
	var body struct {
		Event string `json:"event" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.eventRepo.GetByName(c.Request.Context(), body.Event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown event: " + body.Event})
		return
	}

	if err := h.dispatcher.SetActiveEvent(body.Event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"active_event": body.Event})
}

// DELETE /_simflux/webhook/event
func (h *Handler) ResetEvent(c *gin.Context) {
	if err := h.dispatcher.ResetActiveEvent(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"active_event": nil})
}

// POST /_simflux/reset
func (h *Handler) ResetPayments(c *gin.Context) {
	if err := h.paymentRepo.Clear(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payments cleared"})
}

// GET /_simflux/analytics
func (h *Handler) GetAnalytics(c *gin.Context) {
	ctx := c.Request.Context()

	payments, _ := h.paymentRepo.Get(ctx, "")
	statusCounts := map[string]int{}
	for _, p := range payments {
		statusCounts[p.Status]++
	}

	stats, _ := h.deliveryLog.Stats(ctx)
	recent, _ := h.deliveryLog.Recent(ctx, 20)

	c.JSON(http.StatusOK, gin.H{
		"payments": gin.H{
			"total":     len(payments),
			"by_status": statusCounts,
		},
		"webhook_deliveries": gin.H{
			"total":           stats.Total,
			"successful":      stats.Successful,
			"failed":          stats.Failed,
			"avg_duration_ms": stats.AvgMs,
		},
		"recent_deliveries": recent,
	})
}
