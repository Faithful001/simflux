package payment

import (
	"context"

	"github.com/Faithful001/simflux.git/internal/payment/dto"
	"gorm.io/gorm"
)

type IPaymentRepository interface {
	Create(ctx context.Context, req dto.CreatePaymentRequest) error
	Get(ctx context.Context) ([]Payment, error)
	GetByID(ctx context.Context, id int) (Payment, error)
	GetByReference(ctx context.Context, reference string) (Payment, error)
	GetByStatus(ctx context.Context, status string) ([]Payment, error)
	GetByProvider(ctx context.Context, provider string) ([]Payment, error)
	Update(ctx context.Context, req dto.UpdatePaymentRequest) error
	Delete(ctx context.Context, id int) error
	Clear(ctx context.Context) error
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}


func (r *PaymentRepository) Create(ctx context.Context, req dto.CreatePaymentRequest) error {
	payment := Payment{
        Provider:  req.Provider,
        Reference: req.Reference,
        Status:    req.Status,
        Data:      req.Data,
    }

	return r.db.WithContext(ctx).Create(&payment).Error
}

func (r *PaymentRepository) Get(ctx context.Context) ([]Payment, error) {
	var payments []Payment
	if err := r.db.WithContext(ctx).Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id int) (Payment, error) {
	var payment Payment
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error; err != nil {
		return Payment{}, err
	}

	return payment, nil
}

func (r *PaymentRepository) GetByReference(ctx context.Context, reference string) (Payment, error) {
	var payment Payment
	err := r.db.WithContext(ctx).Where("reference = ?", reference).First(&payment).Error

	return payment, err
}

func (r *PaymentRepository) GetByStatus(ctx context.Context, status string) ([]Payment, error) {
	var payments []Payment
	err := r.db.WithContext(ctx).Where("status = ?", status).First(&payments).Error

	return payments, err
}

func (r *PaymentRepository) GetByProvider(ctx context.Context, provider string) ([]Payment, error) {
	var payments []Payment
	err := r.db.WithContext(ctx).Where("provider = ?", provider).First(&payments).Error

	return payments, err
}

func (r *PaymentRepository) Update(ctx context.Context, req dto.UpdatePaymentRequest) error {
    return r.db.WithContext(ctx).
        Model(&Payment{}).
        Where("reference = ?", req.Reference).
        Updates(map[string]any{
            "provider": req.Provider,
            "status":   req.Status,
            "data":     req.Data,
        }).Error
}

func (r *PaymentRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Payment{}).Error
}

func (r *PaymentRepository) Clear(ctx context.Context) error {
    return r.db.
        WithContext(ctx).
        Session(&gorm.Session{AllowGlobalUpdate: true}).
        Delete(&Payment{}).
        Error
}