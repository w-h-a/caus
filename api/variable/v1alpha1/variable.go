package v1alpha1

import (
	"fmt"
	"slices"
)

var (
	SupportedDimensions = []string{"calls", "duration"}

	SupportedAggregations = map[string][]string{
		"calls":    {"count", "rate_per_sec"},
		"duration": {"avg", "p50", "p95", "p99"},
	}

	SupportedAttributeQueryOperators = []string{"equals", "contains", "isnotnull"}
)

type DiscoveryConfig struct {
	Variables []VariableDefinition `yaml:"variables"`
}

type VariableDefinition struct {
	Name         string            `yaml:"name"`
	Source       string            `yaml:"source"` // e.g., "metrics", "traces"
	MetricsQuery string            `yaml:"metrics_query,omitempty"`
	TraceQuery   TraceQueryDetails `yaml:"trace_query,omitempty"`
}

type TraceQueryDetails struct {
	ServiceName       string           `yaml:"service"`
	SpanName          string           `yaml:"span_name,omitempty"`
	SpanKind          string           `yaml:"span_kind,omitempty"`
	Dimension         string           `yaml:"dimension"` // e.g., "duration", "calls", etc
	AggregationOption string           `yaml:"aggregation,omitempty"`
	AttributeQueries  []AttributeQuery `yaml:"attribute_queries,omitempty"`
}

func (t TraceQueryDetails) Validate() error {
	if t.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}

	if !slices.Contains(SupportedDimensions, t.Dimension) {
		return fmt.Errorf("unsupported dimension '%s'. Supported: %v", t.Dimension, SupportedDimensions)
	}

	validAggs, ok := SupportedAggregations[t.Dimension]
	if !ok {
		return fmt.Errorf("invalid dimension configuration")
	}

	if t.AggregationOption != "" && !slices.Contains(validAggs, t.AggregationOption) {
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
