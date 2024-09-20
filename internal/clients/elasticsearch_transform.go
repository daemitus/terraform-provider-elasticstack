package clients

import (
	"context"
	"net/http"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (c *ElasticsearchClient) GetTransform(ctx context.Context, id string) (*types.TransformSummary, diag.Diagnostics) {
	resp, err := c.API.GetTransform(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneSliceResponse(resp.Output.Transforms, id, "transform", "transforms")
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutTransform(ctx context.Context, id string, req elasticsearch.PutTransformRequest) (*types.TransformSummary, diag.Diagnostics) {
	resp, err := c.API.PutTransform(ctx, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetTransform(ctx, id)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteTransform(ctx context.Context, id string) diag.Diagnostics {
	resp, err := c.API.DeleteTransform(ctx, id)
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
