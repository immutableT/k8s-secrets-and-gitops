#!/usr/bin/env bash
set -e


IV=$(echo "ZKD5DLUyJVhG9T8xnSnMEQ" | ../bing/jose-util b64decode | xxd -p)
CIPHERTEXT=$(echo "JKLXYc7C9ePhFlI53hlnNA" | ../bin/jose-util b64decode)
DEK=8cf6c4affcae6c2c62c235809fa77e65e3222c72332a4279862b45c35e7ab004

echo "${CIPHERTEXT}" | openssl aes-128-cbc -d -K "${DEK}" -iv "${IV}"

