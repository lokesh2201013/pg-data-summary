package local

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
	"gorm.io/gorm"
)

type SummaryRepository interface {
	SaveSummary(summary *domain.Summary) error
	GetSummaries(page, pageSize int) ([]domain.Summary, error)
	GetSummaryByID(id string) (*domain.Summary, error)
}

type summaryRepo struct{}

func NewSummaryRepository() SummaryRepository {
	return &summaryRepo{}
}


func (r *summaryRepo) SaveSummary(summary *domain.Summary) error {
	for i := range summary.Schemas {
		if summary.Schemas[i].ID == "" {
			summary.Schemas[i].ID = uuid.NewString()
		}
		summary.Schemas[i].SummaryID = summary.ID

		for j := range summary.Schemas[i].Tables {
			if summary.Schemas[i].Tables[j].ID == "" {
				summary.Schemas[i].Tables[j].ID = uuid.NewString()
			}
			summary.Schemas[i].Tables[j].SchemaID = summary.Schemas[i].ID
		}
	}

	var existing domain.Summary
	err := dB.Preload("Schemas.Tables").First(&existing, "id = ?", summary.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dB.Create(summary).Error
		}
		return err
	}

	return dB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(summary).Error
}


func (r *summaryRepo) GetSummaries(page, pageSize int) ([]domain.Summary, error) {
    var summaries []domain.Summary
    offset := (page - 1) * pageSize

    
    err := dB.
        Preload("Schemas.Tables").
        Limit(pageSize).
        Offset(offset).
        Find(&summaries).Error

    if err != nil {
        return nil, err
    }

    return summaries, nil
}


func (r *summaryRepo) GetSummaryByID(id string) (*domain.Summary, error) {
    var summary domain.Summary
    if err := dB.Preload("Schemas.Tables").First(&summary, "id = ?", id).Error; err != nil {
        return nil, err
    }
	fmt.Println("Retrieved summary:", summary)
    return &summary, nil
}

