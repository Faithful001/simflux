package main

import (
	"os"

	"github.com/Faithful001/simflux.git/internal/cli"
	"github.com/Faithful001/simflux.git/internal/config/db"
	"github.com/Faithful001/simflux.git/internal/payment"
	"github.com/Faithful001/simflux.git/internal/provider_definition"
	"github.com/Faithful001/simflux.git/internal/router"
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

	router.SetupRouter(r, paymentRepo, providerDefinitionRepo)


	cmd := cli.NewRootCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}