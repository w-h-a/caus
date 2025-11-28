package mock

import (
	"context"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/simulator"
)

type mockSimulator struct {
	options     simulator.Options
	lastRequest *causal.SimulateRequest
}

func (s *mockSimulator) Simulate(ctx context.Context, req *causal.SimulateRequest) (*causal.SimulateResponse, error) {
	s.lastRequest = req
	return &causal.SimulateResponse{
		JsonResults: `{}`,
	}, nil
}

func (d *mockSimulator) LastRequest() *causal.SimulateRequest {
	return d.lastRequest
}

func NewSimulator(opts ...simulator.Option) *mockSimulator {
	options := simulator.NewOptions(opts...)

	md := &mockSimulator{
		options: options,
	}

	return md
}
