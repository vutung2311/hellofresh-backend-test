package httpclient

import (
	"net/http"
	"time"
)

var defaultHttpClient *http.Client

func init() {
	defaultHttpClient = &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       10 * time.Second,
	}
}
