"""Airflow DAG for the Iris ML pipeline: train → export → register → deploy."""

from datetime import datetime, timedelta

from airflow import DAG
from airflow.operators.python import PythonOperator

import sys
import os

# Allow imports from ml/iris relative to repo root
REPO_ROOT = os.environ.get("REPO_ROOT", "/opt/ml-platform-starter")
sys.path.insert(0, REPO_ROOT)

default_args = {
    "owner": "ml-platform",
    "depends_on_past": False,
    "retries": 1,
    "retry_delay": timedelta(minutes=5),
    "email_on_failure": False,
}

with DAG(
    dag_id="ml_pipeline_iris",
    default_args=default_args,
    description="Iris classification: train, export, register, deploy",
    schedule_interval="@daily",
    start_date=datetime(2024, 1, 1),
    catchup=False,
    tags=["iris", "mlops"],
) as dag:

    def _train(**kwargs):
        from ml.iris.train import train
        acc = train()
        kwargs["ti"].xcom_push(key="accuracy", value=acc)

    def _export(**kwargs):
        from ml.iris.export import export
        from ml.iris.train import ARTIFACT_DIR
        import os
        model_path = os.path.join(ARTIFACT_DIR, "model.pkl")
        output_dir = os.path.join(ARTIFACT_DIR, "iris_onnx", "1")
        export(model_path=model_path, output_dir=output_dir)

    def _register(**kwargs):
        from ml.iris.register import register
        from ml.iris.train import ARTIFACT_DIR
        import os
        tracking_uri = os.environ.get("MLFLOW_TRACKING_URI", "http://mlflow:5000")
        model_path = os.path.join(ARTIFACT_DIR, "model.pkl")
        onnx_path = os.path.join(ARTIFACT_DIR, "iris_onnx", "1", "model.onnx")
        register(tracking_uri=tracking_uri, model_path=model_path, onnx_path=onnx_path)

    def _copy_to_triton(**kwargs):
        """Copy the exported ONNX model to the Triton model repository."""
        import shutil
        from ml.iris.train import ARTIFACT_DIR
        triton_repo = os.environ.get("TRITON_MODEL_REPO", "/models")
        src = os.path.join(ARTIFACT_DIR, "iris_onnx")
        dst = os.path.join(triton_repo, "iris_onnx")
        if os.path.exists(src):
            shutil.copytree(src, dst, dirs_exist_ok=True)
            print(f"Copied {src} → {dst}")
        else:
            raise FileNotFoundError(f"ONNX artifacts not found at {src}")

    train_task = PythonOperator(
        task_id="train",
        python_callable=_train,
        provide_context=True,
    )

    export_task = PythonOperator(
        task_id="export",
        python_callable=_export,
        provide_context=True,
    )

    register_task = PythonOperator(
        task_id="register",
        python_callable=_register,
        provide_context=True,
    )

    copy_task = PythonOperator(
        task_id="copy_to_triton",
        python_callable=_copy_to_triton,
        provide_context=True,
    )

    train_task >> export_task >> register_task >> copy_task
