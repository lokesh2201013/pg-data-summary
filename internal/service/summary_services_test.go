package service

import (
	"errors"
	"testing"
	"time"

	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//
// --- Mock types ---
//

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
	return args.Get(0).(*domain.Summary), args.Error(1)
}

type mockClient struct {
	mock.Mock
}

func (m *mockClient) FetchSummary(details domain.ConnectionDetails) (domain.Summary, error) {
	args := m.Called(details)
	return args.Get(0).(domain.Summary), args.Error(1)
}

//
// --- Tests ---
//

func TestUpdateSummary_Success(t *testing.T) {
	repo := new(mockRepo)
	client := new(mockClient)
	svc := NewSummaryService(repo, client)

	details := domain.ConnectionDetails{Host: "localhost", User: "test", Password: "secret"}
	expectedSummary := domain.Summary{ID: "123"}

	client.On("FetchSummary", details).Return(expectedSummary, nil)

	summary, err := svc.UpdateSummary(details)

	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, "123", summary.ID)
	assert.Equal(t, "", summary.SourceInfo.Password) // password should be cleared
	assert.WithinDuration(t, time.Now(), summary.SyncedAt, time.Second)
}

func TestUpdateSummary_FetchError(t *testing.T) {
	repo := new(mockRepo)
	client := new(mockClient)
	svc := NewSummaryService(repo, client)

	details := domain.ConnectionDetails{Host: "localhost"}
	client.On("FetchSummary", details).Return(domain.Summary{}, errors.New("fetch failed"))

	summary, err := svc.UpdateSummary(details)

	assert.Error(t, err)
	assert.Nil(t, summary)
}

func TestGetSummaries(t *testing.T) {
	repo := new(mockRepo)
	client := new(mockClient)
	svc := NewSummaryService(repo, client)

	expected := []domain.Summary{{ID: "1"}, {ID: "2"}}
	repo.On("GetSummaries", 1, 10).Return(expected, nil)

	summaries, err := svc.GetSummaries(1, 10)

	assert.NoError(t, err)
	assert.Len(t, summaries, 2)
	assert.Equal(t, "1", summaries[0].ID)
}

func TestGetSummaryByID(t *testing.T) {
	repo := new(mockRepo)
	client := new(mockClient)
	svc := NewSummaryService(repo, client)

	expected := &domain.Summary{ID: "abc"}
	repo.On("GetSummaryByID", "abc").Return(expected, nil)

	summary, err := svc.GetSummaryByID("abc")

	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, "abc", summary.ID)
}

func TestGetSummaryByID_NotFound(t *testing.T) {
	repo := new(mockRepo)
	client := new(mockClient)
	svc := NewSummaryService(repo, client)

	repo.On("GetSummaryByID", "missing").Return((*domain.Summary)(nil), errors.New("not found"))

	summary, err := svc.GetSummaryByID("missing")

	assert.Error(t, err)
	assert.Nil(t, summary)
}
