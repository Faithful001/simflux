package payment

import "github.com/Faithful001/simflux.git/internal/payment/dto"

func toPaymentResponse(payment Payment) dto.PaymentResponse {
    return dto.PaymentResponse{
        ID:        payment.ID,
        Provider:  payment.Provider,
        Reference: payment.Reference,
        Status:    payment.Status,
        Data:      payment.Data,
        CreatedAt: payment.CreatedAt,
        UpdatedAt: payment.UpdatedAt,
    }
}