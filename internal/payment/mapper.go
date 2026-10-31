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

func toPaymentResponses(payments []Payment) []dto.PaymentResponse {
	responses := make([]dto.PaymentResponse, 0, len(payments))

	for _, p := range payments {
		responses = append(responses, toPaymentResponse(p))
	}

	return responses
}