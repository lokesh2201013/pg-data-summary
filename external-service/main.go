package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/lokesh2201013/postgres-data-summary/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	//"github.com/lokesh2201013/postgres-data-summary/internal/handler"
	"github.com/lokesh2201013/postgres-data-summary/internal/logger"
	//"github.com/lokesh2201013/postgres-data-summary/internal/repository/external"
	"github.com/lokesh2201013/postgres-data-summary/external-service/database"
	"github.com/lokesh2201013/postgres-data-summary/external-service/routes"
	//"github.com/lokesh2201013/postgres-data-summary/internal/repository/local"
	//"github.com/lokesh2201013/postgres-data-summary/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.ConnectDB()

     

	app := fiber.New()
    logger.InitLogger()
	//app.Use(logger.New())
    app.Use(logger.ZapLogger())
	routes.SetupRoutes(app)
     app.Get("/swagger/*", fiberSwagger.WrapHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}

	log.Printf("Server listening on %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
