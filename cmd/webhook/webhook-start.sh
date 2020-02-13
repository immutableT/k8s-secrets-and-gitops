#!/usr/bin/env bash

SECURE_PORT=8083

# Start kube-apiserver
if sudo lsof -ti:${SECURE_PORT}; then
 sudo lsof -ti:${SECURE_PORT} | xargs kill -9
fi

./webhook --secure-port=${SECURE_PORT} --cert-dir=apiserver.local.config/certificates