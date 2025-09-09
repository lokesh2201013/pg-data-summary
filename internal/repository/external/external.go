package external

import (
	"fmt"
	"net/http"
	"context"
	"time"
	"encoding/json"
	"bytes"
	"github.com/lokesh2201013/postgres-data-summary/internal/domain"
)

type SummaryClient interface {
	FetchSummary(details domain.ConnectionDetails) (domain.Summary, error)
}

type summaryClient struct{}

func NewSummaryClient() SummaryClient {
	return &summaryClient{}
}

func (c *summaryClient) FetchSummary(details domain.ConnectionDetails) (domain.Summary, error) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
    jsonData, err := json.Marshal(details)
	fmt.Println("Connection details JSON:", string(jsonData))
	req := fmt.Sprintf("http://127.0.0.1:8000/summarypostgres")
    res, err := http.Post(req, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return domain.Summary{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return domain.Summary{}, fmt.Errorf("failed to fetch summary: %s", res.Status)
	}

	var summary domain.Summary
	if err := json.NewDecoder(res.Body).Decode(&summary); err != nil {
		return domain.Summary{}, err
	}
    fmt.Println("Fetched summary:", summary)
	return summary, nil
}
