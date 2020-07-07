#!/usr/bin/env bash
set -e
# set -x

PUB_KEY_PATH="../certs/cluster/key.pub"
PASSWORD_PLAINTEXT="P@ssword01"

PROJECT="alextc-gke-dev"
LOCATION="us-central1"
KEY_RING="kubecon-eu-demo-ring"
KEY_NAME="kubecon-eu-key"


JWE=$(echo "${PASSWORD_PLAINTEXT}" | ../bin/jose-util encrypt \
      --key "${PUB_KEY_PATH}" \
      --alg RSA-OAEP-256 \
      --enc A128CBC-HS256 \
      --full | jq . )
echo "${JWE}"

gcloud kms keys describe ${KEY_NAME} \
  --project="${PROJECT}" \
  --location="${LOCATION}"  \
  --keyring="${KEY_RING}"

DEK=$(echo "${JWE}" |
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

IV=$(echo "${JWE}" | jq -r .iv | ../bin/jose-util b64decode | xxd -p)
CIPHERTEXT=$(echo "${JWE}" | jq -r .ciphertext | ../bin/jose-util b64decode | xxd -p)

echo "DEK: ${DEK}"
echo "IV: ${IV}"
echo "CIPHERTEXT: ${CIPHERTEXT}"
