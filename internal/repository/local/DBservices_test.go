package local

import (
	"fmt"
	"os"
	"testing"

	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func setupTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "postgres_test"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test postgres: %v", err)
	}


	err = db.Migrator().DropTable(&domain.Summary{})
	if err != nil {
		t.Fatalf("failed to drop schema: %v", err)
	}

	err = db.AutoMigrate(&domain.Summary{})
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	
	dB = db

	return db
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func TestSaveSummary_NewRecord(t *testing.T) {
	setupTestDB(t)
	repo := NewSummaryRepository()

	summary := &domain.Summary{
		ID: "123",
		SourceInfo: domain.ConnectionDetails{
			Host: "localhost", User: "test",
		},
	}

	err := repo.SaveSummary(summary)
	assert.NoError(t, err)

	var found domain.Summary
	err = dB.First(&found, "id = ?", "123").Error
	assert.NoError(t, err)
	assert.Equal(t, "123", found.ID)
}

func TestSaveSummary_UpdateRecord(t *testing.T) {
	setupTestDB(t)
	repo := NewSummaryRepository()

	summary := &domain.Summary{ID: "123", SourceInfo: domain.ConnectionDetails{Host: "localhost"}}
	_ = dB.Create(summary).Error

	summary.SourceInfo.User = "updatedUser"
	err := repo.SaveSummary(summary)
	assert.NoError(t, err)

	var found domain.Summary
	_ = dB.First(&found, "id = ?", "123").Error
	assert.Equal(t, "updatedUser", found.SourceInfo.User)
}

func TestGetSummaries(t *testing.T) {
	setupTestDB(t)
	repo := NewSummaryRepository()

	for i := 1; i <= 5; i++ {
		s := domain.Summary{ID: fmt.Sprintf("id-%d", i)}
		_ = dB.Create(&s).Error
	}

	summaries, err := repo.GetSummaries(1, 2)
	assert.NoError(t, err)
	assert.Len(t, summaries, 2)
}

func TestGetSummaryByID(t *testing.T) {
	setupTestDB(t)
	repo := NewSummaryRepository()

	s := domain.Summary{ID: "abc"}
	_ = dB.Create(&s).Error

	found, err := repo.GetSummaryByID("abc")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "abc", found.ID)

	_, err = repo.GetSummaryByID("notfound")
	assert.Error(t, err)
}
