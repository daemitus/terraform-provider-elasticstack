package elasticsearch

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8/typedapi/enrich/deletepolicy"
	"github.com/elastic/go-elasticsearch/v8/typedapi/enrich/getpolicy"
	"github.com/elastic/go-elasticsearch/v8/typedapi/enrich/putpolicy"
)

func (c *Client) GetEnrichPolicy(ctx context.Context, name string) (*ApiResponse[GetEnrichPolicyResponse], error) {
	return doApiPtr(c.http.Enrich.GetPolicy().Name(name).Do(ctx))
}

func (c *Client) PutEnrichPolicy(ctx context.Context, name string, req PutEnrichPolicyRequest) (*ApiResponse[PutEnrichPolicyResponse], error) {
	return doApiPtr(c.http.Enrich.PutPolicy(name).Request(&req).Do(ctx))
}

func (c *Client) DeleteEnrichPolicy(ctx context.Context, name string) (*ApiResponse[DeleteEnrichPolicyResponse], error) {
	return doApiPtr(c.http.Enrich.DeletePolicy(name).Do(ctx))
}

type (
	GetEnrichPolicyResponse    = getpolicy.Response
	PutEnrichPolicyRequest     = putpolicy.Request
	PutEnrichPolicyResponse    = putpolicy.Response
	DeleteEnrichPolicyResponse = deletepolicy.Response
)
