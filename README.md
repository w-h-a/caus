# caus

## Problem

Your dashboard shows that `cpu_usage` spiked at the same time as `latency`. Did the CPU spike cause the latency? Or did high traffic cause both? Or did a downstream dependency stall, causing threads to pile up, increasing memory, triggering GC, which spiked the CPU? 

All of these sound reasonable. You apply a fix, deploy it, and hope. If it works, you claim victory. If it doesn't, you try the next guess. What if you could test _before_ you start coding, deploying, and hoping?

## Solution

`caus` is a counterfactual prediction engine.

It allows you to test your engineering intuition using historical data. You provide the hypothesis (a causal graph) and variable configuration, and `caus` tells you what *would have happened* if you had intervened on the causes of the effects you care about. It turns your Post-Mortem action items from "Let's try optimizing this SQL query" into "Reducing outbound call duration would have reduced latency by 42%".

**The Question:**
My theory is that the `orders` service is saturating the Redis CPU, slowing down the `payments` service. If we had rate-limited `orders` by 50%, by how much would `payments` latency decrease?

**The Test:**
```bash
caus simulate \
  --graph="/path/to/graph.json" \
  --vars="/path/to/vars.yml" \
  --start="24h" \ # start collecting data from 24 hours ago
  --end="4h" \ # end collecting data once you reach 4 hours ago
  --step="5m" \  # rows of data will be of 5m increments
  --do="orders_redis_calls * 0.5" \
  --effect="payments_latency" \
  --horizon="60" # for how many steps back are we rewriting history?
```

**The Answer:**
```text
--- Simulation Report ---
Intervention: Scaling orders_redis_calls by 50.0%
Effect:       payments_latency
---------------------------------
Baseline Average:     450.00ms
Counterfactual Avg:   120.00ms
Net Impact:           -330.00ms (-73.33%)
---------------------------------
```

**Conclusion:**
Given your hypothesis and the data, rate limiting the orders service is certainly a viable approach!

## The Architecture

* The Orchestrator (Go): Parses your variable configs, fetches aggregated data from your observability backends (e.g., ClickHouse, Datadog, Honeycomb, Prometheus, etc), and manages the simulation workflow.
* The Brain (Python): Fits a linear structural equation model to your data based on your provided causal graph. It then executes a counterfactual simulation to predict the values of your effect given your intervention for some horizon.
