package v1alpha1

import (
	"context"
	"time"

	causal "github.com/w-h-a/caus/api/causal/v1alpha1"
	"github.com/w-h-a/caus/internal/client/discoverer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type v1alpha1Discoverer struct {
	options discoverer.Options
	client  causal.CausalDiscoveryClient
}

func (d *v1alpha1Discoverer) Discover(ctx context.Context, req *causal.DiscoverRequest) (*causal.CausalGraph, error) {
	callCtx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	return d.client.Discover(callCtx, req)
}

func NewDiscoverer(opts ...discoverer.Option) discoverer.Discoverer {
	options := discoverer.NewOptions(opts...)

	// TODO: validate options

	conn, err := grpc.NewClient(options.Location, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	c := causal.NewCausalDiscoveryClient(conn)

	d := &v1alpha1Discoverer{
		options: options,
		client:  c,
	}

	return d
}
