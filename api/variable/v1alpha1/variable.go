package v1alpha1

import (
	"fmt"
	"slices"
)

var (
	SupportedImplementations = map[string][]string{
		"metrics": {"mock", "random", "csv", "prometheus", "datadog"},
		"traces":  {"mock", "random", "csv", "clickhouse", "datadog"}, // TODO: honeycomb
	}

	SupportedDimensions = []string{"calls", "duration"}

	SupportedAggregations = map[string][]string{
		"calls":    {"count"}, // TODO: rate_per_sec
		"duration": {"avg", "p50", "p95", "p99"},
	}

	SupportedAttributeQueryOperators = []string{"equals", "contains", "isnotnull"}
)

type DiscoveryConfig struct {
	Variables []VariableDefinition `yaml:"variables"`
}

func (c *DiscoveryConfig) Validate() error {
	if len(c.Variables) == 0 {
		return fmt.Errorf("no variables defined")
	}

	for i, v := range c.Variables {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("variable[%d] '%s' invalid: %w", i, v.Name, err)
		}
	}

	return nil
}

type VariableDefinition struct {
	Name         string             `yaml:"name"`
	Source       *Source            `yaml:"source"`
	MetricsQuery string             `yaml:"metrics_query,omitempty"`
	TraceQuery   *TraceQueryDetails `yaml:"trace_query,omitempty"`
}

func (v *VariableDefinition) Validate() error {
	if len(v.Name) == 0 {
		return fmt.Errorf("name is required")
	}

	if v.Source == nil {
		return fmt.Errorf("source is required")
	}

	if err := v.Source.Validate(); err != nil {
		return fmt.Errorf("source invalid: %w", err)
	}

	switch v.Source.Type {
	case "metrics":
		if len(v.MetricsQuery) == 0 {
			return fmt.Errorf("metrics_query is required for source type 'metrics'")
		}
	case "traces":
		if v.TraceQuery == nil {
			return fmt.Errorf("trace_query is required for source type 'traces'")
		}
		if err := v.TraceQuery.Validate(); err != nil {
			return fmt.Errorf("trace_query invalid: %w", err)
		}
	default:
		return fmt.Errorf("unknown source type '%s' (supported: metrics, traces)", v.Source.Type)
	}

	return nil
}

type Source struct {
	Type   string `yaml:"type"` // e.g., "metrics", "traces"
	Impl   string `yaml:"impl"` // e.g., "prometheus", "clickhouse", "honeycomb", "datadog", etc
	Loc    string `yaml:"loc"`
	ApiKey string `yaml:"api_key"`
	AppKey string `yaml:"app_key"`
}

func (s *Source) Validate() error {
	if len(s.Type) == 0 {
		return fmt.Errorf("type is required")
	}

	if len(s.Impl) == 0 {
		return fmt.Errorf("impl is required")
	}

	validImpls, ok := SupportedImplementations[s.Type]
	if !ok {
		return fmt.Errorf("unsupported source type '%s'", s.Type)
	}

	if !slices.Contains(validImpls, s.Impl) {
		return fmt.Errorf("unsupported implementation '%s' for type '%s'. Supported: %v", s.Impl, s.Type, validImpls)
	}

	if len(s.Loc) == 0 {
		return fmt.Errorf("loc (location) is required")
	}

	return nil
}

type TraceQueryDetails struct {
	ServiceName       string           `yaml:"service"`
	Dimension         string           `yaml:"dimension"` // e.g., "duration", "calls", etc
	AggregationOption string           `yaml:"aggregation"`
	SpanName          string           `yaml:"span_name,omitempty"`
	SpanKind          string           `yaml:"span_kind,omitempty"`
	AttributeQueries  []AttributeQuery `yaml:"attribute_queries,omitempty"`
}

func (t *TraceQueryDetails) Validate() error {
	if len(t.ServiceName) == 0 {
		return fmt.Errorf("service name is required")
	}

	if !slices.Contains(SupportedDimensions, t.Dimension) {
		return fmt.Errorf("unsupported dimension '%s'. Supported: %v", t.Dimension, SupportedDimensions)
	}

	validAggs, ok := SupportedAggregations[t.Dimension]
	if !ok {
		return fmt.Errorf("invalid dimension configuration")
	}

	if len(t.AggregationOption) > 0 && !slices.Contains(validAggs, t.AggregationOption) {
		return fmt.Errorf("unsupported aggregation '%s' for dimension '%s'. Supported: %v", t.AggregationOption, t.Dimension, validAggs)
	}

	for _, q := range t.AttributeQueries {
		if !slices.Contains(SupportedAttributeQueryOperators, q.Operator) {
			return fmt.Errorf("unsupported operator '%s'. Supported: %v", q.Operator, SupportedAttributeQueryOperators)
		}
	}

	return nil
}

type AttributeQuery struct {
	Key      string `yaml:"key"`
	Value    string `yaml:"value"`
	Operator string `yaml:"operator"`
}
