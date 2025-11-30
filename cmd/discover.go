package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/urfave/cli/v2"
	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/discoverer"
	"github.com/w-h-a/caus/internal/client/discoverer/v1alpha1"
	"github.com/w-h-a/caus/internal/client/simulator/noop"
	"github.com/w-h-a/caus/internal/config"
	"github.com/w-h-a/caus/internal/service/orchestrator"
	"google.golang.org/protobuf/encoding/protojson"
)

func Discover(c *cli.Context) error {
	ctx := c.Context

	// 1. Parse inputs
	configPath := c.String("vars")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	step := c.Duration("step")
	now := time.Now().UTC().Truncate(step)
	start := now.Add(-1 * c.Duration("start"))
	end := now.Add(-1 * c.Duration("end"))

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
	v1alpha1Discoverer := v1alpha1.NewDiscoverer(
		discoverer.WithLocation("localhost:50051"),
	)

	noopSimulator := noop.NewSimulator()

	// 3. Build services
	o := orchestrator.New(fetchers, v1alpha1Discoverer, noopSimulator)

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

	// 5. Print graph (json or pretty)
	if c.Bool("json") {
		opts := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			EmitUnpopulated: true,
		}
		bs, _ := opts.Marshal(graph)
		fmt.Println(string(bs))
	} else {
		printGraph(graph, step)
	}

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
