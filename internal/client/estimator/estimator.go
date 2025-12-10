package estimator

import (
	"context"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
)

type Estimator interface {
	Estimate(ctx context.Context, req *causal.EstimateRequest) (*causal.EstimateResponse, error)
}
