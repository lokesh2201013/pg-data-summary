package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lokesh2201013/postgres-data-summary/internal/handler"
)

func SummaryRoutes(app *fiber.App, h handler.SummaryHandler) {
	api := app.Group("/summary")
	api.Post("/sync", h.SyncSummary)
	api.Get("/summaries", h.GetSummaries)
	api.Get("/summaries/:id", h.GetSummaryByID)
}
