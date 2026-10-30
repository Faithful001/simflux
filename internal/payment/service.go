package payment

import (
	"context"

	"github.com/Faithful001/simflux.git/internal/payment/dto"
)

type IPaymentService interface {
	Create(ctx context.Context, req dto.CreatePaymentRequest) error
	Get(ctx context.Context) ([]dto.PaymentResponse, error)
	GetByID(ctx context.Context, id int) (dto.PaymentResponse, error)
	GetByReference(ctx context.Context, reference string) (dto.PaymentResponse, error)
	GetByStatus(ctx context.Context, status string) ([]dto.PaymentResponse, error)
	GetByProvider(ctx context.Context, provider string) ([]dto.PaymentResponse, error)
	Update(ctx context.Context, req dto.UpdatePaymentRequest) error
	Delete(ctx context.Context, id int) error
	Clear(ctx context.Context) error
}

type PaymentService struct{
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

func (*PaymentService) Verify() {

}

func (s *PaymentService) Get(ctx context.Context) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.Get(ctx)

	if err != nil {
		return nil, err
	}

	responses := make([]dto.PaymentResponse, 0, len(payments))

	for _, payment := range payments {
		responses = append(responses, toPaymentResponse(payment))
	}

	return responses, nil

}

func (s *PaymentService) GetByID(ctx context.Context, id int) (dto.PaymentResponse, error) {
	payment, err := s.repo.GetByID(ctx, id)
	
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	response := toPaymentResponse(payment)

	return response, nil
}

func (s *PaymentService) GetByReference(ctx context.Context, reference string) (dto.PaymentResponse, error) {
	payment, err := s.repo.GetByReference(ctx, reference)

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	response := toPaymentResponse(payment)
	
	return response, nil
}

func (s *PaymentService) GetByStatus(ctx context.Context, status string) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.GetByStatus(ctx, status) 

	if err != nil {
		return nil, err
	}

	responses := make([]dto.PaymentResponse, 0, len(payments))

	for _, payment := range payments {
		responses = append(responses, toPaymentResponse(payment))
	}

	return responses, nil
}
func (s *PaymentService) GetByProvider(ctx context.Context, provider string) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.GetByProvider(ctx, provider) 

	if err != nil {
		return nil, err
	}

	responses := make([]dto.PaymentResponse, 0, len(payments))

	for _, payment := range payments {
		responses = append(responses, toPaymentResponse(payment))
	}

	return responses, nil
}
func (s *PaymentService) Update(ctx context.Context, req dto.UpdatePaymentRequest) error {
	return s.repo.Update(ctx, req) 
}

func (s *PaymentService) Delete(ctx context.Context, id int) error  {
	return s.repo.Delete(ctx, id) 
}

func (s *PaymentService) Clear(ctx context.Context) error {
    return s.repo.Clear(ctx)
}

