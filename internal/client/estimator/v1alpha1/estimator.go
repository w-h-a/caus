package v1alpha1

import (
	"context"
	"time"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/estimator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type v1alpha1Estimator struct {
	options estimator.Options
	client  causal.CausalEstimationClient
}

func (s *v1alpha1Estimator) Estimate(ctx context.Context, req *causal.EstimateRequest) (*causal.EstimateResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	return s.client.Estimate(callCtx, req)
}

func NewEstimator(opts ...estimator.Option) estimator.Estimator {
	options := estimator.NewOptions(opts...)

	// TODO: validate options

	conn, err := grpc.NewClient(options.Location, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	c := causal.NewCausalEstimationClient(conn)

	s := &v1alpha1Estimator{
		options: options,
		client:  c,
	}

	return s
}
