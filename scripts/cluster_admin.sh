#!/usr/bin/env bash
set -e

ETCD_DATA_DIR='../etcd-data'
ETCD_PORT='2379'
ETCD_LOG='../logs/etcd.log'

KAS_SECURE_PORT='8081'
KAS_LOG='../logs/kas.log'
KAS_CERT_DIR='../certs/kas'

WEB_HOOK_PORT=8083
WEB_HOOK_CERT_DIR='../certs/webhook'
WEB_HOOK_LOG="../logs/webhook.log"

PASSWORD_CIPHERTEXT_PATH="../cluster/secrets/db-secret.yaml"

START_WEB_HOOK="${1:-true}"

PORTS=("${ETCD_PORT}" "${KAS_SECURE_PORT}" "${WEB_HOOK_PORT}")
for p in "${PORTS[@]}"; do
  if lsof -ti:"${p}"; then
    lsof -ti:"${p}" | xargs kill -9
  fi
done

rm ${ETCD_DATA_DIR:?}/* -r || true

../bin/etcd \
  --advertise-client-urls http://127.0.0.1:${ETCD_PORT} \
  --data-dir ${ETCD_DATA_DIR} \
  --listen-client-urls http://127.0.0.1:${ETCD_PORT} &> "${ETCD_LOG}" &

../bin/kube-apiserver \
  --secure-port=${KAS_SECURE_PORT} \
  --etcd-servers=http://127.0.0.1:${ETCD_PORT} \
  --storage-backend=etcd3 \
  --cert-dir="${KAS_CERT_DIR}"  \
  --enable-admission-plugins=MutatingAdmissionWebhook \
  --logtostderr=true &> "${KAS_LOG}" &

sleep 3
#openssl s_client \
#  -connect 127.0.0.1:${KAS_SECURE_PORT} \
#  -CAfile ${KAS_CERT_DIR}/apiserver.crt <<< 'Q'



if [[ "${START_WEB_HOOK:-false}" == "true" ]]; then
  ../bin/kubectl \
  --server=127.0.0.1:8080 \
  apply -f ../manifests/deployment/mutating-webhook-registration.yaml

  go build -o ../cmd/webhook/webhook ../cmd/webhook
  ../cmd/webhook/webhook \
    --secure-port="${WEB_HOOK_PORT}" \
    --cert-dir="${WEB_HOOK_CERT_DIR}"  &> "${WEB_HOOK_LOG}" &
fi


sleep 2
#openssl s_client \
#  -connect 127.0.0.1:${WEB_HOOK_PORT} \
#  -CAfile ${WEB_HOOK_CERT_DIR}/secrets-decryption-webhook.crt <<< 'Q'

../bin/kubectl \
  --server=127.0.0.1:8080 \
  apply -f "${PASSWORD_CIPHERTEXT_PATH}"

../bin/kubectl --server=127.0.0.1:8080 get secret db-secret -o yaml
