// handlers_test.go
package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Health)

	// pass in our Request and ResponseRecorder to handler
	handler.ServeHTTP(rr, req)

	// Check the status code is what expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what expect.
	expected := `Welcome! KMS Service Wrapper`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

var testcases = []struct {
	opr  string
	body []byte
}{
	{"createkeyringcrypto", []byte(`{"krid":"testkeyID", "clientid":"test_client", "keyid":"key01", "labels":{"keyname": "global"}}`)},
	{"encrypt", []byte(`{"krid":"7f749e88-ae86-4493-149d-8a80292e7c70", "clientid":"test_client", "keyid":"key01", "dataset":{"label1": "encrypt me"}}`)},
	{"decrypt", []byte(`{"krid":"7f749e88-ae86-4493-149d-8a80292e7c70", "clientid":"test_client", "keyid":"key01", "dataset":{"label1": "CiQAIWT014nsVsS7g82lYuVpR8GIAdmtCVz_kIBwhgK-9gapYuYSMwAufkdgpSHyXgo5vZBDdcrSeMArbWikdMgYhEEmiQAdq98E2LCZqCGhyTkAT7NeJbtQrQ=="}}`)},
	{"keyringrotation", []byte(`{"krid":"7f749e88-ae86-4493-149d-8a80292e7c70", "clientid":"test_client", "keyid":"key01", "labels":{"keyname": "global"}}`)},
}

func TestRequestKMS(t *testing.T) {
	t.Log(len(testcases))
	for _, tt := range testcases {
		t.Run(tt.opr, func(t *testing.T) {

			//var jsonStr = []byte(`{"krid":"testkeyID", "clientid":"test_client", "keyid":"key01"}`)
			req, err := http.NewRequest("POST", "/kms", bytes.NewBuffer(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			q, err := url.ParseQuery(req.URL.RawQuery)
			if err != nil {
				t.Fatal(err)
			}
			// Set Operation Query Param
			q.Set("opr", tt.opr)

			req.URL.RawQuery = q.Encode()

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(RequestKMS)

			// pass in our Request and ResponseRecorder to handler
			handler.ServeHTTP(rr, req)

			// Check the status code is what expect.
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			t.Log(rr.Body.String())
		})
	}
}
