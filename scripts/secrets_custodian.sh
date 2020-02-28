#!/usr/bin/env bash
set -e

PUB_KEY_PATH="../certs/cluster/key.pub"
PASSWORD_PLAINTEXT="P@ssword01"

JWE=$(echo "${PASSWORD_PLAINTEXT}" | ../bin/jose-util encrypt \
      --key "${PUB_KEY_PATH}" \
      --alg RSA-OAEP \
      --enc A128CBC-HS256 \
      --full | jq . )

cat > ../cluster/secrets/db-secret.yaml <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  password: ${JWE}
EOF
