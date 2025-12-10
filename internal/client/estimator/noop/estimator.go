package noop

import (
	"context"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/estimator"
)

type noopEstimator struct{}

func (s *noopEstimator) Estimate(ctx context.Context, req *causal.EstimateRequest) (*causal.EstimateResponse, error) {
	return &causal.EstimateResponse{}, nil
}

func NewEstimator(opts ...estimator.Option) estimator.Estimator {
	return &noopEstimator{}
}
