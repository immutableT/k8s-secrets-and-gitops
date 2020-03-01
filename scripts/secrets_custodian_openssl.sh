#!/usr/bin/env bash
set -e

PASSWORD_PLAINTEXT="P@ssword01"
PUB_KEY_PATH="../certs/cluster/key.pub"


echo "${PASSWORD_PLAINTEXT}" | openssl pkeyutl \
  -encrypt -pubin \
  -inkey "${PUB_KEY_PATH}" \
  -pkeyopt rsa_padding_mode:oaep \
  -pkeyopt rsa_oaep_md:sha256 \
  -pkeyopt rsa_mgf1_md:sha256 | base64 -w 0

