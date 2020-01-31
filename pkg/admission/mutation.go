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

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/klog"
)

var (
	scheme    = runtime.NewScheme()
	codecs    = serializer.NewCodecFactory(scheme)
	reviewGVK = admissionv1beta1.SchemeGroupVersion.WithKind("AdmissionReview")
)

func init() {
	utilruntime.Must(admissionv1beta1.AddToScheme(scheme))
}

func Serve(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responsewriters.InternalError(w, req, fmt.Errorf("failed to read body: %v", err))
		return
	}

	obj, gvk, err := codecs.UniversalDeserializer().Decode(body, &reviewGVK, &admissionv1beta1.AdmissionReview{})
	if err != nil {
		responsewriters.InternalError(w, req, fmt.Errorf("failed to decode body: %v", err))
		return
	}
	review, ok := obj.(*admissionv1beta1.AdmissionReview)
	if !ok {
		responsewriters.InternalError(w, req, fmt.Errorf("unexpected GroupVersionKind: %s", gvk))
		return
	}
	if review.Request == nil {
		responsewriters.InternalError(w, req, errors.New("unexpected nil request"))
		return
	}
	review.Response = &admissionv1beta1.AdmissionResponse{
		UID: review.Request.UID,
	}

	// TODO(immutableT) check if review.Request.Object.Object == nil
	switch secret := review.Request.Object.Object.(type) {
	case *corev1.Secret:
		for k, v := range secret.Data {
			// TODO (immutableT) Add logic to detect encrypted values.
			// TODO (immutableT) Add logic to decrypt values.
			klog.Infof("Processing k:%v, v: %v", k, v)
		}

	default:
		responsewriters.InternalError(w, req, fmt.Errorf("unexpected object type: %v", secret))
		return
	}

	klog.V(2).Infof("Defaulting %s/%s in version %s", review.Request.Namespace, review.Request.Name, gvk)

	// TODO(immutableT) Generate Json path instead of writing the full record.
	// See github.com/appscode/jsonpatch or k8s.io/client-go/util/jsonpath/jsonpath

	review.Response.Allowed = true

	// TODO(immutableT) This could handled more concisely with k8s.io/apiserver/pkg/endpoints/handlers/responsewriters
	resp, err := json.Marshal(review)
	if err != nil {
		klog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}

	if _, err := w.Write(resp); err != nil {
		klog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
