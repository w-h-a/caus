package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/w-h-a/caus/internal/client/causal"
	"github.com/w-h-a/caus/internal/client/causal/grpc"
	mockmetrics "github.com/w-h-a/caus/internal/client/metrics/mock"
	"github.com/w-h-a/caus/internal/service/orchestrator"
)

func main() {
	log.Println("Starting causal discovery...")

	ctx := context.Background()

	log.Println("Building clients...")

	mockMetricsFetcher := mockmetrics.NewFetcher()
	grpcDiscoverer := grpc.NewDiscoverer()

	log.Println("Building core service...")

	o := orchestrator.New(mockMetricsFetcher, grpcDiscoverer)

	log.Println("Running analysis...")

	graph, err := o.RunAnalysis(ctx, []string{}, time.Now().Add(-1*time.Hour), time.Now())
	if err != nil {
		log.Fatalf("Error running analysis: %v", err)
	}

	log.Println("Causal discovery completed")

	printGraph(graph)
}

func printGraph(graph *causal.CausalGraph) {
	fmt.Println("\n--- Causal Graph Results ---")
	fmt.Println("Nodes:")
	for _, node := range graph.Nodes {
		fmt.Printf("  - %s\n", node.Label)
	}
	fmt.Println("\nDiscovered Edges:")
	if len(graph.Edges) == 0 {
		fmt.Println("  No causal edges were found.")
	} else {
		for _, edge := range graph.Edges {
			fmt.Printf("  - %s --> %s (lag: %d)\n", edge.Source, edge.Target, edge.Lag)
		}
	}
	fmt.Println("--------------------------")
}
