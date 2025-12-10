package mock

import (
	"context"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/estimator"
)

type mockEstimator struct {
	options     estimator.Options
	lastRequest *causal.EstimateRequest
}

func (s *mockEstimator) Estimate(ctx context.Context, req *causal.EstimateRequest) (*causal.EstimateResponse, error) {
	s.lastRequest = req
	return &causal.EstimateResponse{}, nil
}

func (d *mockEstimator) LastRequest() *causal.EstimateRequest {
	return d.lastRequest
}

func NewEstimator(opts ...estimator.Option) *mockEstimator {
	options := estimator.NewOptions(opts...)

	md := &mockEstimator{
		options: options,
	}

	return md
}
