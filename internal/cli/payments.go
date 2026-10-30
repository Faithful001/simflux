package cli

import (
	"fmt"

	"github.com/Faithful001/simflux.git/internal/payment"
	"github.com/spf13/cobra"
)

func NewPaymentsCommand(repo payment.IPaymentRepository) *cobra.Command {
	paymentService := payment.NewPaymentService(repo)

	cmd := &cobra.Command{
		Use: "payments",
		Short: "Manage Simflux payments",
		Long: "List, show, enable, disable, and clear simulation payments",
	}

	cmd.AddCommand(newPaymentsListCommand(paymentService))

	return cmd
}

func newPaymentsListCommand(service payment.IPaymentService) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Short: "List all payments",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			payments, err := service.Get(cmd.Context())

			if err != nil {
				return err
			}

			if len(payments) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No payments found")
				return nil
			}

			for _, p := range payments {
				fmt.Fprintf(
					cmd.OutOrStdout(),
					"%d\t%st%s\t%s\t%s\n",
					p.ID,
					p.Provider,
					p.Reference,
					p.Status,
					p.Data,
				)
			}

			return nil
		},
	}
}