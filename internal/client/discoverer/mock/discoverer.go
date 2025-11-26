package mock

import (
	"context"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/discoverer"
)

type mockDiscoverer struct {
	options     discoverer.Options
	lastRequest *causal.DiscoverRequest
}

func (d *mockDiscoverer) Discover(ctx context.Context, req *causal.DiscoverRequest) (*causal.CausalGraph, error) {
	d.lastRequest = req
	return &causal.CausalGraph{
		Nodes: []*causal.Node{{Id: 0, Label: "test"}},
	}, nil
}

func (d *mockDiscoverer) LastRequest() *causal.DiscoverRequest {
	return d.lastRequest
}

func NewDiscoverer(opts ...discoverer.Option) *mockDiscoverer {
	options := discoverer.NewOptions(opts...)

	md := &mockDiscoverer{
		options: options,
	}

	return md
}
