package tracing

import (
	"context"
	"fmt"
	"math/rand"
	"os"
)

type ctxKeyRequestID int

const requestIDKey ctxKeyRequestID = iota

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
}

func GenerateRequestID() string {
	myID := rand.Int()
	return fmt.Sprintf("%s-%06x", hostname, myID)
}

func PutRequestID(requestID string, ctx context.Context) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}
