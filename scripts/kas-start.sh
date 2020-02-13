#!/usr/bin/env bash
SECURE_PORT=8081

if sudo lsof -ti:${SECURE_PORT}; then
 sudo lsof -ti:${SECURE_PORT} | xargs kill -9
fi

../bin/kube-apiserver \
  --secure-port=${SECURE_PORT} \
  --etcd-servers=http://127.0.0.1:2379 \
  --storage-backend=etcd3 \
  --cert-dir=./apiserver.local.config/certificates  \
  --enable-admission-plugins=MutatingAdmissionWebhook \
  --logtostderr=true
  #--admission-control-config-file=../manifests/deployment/mutating-webhook-registration.yaml \
