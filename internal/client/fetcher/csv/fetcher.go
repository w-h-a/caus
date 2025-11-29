package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/fetcher"
)

type csvFetcher struct {
	options fetcher.Options
}

func (f *csvFetcher) Fetch(ctx context.Context, v variable.VariableDefinition, start time.Time, end time.Time, step time.Duration) (map[time.Time]float64, error) {
	// 1. Nab the csv file
	file, err := os.Open(f.options.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file %s: %w", f.options.Location, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv data: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("csv file is empty")
	}

	// 2. Find the column index for this variable
	header := records[0]
	colIdx := -1
	for i, colName := range header {
		if colName == v.Name {
			colIdx = i
			break
		}
	}

	if colIdx == -1 {
		return nil, fmt.Errorf("variable '%s' not found in csv columns: %v", v.Name, header)
	}

	// 3. Parse Data & Synthesize Timestamps
	result := map[time.Time]float64{}
	totalRows := len(records) - 1

	// Start from the last row (most recent)
	curr := end.UTC().Truncate(step)

	for i := totalRows; i > 0; i-- {
		// Parse value
		valStr := records[i][colIdx]
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue
		}

		// Store only if within the requested window
		if !curr.Before(start) && !curr.After(end) {
			result[curr.Truncate(0)] = val
		}

		// Move back in time
		curr = curr.Add(-1 * step)

		// Stop if we go past start
		if curr.Before(start) {
			break
		}
	}

	return result, nil
}

func NewFetcher(opts ...fetcher.Option) fetcher.Fetcher {
	options := fetcher.NewOptions(opts...)

	cf := &csvFetcher{
		options: options,
	}

	return cf
}
