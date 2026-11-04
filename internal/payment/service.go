package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Faithful001/simflux.git/internal/payment/dto"
	"github.com/Faithful001/simflux.git/internal/provider_definition"
	"github.com/Faithful001/simflux.git/internal/webhook"
	webhookDto "github.com/Faithful001/simflux.git/internal/webhook/dto"
)

type IPaymentService interface {
	Create(ctx context.Context, req dto.CreatePaymentRequest) error
	Get(ctx context.Context, status string) ([]dto.PaymentResponse, error)
	GetByID(ctx context.Context, id int64) (dto.PaymentResponse, error)
	GetByReference(ctx context.Context, reference string) (dto.PaymentResponse, error)
	GetByStatus(ctx context.Context, status string) ([]dto.PaymentResponse, error)
	Update(ctx context.Context, req dto.UpdatePaymentRequest) error
	Delete(ctx context.Context, id int64) error
	Clear(ctx context.Context) error
}

type PaymentService struct {
	paymentRepo            IPaymentRepository
	providerDefinitionRepo provider_definition.IProviderDefinitionRepository
	webhookDispatcher      webhook.IWebhookDispatcher
	webhookDelivery        webhook.IWebhookDelivery
}

func NewPaymentService(
	paymentRepo IPaymentRepository,
	providerDefinitionRepo provider_definition.IProviderDefinitionRepository,
	webhookDispatcher webhook.IWebhookDispatcher,
	webhookDelivery webhook.IWebhookDelivery,
) IPaymentService {
	return &PaymentService{
		paymentRepo:            paymentRepo,
		providerDefinitionRepo: providerDefinitionRepo,
		webhookDispatcher:      webhookDispatcher,
		webhookDelivery:        webhookDelivery,
	}
}

func (s *PaymentService) Create(
	ctx context.Context,
	req dto.CreatePaymentRequest,
) error {
	providerDefinition, err := s.providerDefinitionRepo.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("no active provider definition found: %w", err)
	}

	var paymentModel provider_definition.PaymentModel
	if err := json.Unmarshal(providerDefinition.Payment, &paymentModel); err != nil {
		return fmt.Errorf("failed to decode provider payment definition: %w", err)
	}

	definition := provider_definition.Definition{
		Name:     providerDefinition.Name,
		Provider: providerDefinition.Provider,
		Version:  providerDefinition.Version,
		Payment:  paymentModel,
	}

	if err := provider_definition.ValidatePaymentDefinition(definition); err != nil {
		return err
	}

	if err := provider_definition.ValidatePaymentData(definition, req.Fields); err != nil {
		return err
	}

	if err := s.paymentRepo.Create(ctx, req); err != nil {
		return err
	}

	if s.webhookDispatcher != nil && s.webhookDelivery != nil {
		go func() {
			bgCtx := context.WithoutCancel(ctx)

			result, err := s.webhookDispatcher.GetWebhookConfig()
			if err != nil {
				fmt.Printf("Error while retrieving the webhook endpoint: %v\n", err)
				return
			}

			// Resolve event name from the persisted override.
			// the developer must call `simflux webhook set --event <name>`.
			eventName, err := s.webhookDispatcher.GetActiveEvent()
			if err != nil {
				fmt.Printf("Error reading active webhook event: %v\n", err)
				return
			}
			if eventName == "" {
				fmt.Println("Webhook skipped: no active event set. Run `simflux webhook set --event <name>` to enable webhook delivery.")
				return
			}

			payload, err := json.Marshal(webhookDto.WebhookRequest{
				Event:     eventName,
				Reference: req.Reference,
				Status:    req.Status,
				Fields:    req.Fields,
			})
			if err != nil {
				fmt.Printf("Error encoding webhook payload: %v\n", err)
				return
			}

			if err := s.webhookDelivery.Deliver(
				bgCtx,
				result.URL,
				payload,
				result.HeaderKey,
				result.HeaderVal,
			); err != nil {
				fmt.Printf("Error delivering webhook: %v\n", err)
			} else {
				fmt.Printf("Webhook delivered successfully to %s (event: %s)\n", result.URL, eventName)
			}
		}()
	}

	return nil
}

func (s *PaymentService) Get(
	ctx context.Context,
	status string,
) ([]dto.PaymentResponse, error) {
	payments, err := s.paymentRepo.Get(ctx, status)
	if err != nil {
		return nil, err
	}

	return toPaymentResponses(payments)
}

func (s *PaymentService) GetByID(ctx context.Context, id int64) (dto.PaymentResponse, error) {
	p, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return toPaymentResponse(p)
}

func (s *PaymentService) GetByReference(
	ctx context.Context,
	reference string,
) (dto.PaymentResponse, error) {
	p, err := s.paymentRepo.GetByReference(ctx, reference)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return toPaymentResponse(p)
}

func (s *PaymentService) GetByStatus(
	ctx context.Context,
	status string,
) ([]dto.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	return toPaymentResponses(payments)
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
