#!/usr/bin/env bash

ETCD_DATA_DIR='../etcd-data'

# Start etcd
if lsof -ti:2379; then
 lsof -ti:2379 | xargs kill -9
fi
rm ${ETCD_DATA_DIR:?}/* -r
etcd --advertise-client-urls http://127.0.0.1:2379 --data-dir ${ETCD_DATA_DIR} --listen-client-urls http://127.0.0.1:2379 --debug
