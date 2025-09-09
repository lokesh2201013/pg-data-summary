package service

import (
	//"context"
	//"fmt"
	"time"

	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
	"github.com/lokesh2201013/postgres-data-summary/internal/logger"
	"github.com/lokesh2201013/postgres-data-summary/internal/repository/external"
	"github.com/lokesh2201013/postgres-data-summary/internal/repository/local"
	"go.uber.org/zap"
)

type ISummaryService interface {
	UpdateSummary(details domain.ConnectionDetails) (*domain.Summary, error)
	GetSummaries(page, pageSize int) ([]domain.Summary, error)
	GetSummaryByID(id string) (*domain.Summary, error)
}

type SummaryService struct {
	repo     local.SummaryRepository
	exclient external.SummaryClient
	retries  int
	delay    time.Duration
}

func NewSummaryService(repo local.SummaryRepository, client external.SummaryClient, retries int, delay time.Duration) *SummaryService {
	return &SummaryService{
		repo:     repo,
		exclient: client,
		retries:  retries,
		delay:    delay,
	}
}

// UpdateSummary fetches summary from external DB with retries and logs
func (s *SummaryService) UpdateSummary(details domain.ConnectionDetails) (*domain.Summary, error) {
	logger.Log.Info("Starting UpdateSummary", zap.Any("details", details))

	var summary domain.Summary
	var err error

	for attempt := 1; attempt <= s.retries; attempt++ {
		summary, err = s.exclient.FetchSummary(details)
		if err == nil {
			break
		}
		logger.Log.Warn("FetchSummary attempt failed",
			zap.Int("attempt", attempt),
			zap.Any("details", details),
			zap.Error(err),
		)
		time.Sleep(s.delay)
	}

	if err != nil {
		logger.Log.Error("FetchSummary failed after retries", zap.Error(err))
		return nil, err
	}

	// Mask password for logging
	details.Password = ""
	summary.SourceInfo = details
	summary.SyncedAt = time.Now()

	logger.Log.Info("Fetched summary successfully", zap.String("summaryID", summary.ID))

	// Retry DB save with logging
	for attempt := 1; attempt <= s.retries; attempt++ {
		err = s.repo.SaveSummary(&summary)
		if err == nil {
			logger.Log.Info("Saved summary successfully", zap.String("summaryID", summary.ID))
			break
		}
		logger.Log.Warn("SaveSummary attempt failed",
			zap.Int("attempt", attempt),
			zap.String("summaryID", summary.ID),
			zap.Error(err),
		)
		time.Sleep(s.delay)
	}

	if err != nil {
		logger.Log.Error("SaveSummary failed after retries", zap.String("summaryID", summary.ID), zap.Error(err))
		return nil, err
	}

	return &summary, nil
}

func (s *SummaryService) GetSummaries(page, pageSize int) ([]domain.Summary, error) {
	logger.Log.Info("GetSummaries called", zap.Int("page", page), zap.Int("pageSize", pageSize))
	summaries, err := s.repo.GetSummaries(page, pageSize)
	if err != nil {
		logger.Log.Error("GetSummaries failed", zap.Error(err))
		return nil, err
	}
	logger.Log.Info("GetSummaries succeeded", zap.Int("count", len(summaries)))
	return summaries, nil
}

func (s *SummaryService) GetSummaryByID(id string) (*domain.Summary, error) {
	logger.Log.Info("GetSummaryByID called", zap.String("id", id))
	summary, err := s.repo.GetSummaryByID(id)
	if err != nil {
		logger.Log.Error("GetSummaryByID failed", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	if summary != nil {
		logger.Log.Info("GetSummaryByID succeeded", zap.String("id", summary.ID))
	} else {
		logger.Log.Warn("Summary not found", zap.String("id", id))
	}
	return summary, nil
}
