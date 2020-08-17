/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/immutableT/k8s-secrets-and-gitops/pkg/admission"
	"github.com/immutableT/k8s-secrets-and-gitops/pkg/kms/google"
)

func main() {
	var (
		project  = flag.String("project", "", "GCP Project where Cloud KMS key is located.")
		location = flag.String("location", "", "Location where Cloud KMS key-ring is located.")
		ring     = flag.String("ring", "", "Key-Ring where Cloud KMS key is stored.")
		key      = flag.String("key", "", "Name of the Cloud KMS key.")
		ver      = flag.Int("ver", 1, "Key's version.")
	)

	h := &admission.WebHook{
		KMSClient: &google.Client{
			Project:    *project,
			Location:   *location,
			KeyRing:    *ring,
			KeyName:    *key,
			KeyVersion: *ver,
		},
	}
	http.HandleFunc("/secrets", h.Serve)
	err := http.ListenAndServeTLS(
		":8083",
		"../certs/webhook/secrets-decryption-webhook.crt",
		"../certs/webhook/secrets-decryption-webhook.key",
		nil)
	if err != nil {
		log.Fatal(err)
	}

}
