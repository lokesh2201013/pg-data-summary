package service

import (
	"time"
    "fmt"
	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
	"github.com/lokesh2201013/postgres-data-summary/internal/repository/external"
	"github.com/lokesh2201013/postgres-data-summary/internal/repository/local"
)

type ISummaryService interface {
	UpdateSummary(details domain.ConnectionDetails) (*domain.Summary, error)
	GetSummaries(page, pageSize int) ([]domain.Summary, error)
	GetSummaryByID(id string) (*domain.Summary, error)
}

type SummaryService struct {
	repo     local.SummaryRepository
	exclient external.SummaryClient
}

func NewSummaryService(repo local.SummaryRepository, client external.SummaryClient) *SummaryService {
	return &SummaryService{repo: repo, exclient: client}
}

func (s *SummaryService) UpdateSummary(details domain.ConnectionDetails) (*domain.Summary, error) {
	fmt.Println("Updating summary for:", details)
	summary, err := s.exclient.FetchSummary(details)
	if err != nil {
		return nil, err
	}

	details.Password = ""
	summary.SourceInfo = details
	summary.SyncedAt = time.Now()

	return &summary, nil
}

func (s *SummaryService) GetSummaries(page, pageSize int) ([]domain.Summary, error) {
	return s.repo.GetSummaries(page, pageSize)
}

func (s *SummaryService) GetSummaryByID(id string) (*domain.Summary, error) {
	return s.repo.GetSummaryByID(id)
}
