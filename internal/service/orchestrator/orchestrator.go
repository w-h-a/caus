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
	"github.com/w-h-a/caus/internal/client/simulator"
)

type Service struct {
	fetchers   map[string]map[string]fetcher.Fetcher
	discoverer discoverer.Discoverer
	simulator  simulator.Simulator
}

func (s *Service) Discover(
	ctx context.Context,
	vars []variable.VariableDefinition,
	start time.Time,
	end time.Time,
	step time.Duration,
	discoveryArgs DiscoveryArgs,
) (*causal.CausalGraph, error) {
	// 1. fetch and stitch
	csvData, err := s.fetch(ctx, vars, start, end, step)
	if err != nil {
		return nil, err
	}

	// 2. discover direct causes
	graph, err := s.discover(ctx, csvData, discoveryArgs)
	if err != nil {
		return nil, err
	}

	return graph, nil
}

func (s *Service) Simulate(
	ctx context.Context,
	vars []variable.VariableDefinition,
	start time.Time,
	end time.Time,
	step time.Duration,
	simulationArgs SimulationArgs,
) (string, error) {
	// 1. fetch and stitch
	csvData, err := s.fetch(ctx, vars, start, end, step)
	if err != nil {
		return "", err
	}

	// 2. do counterfactual prediction
	result, err := s.simulate(ctx, csvData, simulationArgs)
	if err != nil {
		return "", err
	}

	return result, nil
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

	// keep track of the last known for each column
	lastKnown := make([]string, len(vars))
	for i := range lastKnown {
		lastKnown[i] = "0.0"
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
				zeroOK := v.Source.Type == "traces" && v.TraceQuery != nil && v.TraceQuery.Dimension == "calls"
				if zeroOK {
					row[i] = "0.0"
					lastKnown[i] = "0.0"
				} else {
					row[i] = lastKnown[i]
				}
			} else {
				val := strconv.FormatFloat(val, 'f', 6, 64)
				row[i] = val
				lastKnown[i] = val
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

func (s *Service) discover(ctx context.Context, csvData []byte, discovery DiscoveryArgs) (*causal.CausalGraph, error) {
	req := &causal.DiscoverRequest{
		CsvData: string(csvData),
		MaxLag:  discovery.MaxLag,
		PcAlpha: discovery.PcAlpha,
	}

	graph, err := s.discoverer.Discover(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to discover causes: %w", err)
	}

	return graph, nil
}

func (s *Service) simulate(ctx context.Context, csvData []byte, simulation SimulationArgs) (string, error) {
	req := &causal.SimulateRequest{
		CsvData:         string(csvData),
		Graph:           simulation.Graph,
		Intervention:    simulation.Intervention,
		SimulationSteps: simulation.Horizon,
	}

	rsp, err := s.simulator.Simulate(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to perform counterfactual simulation: %w", err)
	}

	return rsp.JsonResults, nil
}

func New(fs map[string]map[string]fetcher.Fetcher, d discoverer.Discoverer, s simulator.Simulator) *Service {
	return &Service{
		fetchers:   fs,
		discoverer: d,
		simulator:  s,
	}
}
