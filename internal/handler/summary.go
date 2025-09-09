package handler

import (
	//"net/http"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	//"github.com/lokesh2201013/postgres-data-summary/internal/repository/external"
	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
	"github.com/lokesh2201013/postgres-data-summary/internal/logger"
	//"github.com/lokesh2201013/postgres-data-summary/internal/repository/local"
	"github.com/lokesh2201013/postgres-data-summary/internal/service"
)


type SummaryHandler interface {
    SyncSummary(c *fiber.Ctx) error
    GetSummaries(c *fiber.Ctx) error
    GetSummaryByID(c *fiber.Ctx) error
}

type summaryHandlerImpl struct {
    service service.ISummaryService
}

func NewSummaryHandler(service service.ISummaryService) SummaryHandler {
    return &summaryHandlerImpl{service: service}
}


// SyncSummary godoc
// @Summary Sync a new database summary
// @Description Connects to remote PostgreSQL via external API and saves summary
// @Tags summary
// @Accept  json
// @Produce  json
// @Param details body domain.ConnectionDetails true "Remote DB connection"
// @Success 201 {object} domain.Summary
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /summary/sync [post]
func (h *summaryHandlerImpl) SyncSummary(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Log.Info("SyncSummary request received")

	var details domain.ConnectionDetails
	if err := c.BodyParser(&details); err != nil {
		logger.Log.Error("Failed to parse ConnectionDetails",
			zap.Error(err),
			zap.ByteString("body", c.Body()))
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if details.Host == "" || details.Port == nil || details.User == "" || details.DBName == "" {
		logger.Log.Warn("Missing required connection details",
			zap.Any("details", details))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required connection details"})
	}

	summary, err := h.service.UpdateSummary(details)
	if err != nil {
		logger.Log.Error("UpdateSummary failed",
			zap.Any("details", details),
			zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to sync summary")
	}

	// x := local.NewSummaryRepository()
	// if err := x.SaveSummary(summary); err != nil {
	// 	logger.Log.Error("SaveSummary failed",
	// 		zap.Any("summary", summary),
	// 		zap.Error(err))
	// 	return fiber.NewError(fiber.StatusInternalServerError, "Failed to save summary")
	// }

	logger.Log.Info("Summary synced successfully", zap.String("id", summary.ID))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Summary synced successfully",
		"summary": summary,
	})
}


// GetSummaries godoc
// @Summary Get all summaries
// @Description Retrieves paginated summaries
// @Tags summary
// @Produce  json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {array} domain.Summary
// @Failure 500 {object} map[string]string
// @Router /summary/summaries [get]
func (h *summaryHandlerImpl) GetSummaries(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil {
		pageSize = 10
	}

	logger.Log.Info("GetSummaries request received", zap.Int("page", page), zap.Int("pageSize", pageSize))

	summaries, err := h.service.GetSummaries(page, pageSize)
	if err != nil {
		logger.Log.Error("GetSummaries failed", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get summaries")
	}

	logger.Log.Info("GetSummaries succeeded", zap.Int("count", len(summaries)))
	return c.JSON(summaries)
}


// GetSummaryByID godoc
// @Summary Get summary by ID
// @Description Retrieves full summary by ID
// @Tags summary
// @Produce  json
// @Param id path string true "Summary ID"
// @Success 200 {object} domain.Summary
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /summary/summaries/{id} [get]
func (h *summaryHandlerImpl) GetSummaryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	logger.Log.Info("GetSummaryByID request received", zap.String("id", id))

	summary, err := h.service.GetSummaryByID(id)
	if err != nil {
		logger.Log.Error("GetSummaryByID failed", zap.String("id", id), zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get summary")
	}

	if summary == nil {
		logger.Log.Warn("Summary not found", zap.String("id", id))
		return fiber.NewError(fiber.StatusNotFound, "Summary not found")
	}

	logger.Log.Info("GetSummaryByID succeeded", zap.String("id", summary.ID))
	return c.JSON(summary)
}