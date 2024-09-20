package elasticsearch

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/daemitus/terraform-provider-elasticstack/internal/config"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util/logging"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// Client contain the REST client.
type Client struct {
	http *elasticsearch.TypedClient
}

// NewClient creates a new Client.
func NewClient(config config.ServiceConfig) (*Client, error) {
	esConfig := elasticsearch.Config{
		Addresses: []string{config.Endpoint},
		APIKey:    config.ApiKey,
		Username:  config.Username,
		Password:  config.Password,
	}

	if esConfig.Transport == nil {
		esConfig.Transport = http.DefaultTransport.(*http.Transport)
	}
	if esConfig.Transport.(*http.Transport).TLSClientConfig == nil {
		esConfig.Transport.(*http.Transport).TLSClientConfig = &tls.Config{}
	}

	if config.Insecure {
		esConfig.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = config.Insecure
	}

	if logging.IsDebugOrHigher() {
		esConfig.Transport = logging.NewDebugTransport("Elasticsearch", esConfig.Transport)
	}

	esClient, err := elasticsearch.NewTypedClient(esConfig)
	if err != nil {
		return nil, err
	}

	client := Client{http: esClient}
	return &client, nil
}

// ============================================================================

// ApiResponse defines the response model.
type ApiResponse[T any] struct {
	Output T
	Error  *types.ElasticsearchError
	code   int
}

// Body returns the response body as a string.
func (r ApiResponse[T]) Body() string {
	return r.Error.Error()
}

// StatusCode gets the HTTP status code.
func (r ApiResponse[T]) StatusCode() int {
	return r.code
}

// Perform an API request.
func doApi[T any](resp T, err error) (*ApiResponse[T], error) {
	if err != nil {
		if e, ok := err.(*types.ElasticsearchError); ok {
			response := &ApiResponse[T]{
				Error: e,
				code:  e.Status,
			}
			return response, nil
		}
		return nil, err
	}

	response := &ApiResponse[T]{Output: resp, code: 200}
	return response, nil
}

func doApiPtr[T any](resp *T, err error) (*ApiResponse[T], error) {
	if err != nil {
		return nil, err
	}

	return doApi(*resp, err)
}

func doApiResp[T any](resp *http.Response, err error) (*ApiResponse[T], error) {
	if err != nil {
		return nil, err
	}

	response := &ApiResponse[T]{
		code: resp.StatusCode,
	}

	if strings.Contains(resp.Header.Get("Content-Type"), "application/vnd.elasticsearch+json") {
		switch resp.StatusCode {
		case http.StatusOK, http.StatusNoContent:
			var dest T
			bytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(bytes, &dest); err != nil {
				return nil, err
			}
			response.Output = dest
		default:
			var dest types.ElasticsearchError
			bytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(bytes, &dest); err != nil {
				return nil, err
			}
			response.Error = &dest
		}
	}

	return response, nil
}
