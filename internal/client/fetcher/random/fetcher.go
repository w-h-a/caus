package random

import (
	"context"
	"math/rand/v2"
	"time"

	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/fetcher"
)

type randomFetcher struct {
}

func (f *randomFetcher) Fetch(ctx context.Context, v variable.VariableDefinition, start time.Time, end time.Time, step time.Duration) (map[time.Time]float64, error) {
	data := make(map[time.Time]float64)

	curr := start.Truncate(step)
	endTime := end.Truncate(step)

	for !curr.After(endTime) {
		data[curr] = rand.Float64() * 100
		curr = curr.Add(step)
	}

	return data, nil
}

// TODO: options
func NewFetcher() fetcher.Fetcher {
	f := &randomFetcher{}

	return f
}
