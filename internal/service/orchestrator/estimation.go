package orchestrator

import causal "github.com/w-h-a/caus/api/causal/v1alpha1"

type EstimateArgs struct {
	Graph *causal.CausalGraph
}
