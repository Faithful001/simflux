// retrieve provider_definition fields and return it - i need a method for that

package webhook

type IWebhookService interface {
}

type WebhookService struct {
	repo IWebhookRepository
}

func NewWebhookService(repo IWebhookRepository) IWebhookService {
	return &WebhookService{repo: repo}
}


