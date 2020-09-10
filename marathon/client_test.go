package marathon

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func checkRequest(statusCode int, expectError bool, t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
	defer ts.Close()
	m := new(Marathon)
	m.Server = ts.URL
	client := NewClient(m)
	err := client.CreateOrUpdateApplication()
	if (expectError && err == nil) || (!expectError && err != nil) {
		t.Error(err)
	}
}
