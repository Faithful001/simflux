package cli

import (
	"fmt"

	"github.com/Faithful001/simflux.git/internal/event"
	"github.com/Faithful001/simflux.git/internal/webhook"
	"github.com/spf13/cobra"
)

func NewWebhookCommand(
	webhookDispatcher webhook.IWebhookDispatcher,
	eventRepo event.IEventRepository,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Manage webhook event simulation",
		Long:  "Set, inspect, and reset the active webhook event used during payment simulation.",
	}

	cmd.AddCommand(
		newWebhookSetCommand(webhookDispatcher, eventRepo),
		newWebhookGetCommand(webhookDispatcher),
		newWebhookResetCommand(webhookDispatcher),
	)

	return cmd
}

func newWebhookSetCommand(
	dispatcher webhook.IWebhookDispatcher,
	eventRepo event.IEventRepository,
) *cobra.Command {
	var eventName string

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set the active webhook event",
		Long: `Set the webhook event that will be sent on the next payment creation.
The event must be registered (via simflux load) before it can be set.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if eventName == "" {
				return fmt.Errorf("--event is required")
			}

			// Validate that the event is registered.
			if _, err := eventRepo.GetByName(cmd.Context(), eventName); err != nil {
				return fmt.Errorf("event %q is not registered; add it to your simflux.yaml and run simflux load", eventName)
			}

			if err := dispatcher.SetActiveEvent(eventName); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Active webhook event set to: %s\n", eventName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&eventName, "event", "e", "", "Event name to activate (e.g. charge.failed)")
	_ = cmd.MarkFlagRequired("event")

	return cmd
}

func newWebhookGetCommand(dispatcher webhook.IWebhookDispatcher) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Show the active webhook event",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			eventName, err := dispatcher.GetActiveEvent()
			if err != nil {
				return err
			}

			if eventName == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "No active event set. Webhook delivery is disabled until you run `simflux webhook set --event <name>`.")
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Active webhook event: %s\n", eventName)
			return nil
		},
	}
}

func newWebhookResetCommand(dispatcher webhook.IWebhookDispatcher) *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset the active webhook event to default",
		Long:  "Clears the active event override. Webhooks will use the payment status as the event name.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dispatcher.ResetActiveEvent(); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Webhook event override cleared.")
			return nil
		},
	}
}

