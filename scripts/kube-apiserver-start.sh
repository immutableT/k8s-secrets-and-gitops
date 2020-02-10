#!/usr/bin/env bash

SECURE_PORT=8081
if sudo lsof -ti:${SECURE_PORT}; then
 sudo lsof -ti:${SECURE_PORT} | xargs kill -9
fi

../bin/kube-apiserver --secure-port=${SECURE_PORT} --etcd-servers=http://127.0.0.1:2379 --storage-backend=etcd3 --tls-cert-file=../certs/kube-apiserver.crt --tls-private-key-file=../certs/kube-apiserver.key --logtostderr=true