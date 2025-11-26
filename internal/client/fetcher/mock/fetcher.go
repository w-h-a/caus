package mock

import (
	"context"
	"time"

	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/fetcher"
)

type mockFetcher struct {
	options     fetcher.Options
	data        map[string]map[time.Time]float64
	calledStart time.Time
}

func (f *mockFetcher) Fetch(ctx context.Context, v variable.VariableDefinition, start time.Time, end time.Time, step time.Duration) (map[time.Time]float64, error) {
	f.calledStart = start

	if series, ok := f.data[v.Name]; ok {
		return series, nil
	}

	return map[time.Time]float64{}, nil
}

func (f *mockFetcher) CalledStart() time.Time {
	return f.calledStart
}

func NewFetcher(opts ...fetcher.Option) *mockFetcher {
	options := fetcher.NewOptions(opts...)

	data := map[string]map[time.Time]float64{}

	if d, ok := getDataFromCtx(options.Context); ok {
		data = d
	}

	mf := &mockFetcher{
		options: options,
		data:    data,
	}

	return mf
}
