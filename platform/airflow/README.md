# Airflow — ML Pipeline DAG

Airflow DAG that orchestrates the full Iris ML pipeline.

## DAG: ml_pipeline_iris

```
train → export → register → copy_to_triton
```

| Task | Description |
|------|-------------|
| train | Train LogisticRegression on Iris dataset |
| export | Export model to ONNX |
| register | Register run + artifacts in MLflow |
| copy_to_triton | Copy ONNX files to Triton model repository |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| MLFLOW_TRACKING_URI | http://mlflow:5000 | MLflow tracking server URL |
| TRITON_MODEL_REPO | /models | Path to Triton model repository |
| REPO_ROOT | /opt/ml-platform-starter | Path to repo root (for Python imports) |

## Local Testing

```bash
# Install Airflow and dependencies
pip install apache-airflow mlflow scikit-learn skl2onnx

# Set repo root so the DAG can import ml.iris
export REPO_ROOT=$(pwd)
export MLFLOW_TRACKING_URI=http://localhost:5000

# Test individual functions
python -c "from platform.airflow.dags.ml_pipeline_iris import *"

# Or run the full pipeline manually
python ml/iris/train.py
python ml/iris/export.py
python ml/iris/register.py
```

## Running with Airflow standalone

```bash
export AIRFLOW_HOME=$(pwd)/.airflow
export AIRFLOW__CORE__DAGS_FOLDER=$(pwd)/platform/airflow/dags

airflow standalone
# → Open http://localhost:8080
# → Enable and trigger ml_pipeline_iris
```
