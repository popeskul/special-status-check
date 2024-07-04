package main

import (
	"github.com/popeskul/special-status-check/internal/app"
	"github.com/popeskul/special-status-check/internal/config"
	"log"
)

// @title Special Status Check Service API
// @version 1.0
// @description This is a sample api for Special Status Check Service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.LoadConfig([]string{"./config", "."})

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	application.Run()
}
