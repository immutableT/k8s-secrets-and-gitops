#!/usr/bin/env bash
set -e

PROJECT="alextc-gke-dev"
LOCATION="us-central1"
KEY_RING="kubecon-eu-demo-ring"
KEY_NAME="kubecon-eu-key"

SECRET=$(../bin/kubectl --server=127.0.0.1:8080 get secret db-secret -o json | jq -r '.data.password' | base64 -d)
echo "${SECRET}" | jq -r .iv | ../bin/jose-util b64decode | xxd -p
echo "${SECRET}" | jq -r .ciphertext | ../bin/jose-util b64decode | xxd -p

DEK=$(echo "${SECRET}" |
  jq -r .encrypted_key |
  ../bin/jose-util b64decode |
  gcloud kms asymmetric-decrypt \
  --project="${PROJECT}" \
  --location="${LOCATION}"  \
  --keyring="${KEY_RING}" \
  --key="${KEY_NAME}" \
  --version 1 \
  --ciphertext-file=- \
  --plaintext-file=- | xxd -p -c 1000)

IV=$(echo "${SECRET}" | jq -r .iv | ../bin/jose-util b64decode | xxd -p)
CIPHERTEXT=$(echo "${SECRET}" | jq -r .ciphertext | ../bin/jose-util b64decode | xxd -p)

echo DEK:"${DEK}"
echo IV:"${IV}"
echo CIPHERTEXT:"${CIPHERTEXT}"

PLAINTEXT=$(echo "${CIPHERTEXT}" | openssl aes-128-cbc -d -K "${DEK}" -iv "${IV}")
