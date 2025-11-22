package grpc

import (
	"context"
	"time"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	discoverer "github.com/w-h-a/caus/internal/client/causal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcDiscoverer struct {
	client discoverer.CausalDiscoveryClient
}

func (d *grpcDiscoverer) Discover(ctx context.Context, req *causal.DiscoverRequest) (*causal.CausalGraph, error) {
	callCtx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	return d.client.Discover(callCtx, req)
}

// TODO: options (addr)
func NewDiscoverer() discoverer.Discoverer {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	c := discoverer.NewCausalDiscoveryClient(conn)

	d := &grpcDiscoverer{
		client: c,
	}

	return d
}
