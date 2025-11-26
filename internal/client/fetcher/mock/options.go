package mock

import (
	"context"
	"time"

	"github.com/w-h-a/caus/internal/client/fetcher"
)

type dataKey struct{}

func WithData(d map[string]map[time.Time]float64) fetcher.Option {
	return func(o *fetcher.Options) {
		o.Context = context.WithValue(o.Context, dataKey{}, d)
	}
}

func getDataFromCtx(ctx context.Context) (map[string]map[time.Time]float64, bool) {
	d, ok := ctx.Value(dataKey{}).(map[string]map[time.Time]float64)
	return d, ok
}
