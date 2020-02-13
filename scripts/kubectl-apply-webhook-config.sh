#!/usr/bin/env bash

../bin/kubectl --server=127.0.0.1:8080 apply -f ../manifests/deployment/mutating-webhook-registration.yaml

