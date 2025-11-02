package causal

import context "context"

type Discoverer interface {
	Discover(ctx context.Context, req *DiscoverRequest) (*CausalGraph, error)
}
