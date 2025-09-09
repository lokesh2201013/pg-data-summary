package service

import (
	"errors"
	"testing"
	//"time"

	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)
func init() {

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) SaveSummary(summary *domain.Summary) error {
	args := m.Called(summary)
	return args.Error(0)
}

func (m *mockRepo) GetSummaries(page, pageSize int) ([]domain.Summary, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]domain.Summary), args.Error(1)
}

func (m *mockRepo) GetSummaryByID(id string) (*domain.Summary, error) {
	args := m.Called(id)
	if s := args.Get(0); s != nil {
		return s.(*domain.Summary), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- Mock External Client ---
type mockExternalClient struct {
	mock.Mock
}

func (m *mockExternalClient) FetchSummary(details domain.ConnectionDetails) (domain.Summary, error) {
	args := m.Called(details)
	if s := args.Get(0); s != nil {
		return s.(domain.Summary), args.Error(1)
	}
	return domain.Summary{}, args.Error(1)
}

// --- Tests ---
func TestUpdateSummary_Success(t *testing.T) {
	repo := new(mockRepo)
	client := new(mockExternalClient)
	service := NewSummaryService(repo, client, 1, 0)

	port := 5432
	details := domain.ConnectionDetails{
		Host:   "localhost",
		Port:   &port,
		User:   "test",
		DBName: "demo",
	}

	expectedSummary := domain.Summary{ID: "123", Name: "Test Summary"}

	client.On("FetchSummary", details).Return(expectedSummary, nil)
	repo.On("SaveSummary", mock.AnythingOfType("*domain.Summary")).Return(nil)

	summary, err := service.UpdateSummary(details)
	assert.NoError(t, err)
	assert.Equal(t, "123", summary.ID)
	client.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestUpdateSummary_FetchError(t *testing.T) {
	repo := new(mockRepo)
	client := new(mockExternalClient)
	service := NewSummaryService(repo, client, 2, 0)

	details := domain.ConnectionDetails{Host: "badhost"}
	client.On("FetchSummary", details).Return(domain.Summary{}, errors.New("fetch failed"))

	summary, err := service.UpdateSummary(details)
	assert.Nil(t, summary)
	assert.Error(t, err)
}

func TestUpdateSummary_SaveError(t *testing.T) {
	repo := new(mockRepo)
	client := new(mockExternalClient)
	service := NewSummaryService(repo, client, 2, 0)

	details := domain.ConnectionDetails{Host: "localhost"}
	expectedSummary := domain.Summary{ID: "999"}

	client.On("FetchSummary", details).Return(expectedSummary, nil)
	repo.On("SaveSummary", mock.AnythingOfType("*domain.Summary")).Return(errors.New("db error"))

	summary, err := service.UpdateSummary(details)
	assert.Nil(t, summary)
	assert.Error(t, err)
}

func TestGetSummaries_Success(t *testing.T) {
	repo := new(mockRepo)
	service := NewSummaryService(repo, nil, 1, 0)

	summaries := []domain.Summary{
		{ID: "1", Name: "Summary1"},
		{ID: "2", Name: "Summary2"},
	}

	repo.On("GetSummaries", 1, 10).Return(summaries, nil)

	result, err := service.GetSummaries(1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetSummaryByID_Success(t *testing.T) {
	repo := new(mockRepo)
	service := NewSummaryService(repo, nil, 1, 0)

	expected := &domain.Summary{ID: "123", Name: "Demo Summary"}
	repo.On("GetSummaryByID", "123").Return(expected, nil)

	summary, err := service.GetSummaryByID("123")
	assert.NoError(t, err)
	assert.Equal(t, "123", summary.ID)
}

func TestGetSummaryByID_NotFound(t *testing.T) {
	repo := new(mockRepo)
	service := NewSummaryService(repo, nil, 1, 0)

	repo.On("GetSummaryByID", "999").Return(nil, nil)

	summary, err := service.GetSummaryByID("999")
	assert.NoError(t, err)
	assert.Nil(t, summary)
}
