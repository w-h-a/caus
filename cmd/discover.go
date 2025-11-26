package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/urfave/cli/v2"
	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/discoverer"
	"github.com/w-h-a/caus/internal/client/discoverer/grpc"
	"github.com/w-h-a/caus/internal/client/fetcher"
	"github.com/w-h-a/caus/internal/client/fetcher/clickhouse"
	"github.com/w-h-a/caus/internal/client/fetcher/prometheus"
	"github.com/w-h-a/caus/internal/client/fetcher/random"
	"github.com/w-h-a/caus/internal/config"
	"github.com/w-h-a/caus/internal/service/orchestrator"
)

func Discover(c *cli.Context) error {
	ctx := c.Context

	// 1. Parse vars.yml
	configPath := c.String("vars")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	now := time.Now().UTC()
	start := now.Add(-1 * c.Duration("start"))
	end := now.Add(-1 * c.Duration("end"))
	step := c.Duration("step")

	args := orchestrator.DiscoveryArgs{
		MaxLag:  int32(c.Int("lag")),
		PcAlpha: float32(c.Float64("alpha")),
	}

	log.Printf("Starting Discovery on %d variables...", len(cfg.Variables))
	log.Printf("Window: %s -> %s (Step: %s)", start.Format(time.RFC3339), end.Format(time.RFC3339), step)

	// 2. Build clients
	fetchers, err := initFetchers(cfg)
	if err != nil {
		return err
	}

	// TODO: pass in discoverer config and location via cli or expand variable cfg
	grpcDiscoverer := grpc.NewDiscoverer(
		discoverer.WithLocation("localhost:50051"),
	)

	// 3. Build services
	o := orchestrator.New(fetchers, grpcDiscoverer)

	// 4. Run Discover
	graph, err := o.Discover(
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

func initFetchers(cfg *variable.DiscoveryConfig) (map[string]map[string]fetcher.Fetcher, error) {
	fetchers := map[string]map[string]fetcher.Fetcher{
		"metrics": {},
		"traces":  {},
	}

	factories := map[string]map[string]func(string) fetcher.Fetcher{
		"metrics": {
			"random":     func(_ string) fetcher.Fetcher { return random.NewFetcher() },
			"prometheus": func(loc string) fetcher.Fetcher { return prometheus.NewFetcher(fetcher.WithLocation(loc)) },
		},
		"traces": {
			"random":     func(_ string) fetcher.Fetcher { return random.NewFetcher() },
			"clickhouse": func(loc string) fetcher.Fetcher { return clickhouse.NewFetcher(fetcher.WithLocation(loc)) },
		},
	}

	for _, v := range cfg.Variables {
		impls, typeOk := factories[v.Source.Type]
		if !typeOk {
			return nil, fmt.Errorf("unsupported source type: %s", v.Source.Type)
		}

		factory, implOk := impls[v.Source.Impl]
		if !implOk {
			return nil, fmt.Errorf("unsupported implementation '%s' for type '%s'", v.Source.Impl, v.Source.Type)
		}

		if _, exists := fetchers[v.Source.Type][v.Source.Impl]; exists {
			continue
		}

		fetchers[v.Source.Type][v.Source.Impl] = factory(v.Source.Loc)
	}

	return fetchers, nil
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
