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
	"errors"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/klog"
)

var (
	scheme    = runtime.NewScheme()
	codecs    = serializer.NewCodecFactory(scheme)
	reviewGVK = admissionv1.SchemeGroupVersion.WithKind("AdmissionReview")
)

func init() {
	utilruntime.Must(admissionv1.AddToScheme(scheme))
}

func Serve(w http.ResponseWriter, req *http.Request) {
	review, gvk, err := validateRequest(req)
	if err != nil {
		responsewriters.InternalError(w, req, err)
		return
	}

	secret := validateReview(review)
	if review.Response.Allowed == false {
		responsewriters.WriteObjectNegotiated(codecs, nil, gvk.GroupVersion(), w, req, http.StatusOK, review)
		return
	}

	for k, v := range secret.Data {
		// TODO (immutableT) Add logic to detect encrypted values.
		// TODO (immutableT) Add logic to decrypt values.
		klog.Infof("Processing k:%v, v: %v", k, v)
	}

	// TODO(immutableT) Generate Json path - this is what has to be attached to the response.
	// See github.com/appscode/jsonpatch or k8s.io/client-go/util/jsonpath/jsonpath

	review.Response.Allowed = true
	responsewriters.WriteObjectNegotiated(codecs, nil, gvk.GroupVersion(), w, req, http.StatusOK, review)
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
	if review.Request == nil {
		return nil, nil, errors.New("unexpected nil request")
	}
	return review, gvk, nil
}

func validateReview(review *admissionv1.AdmissionReview) *corev1.Secret {
	review.Response = &admissionv1.AdmissionResponse{
		UID: review.Request.UID,
	}

	if review.Request.Object.Object == nil {
		var err error
		review.Request.Object.Object, _, err = codecs.UniversalDeserializer().Decode(review.Request.Object.Raw, nil, nil)
		if err != nil {
			review.Response.Result = &metav1.Status{
				Message: fmt.Sprintf("Request.Object.Object is nil, and the attempt to deserialize Request.Object.Raw failed with the error: %v", err),
				Status:  metav1.StatusFailure,
			}
			return nil
		}
	}

	if secret, ok := review.Request.Object.Object.(*corev1.Secret); !ok {
		review.Response.Result = &metav1.Status{
			Message: "AdmissionReview does not contain coverv1.Secret",
			Status:  metav1.StatusFailure,
		}
		return nil
	} else {
		return secret
	}
}
