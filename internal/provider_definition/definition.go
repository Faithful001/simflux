package provider_definition

type Definition struct {
	Name     string       `yaml:"name" json:"name"`
	Provider string       `yaml:"provider" json:"provider"`
	Version  string       `yaml:"version" json:"version"`
	Payment  PaymentModel `yaml:"payment" json:"payment"`
	Webhook  WebhookModel `yaml:"webhook" json:"webhook"`
}

type WebhookModel struct {
	URL       string         `yaml:"url" json:"url"`
	Signature SignatureModel `yaml:"signature" json:"signature"`
	Events    []string       `yaml:"events" json:"events"`
}

type SignatureModel struct {
	Key   string `yaml:"key" json:"key"`
	Value string `yaml:"value" json:"value"`
}

type PaymentModel struct {
	Fields map[string]FieldDefinition `yaml:"fields" json:"fields"`
}

type FieldDefinition struct {
	Type     string `yaml:"type" json:"type"`
	Required bool   `yaml:"required" json:"required"`
}
