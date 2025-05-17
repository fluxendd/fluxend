package main

import (
	"fluxton/cmd"
	_ "fluxton/docs"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

// @title Fluxton API
// @version 1.0
// @description Fluxton is backend as-a-service platform that allows you to build, deploy, and scale applications without managing infrastructure.

// @contact.name API Support
// @contact.url http://github.com/fluxton-io/fluxton
// @contact.email chief@fluxton.io

// @host fluxton.io/api
// @BasePath /v2
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cmd.Execute()
}
