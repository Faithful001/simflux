package main

import (
	"github.com/Faithful001/simflux.git/internal/config/db"
	"github.com/Faithful001/simflux.git/internal/scenario"
)

func main() { 
	db.InitDB()
	
	migrationModels := []interface{}{scenario.Scenario{}}

	db.AutoMigrate(migrationModels)
	
}