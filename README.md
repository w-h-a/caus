# caus

## Problem

Your post-mortems are fiction. You're staring at a dozen dashboards that _look_ correlated to figure out what to blame. The action items you write are educated guesses. You fixed the symptoms, but you never found the disease.

Your on-call is a guessing game. The next time PagerDuty screams, you play the same game. You'll see the same sea of red, the same storm of alerts, and you're just as blind as last time. Your MTTR is a function of how much coffee you can chug, how many times you've seen this in the past, and how lucky you are.

Your chaos engineering is just chaos. You're not a scientist testing a hypothesis; you're just a demolition expert with a blindfold. You break things to see what happens because you have no way to predict what will happen. 

## Solution

`caus` discovers the causal structure of aggregated trace telemetry and metrics. Once you have this, it unlocks so much.

### The Flight Recorder

We built an offline, CLI for after the incident. You feed it data from during the incident, and it will present the causal structure of what burned down your system. It will tell you that the `checkout` latency spike didn't just happen. It was caused by the `payment` error rate, which was caused by a spike in `auth` latency 4 minutes earlier. 

The `caus` CLI:
```bash
caus discover \
  --vars="/path/to/vars.yml" \
  --start="2h" \
  --end="0" \
  --step="1m" \
  --max-lag=3 \
  --alpha=0.05
```
The output is a textual representation of the causal graph.

### The Flight Simulator

We also built a counterfactual prediction engine. You feed it the causal graph, select a historical time window, and ask a "What Would Happen If?" question. It allows you to test remediation strategies without touching production. It will tell you: "if we had increased `frontend` calls by 20% during yesterday's traffic spike, `orders` cpu would have spiked to 100% and crashed. This turns our post-mortem action items from guesses into rigorous predictions.

The `caus` CLI:
```bash
caus simulate \
  --graph="/path/to/graph.json" \
  --do="service_a_calls * 1.2" \
  --vars="/path/to/vars.yml" \
  --start="2h" \
  --end="0" \
  --step="1m" \  
  --effect="service_b_cpu" \
  --horizon="60"
```
The output is a quantitative answer to your counterfactual question.

## The Architecture

* The Go gRPC Client (Orchestrator): It takes a request with a list of your aggregated trace and metrics data and a past time window. It queries for metrics and for aggregated trace data and fires off requests to perform discovery or simulations to the Python gRPC Server.
* The Python gRPC Server (Worker): It receives the data and query and performs the heavy lifting:
  * For Discovery, it runs the PCMCI time-series causal discovery algorithm
  * For Simulation, it fits a linear causal model based on the discovered graph and data to make counterfactual predictions.

## The Data: Map vs Street

`caus` works on aggregated data, whether that's aggregates from trace telemetry or metrics. The observability world is rightly obsessed with trace data---the raw, wide events that tell the story of a single request. `caus` does not replace that. 

Traces are the street view; `caus` is the map. 

The process of aggregation is inherently lossy. The full context for why user #54321's request failed lives in the trace. But when the whole system is on fire, you don't need to look at a single house; you need a map of the city to find where the fire started.

That is what `caus` does. It analyzes the macro-level signals to build a map of the system's causal structure. It gives you the hypothesis that tells you which street to look at. Then you can pull up your tracing tool and get the high-fidelity view of the specific requests if necessary. 