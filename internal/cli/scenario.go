package cli

import "github.com/spf13/cobra"

func NewScenarioCommand() *cobra.Command {
	scenarioCommand := cobra.Command{
		Use: "scenario",
		Short: "Handle the scenario of simflux",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return &scenarioCommand
}