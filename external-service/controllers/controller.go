package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lokesh2201013/postgres-data-summary/external-service/models"
)
func GetSummaryPostgres(c *fiber.Ctx) error {
    var connDetails models.ConnectionDetails
    if err := c.BodyParser(&connDetails); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }
     fmt.Println("Received connection details:", connDetails)
    summary := models.Summary{
        ID:         "sum123",
        SourceInfo: connDetails,
      //  SyncedAt:   time.Now(),
        Schemas: []models.Schema{
            {
                Name: "public",
                Tables: []models.Table{
                    {Name: "users", RowCount: 0, SizeMB: 12.5},
                    {Name: "orders", RowCount: 0, SizeMB: 8.3},
                },
            },
            {
                Name: "sales",
                Tables: []models.Table{
                    {Name: "invoices", RowCount: 300, SizeMB: 150},
                    {Name: "customers", RowCount: 1500, SizeMB: 100},
                },
            },
        },
    }
    fmt.Println("Returning summary:", summary)
    return c.JSON(summary)
}
