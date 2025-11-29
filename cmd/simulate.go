package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/discoverer/noop"
	"github.com/w-h-a/caus/internal/client/simulator"
	"github.com/w-h-a/caus/internal/client/simulator/v1alpha1"
	"github.com/w-h-a/caus/internal/config"
	"github.com/w-h-a/caus/internal/service/orchestrator"
	"google.golang.org/protobuf/encoding/protojson"
)

func Simulate(c *cli.Context) error {
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

	intervention, err := parseIntervention(c.String("do"))
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	start := now.Add(-1 * c.Duration("start"))
	end := now.Add(-1 * c.Duration("end"))
	step := c.Duration("step")

	args := orchestrator.SimulationArgs{
		Graph:        &graph,
		Intervention: intervention,
		Horizon:      int32(c.Int("horizon")),
	}

	// 2. Build clients
	fetchers, err := initFetchers(cfg)
	if err != nil {
		return err
	}

	noopDiscoverer := noop.NewDiscoverer()

	// TODO: pass in discoverer config and location via cli or expand variable cfg
	v1alpha1Simulator := v1alpha1.NewSimulator(
		simulator.WithLocation("localhost:50051"),
	)

	// 3. Build services
	o := orchestrator.New(fetchers, noopDiscoverer, v1alpha1Simulator)

	// 4. Run Simulate
	json, err := o.Simulate(
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
	if len(c.String("effect")) == 0 {
		fmt.Println(json)
		return nil
	}

	return printEffectSpecificResults(json, intervention, c.String("effect"))
}

func parseIntervention(input string) (*causal.Intervention, error) {
	// looking for service[_.-]metric * 1.2 or service[_.-]metric = 500
	re := regexp.MustCompile(`^([a-zA-Z0-9_.-]+)\s*(\*|=)\s*([0-9\.]+)$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(input))
	// matches is:
	// [
	//   "front_calls * 1.2",  // Index 0: The whole match
	//   "front_calls",        // Index 1: The Variable Name
	//   "*",                  // Index 2: The Operator
	//   "1.2"                 // Index 3: The Value
	// ]
	if len(matches) != 4 {
		return nil, fmt.Errorf("Got invalid intervention format '%s'. Wanted 'variable * 1.2' or 'variable = 500'", input)
	}

	target := matches[1]
	op := matches[2]

	val, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", matches[3])
	}

	i := &causal.Intervention{
		TargetNode: target,
	}

	switch op {
	case "*":
		i.Action = "INCREASE_BY_PERCENT"
		// making this come out right in python worker
		// new_val = original_val * (1.0 + intervention.value)
		i.Value = val - 1.0
	case "=":
		i.Action = "SET_TO_FIXED"
		i.Value = val
	}

	return i, nil
}

func printEffectSpecificResults(JSON string, intervention *causal.Intervention, effect string) error {
	var result struct {
		Metrics map[string]struct {
			Original  []float64 `json:"original"`
			Simulated []float64 `json:"simulated"`
		} `json:"metrics"`
	}

	if err := json.Unmarshal([]byte(JSON), &result); err != nil {
		return fmt.Errorf("failed to parse simulation results: %w", err)
	}

	data, ok := result.Metrics[effect]
	if !ok {
		return fmt.Errorf("effect variable '%s' was not found in simulation results", effect)
	}

	var totalOrig, totalSim float64
	count := float64(len(data.Original))

	for i := range data.Original {
		totalOrig += data.Original[i]
		totalSim += data.Simulated[i]
	}

	if count == 0 {
		return fmt.Errorf("no data was found for effect variable '%s'", effect)
	}

	avgOrig := totalOrig / count
	avgSim := totalSim / count
	delta := avgSim - avgOrig

	var pctChangeStr string

	if math.Abs(avgOrig) < 1e-9 {
		if math.Abs(delta) < 1e-9 {
			pctChangeStr = "0.00%" // 0 -> 0 is 0% change
		} else {
			pctChangeStr = "N/A" // 0 -> x (where x != 0) is undefined % change (infinite)
		}
	} else {
		pctChange := (delta / avgOrig) * 100
		pctChangeStr = fmt.Sprintf("%.2f%%", pctChange)
	}

	fmt.Printf("\n--- Simulation Report ---\n")

	action := ""
	switch intervention.Action {
	case "INCREASE_BY_PERCENT":
		action = fmt.Sprintf("Scaling %s by %.1f%%", intervention.TargetNode, intervention.Value*100)
	case "SET_TO_FIXED":
		action = fmt.Sprintf("Setting %s to %.2f", intervention.TargetNode, intervention.Value)
	}
	fmt.Printf("Intervention: %s\n", action)
	fmt.Printf("Effect:       %s\n", effect)
	fmt.Println("---------------------------------")

	fmt.Printf("Baseline Average:     %.2f\n", avgOrig)
	fmt.Printf("Counterfactual Avg:   %.2f\n", avgSim)

	sign := ""
	if delta > 0 {
		sign = "+"
	}
	fmt.Printf("Net Impact:           %s%.2f (%s)\n", sign, delta, pctChangeStr)
	fmt.Println("---------------------------------")

	return nil
}
