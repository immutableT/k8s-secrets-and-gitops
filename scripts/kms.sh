#!/usr/bin/env bash
set -e

PROJECT="alextc-gke-dev"
LOCATION="us-central1"
KEY_RING="kubecon-eu-demo-ring"
KEY_NAME="kubecon-eu-key"

gcloud kms keys create "${KEY_NAME}" \
  --project "${PROJECT}" \
  --location "${LOCATION}" \
  --keyring "${KEY_RING}" \
  --purpose asymmetric-encryption \
  --default-algorithm rsa-decrypt-oaep-4096-sha256

gcloud kms keys versions  \
  get-public-key 1 \
  --project "${PROJECT}" \
  --location "$LOCATION" \
  --keyring "$KEY_RING" \
  --key "$KEY_NAME" \
  --output-file ../certs/cluster/key.pub
