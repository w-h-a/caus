package fetcher

import (
	"context"
	"time"
)

type Fetcher interface {
	Fetch(ctx context.Context, metrics []string, start time.Time, end time.Time) ([]byte, error)
}
