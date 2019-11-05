package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vutung2311-golang-test/pkg/router/middleware"

	"github.com/sirupsen/logrus"
)

func TestRequestResponseLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(new(logrus.JSONFormatter))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data": "Hello World"}`))
	})
	finalHandler := middleware.RequestResponseLogger(logger)(handler)
	request := httptest.NewRequest("POST", "http://hello-world.com/api/hello-world", bytes.NewReader([]byte(`{"message": "Hello World"}`)))
	responseWriter := httptest.NewRecorder()
	finalHandler.ServeHTTP(responseWriter, request)

	if buf.Len() == 0 {
		t.Error("expect request response logger should create output")
	}
	if !strings.Contains(buf.String(), `{\"message\": \"Hello World\"}`) {
		t.Error("expect request response log to contain request payload")
	}
	if !strings.Contains(buf.String(), `{\"data\": \"Hello World\"}`) {
		t.Error("expect request response log to contain response body")
	}
	if !strings.Contains(buf.String(), "POST http://hello-world.com/api/hello-world") {
		t.Error("method name and path should be logged")
	}
	if !strings.Contains(buf.String(), "Incoming request and its response") {
		t.Error("should have explanation message")
	}
}
