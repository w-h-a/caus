package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/w-h-a/caus/internal/client/causal"
	"github.com/w-h-a/caus/internal/client/metrics"
)

type Service struct {
	metricsFetcher metrics.Fetcher
	discoverer     causal.Discoverer
}

func (s *Service) RunAnalysis(ctx context.Context, metrics []string, start time.Time, end time.Time) (*causal.CausalGraph, error) {
	csvData, err := s.metricsFetcher.Fetch(ctx, metrics, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}

	req := &causal.DiscoverRequest{
		CsvData: string(csvData),
		MaxLag:  1,
		PcAlpha: 0,
	}

	graph, err := s.discoverer.Discover(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to discover causes: %w", err)
	}

	return graph, nil
}

func New(f metrics.Fetcher, d causal.Discoverer) *Service {
	return &Service{
		metricsFetcher: f,
		discoverer:     d,
	}
}
