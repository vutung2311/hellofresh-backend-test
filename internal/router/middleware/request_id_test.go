package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"vutung2311-golang-test/internal/router/middleware"
	"vutung2311-golang-test/pkg/tracing"
)

func TestRequestID(t *testing.T) {
	var actualRequestID string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequestID = tracing.GetReqID(r.Context())
		_, _ = w.Write([]byte(`Hello World`))
	})
	finalHandler := middleware.RequestID(handler)
	request := httptest.NewRequest("POST", "/api/hello-world", nil)
	responseWriter := httptest.NewRecorder()
	finalHandler.ServeHTTP(responseWriter, request)

	if len(actualRequestID) == 0 {
		t.Error("expect request id to be not empty")
	}
	currentHostname, err := os.Hostname()
	if err != nil {
		t.Errorf("expected no error when getting host name. Got %v", err)
	}
	if !strings.HasPrefix(actualRequestID, currentHostname) {
		t.Error("expect request id to start with hostname")
	}
}
