package main

import (
	"net/http"
	"os"

	"github.com/Faithful001/simflux.git/internal/cli"
	"github.com/Faithful001/simflux.git/internal/config/db"
	"github.com/Faithful001/simflux.git/internal/payment"
	"github.com/Faithful001/simflux.git/internal/provider_definition"
	"github.com/Faithful001/simflux.git/internal/router"
	"github.com/Faithful001/simflux.git/internal/webhook"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	database := db.GetDB()

	if err := db.AutoMigrate(); err != nil {
		os.Exit(1)
	}

	r := gin.Default()

	paymentRepo := payment.NewPaymentRepository(database)
	providerDefinitionRepo := provider_definition.NewProviderDefinitionRepository(database)
	webhookDispatcher := webhook.NewWebhookDispatcher(database)
	webhookDelivery := webhook.NewWebhookDelivery(http.DefaultClient)

	router.SetupRouter(
		r,
		paymentRepo,
		providerDefinitionRepo,
		webhookDispatcher,
		webhookDelivery,
	)

	cmd := cli.NewRootCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
