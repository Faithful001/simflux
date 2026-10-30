package payment

import (
	"context"

	"github.com/Faithful001/simflux.git/internal/payment/dto"
)

type IPaymentService interface {
	Create(ctx context.Context, req dto.CreatePaymentDto) error
	Get(ctx context.Context) ([]Payment, error)
	GetByID(ctx context.Context, id int) (Payment, error)
	GetByReference(ctx context.Context, reference string) (Payment, error)
	GetByStatus(ctx context.Context, status string) ([]Payment, error)
	GetByProvider(ctx context.Context, provider string) ([]Payment, error)
	Update(ctx context.Context, req dto.UpdatePaymentDto) error
	Delete(ctx context.Context, id int) error 
}

type PaymentService struct{
	repo IPaymentRepository
}

func NewPaymentService(repo IPaymentRepository) IPaymentService {
	return &PaymentService{
		repo: repo,
	}
}

func (s *PaymentService) Create(ctx context.Context, req dto.CreatePaymentDto) error {
	return s.repo.Create(ctx, req)
}

func (*PaymentService) Verify() {

}

func (s *PaymentService) Get(ctx context.Context) ([]Payment, error) {
	return s.repo.Get(ctx) 
}

func (s *PaymentService) GetByID(ctx context.Context, id int) (Payment, error) {
	return s.repo.GetByID(ctx, id) 
}

func (s *PaymentService) GetByReference(ctx context.Context, reference string) (Payment, error) {
	return s.repo.GetByReference(ctx, reference)
}

func (s *PaymentService) GetByStatus(ctx context.Context, status string) ([]Payment, error) {
	return s.repo.GetByStatus(ctx, status) 
}
func (s *PaymentService) GetByProvider(ctx context.Context, provider string) ([]Payment, error) {
	return s.repo.GetByProvider(ctx, provider) 
}
func (s *PaymentService) Update(ctx context.Context, req dto.UpdatePaymentDto) error {
	return s.repo.Update(ctx, req) 
}

func (s *PaymentService) Delete(ctx context.Context, id int) error  {
	return s.repo.Delete(ctx, id) 
}

