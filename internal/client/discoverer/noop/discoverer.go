package noop

import (
	"context"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/discoverer"
)

type noopDiscoverer struct{}

func (d *noopDiscoverer) Discover(ctx context.Context, req *causal.DiscoverRequest) (*causal.CausalGraph, error) {
	return &causal.CausalGraph{}, nil
}

func NewDiscoverer(opts ...discoverer.Options) discoverer.Discoverer {
	return &noopDiscoverer{}
}
