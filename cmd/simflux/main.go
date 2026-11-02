package main

import (
	"net/http"
	"os"

	"github.com/Faithful001/simflux.git/internal/cli"
	"github.com/Faithful001/simflux.git/internal/config/db"
	"github.com/Faithful001/simflux.git/internal/event"
	"github.com/Faithful001/simflux.git/internal/payment"
	"github.com/Faithful001/simflux.git/internal/provider_definition"
	"github.com/Faithful001/simflux.git/internal/router"
	"github.com/Faithful001/simflux.git/internal/scenario"
	"github.com/Faithful001/simflux.git/internal/webhook"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	database := db.GetDB()

	if err := db.AutoMigrate(
		&payment.Payment{},
		&provider_definition.ProviderDefinition{},
		&event.Event{},
		&webhook.Webhook{},
		&webhook.WebhookConfig{},
		&webhook.WebhookState{},
		&scenario.Scenario{},
	); err != nil {
		os.Exit(1)
	}

	r := gin.Default()

	paymentRepo := payment.NewPaymentRepository(database)
	providerDefinitionRepo := provider_definition.NewProviderDefinitionRepository(database)
	providerDefinitionService := provider_definition.NewProviderDefinitionService(providerDefinitionRepo)
	eventRepo := event.NewEventRepository(database)
	webhookDispatcher := webhook.NewWebhookDispatcher(database)
	webhookDelivery := webhook.NewWebhookDelivery(http.DefaultClient)

	router.SetupRouter(
		r,
		paymentRepo,
		providerDefinitionRepo,
		webhookDispatcher,
		webhookDelivery,
	)

	scenarioRepo, _ := scenario.NewScenarioRepository(database)

	cmd := cli.NewRootCommand()
	cmd.AddCommand(
		cli.NewLoadCommand(providerDefinitionService, eventRepo),
		cli.NewPaymentsCommand(paymentRepo, providerDefinitionRepo, webhookDispatcher, webhookDelivery),
		cli.NewWebhookCommand(webhookDispatcher, eventRepo),
		cli.NewScenarioCommand(scenarioRepo),
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
