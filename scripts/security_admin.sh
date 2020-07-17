#!/usr/bin/env bash
set -e

PROJECT="alextc-gke-dev"
LOCATION="us-central1"
KEY_RING="kubecon-eu-demo-ring"
KEY_NAME="kubecon-eu-key"
SA_NAME="secrets-decrypter-sa"
SA="${SA_NAME}@${PROJECT}.iam.gserviceaccount.com"
PUB_KEY_PATH="../certs/cluster/key.pub"
#CLUSTER_ADMIN="alextc@google.com"

gcloud config set project "${PROJECT}"

gcloud kms keys create "${KEY_NAME}" \
  --location "${LOCATION}" \
  --keyring "${KEY_RING}" \
  --purpose asymmetric-encryption \
  --default-algorithm rsa-decrypt-oaep-4096-sha256

gcloud iam service-accounts create "${SA_NAME}" \
    --description "SA used for decrypting secrets" \
    --display-name "secrets-decrypter-sa"

gcloud kms keys add-iam-policy-binding "${KEY_NAME}" \
  --location "${LOCATION}" \
  --keyring "${KEY_RING}" \
  --member "serviceAccount:${SA}" \
  --role roles/cloudkms.cryptoKeyDecrypter

#gcloud iam service-accounts add-iam-policy-binding "${SA}" \
#  --member "${CLUSTER_ADMIN}" \
#  --role roles/iam.serviceAccountUser

gcloud kms keys versions  \
  get-public-key 1 \
  --location "${LOCATION}" \
  --keyring "${KEY_RING}" \
  --key "${KEY_NAME}" \
  --output-file "${PUB_KEY_PATH}"
