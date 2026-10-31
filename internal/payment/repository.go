package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Faithful001/simflux.git/internal/payment/dto"
	"gorm.io/gorm"
)

type IPaymentRepository interface {
	Create(ctx context.Context, req dto.CreatePaymentRequest) error

	Get(ctx context.Context) ([]Payment, error)
	GetWithProviderDefinition(ctx context.Context) ([]Payment, error)

	GetByID(ctx context.Context, id int64) (Payment, error)
	GetByIDWithProviderDefinition(ctx context.Context, id int64) (Payment, error)

	GetByReference(ctx context.Context, providerDefinitionID int64, reference string) (Payment, error)
	GetByReferenceWithProviderDefinition(ctx context.Context, providerDefinitionID int64, reference string) (Payment, error)

	GetByStatus(ctx context.Context, status string) ([]Payment, error)
	GetByStatusWithProviderDefinition(ctx context.Context, status string) ([]Payment, error)

	GetByProviderDefinitionID(ctx context.Context, providerDefinitionID int64) ([]Payment, error)
	GetByProviderDefinitionIDWithProviderDefinition(ctx context.Context, providerDefinitionID int64) ([]Payment, error)

	Update(ctx context.Context, req dto.UpdatePaymentRequest) error
	Delete(ctx context.Context, id int64) error
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

func (r *PaymentRepository) Create(
    ctx context.Context,
    req dto.CreatePaymentRequest,
) error {
	fields, err := json.Marshal(req.Fields)
	if err != nil {
		return fmt.Errorf("failed to encode payment fields: %w", err)
	}

    payment := Payment{
        ProviderDefinitionID: req.ProviderDefinitionID,
        Reference:            req.Reference,
        Status:               req.Status,
        Fields:            	  fields,
    }

    return r.db.WithContext(ctx).Create(&payment).Error
}

func (r *PaymentRepository) Get(ctx context.Context) ([]Payment, error) {
	var payments []Payment

	err := r.db.WithContext(ctx).Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) GetWithProviderDefinition(ctx context.Context) ([]Payment, error) {
	var payments []Payment

	err := r.db.WithContext(ctx).
		Preload("ProviderDefinition").
		Find(&payments).Error

	return payments, err
}

func (r *PaymentRepository) GetByID(ctx context.Context, id int64) (Payment, error) {
	var payment Payment

	err := r.db.WithContext(ctx).First(&payment, id).Error
	return payment, err
}

func (r *PaymentRepository) GetByIDWithProviderDefinition(ctx context.Context, id int64) (Payment, error) {
	var payment Payment

	err := r.db.WithContext(ctx).
		Preload("ProviderDefinition").
		First(&payment, id).Error

	return payment, err
}

func (r *PaymentRepository) GetByReference(
	ctx context.Context,
	providerDefinitionID int64,
	reference string,
) (Payment, error) {
	var payment Payment

	err := r.db.WithContext(ctx).
		Where(
			"provider_definition_id = ? AND reference = ?",
			providerDefinitionID,
			reference,
		).
		First(&payment).Error

	return payment, err
}

func (r *PaymentRepository) GetByReferenceWithProviderDefinition(
	ctx context.Context,
	providerDefinitionID int64,
	reference string,
) (Payment, error) {
	var payment Payment

	err := r.db.WithContext(ctx).
		Preload("ProviderDefinition").
		Where(
			"provider_definition_id = ? AND reference = ?",
			providerDefinitionID,
			reference,
		).
		First(&payment).Error

	return payment, err
}

func (r *PaymentRepository) GetByStatus(ctx context.Context, status string) ([]Payment, error) {
	var payments []Payment

	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Find(&payments).Error

	return payments, err
}

func (r *PaymentRepository) GetByStatusWithProviderDefinition(
	ctx context.Context,
	status string,
) ([]Payment, error) {
	var payments []Payment

	err := r.db.WithContext(ctx).
		Preload("ProviderDefinition").
		Where("status = ?", status).
		Find(&payments).Error

	return payments, err
}

func (r *PaymentRepository) GetByProviderDefinitionID(
	ctx context.Context,
	providerDefinitionID int64,
) ([]Payment, error) {
	var payments []Payment

	err := r.db.WithContext(ctx).
		Where("provider_definition_id = ?", providerDefinitionID).
		Find(&payments).Error

	return payments, err
}

func (r *PaymentRepository) GetByProviderDefinitionIDWithProviderDefinition(
	ctx context.Context,
	providerDefinitionID int64,
) ([]Payment, error) {
	var payments []Payment

	err := r.db.WithContext(ctx).
		Preload("ProviderDefinition").
		Where("provider_definition_id = ?", providerDefinitionID).
		Find(&payments).Error

	return payments, err
}

func (r *PaymentRepository) Update(ctx context.Context, req dto.UpdatePaymentRequest) error {
	result := r.db.WithContext(ctx).
		Model(&Payment{}).
		Where(
			"provider_definition_id = ? AND reference = ?",
			req.ProviderDefinitionID,
			req.Reference,
		).
		Updates(map[string]any{
			"status": req.Status,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PaymentRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&Payment{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PaymentRepository) Clear(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&Payment{}).Error
}
