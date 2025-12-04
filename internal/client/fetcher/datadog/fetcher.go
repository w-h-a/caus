package datadog

import (
	"context"
	"fmt"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/fetcher"
)

type datadogFetcher struct {
	options fetcher.Options
	client  *datadogV1.MetricsApi
}

func (f *datadogFetcher) Fetch(ctx context.Context, v variable.VariableDefinition, start time.Time, end time.Time, step time.Duration) (map[time.Time]float64, error) {
	interval := int(step.Seconds())
	if interval < 1 {
		interval = 60
	}

	var query string

	if v.Source.Type == "traces" {
		query = fmt.Sprintf("%s:trace.%s.%s.rollup(%s, %d)", v.TraceQuery.AggregationOption, v.TraceQuery.ServiceName, v.TraceQuery.Dimension, v.TraceQuery.AggregationOption, interval)
	} else {
		query = fmt.Sprintf("%s.rollup(avg, %d)", v.MetricsQuery, interval)
	}

	rsp, _, err := f.client.QueryMetrics(f.withContext(ctx), start.Unix(), end.Unix(), query)
	if err != nil {
		return nil, fmt.Errorf("datadog query failed: %w", err)
	}

	if len(rsp.Series) == 0 {
		return map[time.Time]float64{}, nil
	}
	if len(rsp.Series) > 1 {
		return nil, fmt.Errorf("query returned %d series, expected 1", len(rsp.Series))
	}

	result := map[time.Time]float64{}

	for _, point := range rsp.Series[0].Pointlist {
		if len(point) < 2 || point[0] == nil || point[1] == nil {
			continue
		}
		tsSeconds := int64(*point[0]) / 1000
		t := time.Unix(tsSeconds, 0)
		result[t.UTC().Truncate(step)] = *point[1]
	}

	return result, nil
}

func (f *datadogFetcher) withContext(ctx context.Context) context.Context {
	ctx = context.WithValue(
		ctx,
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {Key: f.options.ApiKey},
			"appKeyAuth": {Key: f.options.AppKey},
		},
	)

	return ctx
}

func NewFetcher(opts ...fetcher.Option) fetcher.Fetcher {
	options := fetcher.NewOptions(opts...)

	// TODO: validate for dadog

	df := &datadogFetcher{
		options: options,
	}

	cfg := datadog.NewConfiguration()
	cfg.Servers = datadog.ServerConfigurations{
		{URL: options.Location},
	}

	client := datadog.NewAPIClient(cfg)

	df.client = datadogV1.NewMetricsApi(client)

	return df
}
