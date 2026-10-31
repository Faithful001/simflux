package payment

import (
	"net/http"
	"strconv"

	"github.com/Faithful001/simflux.git/internal/payment/dto"
	"github.com/gin-gonic/gin"
)

type IPaymentHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetByID(c *gin.Context)
	GetByReference(c *gin.Context)
	GetByStatus(c *gin.Context)
	GetByProviderDefinitionID(c *gin.Context)
}

type PaymentHandler struct {
	service IPaymentService
}

func NewPaymentHandler(service IPaymentService) IPaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var req dto.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Payment created successfully"})
}

func (h *PaymentHandler) Get(c *gin.Context) {
	payments, err := h.service.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetByReference(c *gin.Context) {
	providerDefinitionID, err := strconv.ParseInt(c.Query("provider_definition_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider definition ID"})
		return
	}

	reference := c.Query("reference")
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reference is required"})
		return
	}

	payment, err := h.service.GetByReference(
		c.Request.Context(),
		providerDefinitionID,
		reference,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	payments, err := h.service.GetByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) GetByProviderDefinitionID(c *gin.Context) {
	providerDefinitionID, err := strconv.ParseInt(c.Param("provider_definition_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider definition ID"})
		return
	}

	payments, err := h.service.GetByProviderDefinitionID(
		c.Request.Context(),
		providerDefinitionID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}