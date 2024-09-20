package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const logReqMsg = `%s API Request for [%s] Details:
---[ REQUEST ]---------------------------------------
%s
-----------------------------------------------------`

const logRespMsg = `%s API Response for [%s] Details:
---[ RESPONSE ]--------------------------------------
%s
-----------------------------------------------------`

var _ http.RoundTripper = &debugRoundTripper{}

type debugRoundTripper struct {
	name      string
	transport http.RoundTripper
}

func NewDebugTransport(name string, transport http.RoundTripper) *debugRoundTripper {
	return &debugRoundTripper{
		name:      name,
		transport: transport,
	}
}

func (d *debugRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	requestId := "<nil>"
	if req != nil {
		requestId = fmt.Sprintf("%s %s", req.Method, req.URL)
	}

	d.logRequest(ctx, req, requestId)

	resp, err := d.transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	d.logResponse(ctx, resp, requestId)

	return resp, nil
}

func (d *debugRoundTripper) logRequest(ctx context.Context, req *http.Request, requestId string) {
	data, err := httputil.DumpRequestOut(req, true)
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf(logReqMsg, d.name, requestId, prettyPrint(data)))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("%s API request dump error: %#v", d.name, err))
	}
}

func (d *debugRoundTripper) logResponse(ctx context.Context, resp *http.Response, requestId string) {
	data, err := httputil.DumpResponse(resp, true)
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf(logRespMsg, d.name, requestId, prettyPrint(data)))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("%s API response dump error: %#v", d.name, err))
	}
}

// prettyPrint iterates through a []byte line-by-line,
// transforming any lines that are complete json into pretty-printed json.
func prettyPrint(b []byte) string {
	parts := strings.Split(string(b), "\n")
	for i, p := range parts {
		if b := []byte(p); json.Valid(b) {
			var out bytes.Buffer
			if err := json.Indent(&out, b, "", " "); err != nil {
				continue
			}
			parts[i] = out.String()
		}
		// Mask Authorization header value
		if strings.Contains(strings.ToLower(p), "authorization:") {
			kv := strings.Split(p, ": ")
			if len(kv) != 2 {
				continue
			}
			kv[1] = strings.Repeat("*", len(kv[1]))
			parts[i] = strings.Join(kv, ": ")
		}
	}
	return strings.Join(parts, "\n")
}
