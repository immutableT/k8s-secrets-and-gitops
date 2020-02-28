#!/usr/bin/env bash
set -e

../bin/kubectl --server=127.0.0.1:8080 get secret my-secret-01 -o yaml
