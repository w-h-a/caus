package cmd

import (
	"fmt"

	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/fetcher"
	"github.com/w-h-a/caus/internal/client/fetcher/clickhouse"
	"github.com/w-h-a/caus/internal/client/fetcher/csv"
	"github.com/w-h-a/caus/internal/client/fetcher/datadog"
	"github.com/w-h-a/caus/internal/client/fetcher/prometheus"
	"github.com/w-h-a/caus/internal/client/fetcher/random"
)

func initFetchers(cfg *variable.DiscoveryConfig) (map[string]map[string]fetcher.Fetcher, error) {
	fetchers := map[string]map[string]fetcher.Fetcher{
		"metrics": {},
		"traces":  {},
	}

	factories := map[string]map[string]func(loc, apiKey, appKey string) fetcher.Fetcher{
		"metrics": {
			"random": func(_, _, _ string) fetcher.Fetcher { return random.NewFetcher() },
			"csv": func(loc, _, _ string) fetcher.Fetcher {
				return csv.NewFetcher(fetcher.WithLocation(loc))
			},
			"prometheus": func(loc, _, _ string) fetcher.Fetcher { return prometheus.NewFetcher(fetcher.WithLocation(loc)) },
			"datadog": func(loc, apiKey, appKey string) fetcher.Fetcher {
				return datadog.NewFetcher(fetcher.WithLocation(loc), fetcher.WithApiKey(apiKey), fetcher.WithAppKey(appKey))
			},
		},
		"traces": {
			"random": func(_, _, _ string) fetcher.Fetcher { return random.NewFetcher() },
			"csv": func(loc, _, _ string) fetcher.Fetcher {
				return csv.NewFetcher(fetcher.WithLocation(loc))
			},
			"clickhouse": func(loc, _, _ string) fetcher.Fetcher { return clickhouse.NewFetcher(fetcher.WithLocation(loc)) },
			"datadog": func(loc, apiKey, appKey string) fetcher.Fetcher {
				return datadog.NewFetcher(fetcher.WithLocation(loc), fetcher.WithApiKey(apiKey), fetcher.WithAppKey(appKey))
			},
		},
	}

	for _, v := range cfg.Variables {
		impls, typeOk := factories[v.Source.Type]
		if !typeOk {
			return nil, fmt.Errorf("unsupported source type: %s", v.Source.Type)
		}

		factory, implOk := impls[v.Source.Impl]
		if !implOk {
			return nil, fmt.Errorf("unsupported implementation '%s' for type '%s'", v.Source.Impl, v.Source.Type)
		}

		if _, exists := fetchers[v.Source.Type][v.Source.Impl]; exists {
			continue
		}

		fetchers[v.Source.Type][v.Source.Impl] = factory(v.Source.Loc, v.Source.ApiKey, v.Source.AppKey)
	}

	return fetchers, nil
}
