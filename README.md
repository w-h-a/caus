# caus

## Problem

Your dashboard shows that `cpu_usage` spiked at the same time as `latency`. 
* Did the CPU spike cause the latency? 
* Did high traffic cause both? 
* Did a downstream dependency stall, causing threads to pile up, increasing memory, triggering GC, which spiked the CPU? 

All of these sound reasonable. You apply a fix, deploy it, and hope. If it works, you claim victory. If it doesn't, you try the next guess. What if you could test _before_ you start coding, deploying, and hoping?

## Solution

`caus` is a causal estimation engine.

It allows you to **quantitatively test your engineering intuition** using historical data. You provide the hypothesis (the causal graph), and `caus` fits a Structural Equation Model to your data to tell you the "physics" of your system.

It turns Post-Mortem debates from "I think it's the database" into "The data and causal assumptions indicate that Database Wait Time drives Latency with a 10x multiplier."

### Usage

**The Question:**
My theory is that the `orders` service is saturating the Redis CPU, slowing down the `payments` service. If we had rate-limited `orders` by 50%, by how much would `payments` latency decrease?

**The Test:**
```bash
caus estimate \
  --graph="/path/to/graph.json" \
  --vars="/path/to/vars.yml" \
  --start="24h" \ # start collecting data from 24 hours ago
  --end="4h" \ # end collecting data once you reach 4 hours ago
  --step="5m" \  # rows of data will be of 5m increments
```

**The Answer:**
```text
--- Causal Physics (Discovered Coefficients) ---
Node: publish_latency
  Intercept: -49.99
  -> publish_latency_lag1: -0.049
  -> node_loop_lag_lag0: 10.95 (STRONG)

Node: node_loop_lag
  Intercept: 5.06
  -> node_loop_lag_lag1: -0.013
  -> container_cpu_lag0: 0.005
```

**Interpretation:**
* The "Strong" Signal: node_loop_lag has a coefficient of 10.95 on publish_latency.

* Meaning: For every 1ms the Event Loop lags, user latency increases by ~11ms.

* Conclusion: The system is CPU-bound. Switching to a multi-threaded runtime (Go) or offloading compute will yield massive gains.

### The Architecture

* The Orchestrator (Go): Parses your variable configs, fetches aggregated data from your observability backends (e.g., ClickHouse, Datadog, Honeycomb, Prometheus, etc), and manages the causal inference workflow.
* The Brain (Python): Fits a linear structural equation model to your data based on your provided causal graph. It calculates the coefficients that quantify the strength and direction of relationships between your variables.
