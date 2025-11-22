package discoverer

import (
	"context"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
)

type Discoverer interface {
	Discover(ctx context.Context, req *causal.DiscoverRequest) (*causal.CausalGraph, error)
}
