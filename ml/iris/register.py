"""Register Iris model metadata in MLflow."""

import argparse
import logging
import os
import pickle

import mlflow
import numpy as np
from sklearn.metrics import accuracy_score
from sklearn.datasets import load_iris
from sklearn.model_selection import train_test_split

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
log = logging.getLogger(__name__)

ARTIFACT_DIR = os.path.join(os.path.dirname(__file__), "artifacts")
EXPERIMENT_NAME = "iris_classification"


def parse_args():
    parser = argparse.ArgumentParser(description="Register Iris model in MLflow")
    parser.add_argument("--tracking-uri", default=os.getenv("MLFLOW_TRACKING_URI", "http://mlflow:5000"))
    parser.add_argument("--model-path", default=os.path.join(ARTIFACT_DIR, "model.pkl"))
    parser.add_argument("--onnx-path", default=os.path.join(ARTIFACT_DIR, "iris_onnx", "1", "model.onnx"))
    return parser.parse_args()


def register(tracking_uri: str, model_path: str, onnx_path: str):
    mlflow.set_tracking_uri(tracking_uri)
    mlflow.set_experiment(EXPERIMENT_NAME)

    with open(model_path, "rb") as f:
        model = pickle.load(f)

    iris = load_iris()
    _, X_test, _, y_test = train_test_split(iris.data, iris.target, test_size=0.2, random_state=42)
    preds = model.predict(X_test)
    acc = accuracy_score(y_test, preds)

    with mlflow.start_run() as run:
        mlflow.log_param("model_type", "LogisticRegression")
        mlflow.log_param("max_iter", getattr(model, "max_iter", 200))
        mlflow.log_metric("accuracy", acc)
        mlflow.log_artifact(model_path, artifact_path="model")
        if os.path.exists(onnx_path):
            mlflow.log_artifact(onnx_path, artifact_path="onnx")
        log.info("Registered run %s with accuracy=%.4f", run.info.run_id, acc)


if __name__ == "__main__":
    args = parse_args()
    register(tracking_uri=args.tracking_uri, model_path=args.model_path, onnx_path=args.onnx_path)
