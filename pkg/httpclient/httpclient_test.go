package httpclient_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"vutung2311-golang-test/pkg/httpclient"

	"github.com/sirupsen/logrus"
)

func TestClient_WithRequestResponseLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(new(logrus.JSONFormatter))

	loggerCreator := func(_ context.Context) logrus.FieldLogger {
		return logger
	}
	testServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Hello World"))
		}),
	)
	client := httpclient.New(10 * time.Second).WithRequestResponseLogger(loggerCreator)
	_, err := client.Get(context.Background(), testServer.URL)
	if err != nil {
		t.Errorf("error should be nil, got %v", err)
	}
	if len(buf.String()) == 0 {
		t.Error("should have output")
	}
	if !strings.Contains(buf.String(), `request`) ||
		!strings.Contains(buf.String(), `GET / HTTP/1.1`) ||
		!strings.Contains(buf.String(), `url`) ||
		!strings.Contains(buf.String(), testServer.URL) ||
		!strings.Contains(buf.String(), `response`) ||
		!strings.Contains(buf.String(), `Hello World`) ||
		!strings.Contains(buf.String(), `status`) ||
		!strings.Contains(buf.String(), `200 OK`) {
		t.Error("logged message missed some part")
	}
}
