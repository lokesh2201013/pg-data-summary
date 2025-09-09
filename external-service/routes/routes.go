package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lokesh2201013/postgres-data-summary/external-service/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/summarypostgres",controllers.GetSummaryPostgres)
}