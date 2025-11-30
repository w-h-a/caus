package unit

import (
	"context"
	"encoding/csv"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	mockdiscoverer "github.com/w-h-a/caus/internal/client/discoverer/mock"
	noopdisc "github.com/w-h-a/caus/internal/client/discoverer/noop"
	"github.com/w-h-a/caus/internal/client/fetcher"
	mockfetcher "github.com/w-h-a/caus/internal/client/fetcher/mock"
	mocksimulator "github.com/w-h-a/caus/internal/client/simulator/mock"
	noopsim "github.com/w-h-a/caus/internal/client/simulator/noop"
	"github.com/w-h-a/caus/internal/service/orchestrator"
)

func TestOrchestrator_Discover(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	// Arrange
	pst, _ := time.LoadLocation("America/Los_Angeles")
	start := time.Date(2023, 10, 1, 10, 0, 0, 0, pst) // 10:00 PST
	end := start.Add(2 * time.Minute)                 // 10:02 PST
	step := time.Minute

	t0 := start.UTC().Truncate(step)
	t1 := t0.Add(step)
	t2 := t1.Add(step)

	mockData := map[string]map[time.Time]float64{
		"var_a": {
			t0: 10.0,
			t1: 11.0,
			t2: 12.0,
		},
		"var_b": {
			t0: 20.0,
			// t1 is MISSING! (Gap)
			t2: 22.0,
		},
	}

	mFetcher := mockfetcher.NewFetcher(
		mockfetcher.WithData(mockData),
	)

	mDiscoverer := mockdiscoverer.NewDiscoverer()

	nSimulator := noopsim.NewSimulator()

	fetchers := map[string]map[string]fetcher.Fetcher{
		"metrics": {"mock": mFetcher},
	}

	svc := orchestrator.New(fetchers, mDiscoverer, nSimulator)

	// Act
	vars := []variable.VariableDefinition{
		{Name: "var_a", Source: &variable.Source{Type: "metrics", Impl: "mock"}},
		{Name: "var_b", Source: &variable.Source{Type: "metrics", Impl: "mock"}},
	}

	_, err := svc.Discover(context.Background(), vars, start, end, step, orchestrator.DiscoveryArgs{MaxLag: 1, PcAlpha: 0.05})
	require.NoError(t, err)

	// Assert
	req := mDiscoverer.LastRequest()
	assert.NotNil(t, req)
	csvStr := req.CsvData
	reader := csv.NewReader(strings.NewReader(csvStr))
	rows, _ := reader.ReadAll()

	/* Expected CSV Structure:
	header: var_a, var_b
	row 1:  10.0,  20.0
	row 2:  11.0,  20.0   <-- var_b filled with last known since it's metrics
	row 3:  12.0,  22.0
	*/

	assert.Equal(t, 4, len(rows)) // Header + 3 rows
	assert.Equal(t, "var_a", rows[0][0])
	assert.Equal(t, "var_b", rows[0][1])
	assert.Equal(t, "11.000000", rows[2][0]) // var_a should have data
	assert.Equal(t, "20.000000", rows[2][1]) // var_b is missing data here so should be filled with last known
	assert.Equal(t, int32(1), req.MaxLag)
	assert.Equal(t, float32(0.05), req.PcAlpha)
}

func TestOrchestrator_Simulate(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	// Arrange
	start := time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC) // 10:00
	end := start.Add(2 * time.Minute)                      // 10:02
	step := time.Minute

	t0 := start
	t1 := t0.Add(step)
	t2 := t1.Add(step)

	mockData := map[string]map[time.Time]float64{
		"var_a": {
			t0: 10.0,
			t1: 11.0,
			t2: 12.0,
		},
	}

	mFetcher := mockfetcher.NewFetcher(
		mockfetcher.WithData(mockData),
	)

	nDiscoverer := noopdisc.NewDiscoverer()

	mSimulator := mocksimulator.NewSimulator()

	fetchers := map[string]map[string]fetcher.Fetcher{
		"metrics": {"mock": mFetcher},
	}

	svc := orchestrator.New(fetchers, nDiscoverer, mSimulator)

	// Act
	vars := []variable.VariableDefinition{
		{Name: "var_a", Source: &variable.Source{Type: "metrics", Impl: "mock"}},
	}

	inputGraph := &causal.CausalGraph{
		Nodes: []*causal.Node{{Id: 0, Label: "var_a"}},
	}

	intervention := &causal.Intervention{
		TargetNode: "var_a",
		Action:     "SET_TO_FIXED",
		Value:      100.0,
	}

	_, err := svc.Simulate(context.Background(), vars, start, end, step, orchestrator.SimulationArgs{Graph: inputGraph, Intervention: intervention, Horizon: 60})
	require.NoError(t, err)

	// Assert
	req := mSimulator.LastRequest()
	assert.NotNil(t, req)
	assert.True(t, len(req.CsvData) > 0)
	assert.True(t, len(req.Graph.Nodes) == 1)
	assert.Equal(t, "var_a", req.Graph.Nodes[0].Label)
	assert.Equal(t, "var_a", req.Intervention.TargetNode)
	assert.Equal(t, 100.0, req.Intervention.Value)
	assert.Equal(t, int32(60), req.SimulationSteps)
}

func TestOrchestrator_FetchAlignment(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	// Arrange
	start := time.Date(2023, 1, 1, 10, 0, 47, 0, time.UTC) // 10:00:47
	end := start.Add(1 * time.Minute)
	step := time.Minute

	mFetcher := mockfetcher.NewFetcher()

	mDiscoverer := mockdiscoverer.NewDiscoverer()

	nSimulator := noopsim.NewSimulator()

	fetchers := map[string]map[string]fetcher.Fetcher{
		"metrics": {"mock": mFetcher},
	}

	svc := orchestrator.New(fetchers, mDiscoverer, nSimulator)

	// Act
	vars := []variable.VariableDefinition{
		{Name: "test", Source: &variable.Source{Type: "metrics", Impl: "mock"}},
	}

	_, err := svc.Discover(context.Background(), vars, start, end, step, orchestrator.DiscoveryArgs{})
	require.NoError(t, err)

	// Assert
	expectedStart := start.Truncate(step)
	assert.Equal(t, expectedStart, mFetcher.CalledStart()) //Fetcher received 10:00:00, NOT 10:00:47
}
