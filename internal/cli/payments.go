package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Faithful001/simflux.git/internal/payment"
	"github.com/spf13/cobra"
)

func NewPaymentsCommand(repo payment.IPaymentRepository) *cobra.Command {
	paymentService := payment.NewPaymentService(repo)

	cmd := &cobra.Command{
		Use:   "payments",
		Short: "Manage Simflux payments",
		Long:  "List, inspect, and clear simulated payment records.",
	}

	cmd.AddCommand(
		newPaymentsListCommand(paymentService),
		newPaymentsShowCommand(paymentService),
		newPaymentsClearCommand(paymentService),
	)

	return cmd
}

func newPaymentsListCommand(service payment.IPaymentService) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all payments",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			payments, err := service.Get(cmd.Context())
			if err != nil {
				return err
			}

			if len(payments) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No payments found.")
				return nil
			}

			fmt.Fprintf(
				cmd.OutOrStdout(),
				"%-8s %-25s %-25s %-15s\n",
				"ID",
				"PROVIDER DEFINITION ID",
				"REFERENCE",
				"STATUS",
			)

			for _, p := range payments {
				fmt.Fprintf(
					cmd.OutOrStdout(),
					"%-8d %-25d %-25s %-15s\n",
					p.ID,
					p.ProviderDefinitionID,
					p.Reference,
					p.Status,
				)
			}

			return nil
		},
	}
}

func newPaymentsShowCommand(service payment.IPaymentService) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show payment details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid payment ID: %s", args[0])
			}

			p, err := service.GetByID(cmd.Context(), id)
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")

			return encoder.Encode(p)
		},
	}
}

func newPaymentsClearCommand(service payment.IPaymentService) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all simulated payments",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				return fmt.Errorf("clearing payments requires --force")
			}

			if err := service.Clear(cmd.Context()); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "All payments cleared successfully.")
			return nil
		},
	}

	cmd.Flags().BoolVar(
		&force,
		"force",
		false,
		"Skip confirmation when clearing payments",
	)

	return cmd
}
