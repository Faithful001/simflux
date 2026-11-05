package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Faithful001/simflux.git/internal/delivery_log"
	"github.com/Faithful001/simflux.git/internal/event"
	"github.com/Faithful001/simflux.git/internal/payment"
	"github.com/Faithful001/simflux.git/internal/provider_definition"
	"github.com/Faithful001/simflux.git/internal/router"
	"github.com/Faithful001/simflux.git/internal/webhook"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func NewServeCommand(
	paymentRepo payment.IPaymentRepository,
	providerDefinitionRepo provider_definition.IProviderDefinitionRepository,
	webhookDispatcher webhook.IWebhookDispatcher,
	webhookDelivery webhook.IWebhookDelivery,
	eventRepo event.IEventRepository,
	deliveryLog delivery_log.IDeliveryLogRepository,
) *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the Simflux HTTP server",
		Long:  "Start the Simflux HTTP server and listen for payment simulation requests.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := gin.Default()

			router.SetupRouter(
				r,
				paymentRepo,
				providerDefinitionRepo,
				webhookDispatcher,
				webhookDelivery,
				eventRepo,
				deliveryLog,
			)

			addr := fmt.Sprintf(":%d", port)
			srv := &http.Server{
				Addr:    addr,
				Handler: r,
			}

			// Listen for OS shutdown signals in the background.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				<-quit
				fmt.Fprintln(cmd.OutOrStdout(), "\nShutting down server...")

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				if err := srv.Shutdown(ctx); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Server forced to shutdown: %v\n", err)
				}
			}()

			fmt.Fprintf(cmd.OutOrStdout(), "Simflux server listening on :%d\n", port)
			fmt.Fprintf(cmd.OutOrStdout(), "Dashboard:  http://localhost:%d/ui\n", port)

			if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				return err
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")

	return cmd
}
