package router

import (
	"github.com/Faithful001/simflux.git/internal/payment"
	"github.com/Faithful001/simflux.git/internal/provider_definition"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	r *gin.Engine,
	paymentRepo payment.IPaymentRepository,
	providerDefinitionRepo provider_definition.IProviderDefinitionRepository,
	) {
		
	paymentService := payment.NewPaymentService(
		paymentRepo,
		providerDefinitionRepo,
	)
	paymentHandler := payment.NewPaymentHandler(paymentService)

	paymentRoutes := r.Group("/payments")
	{
		paymentRoutes.POST("", paymentHandler.Create)
		paymentRoutes.GET("", paymentHandler.Get)
		paymentRoutes.GET("/:id", paymentHandler.GetByID)
		paymentRoutes.GET("/reference", paymentHandler.GetByReference)
		paymentRoutes.GET("/status", paymentHandler.GetByStatus)
		paymentRoutes.GET("/provider/:provider_definition_id", paymentHandler.GetByProviderDefinitionID)
	}
}