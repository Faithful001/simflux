package provider_definition

type Definition struct {
	Name    string       `yaml:"name"`
	Version string       `yaml:"version"`
	Payment PaymentModel `yaml:"payment"`
	Webhook WebhookModel `yaml:"webhook"`
}

type WebhookModel struct {
	URL       string         `yaml:"url"`
	Signature SignatureModel `yaml:"signature"`
	events    []string       `yaml:"events"`
}

type SignatureModel struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type PaymentModel struct {
	Fields map[string]FieldDefinition `yaml:"fields"`
}

type FieldDefinition struct {
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}