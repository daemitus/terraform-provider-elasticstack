package kibana

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/daemitus/terraform-provider-elasticstack/internal/config"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util/logging"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

// Client contain the REST client.
type Client struct {
	http *resty.Client
}

// NewClient creates a new Client.
func NewClient(cfg config.ServiceConfig) (*Client, error) {
	serverURL := cfg.Endpoint
	if serverURL == "" {
		serverURL = "http://localhost:5601"
	}
	serverURL = strings.TrimSuffix(serverURL, "/")

	var roundTripper http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.Insecure,
		},
	}

	if logging.IsDebugOrHigher() {
		roundTripper = logging.NewDebugTransport("Kibana", roundTripper)
	}

	httpClient := &http.Client{Transport: roundTripper}
	restyClient := resty.
		NewWithClient(httpClient).
		SetBaseURL(serverURL).
		SetHeader("kbn-xsrf", "true").
		SetHeader("Content-Type", "application/json").
		SetDisableWarn(true)

	if cfg.ApiKey != "" {
		restyClient.SetAuthScheme("ApiKey").SetAuthToken(cfg.ApiKey)
	} else {
		restyClient.SetBasicAuth(cfg.Username, cfg.Password)
	}

	client := Client{http: restyClient}
	return &client, nil
}

// ApiResponse defines the response model.
type ApiResponse[T any] struct {
	Contents []byte
	Response *resty.Response
	Output   T
	Error    *Error
}

// Body returns the response body as a string.
func (r ApiResponse[T]) Body() string {
	return string(r.Contents)
}

// StatusCode gets the HTTP status code.
func (r ApiResponse[T]) StatusCode() int {
	if r.Response != nil {
		return r.Response.StatusCode()
	}
	return 0
}

// EmptyResponse is a empty placeholder for 204 responses.
type EmptyResponse struct{}

// Error defines the model for error.
type Error struct {
	Error      *string  `json:"error,omitempty"`
	Message    *string  `json:"message,omitempty"`
	StatusCode *float32 `json:"statusCode,omitempty"`
}

// Perform an API request.
func doAPI[T any](c *Client, ctx context.Context, method string, path string, pathParams map[string]string, body any, queryParams any) (*ApiResponse[T], error) {
	req := c.http.R().SetContext(ctx)

	if pathParams != nil {
		if space, ok := pathParams["space"]; ok && space == "" {
			pathParams["space"] = "default"
		}
		req.SetPathParams(pathParams)
	}

	if body != nil {
		req.SetBody(body)
	}

	if queryParams != nil {
		values, err := query.Values(queryParams)
		if err != nil {
			return nil, err
		}
		req.SetQueryParamsFromValues(values)
	}

	resp, err := req.Execute(method, path)
	if err != nil {
		return nil, err
	}

	response := &ApiResponse[T]{
		Contents: resp.Body(),
		Response: resp,
	}

	if strings.Contains(resp.Header().Get("Content-Type"), "application/json") {
		switch resp.StatusCode() {
		case http.StatusOK, http.StatusNoContent:
			var dest T
			if err := json.Unmarshal(response.Contents, &dest); err != nil {
				return nil, err
			}
			response.Output = dest
		default:
			var dest Error
			if err := json.Unmarshal(response.Contents, &dest); err != nil {
				return nil, err
			}
			response.Error = &dest
		}
	}

	return response, nil
}
