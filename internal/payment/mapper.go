package payment

import (
	"encoding/json"
	"fmt"

	"github.com/Faithful001/simflux.git/internal/payment/dto"
	"github.com/Faithful001/simflux.git/internal/provider_definition"
)

func toPaymentResponse(payment Payment) (dto.PaymentResponse, error) {
	fields := make(map[string]any)

	if len(payment.Fields) > 0 {
		if err := json.Unmarshal(payment.Fields, &fields); err != nil {
			return dto.PaymentResponse{}, fmt.Errorf("failed to decode payment attributes: %w", err)
		}
	}

	filteredFields := make(map[string]any)

	for key, value := range fields {
		if !provider_definition.IsReservedPaymentField(key) {
			filteredFields[key] = value
		}
	}

	return dto.PaymentResponse{
		ID:                   payment.ID,
		ProviderDefinitionID: payment.ProviderDefinitionID,
		Reference:            payment.Reference,
		Status:               payment.Status,
		CreatedAt:            payment.CreatedAt,
		UpdatedAt:            payment.UpdatedAt,
		Fields:               filteredFields,
	}, nil
}

func toPaymentResponses(payments []Payment) ([]dto.PaymentResponse, error) {
	responses := make([]dto.PaymentResponse, 0, len(payments))

	for _, p := range payments {
		response, err := toPaymentResponse(p)
		if err != nil {
			return nil, err
		}

		responses = append(responses, response)
	}

	return responses, nil
}