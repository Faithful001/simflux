package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Faithful001/simflux.git/internal/scenario"
	"github.com/spf13/cobra"
)

func NewScenarioCommand(repo scenario.IScenarioRepository) *cobra.Command {
	scenarioService := scenario.NewScenarioService(repo)

	cmd := &cobra.Command{
		Use:   "scenario",
		Short: "Manage Simflux scenarios",
		Long:  "Create, list, inspect, update, and delete payment simulation scenarios.",
	}

	cmd.AddCommand(
		newScenarioListCommand(scenarioService),
		newScenarioCreateCommand(scenarioService),
		newScenarioShowCommand(scenarioService),
		newScenarioUpdateCommand(scenarioService),
		newScenarioDeleteCommand(scenarioService),
	)

	return cmd
}

func newScenarioListCommand(service scenario.IScenarioService) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all scenarios",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			scenarios, err := service.Get(cmd.Context())
			if err != nil {
				return err
			}

			if len(scenarios) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No scenarios found.")
				return nil
			}

			for _, s := range scenarios {
				fmt.Fprintf(
					cmd.OutOrStdout(),
					"%d\t%s\t%s\n",
					s.ID,
					s.Name,
					s.Status,
				)
			}

			return nil
		},
	}
}

func newScenarioCreateCommand(service scenario.IScenarioService) *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new scenario",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" {
				return fmt.Errorf("scenario definition file is required")
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read scenario file: %w", err)
			}

			var s scenario.Scenario
			if err := json.Unmarshal(content, &s); err != nil {
				return fmt.Errorf("failed to parse scenario definition: %w", err)
			}

			if err := service.Create(cmd.Context(), s); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Scenario created successfully.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the scenario definition JSON file")

	return cmd
}

func newScenarioShowCommand(service scenario.IScenarioService) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show scenario details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid scenario ID: %s", args[0])
			}

			s, err := service.GetByID(cmd.Context(), id)
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")

			return encoder.Encode(s)
		},
	}
}

func newScenarioUpdateCommand(service scenario.IScenarioService) *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an existing scenario",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid scenario ID: %s", args[0])
			}

			if filePath == "" {
				return fmt.Errorf("scenario definition file is required")
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read scenario file: %w", err)
			}

			var s scenario.Scenario
			if err := json.Unmarshal(content, &s); err != nil {
				return fmt.Errorf("failed to parse scenario definition: %w", err)
			}

			s.ID = id

			if err := service.Update(cmd.Context(), s); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Scenario updated successfully.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the scenario definition JSON file")

	return cmd
}

func newScenarioDeleteCommand(service scenario.IScenarioService) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a scenario",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid scenario ID: %s", args[0])
			}

			if err := service.Delete(cmd.Context(), id); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Scenario deleted successfully.")
			return nil
		},
	}
}
