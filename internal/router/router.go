package router

import (
	"github.com/Faithful001/simflux.git/internal/delivery_log"
	"github.com/Faithful001/simflux.git/internal/control"
	"github.com/Faithful001/simflux.git/internal/event"
	"github.com/Faithful001/simflux.git/internal/payment"
	"github.com/Faithful001/simflux.git/internal/provider_definition"
	"github.com/Faithful001/simflux.git/internal/ui"
	"github.com/Faithful001/simflux.git/internal/webhook"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	r *gin.Engine,
	paymentRepo payment.IPaymentRepository,
	providerDefinitionRepo provider_definition.IProviderDefinitionRepository,
	webhookDispatcher webhook.IWebhookDispatcher,
	webhookDelivery webhook.IWebhookDelivery,
	eventRepo event.IEventRepository,
	deliveryLog delivery_log.IDeliveryLogRepository,
) {
	paymentService := payment.NewPaymentService(
		paymentRepo,
		providerDefinitionRepo,
		webhookDispatcher,
		webhookDelivery,
	)
	paymentHandler := payment.NewPaymentHandler(paymentService)

	paymentRoutes := r.Group("/payments")
	{
		paymentRoutes.POST("", paymentHandler.Create)
		paymentRoutes.GET("", paymentHandler.Get)
		paymentRoutes.GET("/:id", paymentHandler.GetByID)
		paymentRoutes.GET("/:id/status", paymentHandler.GetStatusByID)
		paymentRoutes.DELETE("/:id", paymentHandler.Delete)
		paymentRoutes.GET("/reference/:reference", paymentHandler.GetByReference)
		paymentRoutes.GET("/reference/:reference/status", paymentHandler.GetStatusByReference)
		paymentRoutes.PATCH("/reference/:reference", paymentHandler.Update)
	}

	controlHandler := control.NewHandler(webhookDispatcher, eventRepo, paymentRepo, deliveryLog)
	simflux := r.Group("/_simflux")
	{
		simflux.GET("/state", controlHandler.GetState)
		simflux.POST("/webhook/event", controlHandler.SetEvent)
		simflux.DELETE("/webhook/event", controlHandler.ResetEvent)
		simflux.POST("/reset", controlHandler.ResetPayments)
		simflux.GET("/analytics", controlHandler.GetAnalytics)
	}

	ui.RegisterRoutes(r)
}
