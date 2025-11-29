import numpy as np
import pandas as pd
import json
from tigramite.toymodels import structural_causal_processes as toys

print("Generating ground truth data...")

# 1. Define a simple, known causal structure
#    Format: ((parent_index, lag), coefficient, function)
#    Structure:
#    X0 --> X1 (at lag 1)
#    X1 --> X2 (at lag 2)
def lin_f(x): 
    return x

links = {
    0: [((0, -1), 0.5, lin_f)],                   
    1: [((1, -1), 0.5, lin_f), ((0, -1), 0.5, lin_f)],
    2: [((2, -1), 0.5, lin_f), ((1, -2), 0.5, lin_f)] 
}

# 2. Generate 1000 data points from this structure
T = 1000
data, _ = toys.structural_causal_process(links, T=T)

# 3. Save to a CSV file inside the container at /app/ground_truth.csv
var_names = ['service_a', 'service_b', 'service_c']
df = pd.DataFrame(data, columns=var_names)
df.to_csv('/app/ground_truth.csv', index=False)

# 4. Save the graph
graph_structure = {
    "nodes": [
        {"id": 0, "label": "service_a"},
        {"id": 1, "label": "service_b"},
        {"id": 2, "label": "service_c"}
    ],
    "edges": [
        {"source": "service_a", "target": "service_a", "type": "directed", "lag": 1},
        {"source": "service_b", "target": "service_b", "type": "directed", "lag": 1},
        {"source": "service_c", "target": "service_c", "type": "directed", "lag": 1},
        {"source": "service_a", "target": "service_b", "type": "directed", "lag": 1},
        {"source": "service_b", "target": "service_c", "type": "directed", "lag": 2}
    ]
}

with open('/app/ground_truth_graph.json', 'w') as f:
    json.dump(graph_structure, f, indent=2)

# 5. Print
print("Files generated: ground_truth.csv, ground_truth_graph.json")
