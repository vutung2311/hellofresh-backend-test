package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"sync"
	"vutung2311-golang-test/pkg/tracing"

	"github.com/sirupsen/logrus"
)

type requestResponseLogger struct {
	requestPayload []byte
	responseData   []byte

	http.ResponseWriter
}

func (l *requestResponseLogger) WithRequestAndResponse(w http.ResponseWriter, r *http.Request) error {
	var err error
	l.ResponseWriter = w
	l.requestPayload, err = httputil.DumpRequest(r, true)
	return err
}

func (l *requestResponseLogger) Empty() *requestResponseLogger {
	l.responseData = nil
	l.requestPayload = nil
	l.ResponseWriter = nil
	return l
}

func (l *requestResponseLogger) Header() http.Header {
	return l.ResponseWriter.Header()
}

func (l *requestResponseLogger) Write(b []byte) (int, error) {
	l.responseData = append(l.responseData, b...)
	return l.ResponseWriter.Write(b)
}

func (l *requestResponseLogger) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
}

var loggerPool *sync.Pool

func init() {
	loggerPool = &sync.Pool{
		New: func() interface{} {
			return new(requestResponseLogger)
		},
	}
}

func RequestResponseLogger(logger logrus.FieldLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			l := loggerPool.Get().(*requestResponseLogger)
			err := l.WithRequestAndResponse(w, r)
			if err != nil {
				http.Error(w, fmt.Sprintf("error while collecting request and response: %v", err), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(l, r)

			logger.WithFields(
				logrus.Fields{
					"contextId": tracing.GetReqID(r.Context()),
					"request":   string(l.requestPayload),
					"response":  string(l.responseData),
				},
			).Info("Incoming request and its response")
			loggerPool.Put(l.Empty())
		}
		return http.HandlerFunc(fn)
	}
}
