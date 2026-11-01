package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Faithful001/simflux.git/internal/payment/dto"
	"github.com/Faithful001/simflux.git/internal/provider_definition"
)

type IPaymentService interface {
	Create(ctx context.Context, req dto.CreatePaymentRequest) error
	Get(ctx context.Context,
	status string,
	reference string,
	) ([]dto.PaymentResponse, error)
	GetByID(ctx context.Context, id int64) (dto.PaymentResponse, error)
	GetByReference(ctx context.Context, providerDefinitionID int64, reference string) (dto.PaymentResponse, error)
	GetByStatus(ctx context.Context, status string) ([]dto.PaymentResponse, error)
	GetByProviderDefinitionID(ctx context.Context, providerDefinitionID int64) ([]dto.PaymentResponse, error)
	Update(ctx context.Context, req dto.UpdatePaymentRequest) error
	Delete(ctx context.Context, id int64) error
	Clear(ctx context.Context) error
}

type PaymentService struct {
    paymentRepo            IPaymentRepository
    providerDefinitionRepo provider_definition.IProviderDefinitionRepository
}

func NewPaymentService(
    paymentRepo IPaymentRepository,
    providerDefinitionRepo provider_definition.IProviderDefinitionRepository,
) IPaymentService {
    return &PaymentService{
        paymentRepo:            paymentRepo,
        providerDefinitionRepo: providerDefinitionRepo,
    }
}

func (s *PaymentService) Create(
    ctx context.Context,
    req dto.CreatePaymentRequest,
) error {
    providerDefinition, err := s.providerDefinitionRepo.GetByID(
        ctx,
        req.ProviderDefinitionID,
    )
    if err != nil {
        return err
    }

    var definition provider_definition.Definition

    if err := json.Unmarshal(providerDefinition.Fields, &definition); err != nil {
        return fmt.Errorf("failed to decode provider definition: %w", err)
    }

    if err := provider_definition.ValidatePaymentDefinition(definition); err != nil {
        return err
    }

    if err := provider_definition.ValidatePaymentData(definition, req.Fields); err != nil {
        return err
    }

    return s.paymentRepo.Create(ctx, req)
}

func (s *PaymentService) Get(
	ctx context.Context,
	status string,
	reference string,
	) ([]dto.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetWithProviderDefinition(ctx, status, reference)
	if err != nil {
		return nil, err
	}

	paymentResponse, err := toPaymentResponses(payments)
	
	if err != nil {
		return nil, err
	}

	return paymentResponse, nil
}

func (s *PaymentService) GetByID(ctx context.Context, id int64) (dto.PaymentResponse, error) {
	p, err := s.paymentRepo.GetByIDWithProviderDefinition(ctx, id)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	payment, err := toPaymentResponse(p)
	
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return payment, nil
}

func (s *PaymentService) GetByReference(
	ctx context.Context,
	providerDefinitionID int64,
	reference string,
) (dto.PaymentResponse, error) {
	p, err := s.paymentRepo.GetByReferenceWithProviderDefinition(
		ctx,
		providerDefinitionID,
		reference,
	)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	payment, err := toPaymentResponse(p)
	
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return payment, nil
}

func (s *PaymentService) GetByStatus(
	ctx context.Context,
	status string,
) ([]dto.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetByStatusWithProviderDefinition(ctx, status)
	if err != nil {
		return nil, err
	}

	paymentResponse, err := toPaymentResponses(payments)
	
	if err != nil {
		return nil, err
	}

	return paymentResponse, nil
}

func (s *PaymentService) GetByProviderDefinitionID(
	ctx context.Context,
	providerDefinitionID int64,
) ([]dto.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetByProviderDefinitionIDWithProviderDefinition(
		ctx,
		providerDefinitionID,
	)
	if err != nil {
		return nil, err
	}

	paymentResponse, err := toPaymentResponses(payments)
	
	if err != nil {
		return nil, err
	}

	return paymentResponse, nil
}

func (s *PaymentService) Update(ctx context.Context, req dto.UpdatePaymentRequest) error {
	return s.paymentRepo.Update(ctx, req)
}

func (s *PaymentService) Delete(ctx context.Context, id int64) error {
	return s.paymentRepo.Delete(ctx, id)
}

func (s *PaymentService) Clear(ctx context.Context) error {
	return s.paymentRepo.Clear(ctx)
}
