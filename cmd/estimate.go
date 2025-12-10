package cmd

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/discoverer/noop"
	"github.com/w-h-a/caus/internal/client/estimator"
	"github.com/w-h-a/caus/internal/client/estimator/v1alpha1"
	"github.com/w-h-a/caus/internal/config"
	"github.com/w-h-a/caus/internal/service/orchestrator"
	"google.golang.org/protobuf/encoding/protojson"
)

func Estimate(c *cli.Context) error {
	ctx := c.Context

	// 1. Parse inputs
	configPath := c.String("vars")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	bs, err := os.ReadFile(c.String("graph"))
	if err != nil {
		return fmt.Errorf("failed to read graph: %w", err)
	}
	var graph causal.CausalGraph
	if err := protojson.Unmarshal(bs, &graph); err != nil {
		return fmt.Errorf("invalid graph: %w", err)
	}

	step := c.Duration("step")
	now := time.Now().UTC().Truncate(step)
	start := now.Add(-1 * c.Duration("start"))
	end := now.Add(-1 * c.Duration("end"))

	args := orchestrator.EstimateArgs{
		Graph: &graph,
	}

	log.Printf("Starting Estimation on %d variables...", len(cfg.Variables))
	log.Printf("Window: %s -> %s (Step: %s)", start.Format(time.RFC3339), end.Format(time.RFC3339), step)

	// 2. Build clients
	fetchers, err := initFetchers(cfg)
	if err != nil {
		return err
	}

	noopDiscoverer := noop.NewDiscoverer()

	// TODO: pass in discoverer config and location via cli or expand variable cfg
	v1alpha1Estimator := v1alpha1.NewEstimator(
		estimator.WithLocation("localhost:50051"),
	)

	// 3. Build services
	o := orchestrator.New(fetchers, noopDiscoverer, v1alpha1Estimator)

	// 4. Run Estimate
	results, err := o.Estimate(
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

	// 5. Display results
	return printEstimationResults(results)
}

func printEstimationResults(results *causal.EstimateResponse) error {
	fmt.Printf("\n--- Causal Physics (Discovered Coefficients) ---\n")

	for node, model := range results.Models {
		fmt.Printf("Node: %s\n", node)
		fmt.Printf("  Intercept: %.4f\n", model.Intercept)

		for i, feature := range model.Features {
			coeff := model.Coefficients[i]
			strength := ""
			if math.Abs(float64(coeff)) > 1.0 {
				strength = " (STRONG)"
			}
			fmt.Printf("  -> %s: %.4f%s\n", feature, coeff, strength)
		}

		fmt.Println("")
	}

	return nil
}
