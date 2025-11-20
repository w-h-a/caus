# caus

## Problem

Your post-mortems are fiction. You're staring at a dozen dashboards that _look_ correlated to figure out what to blame. The action items you write are educated guesses. You fixed the symptoms, but you never found the disease.

Your on-call is a guessing game. The next time PagerDuty screams, you play the same game. You'll see the same sea of red, the same storm of alerts, and you're just as blind as last time. Your MTTR is a function of how much coffee you can chug, how many times you've seen this in the past, and how lucky you are.

Your chaos engineering is just chaos. You're not a scientist testing a hypothesis; you're just a demolition expert with a blindfold. You break things to see what happens because you have no way to predict what will happen. 

## Solution

`caus` discovers the causal structure of aggregated trace telemetry and metrics. Once you have this, it unlocks so much.

### Stage 1: The Flight Recorder (WIP)

We are building an offline, command-line tool for after the incident. You feed it data from during the incident, you go make coffee, and when you come back it will present the causal structure of what burned down your system.

It will tell you that the `checkout` latency spike didn't just happen. It was caused by the `payment` error rate, which was caused by a spike in `auth` latency 4 minutes earlier. 

The Architecture

* The Go gRPC Client: It takes a request with a list of your aggregated trace and metrics data and a past time window. It queries Prometheus for metrics and Clickhouse for trace data and fires it off to the Python gRPC Server.
* The Python gRPC Server: It gets the data and runs a time-series causal discovery algorithm (PCMCI) on it. It churns and sends back a graphical representation of your causal structure.
* The `caus` CLI: The engineer runs:
```bash
caus discover \
  --metrics="redis_cpu" \
  --traces="checkout_latency,payment_errors,auth_latency" \
  --start="2025-10-29T03:00:00Z" \
  --end="2025-10-29T05:00:00Z" \
  --max-lag=2 \
  --alpha=0.05
```
The output is a JSON object representing a causal graph.

### Stages 2 & 3: The Co-Pilot & Flight Simulator

Having the causal structure unlocks a roadmap:

* We can use causal graphs to structure probabilistic inference to a few important causes of an ongoing incident instead of squinting at 'correlated' dashboards.
* We can use causal graphs to ask 'what if' questions to structure our chaos engineering and determine our post-mortem action items instead of guessing.

## The Data: Map vs Street

`caus` works on aggregated data, whether that's aggregates from trace telemetry or metrics. The observability world is rightly obsessed with trace data---the raw, wide events that tell the story of a single request. `caus` does not replace that. 

Traces are the street view; `caus` is the map. The process of aggregation is inherently lossy. The full context for why user #54321's request failed lives in the trace. But when the whole system is on fire, you don't need to look at a single house; you need a map of the city to find where the fire started.

That is what `caus` does. It analyzes the macro-level signals to build a map of the system's causal structure. It gives you the hypothesis that tells you which street to look at. Then you can pull up your tracing tool and get the high-fidelity view of the specific requests if necessary. 