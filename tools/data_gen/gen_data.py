import numpy as np
import pandas as pd
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
    1: [((1, -1), 0.5, lin_f), ((0, -1), 0.7, lin_f)],
    2: [((2, -1), 0.5, lin_f), ((1, -2), 0.8, lin_f)] 
}

# 2. Generate 500 data points from this structure
T = 500
data, _ = toys.structural_causal_process(links, T=T)

# 3. Save to a CSV file inside the container at /app/ground_truth.csv
var_names = ['service_a', 'service_b', 'service_c']
df = pd.DataFrame(data, columns=var_names)
df.to_csv('/app/ground_truth.csv', index=False)

# 4. Print to console
print("... Data saved to 'ground_truth.csv'")
print("\n--- ANSWER KEY (GROUND TRUTH) ---")
print("  - service_a --> service_b (lag: 1)")
print("  - service_b --> service_c (lag: 2)")
print("---------------------------------")