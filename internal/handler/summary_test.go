package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
	"github.com/lokesh2201013/postgres-data-summary/internal/handler"
)

// ---- Mock Service ----
type mockSummaryService struct {
	mock.Mock
}

func (m *mockSummaryService) UpdateSummary(details domain.ConnectionDetails) (*domain.Summary, error) {
	args := m.Called(details)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Summary), args.Error(1)
}

func (m *mockSummaryService) GetSummaries(page, pageSize int) ([]domain.Summary, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]domain.Summary), args.Error(1)
}

func (m *mockSummaryService) GetSummaryByID(id string) (*domain.Summary, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Summary), args.Error(1)
}

// ---- Tests ----

func setupApp(svc *mockSummaryService) *fiber.App {
	h := handler.NewSummaryHandler(svc)
	app := fiber.New()
	api := app.Group("/summary")
	api.Post("/sync", h.SyncSummary)
	api.Get("/summaries", h.GetSummaries)
	api.Get("/summaries/:id", h.GetSummaryByID)
	return app
}

func TestSyncSummary_Success(t *testing.T) {
	svc := new(mockSummaryService)
	app := setupApp(svc)

	details := domain.ConnectionDetails{
		Host: "localhost", Port: new(int), User: "test", DBName: "demo",
	}
	*details.Port = 5432

	expected := &domain.Summary{
		ID:       "123",
		SyncedAt: time.Now(),
		SourceInfo: details,
	}

	svc.On("UpdateSummary", details).Return(expected, nil)

	body, _ := json.Marshal(details)
	req := httptest.NewRequest(http.MethodPost, "/summary/sync", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestSyncSummary_InvalidBody(t *testing.T) {
	svc := new(mockSummaryService)
	app := setupApp(svc)

	req := httptest.NewRequest(http.MethodPost, "/summary/sync", bytes.NewReader([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetSummaries_Success(t *testing.T) {
	svc := new(mockSummaryService)
	app := setupApp(svc)

	expected := []domain.Summary{
		{ID: "1", SyncedAt: time.Now()},
		{ID: "2", SyncedAt: time.Now()},
	}

	svc.On("GetSummaries", 1, 10).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/summary/summaries?page=1&pageSize=10", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetSummaryByID_NotFound(t *testing.T) {
	svc := new(mockSummaryService)
	app := setupApp(svc)

	svc.On("GetSummaryByID", "999").Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/summary/summaries/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetSummaryByID_Error(t *testing.T) {
	svc := new(mockSummaryService)
	app := setupApp(svc)

	svc.On("GetSummaryByID", "err").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/summary/summaries/err", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
