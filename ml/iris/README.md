# Iris Classifier

Minimal Iris classification example using scikit-learn, ONNX export, and MLflow registration.

## Directory Layout

```
ml/iris/
├── train.py       # Train logistic regression on Iris dataset
├── export.py      # Export model to ONNX (Triton-compatible)
├── register.py    # Register metadata in MLflow
├── requirements.txt
└── artifacts/     # Generated: model.pkl, iris_onnx/1/model.onnx
```

## Commands

```bash
pip install -r requirements.txt

# Train
python train.py

# Export to ONNX
python export.py

# Register in MLflow (set MLFLOW_TRACKING_URI first)
export MLFLOW_TRACKING_URI=http://localhost:5000
python register.py
```

## Triton Model Repository Output

After running train + export:
```
artifacts/
└── iris_onnx/
    └── 1/
        └── model.onnx
```

Copy `artifacts/iris_onnx/` to the Triton model repository.
