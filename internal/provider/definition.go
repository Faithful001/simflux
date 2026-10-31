package provider

type ProviderDefinition struct {
	Name    string       `yaml:"name"`
	Version string       `yaml:"version"`
	Payment PaymentModel `yaml:"payment"`
}

type PaymentModel struct {
	Fields map[string]FieldDefinition `yaml:"fields"`
}

type FieldDefinition struct {
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}