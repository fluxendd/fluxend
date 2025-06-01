package main

import (
	_ "fluxend/docs"
	"fluxend/internal/app/commands"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

// @title Fluxend API
// @version 1.0
// @description Fluxend is backend as-a-service platform that allows you to build, deploy, and scale applications without managing infrastructure.

// @contact.name API Support
// @contact.url http://github.com/fluxend/fluxend
// @contact.email hello@fluxend.app

// @host api.fluxend.app
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	commands.Execute()
}
