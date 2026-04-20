"""Export trained Iris model to ONNX format for Triton."""

import argparse
import logging
import os
import pickle

import numpy as np
from skl2onnx import convert_sklearn
from skl2onnx.common.data_types import FloatTensorType

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
log = logging.getLogger(__name__)

ARTIFACT_DIR = os.path.join(os.path.dirname(__file__), "artifacts")
TRITON_MODEL_DIR = os.path.join(ARTIFACT_DIR, "iris_onnx", "1")


def parse_args():
    parser = argparse.ArgumentParser(description="Export Iris model to ONNX")
    parser.add_argument("--model-path", default=os.path.join(ARTIFACT_DIR, "model.pkl"))
    parser.add_argument("--output-dir", default=TRITON_MODEL_DIR)
    return parser.parse_args()


def export(model_path: str, output_dir: str):
    with open(model_path, "rb") as f:
        model = pickle.load(f)

    initial_type = [("float_input", FloatTensorType([None, 4]))]
    onnx_model = convert_sklearn(model, initial_types=initial_type)

    os.makedirs(output_dir, exist_ok=True)
    onnx_path = os.path.join(output_dir, "model.onnx")
    with open(onnx_path, "wb") as f:
        f.write(onnx_model.SerializeToString())
    log.info("ONNX model saved to %s", onnx_path)


if __name__ == "__main__":
    args = parse_args()
    export(model_path=args.model_path, output_dir=args.output_dir)
