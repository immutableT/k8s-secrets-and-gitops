package admission

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/types"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
)

func TestMutation(t *testing.T) {
	testCases := []struct {
		desc           string
		request        *admissionv1beta1.AdmissionReview
		wantStatusCode int
		wantError      string
	}{
		{
			desc:           "nil body",
			wantError:      `Internal Server Error: "/secrets": unexpected nil request`,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			desc:           "empty AdmissionReview",
			wantError:      `Internal Server Error: "/secrets": unexpected nil request`,
			request:        &admissionv1beta1.AdmissionReview{},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			desc:      "AdmissionReview with empty AdmissionRequest",
			wantError: `Internal Server Error: "/secrets": unexpected object type: &lt;nil&gt;`,
			request: &admissionv1beta1.AdmissionReview{
				Request: &admissionv1beta1.AdmissionRequest{},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			desc:      "AdmissionRequest with UID",
			wantError: `Internal Server Error: "/secrets": unexpected object type: &lt;nil&gt;`,
			request: &admissionv1beta1.AdmissionReview{
				Request: &admissionv1beta1.AdmissionRequest{
					UID: types.UID("1"),
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			var (
				requestBody  []byte
				responseBody []byte
				err          error
			)

			if tt.request != nil {
				requestBody, err = json.Marshal(tt.request)
				if err != nil {
					t.Fatalf("Failed to marshal AdmissionReview, err: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodGet, "/secrets", bytes.NewReader(requestBody))
			w := httptest.NewRecorder()
			Serve(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatusCode {
				t.Fatalf("Got StatusCode %v want %v", resp.StatusCode, tt.wantStatusCode)
			}

			responseBody, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read resp.Body, err: %v", err)
			}

			if resp.StatusCode != http.StatusOK && tt.wantError != strings.TrimSuffix(string(responseBody), "\n") {
				t.Fatalf("Expeted the body to contain an error. Got: %s want: %s", responseBody, tt.wantError)
			}

			// We don't expect to receive anything parsable to AdmissionReview in the body if an
			// internal error occurred, so we stop testing here.
			if resp.StatusCode == http.StatusInternalServerError {
				return
			}

			got := &admissionv1beta1.AdmissionReview{}
			err = json.Unmarshal(responseBody, got)
			if err != nil {
				t.Fatalf("Failed to unmarshal AdmissionReview, err: %v", err)
			}
		})
	}
}
