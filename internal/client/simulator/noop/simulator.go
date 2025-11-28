package noop

import (
	"context"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/simulator"
)

type noopSimulator struct{}

func (s *noopSimulator) Simulate(ctx context.Context, req *causal.SimulateRequest) (*causal.SimulateResponse, error) {
	return &causal.SimulateResponse{}, nil
}

func NewSimulator(opts ...simulator.Option) simulator.Simulator {
	return &noopSimulator{}
}
