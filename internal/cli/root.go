package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	rootCmd := cobra.Command{
		Use: "simflux",
		Short: "Test payment integrations locally",
		Long: "Simflux helps you simulate payment providers, test webhooks, and reproduce payment failures locally.",
	}

	return &rootCmd
}
