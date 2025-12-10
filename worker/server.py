import sys
import pandas as pd
import numpy as np
import json
import io
import logging
from concurrent import futures
from sklearn.linear_model import LinearRegression

import grpc
import causal_pb2 as pb
import causal_pb2_grpc

from tigramite import data_processing as pp
from tigramite.pcmci import PCMCI
from tigramite.independence_tests.parcorr import ParCorr

def perform_causal_discovery(csv_data_string: str, max_lag: int, pc_alpha: float) -> pb.CausalGraph:
    """
    Takes CSV, runs PCMCI, and returns a pb.CausalGraph *struct*.
    """
    try:
        # 1. Read data
        raw_data = pd.read_csv(io.StringIO(csv_data_string))
        labels = raw_data.columns.tolist()
        data_values_float = raw_data.values.astype(np.float64)
        dataframe = pp.DataFrame(data_values_float, var_names=labels)

        # 2. Initialize PCMCI
        parcorr = ParCorr(significance='analytic')
        pcmci = PCMCI(dataframe=dataframe, cond_ind_test=parcorr, verbosity=0)

        # 3. Run PCMCI using params from the request
        run_alpha = None
        if pc_alpha > 0:
            run_alpha = pc_alpha
        if max_lag <= 0:
            max_lag = 3 # default

        logging.info(f"Running PCMCI with max_lag={max_lag} and pc_alpha={run_alpha}")
        results = pcmci.run_pcmci(tau_max=max_lag, pc_alpha=run_alpha)
        
        # 4. Build the response
        graph_matrix = results['graph']
        pb_nodes = [pb.Node(id=i, label=label) for i, label in enumerate(labels)]
        pb_edges = []
        for i in range(len(labels)):      # Source
            for j in range(len(labels)):  # Target
                for tau in range(max_lag + 1): # Lag
                    if graph_matrix[i, j, tau] == '-->':
                        pb_edges.append(pb.Edge(
                            source=labels[i],
                            target=labels[j],
                            type="directed",
                            lag=tau
                        ))
        
        # 5. Return the full pb.CausalGraph struct
        return pb.CausalGraph(nodes=pb_nodes, edges=pb_edges)

    except Exception as e:
        logging.error(f"Causal discovery failed: {e}")
        raise

def perform_estimation(csv_data: str, graph_proto: pb.CausalGraph) -> dict[str, pb.ModelInfo]:
    """
    Fits SCM.
    """
    try:
        # 1. Load Data
        df = pd.read_csv(io.StringIO(csv_data))
        df = df.fillna(0)
        
        # 2. Parse Graph into Parents Lookup
        parents = {col: [] for col in df.columns}
        for edge in graph_proto.edges:
            # Graph is Source -> Target. We map Target -> Source(s).
            if edge.target in parents:
                parents[edge.target].append((edge.source, edge.lag))

        # 3. Fit SCM
        models = {}
        
        for node in df.columns:
            node_parents = parents[node]
            if not node_parents:
                continue 

            X_features = []
            feature_names = []
            
            for p_name, p_lag in node_parents:
                if p_lag >= 0:
                    X_features.append(df[p_name].shift(p_lag))
                    feature_names.append(f"{p_name}_lag{p_lag}")
            
            if not X_features:
                continue

            X = pd.concat(X_features, axis=1)
            X.columns = feature_names
            y = df[node]
            
            valid_idx = X.dropna().index
            X = X.loc[valid_idx]
            y = y.loc[valid_idx]
            
            model = LinearRegression()
            model.fit(X, y)

            logging.info(f"Model for {node}: Coeffs={model.coef_} Intercept={model.intercept_} Features={feature_names}")
            models[node] = { "model": model, "features": feature_names }

        # 4. Format Results
        pb_models = {}
        for node, info in models.items():
            pb_models[node] = pb.ModelInfo(
                features=info["features"],
                coefficients=info["model"].coef_.tolist(),
                intercept=float(info["model"].intercept_)
            )

        return pb_models

    except Exception as e:
        logging.error(f"Estimation failed: {e}")
        raise

class CausalDiscoveryServicer(causal_pb2_grpc.CausalDiscoveryServicer):
    def Discover(self, request: pb.DiscoverRequest, context):
        try:
            logging.info("Received causal discovery request.")
            pb_graph = perform_causal_discovery(
                request.csv_data,
                request.max_lag,
                request.pc_alpha
            )
            logging.info("Causal discovery complete.")
            return pb_graph
        
        except Exception as e:
            logging.error(f"Error processing request: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Python error: {e}")
            return pb.CausalGraph()

class CausalEstimationServicer(causal_pb2_grpc.CausalEstimationServicer):
    def Estimate(self, request: pb.EstimateRequest, context):
        try:
            logging.info(f"Received Estimation request.")
            models_map = perform_estimation(
                request.csv_data, 
                request.graph, 
            )
            logging.info("Estimation complete.")
            return pb.EstimateResponse(models=models_map)
        except Exception as e:
            logging.error(f"Error in Estimate: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Python error: {e}")
            return pb.EstimateResponse()

def serve():
    """Starts the gRPC server and waits for connections."""
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    causal_pb2_grpc.add_CausalDiscoveryServicer_to_server(
        CausalDiscoveryServicer(), server
    )
    causal_pb2_grpc.add_CausalEstimationServicer_to_server(
        CausalEstimationServicer(), server
    )

    port = "50051"
    server.add_insecure_port(f"[::]:{port}")
    server.start()
    print(f"Python gRPC server started, listening on port {port}...")
    logging.info(f"Server started on port {port}")
    server.wait_for_termination()

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    serve()