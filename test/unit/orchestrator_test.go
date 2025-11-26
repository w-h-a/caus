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
	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	mockdiscoverer "github.com/w-h-a/caus/internal/client/discoverer/mock"
	"github.com/w-h-a/caus/internal/client/fetcher"
	mockfetcher "github.com/w-h-a/caus/internal/client/fetcher/mock"
	"github.com/w-h-a/caus/internal/service/orchestrator"
)

func TestOrchestrator_Stiching(t *testing.T) {
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

	fetchers := map[string]map[string]fetcher.Fetcher{
		"metrics": {"mock": mFetcher},
	}

	svc := orchestrator.New(fetchers, mDiscoverer)

	// Act
	vars := []variable.VariableDefinition{
		{Name: "var_a", Source: &variable.Source{Type: "metrics", Impl: "mock"}},
		{Name: "var_b", Source: &variable.Source{Type: "metrics", Impl: "mock"}},
	}

	_, err := svc.Discover(context.Background(), vars, start, end, step, orchestrator.DiscoveryArgs{MaxLag: 1, PcAlpha: 0.05})
	require.NoError(t, err)

	// Assert
	assert.NotNil(t, mDiscoverer.LastRequest())
	csvStr := mDiscoverer.LastRequest().CsvData
	reader := csv.NewReader(strings.NewReader(csvStr))
	rows, _ := reader.ReadAll()

	/* Expected CSV Structure:
	header: var_a, var_b
	row 1:  10.0,  20.0
	row 2:  11.0,  0.0   <-- var_b filled with 0.0
	row 3:  12.0,  22.0
	*/

	assert.Equal(t, 4, len(rows)) // Header + 3 rows
	assert.Equal(t, "var_a", rows[0][0])
	assert.Equal(t, "var_b", rows[0][1])
	assert.Equal(t, "11.000000", rows[2][0]) // var_a should have data
	assert.Equal(t, "0.0", rows[2][1])       // var_b is missing data here
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

	fetchers := map[string]map[string]fetcher.Fetcher{
		"metrics": {"mock": mFetcher},
	}

	svc := orchestrator.New(fetchers, mDiscoverer)

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
