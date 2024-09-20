package clients

import (
	"context"
	"net/http"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (c *ElasticsearchClient) GetWatch(ctx context.Context, id string) (*elasticsearch.GetWatchResponse, diag.Diagnostics) {
	resp, err := c.API.GetWatch(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutWatch(ctx context.Context, id string, active bool, req elasticsearch.PutWatchRequest) (*elasticsearch.GetWatchResponse, diag.Diagnostics) {
	resp, err := c.API.PutWatch(ctx, id, active, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetWatch(ctx, id)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteWatch(ctx context.Context, id string) diag.Diagnostics {
	resp, err := c.API.DeleteWatch(ctx, id)
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
