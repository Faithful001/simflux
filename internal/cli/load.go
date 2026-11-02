package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Faithful001/simflux.git/internal/event"
	provider "github.com/Faithful001/simflux.git/internal/provider_definition"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewLoadCommand(
	pdService provider.IProviderDefinitionService,
	eventRepo event.IEventRepository,
) *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load the Simflux provider definition",
		Long:  "Load and validate the Simflux provider definition YAML file.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := resolveDefinitionPath(filePath)
			if err != nil {
				return err
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read definition file: %w", err)
			}

			var definition provider.Definition

			if err := yaml.Unmarshal(content, &definition); err != nil {
				return fmt.Errorf("invalid YAML definition: %w", err)
			}

			if definition.Name == "" {
				return fmt.Errorf("provider name is required")
			}

			if definition.Version == "" {
				return fmt.Errorf("provider version is required")
			}

			if len(definition.Payment.Fields) == 0 {
				return fmt.Errorf("at least one payment field is required")
			}

			if err := provider.ValidatePaymentDefinition(definition); err != nil {
				return err
			}

			paymentBytes, err := json.Marshal(definition.Payment)
			if err != nil {
				return fmt.Errorf("failed to marshal payment definition: %w", err)
			}

			webhookBytes, err := json.Marshal(definition.Webhook)
			if err != nil {
				return fmt.Errorf("failed to marshal webhook definition: %w", err)
			}

			providerDefinition := provider.ProviderDefinition{
				Name:     definition.Name,
				Provider: definition.Provider,
				Version:  definition.Version,
				Payment:  paymentBytes,
				Webhook:  webhookBytes,
			}

			if err := pdService.SaveActive(cmd.Context(), &providerDefinition); err != nil {
				return fmt.Errorf("failed to save provider definition: %w", err)
			}

			if eventRepo != nil {
				for _, eventName := range definition.Webhook.Events {
					trimmed := strings.TrimSpace(eventName)
					if trimmed == "" {
						continue
					}
					if _, err := eventRepo.GetByName(cmd.Context(), trimmed); err != nil {
						if err := eventRepo.Create(cmd.Context(), &event.Event{Name: trimmed}); err != nil {
							return fmt.Errorf("failed to register event %q: %w", trimmed, err)
						}
					}
				}
			}

			fmt.Fprintf(
				cmd.OutOrStdout(),
				"Successfully loaded provider definition: %s (v%s)\n",
				definition.Name,
				definition.Version,
			)

			return nil
		},
	}

	cmd.Flags().StringVarP(
		&filePath,
		"file",
		"f",
		"",
		"Path to the Simflux definition YAML file",
	)

	return cmd
}

func resolveDefinitionPath(filePath string) (string, error) {
	if filePath != "" {
		return filePath, nil
	}

	for _, name := range []string{"simflux.yaml", "simflux.yml"} {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to access %s: %w", name, err)
		}
	}

	return "", fmt.Errorf("no Simflux definition found; expected simflux.yaml or simflux.yml")
}
