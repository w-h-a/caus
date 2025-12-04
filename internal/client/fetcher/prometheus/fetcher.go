package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/fetcher"
)

type prometheusFetcher struct {
	options fetcher.Options
	api     v1.API
}

func (f *prometheusFetcher) Fetch(ctx context.Context, v variable.VariableDefinition, start time.Time, end time.Time, step time.Duration) (map[time.Time]float64, error) {
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	val, _, err := f.api.QueryRange(ctx, v.MetricsQuery, r)
	if err != nil {
		return nil, fmt.Errorf("prometheus query failed: %w", err)
	}

	matrix, ok := val.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("expected prometheus Matrix response, got %T", val)
	}

	if len(matrix) == 0 {
		return map[time.Time]float64{}, nil
	}
	if len(matrix) > 1 {
		return nil, fmt.Errorf("query returned %d series, expected 1", len(matrix))
	}

	result := map[time.Time]float64{}
	stream := matrix[0]

	for _, pair := range stream.Values {
		t := pair.Timestamp.Time()
		result[t.UTC().Truncate(step)] = float64(pair.Value)
	}

	return result, nil
}

func NewFetcher(opts ...fetcher.Option) fetcher.Fetcher {
	options := fetcher.NewOptions(opts...)

	// TODO: validate options

	pf := &prometheusFetcher{
		options: options,
	}

	client, err := api.NewClient(api.Config{
		Address: options.Location,
	})
	if err != nil {
		detail := fmt.Sprintf("prometheus fetcher failed to create client: %v", err)
		panic(detail)
	}

	pf.api = v1.NewAPI(client)

	return pf
}
