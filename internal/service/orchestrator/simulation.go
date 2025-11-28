package orchestrator

import causal "github.com/w-h-a/caus/api/causal/v1alpha1"

type SimulationArgs struct {
	Graph        *causal.CausalGraph
	Intervention *causal.Intervention
	Horizon      int32
}
