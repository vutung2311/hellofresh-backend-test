package httpclient

import (
	"context"
	"net/http"
	"time"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type Client struct {
	realClient *http.Client
}

func New(readTimeout time.Duration) *Client {
	c := &Client{
		realClient: &http.Client{
			Timeout: readTimeout,
		},
	}
	return c
}

func (c *Client) WithRoundTrippers(trippers ...func(http.RoundTripper) http.RoundTripper) *Client {
	finalTransport := http.DefaultTransport
	if c.realClient.Transport != nil {
		finalTransport = c.realClient.Transport
	}
	for _, tripper := range trippers {
		finalTransport = tripper(finalTransport)
	}
	c.realClient.Transport = finalTransport

	return c
}

func (c *Client) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.realClient.Do(req)
}
