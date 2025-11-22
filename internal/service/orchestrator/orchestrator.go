package orchestrator

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"time"

	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/causal"
	"github.com/w-h-a/caus/internal/client/fetcher"
)

type Service struct {
	metricsFetcher fetcher.Fetcher
	tracesFetcher  fetcher.Fetcher
	discoverer     causal.Discoverer
}

func (s *Service) Do(
	ctx context.Context,
	vars []variable.VariableDefinition,
	start time.Time,
	end time.Time,
	step time.Duration,
	analysisArgs AnalysisArgs,
) (*causal.CausalGraph, error) {
	csvData, err := s.fetch(ctx, vars, start, end, step)
	if err != nil {
		return nil, err
	}

	graph, err := s.discover(ctx, csvData, analysisArgs)
	if err != nil {
		return nil, err
	}

	return graph, nil
}

func (s *Service) fetch(ctx context.Context, vars []variable.VariableDefinition, start time.Time, end time.Time, step time.Duration) ([]byte, error) {
	// 1. Scatter: Fetch all series
	results := make(map[string]map[time.Time]float64)

	for _, v := range vars {
		log.Printf("ORCHESTRATOR: Fetching '%s' from %s...", v.Name, v.Source)

		var dataFetcher fetcher.Fetcher
		switch v.Source {
		case "metrics":
			dataFetcher = s.metricsFetcher
		case "traces":
			dataFetcher = s.tracesFetcher
		default:
			return nil, fmt.Errorf("unknown source type '%s' for variable '%s'", v.Source, v.Name)
		}

		series, err := dataFetcher.Fetch(ctx, v, start, end, step)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch '%s': %w", v.Name, err)
		}

		results[v.Name] = series
	}

	// 2. Gather & Stitch: Align everything to the exact timestamps requested
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// make the header row
	header := make([]string, len(vars))
	for i, v := range vars {
		header[i] = v.Name
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// iterate over steps
	current := start.Truncate(step)
	endTime := end.Truncate(step)

	for !current.After(endTime) {
		// make the rows
		row := make([]string, len(vars))
		for i, v := range vars {
			val, ok := results[v.Name][current]
			if !ok {
				// TODO: Handle missing data
				row[i] = "0.0"
			} else {
				row[i] = strconv.FormatFloat(val, 'f', 6, 64)
			}
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
		// add a step
		current = current.Add(step)
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *Service) discover(ctx context.Context, csvData []byte, analysis AnalysisArgs) (*causal.CausalGraph, error) {
	req := &causal.DiscoverRequest{
		CsvData: string(csvData),
		MaxLag:  analysis.MaxLag,
		PcAlpha: analysis.PcAlpha,
	}

	graph, err := s.discoverer.Discover(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to discover causes: %w", err)
	}

	return graph, nil
}

func New(m fetcher.Fetcher, t fetcher.Fetcher, d causal.Discoverer) *Service {
	return &Service{
		metricsFetcher: m,
		tracesFetcher:  t,
		discoverer:     d,
	}
}
