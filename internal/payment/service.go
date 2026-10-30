package payment

import "github.com/Faithful001/simflux.git/internal/payment/dto"

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (*PaymentService) Create(req dto.CreatePaymentDto) {

}

func (*PaymentService) Verify() {

}

func (*PaymentService) Get() {

}