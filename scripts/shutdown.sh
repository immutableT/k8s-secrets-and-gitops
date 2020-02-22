#!/usr/bin/env bash
set -e

ETCD_DATA_DIR='../etcd-data'
ETCD_PORT='2379'

KAS_SECURE_PORT='8081'
WEB_HOOK_PORT=8083

PORTS=("${ETCD_PORT}" "${KAS_SECURE_PORT}" "${WEB_HOOK_PORT}")
for p in "${PORTS[@]}"; do
  if lsof -ti:"${p}"; then
    lsof -ti:"${p}" | xargs kill -9
  fi
done

rm ${ETCD_DATA_DIR:?}/* -r
rm ../logs/* -r
