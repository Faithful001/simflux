package payment

import "github.com/Faithful001/simflux.git/internal/payment/dto"

func toPaymentResponse(payment Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:                   payment.ID,
		ProviderDefinitionID: payment.ProviderDefinitionID,
		Reference:            payment.Reference,
		Status:               payment.Status,
		CreatedAt:            payment.CreatedAt,
		UpdatedAt:            payment.UpdatedAt,
	}
}