"""Train an Iris classifier and save the model artifact."""

import argparse
import logging
import os
import pickle

import numpy as np
from sklearn.datasets import load_iris
from sklearn.linear_model import LogisticRegression
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
log = logging.getLogger(__name__)

ARTIFACT_DIR = os.path.join(os.path.dirname(__file__), "artifacts")


def parse_args():
    parser = argparse.ArgumentParser(description="Train Iris classifier")
    parser.add_argument("--max-iter", type=int, default=200, help="Logistic regression max iterations")
    parser.add_argument("--test-size", type=float, default=0.2, help="Test split fraction")
    parser.add_argument("--random-state", type=int, default=42, help="Random seed")
    return parser.parse_args()


def train(max_iter: int = 200, test_size: float = 0.2, random_state: int = 42):
    iris = load_iris()
    X_train, X_test, y_train, y_test = train_test_split(
        iris.data, iris.target, test_size=test_size, random_state=random_state
    )

    model = LogisticRegression(max_iter=max_iter, random_state=random_state)
    model.fit(X_train, y_train)

    preds = model.predict(X_test)
    acc = accuracy_score(y_test, preds)
    log.info("Accuracy: %.4f", acc)

    os.makedirs(ARTIFACT_DIR, exist_ok=True)
    model_path = os.path.join(ARTIFACT_DIR, "model.pkl")
    with open(model_path, "wb") as f:
        pickle.dump(model, f)
    log.info("Model saved to %s", model_path)

    params_path = os.path.join(ARTIFACT_DIR, "params.npy")
    np.save(params_path, np.array([max_iter, test_size, random_state, acc]))
    log.info("Params saved to %s", params_path)

    return acc


if __name__ == "__main__":
    args = parse_args()
    train(max_iter=args.max_iter, test_size=args.test_size, random_state=args.random_state)
