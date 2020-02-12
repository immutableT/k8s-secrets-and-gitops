#!/usr/bin/env bash

../bin/kubectl --server=127.0.0.1:8080 create secret generic my-secret --from-literal=username=dev-user --from-literal=password=P@ssword01

