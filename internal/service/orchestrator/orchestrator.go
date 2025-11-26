package orchestrator

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"time"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/discoverer"
	"github.com/w-h-a/caus/internal/client/fetcher"
)

type Service struct {
	fetchers   map[string]map[string]fetcher.Fetcher
	discoverer discoverer.Discoverer
}

func (s *Service) Discover(
	ctx context.Context,
	vars []variable.VariableDefinition,
	start time.Time,
	end time.Time,
	step time.Duration,
	analysisArgs DiscoveryArgs,
) (*causal.CausalGraph, error) {
	// 1. fetch and stitch
	csvData, err := s.fetch(ctx, vars, start, end, step)
	if err != nil {
		return nil, err
	}

	// 2. discover direct causes
	graph, err := s.discover(ctx, csvData, analysisArgs)
	if err != nil {
		return nil, err
	}

	return graph, nil
}

func (s *Service) fetch(ctx context.Context, vars []variable.VariableDefinition, start time.Time, end time.Time, step time.Duration) ([]byte, error) {
	start = start.UTC().Truncate(step).Truncate(0)
	end = end.UTC().Truncate(step).Truncate(0)

	results := make(map[string]map[time.Time]float64)

	// 1. scatter
	for _, v := range vars {
		log.Printf("ORCHESTRATOR: Fetching '%s' from %s...", v.Name, v.Source.Loc)

		var dataFetcher fetcher.Fetcher

		impls, ok := s.fetchers[v.Source.Type]
		if !ok {
			return nil, fmt.Errorf("unknown source type '%s' for variable '%s'", v.Source.Type, v.Name)
		}

		dataFetcher, ok = impls[v.Source.Impl]
		if !ok {
			return nil, fmt.Errorf("unknown %s implementation '%s' for variable '%s'", v.Source.Type, v.Source.Impl, v.Name)
		}

		series, err := dataFetcher.Fetch(ctx, v, start, end, step)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch '%s': %w", v.Name, err)
		}

		results[v.Name] = series
	}

	// 2. stitch
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
	bs := buf.Bytes()

	return bs, nil
}

func (s *Service) discover(ctx context.Context, csvData []byte, analysis DiscoveryArgs) (*causal.CausalGraph, error) {
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

func New(fs map[string]map[string]fetcher.Fetcher, d discoverer.Discoverer) *Service {
	return &Service{
		fetchers:   fs,
		discoverer: d,
	}
}
