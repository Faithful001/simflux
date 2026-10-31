package payment

import (
	"context"

	"github.com/Faithful001/simflux.git/internal/payment/dto"
)

type IPaymentService interface {
	Create(ctx context.Context, req dto.CreatePaymentRequest) error
	Get(ctx context.Context) ([]dto.PaymentResponse, error)
	GetByID(ctx context.Context, id int64) (dto.PaymentResponse, error)
	GetByReference(ctx context.Context, providerDefinitionID int64, reference string) (dto.PaymentResponse, error)
	GetByStatus(ctx context.Context, status string) ([]dto.PaymentResponse, error)
	GetByProviderDefinitionID(ctx context.Context, providerDefinitionID int64) ([]dto.PaymentResponse, error)
	Update(ctx context.Context, req dto.UpdatePaymentRequest) error
	Delete(ctx context.Context, id int64) error
	Clear(ctx context.Context) error
}

type PaymentService struct {
	repo IPaymentRepository
}

func NewPaymentService(repo IPaymentRepository) IPaymentService {
	return &PaymentService{
		repo: repo,
	}
}

func (s *PaymentService) Create(ctx context.Context, req dto.CreatePaymentRequest) error {
	return s.repo.Create(ctx, req)
}

func (s *PaymentService) Get(ctx context.Context) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.GetWithProviderDefinition(ctx)
	if err != nil {
		return nil, err
	}

	return toPaymentResponses(payments), nil
}

func (s *PaymentService) GetByID(ctx context.Context, id int64) (dto.PaymentResponse, error) {
	p, err := s.repo.GetByIDWithProviderDefinition(ctx, id)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return toPaymentResponse(p), nil
}

func (s *PaymentService) GetByReference(
	ctx context.Context,
	providerDefinitionID int64,
	reference string,
) (dto.PaymentResponse, error) {
	p, err := s.repo.GetByReferenceWithProviderDefinition(
		ctx,
		providerDefinitionID,
		reference,
	)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return toPaymentResponse(p), nil
}

func (s *PaymentService) GetByStatus(
	ctx context.Context,
	status string,
) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.GetByStatusWithProviderDefinition(ctx, status)
	if err != nil {
		return nil, err
	}

	return toPaymentResponses(payments), nil
}

func (s *PaymentService) GetByProviderDefinitionID(
	ctx context.Context,
	providerDefinitionID int64,
) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.GetByProviderDefinitionIDWithProviderDefinition(
		ctx,
		providerDefinitionID,
	)
	if err != nil {
		return nil, err
	}

	return toPaymentResponses(payments), nil
}

func (s *PaymentService) Update(ctx context.Context, req dto.UpdatePaymentRequest) error {
	return s.repo.Update(ctx, req)
}

func (s *PaymentService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *PaymentService) Clear(ctx context.Context) error {
	return s.repo.Clear(ctx)
}
