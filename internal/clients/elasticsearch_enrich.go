package clients

import (
	"context"
	"net/http"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (c *ElasticsearchClient) GetEnrichPolicy(ctx context.Context, name string) (*types.Summary, diag.Diagnostics) {
	resp, err := c.API.GetEnrichPolicy(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneSliceResponse(resp.Output.Policies, name, "enrich policy", "enrich policies")
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutEnrichPolicy(ctx context.Context, name string, req elasticsearch.PutEnrichPolicyRequest) (*types.Summary, diag.Diagnostics) {
	resp, err := c.API.PutEnrichPolicy(ctx, name, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetEnrichPolicy(ctx, name)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteEnrichPolicy(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteEnrichPolicy(ctx, name)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}
