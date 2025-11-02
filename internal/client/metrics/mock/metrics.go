package mock

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/w-h-a/caus/internal/client/metrics"
)

type mockFetcher struct {
}

func (f *mockFetcher) Fetch(ctx context.Context, metrics []string, start time.Time, end time.Time) ([]byte, error) {
	csvData, err := os.ReadFile("test/testdata/sample_data.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return csvData, nil
}

// TODO: options
func NewFetcher() metrics.Fetcher {
	f := &mockFetcher{}

	return f
}
