package google

import (
	"github.com/google/go-cmp/cmp"
	"testing"

	"github.com/kr/pretty"
	"github.com/square/go-jose"
)

// gcloud auth application-default login - to refresh
func TestRoundTrip(t *testing.T) {
	envelope := `
{
	"protected": "eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0",
	"encrypted_key": "W55XJrzI_rxTnBBtMK5Alg-WCwz3tFBL4JCQTT26o0VT8NdncbfQ_ksonHb3OVjOQieCeCMjNAXTsB2Vv6WkyhElhiRt4TjqGhQSGyBWHtd5o3NVdcHrRXMy5HRhKMJP_idBRr21IxoWtXTaELYLkEqDnANHVZkkvRl69nA3OvKKe4C5n9LQ9nYfwQdYBDwvKiVbGOIRbiTAWFUmAJUzui6YemRAbbv_Q6D2yVV4bK8pSwYuJ9_Hg7Q66lrRgayZJbhRrWAGCMXnMvWhCWrapvQp_oScTh1wanm4mB9deZpkEO-nJp0bloIrYAReZ1BAe05FyBtQUkNyNnB0ATuONmObmVP3aml6aVimgZE6Xsef9K2khofZuy0j6GrUMBf74bIc0ZwLI8nBz6otuQe27UD9mxlkDzZsPoJ-H2CDDkoO-Bu1zp0yfr0uQ0PQ8BlpK5krqBeiqk5jjgshIhl_qj-w2Aa8r-OdpKgny-7g6NpL6QMkYRpN1RN2lm2lUBsnKilRSdWgRkgj_2XneabQnwRYJfFF2PumXnTfcReKDHbd9SjCynxwRUg09uKQKuyyDtcd1YXkSWYeSOhSAZ0qdIPrvFdmsQAIVGBlZGr3Z7xgrzgq9Zc_j7YPWY2KkkTDYPGEhmVV5MtSsgOirU3kWJ6QjLowW16QnQXsNtnfDFQ",
	"iv": "ZKD5DLUyJVhG9T8xnSnMEQ",
	"ciphertext": "JKLXYc7C9ePhFlI53hlnNA",
	"tag": "iZnl_VDdhjMPv6CmmvhEuw"
}
`
	// TODO (immutableT) Is newline expected here?
	want := []byte("P@ssword01\n")

	jwe, err := jose.ParseEncrypted(string([]byte(envelope)))
	if err != nil {
		t.Fatalf("Failed to parse JWE Envelope, err: %v", err)
	}

	t.Logf("%v", pretty.Sprint(jwe))

	kmsClient := &Client{
		Project:    "alextc-gke-dev",
		Location:   "us-central1",
		KeyRing:    "kubecon-eu-demo-ring",
		KeyName:    "kubecon-eu-key",
		KeyVersion: 1,
	}

	got, err := jwe.Decrypt(kmsClient)
	if err != nil {
		t.Fatalf("Failed to decrypt envelope, err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Logf("Got:%v", pretty.Sprint(got))
		t.Fatalf("Mismatch in secret (-want, +got)\n%s", diff)
	}
}
