package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
    fiberSwagger "github.com/swaggo/fiber-swagger"
    _ "github.com/lokesh2201013/postgres-data-summary/docs" 
	"github.com/lokesh2201013/postgres-data-summary/internal/handler"
	"github.com/lokesh2201013/postgres-data-summary/internal/logger"
	"github.com/lokesh2201013/postgres-data-summary/internal/repository/external"
	"github.com/lokesh2201013/postgres-data-summary/internal/repository/local"
	"github.com/lokesh2201013/postgres-data-summary/internal/router"
	"github.com/lokesh2201013/postgres-data-summary/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	local.ConnectDB()

	repo := local.NewSummaryRepository()

	client := external.NewSummaryClient()
	summarySvc := service.NewSummaryService(repo, client)
	h := handler.NewSummaryHandler(summarySvc)

	app := fiber.New()
    logger.InitLogger()
	//app.Use(logger.New())
    app.Use(logger.ZapLogger())
	router.SummaryRoutes(app, h)
     app.Get("/swagger/*", fiberSwagger.WrapHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	log.Printf("Server listening on %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
