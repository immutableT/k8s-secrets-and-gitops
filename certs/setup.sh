#!/usr/bin/env bash

# Root CA
# openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout RootCA.key -out RootCA.pem -subj "/C=US/CN=Kubecon-EU-2020-Demo-Root-CA"
# openssl x509 -outform pem -in RootCA.pem -out RootCA.crt

# Kube-apiserver
openssl req -new -nodes -newkey rsa:2048 -keyout kube-apiserver.key -out kube-apiserver.csr -subj "/C=US/ST=WA/L=Seattle/O=Kubecon-EU-2020-Demo/CN=localhost.local"
openssl x509 -req -sha256 -days 1024 -in kube-apiserver.csr -CA RootCA.pem -CAkey RootCA.key -CAcreateserial -extfile localhost.ext -out kube-apiserver.crt