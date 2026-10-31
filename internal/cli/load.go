package cli

import (
	"fmt"
	"os"

	"github.com/Faithful001/simflux.git/internal/provider"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewLoadCommand() *cobra.Command {
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

			var definition provider.ProviderDefinition

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