package v1alpha1

import (
	"context"
	"time"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/simulator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type v1alpha1Simulator struct {
	options simulator.Options
	client  causal.CausalSimulationClient
}

func (s *v1alpha1Simulator) Simulate(ctx context.Context, req *causal.SimulateRequest) (*causal.SimulateResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	return s.client.Simulate(callCtx, req)
}

func NewSimulator(opts ...simulator.Option) simulator.Simulator {
	options := simulator.NewOptions(opts...)

	// TODO: validate options

	conn, err := grpc.NewClient(options.Location, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	c := causal.NewCausalSimulationClient(conn)

	s := &v1alpha1Simulator{
		options: options,
		client:  c,
	}

	return s
}
