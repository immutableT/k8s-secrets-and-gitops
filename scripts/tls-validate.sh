#!/usr/bin/env bash


openssl s_client \
  -connect 127.0.0.1:8083 \
  -CAfile ../cmd/webhook/apiserver.local.config/certificates/secrets-decryption-webhook.crt <<< 'Q'

openssl s_client \
  -connect 127.0.0.1:8081 \
  -CAfile apiserver.local.config/certificates/apiserver.crt <<< 'Q'