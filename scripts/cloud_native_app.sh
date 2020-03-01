#!/usr/bin/env bash
set -e

PROJECT="alextc-gke-dev"
LOCATION="us-central1"
KEY_RING="kubecon-eu-demo-ring"
KEY_NAME="kubecon-eu-key"

SECRET=$(../bin/kubectl --server=127.0.0.1:8080 get secret db-secret -o json | jq -r '.data.password' | base64 -d)

KEK=$(echo "${SECRET}" | jq -r .encrypted_key | ../bin/jose-util b64decode | gcloud kms asymmetric-decrypt \
  --project="${PROJECT}" \
  --location="${LOCATION}"  \
  --keyring="${KEY_RING}" \
  --key="${KEY_NAME}" \
  --version 1 \
  --ciphertext-file=- \
  --plaintext-file=-)

echo "${KEK}" | base64
