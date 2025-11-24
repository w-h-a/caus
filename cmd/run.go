package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/urfave/cli/v2"
	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/discoverer/grpc"
	"github.com/w-h-a/caus/internal/client/fetcher"
	"github.com/w-h-a/caus/internal/client/fetcher/mock"
	"github.com/w-h-a/caus/internal/config"
	"github.com/w-h-a/caus/internal/service/orchestrator"
)

func Run(c *cli.Context) error {
	ctx := c.Context

	// 1. Parse vars.yml
	configPath := c.String("vars")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	start, _ := time.Parse(time.RFC3339, c.String("start"))
	end, _ := time.Parse(time.RFC3339, c.String("end"))
	step := c.Duration("step")

	args := orchestrator.AnalysisArgs{
		MaxLag:  int32(c.Int("lag")),
		PcAlpha: float32(c.Float64("alpha")),
	}

	log.Printf("Starting Discovery on %d variables...", len(cfg.Variables))
	log.Printf("Window: %s -> %s (Step: %s)", start.Format(time.TimeOnly), end.Format(time.TimeOnly), step)

	// 2. Build clients
	// TODO: choose based on config
	mockFetcher := mock.NewFetcher()

	fetchers := map[string]map[string]fetcher.Fetcher{
		"metrics": {
			"mock":       mockFetcher,
			"prometheus": mockFetcher,
		},
		"traces": {
			"mock":       mockFetcher,
			"clickhouse": mockFetcher,
		},
	}

	grpcDiscoverer := grpc.NewDiscoverer()

	// 3. Build services
	o := orchestrator.New(fetchers, grpcDiscoverer)

	// 4. Do it
	graph, err := o.Do(
		ctx,
		cfg.Variables,
		start,
		end,
		step,
		args,
	)
	if err != nil {
		return err
	}

	// 5. Print the graph
	printGraph(graph, step)

	return nil
}

func printGraph(graph *causal.CausalGraph, step time.Duration) {
	fmt.Println("\n--- Causal Graph Results ---")
	fmt.Println("Nodes:")
	for _, node := range graph.Nodes {
		fmt.Printf("  - %s\n", node.Label)
	}
	fmt.Println("\nDiscovered Edges:")
	if len(graph.Edges) == 0 {
		fmt.Println("  No causal edges were found.")
	} else {
		for _, edge := range graph.Edges {
			lagTime := time.Duration(edge.Lag) * step
			fmt.Printf("  - %s --> %s (lag: %d = %s)\n", edge.Source, edge.Target, edge.Lag, lagTime)
		}
	}
	fmt.Println("--------------------------")
}
