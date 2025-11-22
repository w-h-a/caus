package fetcher

import (
	"context"
	"time"

	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
)

type Fetcher interface {
	Fetch(ctx context.Context, v variable.VariableDefinition, start time.Time, end time.Time, step time.Duration) (map[time.Time]float64, error)
}
