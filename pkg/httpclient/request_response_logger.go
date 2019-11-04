package httpclient

import (
	"context"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

func (c *Client) WithRequestResponseLogger(loggerCreator func(context.Context) logrus.FieldLogger) *Client {
	if loggerCreator == nil {
		loggerCreator = func(_ context.Context) logrus.FieldLogger {
			l := logrus.StandardLogger()
			l.SetFormatter(new(logrus.JSONFormatter))
			return l
		}
	}

	fn := func(tripper http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			var (
				requestDump, responseDump []byte
				reqRespLog                = loggerCreator(req.Context())
				err                       error
				resp                      *http.Response
			)

			fields := logrus.Fields{}
			requestDump, err = httputil.DumpRequestOut(req, true)
			if err != nil {
				reqRespLog.Errorf("can't dump request because of error: %v", err)
			} else {
				fields["request"] = string(requestDump)
				fields["url"] = req.URL.String()
			}

			resp, err = tripper.RoundTrip(req)
			if err != nil {
				reqRespLog.Errorf("http round trip error: %v", err)
				return resp, err
			}

			responseDump, err = httputil.DumpResponse(resp, true)
			if err != nil {
				reqRespLog.Errorf("can't dump response because of error: %v", err)
			} else {
				fields["response"] = string(responseDump)
				fields["status"] = resp.Status
			}

			if fields["request"] != nil && fields["response"] != nil {
				reqRespLog.WithFields(fields).Info("Outgoing request and its response")
			}

			return resp, err
		})
	}

	return c.WithRoundTrippers(fn)
}
