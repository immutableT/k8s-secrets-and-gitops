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

package admission

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kr/pretty"
	"github.com/square/go-jose"
	"gomodules.xyz/jsonpatch/v2"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/klog"
)

var (
	scheme            = runtime.NewScheme()
	codecs            = serializer.NewCodecFactory(scheme)
	reviewGVK         = admissionv1.SchemeGroupVersion.WithKind("AdmissionReview")
	jsonPatchResponse = admissionv1.PatchTypeJSONPatch
)

func init() {
	utilruntime.Must(admissionv1.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
}

func Serve(w http.ResponseWriter, req *http.Request) {
	klog.V(4).Infof("Received request %v", pretty.Sprint(req))

	review, gvk, err := validateRequest(req)
	if err != nil {
		responsewriters.InternalError(w, req, err)
		return
	}
	klog.V(4).Infof("Received request with an AdmissionReview object %v", pretty.Sprint(review))

	secretToPatch, err := secretToReview(review)
	if err != nil {
		responsewriters.InternalError(w, req, err)
		return
	}
	klog.Infof("AdmissionReview request contains a valid secret %v", pretty.Sprint(secretToPatch))

	review.Response = &admissionv1.AdmissionResponse{
		UID:       review.Request.UID,
		PatchType: &jsonPatchResponse,
	}

	beforePatch := review.Request.Object.Raw
	var afterPatch []byte

	for k, v := range secretToPatch.Data {
		if shouldMutate(v) {
			// TODO (immutableT) Add logic to decrypt values.
			secretToPatch.Data[k] = []byte("foo")
			klog.Infof("Patching k:%v, v: %v", k, pretty.Sprint(v))
		} else {
			klog.Infof("Skipping key: %v with value %v, as it does not seem to be enveloped", k, v)
		}
	}

	afterPatch, err = json.Marshal(secretToPatch)
	if err != nil {
		responsewriters.InternalError(w, req, fmt.Errorf("unexpected encoding error: %v", err))
		return
	}
	klog.Infof("Patched and marshalled secret: %s", pretty.Sprint(afterPatch))

	patch, err := jsonpatch.CreatePatch(beforePatch, afterPatch)
	if err != nil {
		responsewriters.InternalError(w, req, fmt.Errorf("unexpected diff error: %v", err))
		return
	}
	klog.Infof("Generated patch: %v", patch)
	review.Response.Patch, err = json.Marshal(patch)
	if err != nil {
		responsewriters.InternalError(w, req, fmt.Errorf("unexpected patch encoding error: %v", err))
		return
	}

	review.Response.Allowed = true
	responsewriters.WriteObjectNegotiated(
		codecs,
		negotiation.DefaultEndpointRestrictions,
		gvk.GroupVersion(),
		w,
		req,
		http.StatusOK,
		review,
	)
}

func validateRequest(req *http.Request) (*admissionv1.AdmissionReview, *schema.GroupVersionKind, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read body: %v", err)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, nil, fmt.Errorf("contentType=%s, expect application/json", contentType)
	}

	obj, gvk, err := codecs.UniversalDeserializer().Decode(body, &reviewGVK, &admissionv1.AdmissionReview{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode body: %v", err)
	}
	review, ok := obj.(*admissionv1.AdmissionReview)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected GroupVersionKind: %s", gvk)
	}
	klog.Infof("Decoded an AdmissionReview object: %s", pretty.Sprint(review))

	if review.Request == nil {
		return nil, nil, errors.New("unexpected nil request")
	}
	return review, gvk, nil
}

func secretToReview(review *admissionv1.AdmissionReview) (*corev1.Secret, error) {
	if review.Request.Object.Object == nil {
		var err error
		review.Request.Object.Object, _, err = codecs.UniversalDeserializer().Decode(review.Request.Object.Raw, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("Request.Object.Object is nil, and the attempt to deserialize Request.Object.Raw failed with the error: %v", err)
		}
	}

	secret, ok := review.Request.Object.Object.(*corev1.Secret)
	if !ok {
		return nil, errors.New("AdmissionReview does not contain corev1.Secret")
	}

	return secret, nil
}

func shouldMutate(secret []byte) bool {
	jwe, err := jose.ParseEncrypted(string(secret))
	if err == nil {
		klog.Infof("Found JWE envelope: %v", pretty.Sprint(jwe))
		return true
	}

	return false
}
