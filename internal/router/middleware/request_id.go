package middleware

import (
	"net/http"
	"vutung2311-golang-test/pkg/tracing"
)

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get("X-Request-Id")
		if len(requestID) == 0 {
			requestID = tracing.GenerateRequestID()
		}
		next.ServeHTTP(w, r.WithContext(tracing.PutRequestID(requestID, ctx)))
	}
	return http.HandlerFunc(fn)
}
