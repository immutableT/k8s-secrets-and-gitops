#!/usr/bin/env bash
set -e

PUB_KEY_PATH="../certs/cluster/key.pub"
PASSWORD_PLAINTEXT="P@ssword01"
PASSWORD_CIPHERTEXT_PATH="../cluster/secrets/db-secret.yaml"

JWE=$(echo "${PASSWORD_PLAINTEXT}" | ../bin/jose-util encrypt \
      --key "${PUB_KEY_PATH}" \
      --alg RSA-OAEP-256 \
      --enc A128CBC-HS256 \
      --full | jq . )
echo "${JWE}"

cat > "${PASSWORD_CIPHERTEXT_PATH}" <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
stringData:
  password: |-
         ${JWE}
EOF
