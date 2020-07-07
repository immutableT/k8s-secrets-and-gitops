package admission

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/kr/pretty"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
)

func TestMutation(t *testing.T) {
	testCases := []struct {
		desc           string
		request        *admissionv1.AdmissionReview
		response       *admissionv1.AdmissionResponse
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
			request:        &admissionv1.AdmissionReview{},
			wantStatusCode: http.StatusInternalServerError,
			wantError:      `Internal Server Error: "/secrets": unexpected nil request`,
		},
		{
			desc: "AdmissionReview with empty AdmissionRequest",
			request: &admissionv1.AdmissionReview{
				Request: &admissionv1.AdmissionRequest{},
			},
			wantError:      `Internal Server Error: "/secrets": Request.Object.Object is nil, and the attempt to deserialize Request.Object.Raw failed with the error: Object 'Kind' is missing in ''`,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			desc: "AdmissionRequest with empty Object",
			request: &admissionv1.AdmissionReview{
				Request: &admissionv1.AdmissionRequest{
					UID:    "705ab4f5-6393-11e8-b7cc-42010a800002",
					Object: runtime.RawExtension{},
				},
			},
			wantError:      `Internal Server Error: "/secrets": Request.Object.Object is nil, and the attempt to deserialize Request.Object.Raw failed with the error: Object 'Kind' is missing in ''`,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			desc: "valid un-enveloped AdmissionRequest should not generate a patch",
			request: &admissionv1.AdmissionReview{
				TypeMeta: v1.TypeMeta{Kind: "AdmissionReview", APIVersion: "admission.k8s.io/v1"},
				Request: &admissionv1.AdmissionRequest{
					UID:                "07d88007-f4e9-4a15-8455-ef1685040976",
					Kind:               v1.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"},
					Resource:           v1.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"},
					SubResource:        "",
					RequestKind:        &v1.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"},
					RequestResource:    &v1.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"},
					RequestSubResource: "",
					Name:               "my-secret-01",
					Namespace:          "default",
					Operation:          "CREATE",
					Object: runtime.RawExtension{
						Raw:    []byte{0x7b, 0x22, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x3a, 0x22, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x2c, 0x22, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x76, 0x31, 0x22, 0x2c, 0x22, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3a, 0x7b, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x22, 0x6d, 0x79, 0x2d, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x2d, 0x30, 0x31, 0x22, 0x2c, 0x22, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x3a, 0x6e, 0x75, 0x6c, 0x6c, 0x7d, 0x2c, 0x22, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3a, 0x7b, 0x22, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x3a, 0x22, 0x55, 0x45, 0x42, 0x7a, 0x63, 0x33, 0x64, 0x76, 0x63, 0x6d, 0x51, 0x77, 0x4d, 0x51, 0x3d, 0x3d, 0x22, 0x2c, 0x22, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x22, 0x5a, 0x47, 0x56, 0x32, 0x4c, 0x58, 0x56, 0x7a, 0x5a, 0x58, 0x49, 0x3d, 0x22, 0x7d, 0x2c, 0x22, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x22, 0x4f, 0x70, 0x61, 0x71, 0x75, 0x65, 0x22, 0x7d},
						Object: nil,
					},
					OldObject: runtime.RawExtension{},
					Options: runtime.RawExtension{
						Raw:    []byte{0x7b, 0x22, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x3a, 0x22, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x2c, 0x22, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x31, 0x22, 0x7d},
						Object: nil,
					},
				},
				Response: (*admissionv1.AdmissionResponse)(nil),
			},
			response: &admissionv1.AdmissionResponse{
				UID:       "07d88007-f4e9-4a15-8455-ef1685040976",
				Allowed:   true,
				Result:    (*v1.Status)(nil),
				Patch:     []byte("[]"),
				PatchType: &jsonPatchResponse,
			},
			wantStatusCode: http.StatusOK,
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

			req := httptest.NewRequest(
				http.MethodPost,
				"/secrets",
				bytes.NewReader(requestBody))
			req.Header.Add("Content-Type", "application/json")
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

			got := &admissionv1.AdmissionReview{}
			err = json.Unmarshal(responseBody, got)
			if err != nil {
				t.Fatalf("Failed to unmarshal AdmissionReview, err: %v", err)
			}

			if diff := cmp.Diff(tt.response, got.Response); diff != "" {
				t.Logf("Got:%v", pretty.Sprint(got))
				t.Fatalf("Mismatch in AdmissionReview (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestShouldMutate(t *testing.T) {
	testCases := []struct {
		desc string
		in   []byte
		want bool
	}{
		{
			desc: "valid jwe envelope",
			in: []byte(`
         {
          "protected": "eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0",
          "encrypted_key": "W55XJrzI_rxTnBBtMK5Alg-WCwz3tFBL4JCQTT26o0VT8NdncbfQ_ksonHb3OVjOQieCeCMjNAXTsB2Vv6WkyhElhiRt4TjqGhQSGyBWHtd5o3NVdcHrRXMy5HRhKMJP_idBRr21IxoWtXTaELYLkEqDnANHVZkkvRl69nA3OvKKe4C5n9LQ9nYfwQdYBDwvKiVbGOIRbiTAWFUmAJUzui6YemRAbbv_Q6D2yVV4bK8pSwYuJ9_Hg7Q66lrRgayZJbhRrWAGCMXnMvWhCWrapvQp_oScTh1wanm4mB9deZpkEO-nJp0bloIrYAReZ1BAe05FyBtQUkNyNnB0ATuONmObmVP3aml6aVimgZE6Xsef9K2khofZuy0j6GrUMBf74bIc0ZwLI8nBz6otuQe27UD9mxlkDzZsPoJ-H2CDDkoO-Bu1zp0yfr0uQ0PQ8BlpK5krqBeiqk5jjgshIhl_qj-w2Aa8r-OdpKgny-7g6NpL6QMkYRpN1RN2lm2lUBsnKilRSdWgRkgj_2XneabQnwRYJfFF2PumXnTfcReKDHbd9SjCynxwRUg09uKQKuyyDtcd1YXkSWYeSOhSAZ0qdIPrvFdmsQAIVGBlZGr3Z7xgrzgq9Zc_j7YPWY2KkkTDYPGEhmVV5MtSsgOirU3kWJ6QjLowW16QnQXsNtnfDFQ",
          "iv": "ZKD5DLUyJVhG9T8xnSnMEQ",
          "ciphertext": "JKLXYc7C9ePhFlI53hlnNA",
          "tag": "iZnl_VDdhjMPv6CmmvhEuw"
         }`),
			want: true,
		},
		{
			desc: "no jwe envelope",
			in:   []byte("P@ssword01"),
			want: false,
		},
		// TODO(immutableT) Add tests for error conditions.
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			got := shouldMutate([]byte(tt.in))
			if got != tt.want {
				t.Fatalf("Got %v want %v", got, tt.want)
			}
		})
	}
}
