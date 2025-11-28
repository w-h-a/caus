package simulator

import (
	"context"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
)

type Simulator interface {
	Simulate(ctx context.Context, req *causal.SimulateRequest) (*causal.SimulateResponse, error)
}
