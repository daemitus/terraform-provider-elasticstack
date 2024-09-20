package clients

import (
	"context"
	"net/http"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (c *ElasticsearchClient) GetIndexTemplate(ctx context.Context, name string) (*types.IndexTemplateItem, diag.Diagnostics) {
	resp, err := c.API.GetIndexTemplate(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneSliceResponse(resp.Output.IndexTemplates, name, "index template", "index templates")
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutIndexTemplate(ctx context.Context, name string, isCreate bool, req elasticsearch.PutIndexTemplateRequest) (*types.IndexTemplateItem, diag.Diagnostics) {
	resp, err := c.API.PutIndexTemplate(ctx, name, isCreate, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetIndexTemplate(ctx, name)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteIndexTemplate(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteIndexTemplate(ctx, name)
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

// ============================================================================

func (c *ElasticsearchClient) GetComponentTemplate(ctx context.Context, name string) (*types.ClusterComponentTemplate, diag.Diagnostics) {
	resp, err := c.API.GetComponentTemplate(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneSliceResponse(resp.Output.ComponentTemplates, name, "component template", "component templates")
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutComponentTemplate(ctx context.Context, name string, req elasticsearch.PutComponentTemplateRequest) (*types.ClusterComponentTemplate, diag.Diagnostics) {
	resp, err := c.API.PutComponentTemplate(ctx, name, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetComponentTemplate(ctx, name)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteComponentTemplate(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteComponentTemplate(ctx, name)
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

// ============================================================================

func (c *ElasticsearchClient) GetIlmPolicy(ctx context.Context, name string) (*types.Lifecycle, diag.Diagnostics) {
	resp, err := c.API.GetIlmPolicy(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneMapResponse(resp.Output, name, "ILM policy", "ILM policies")
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutIlmPolicy(ctx context.Context, name string, req elasticsearch.PutIlmPolicyRequest) (*types.Lifecycle, diag.Diagnostics) {
	resp, err := c.API.PutIlmPolicy(ctx, name, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetIlmPolicy(ctx, name)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteIlmPolicy(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteIlmPolicy(ctx, name)
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

// ============================================================================

func (c *ElasticsearchClient) GetDataStream(ctx context.Context, name string) (*types.DataStream, diag.Diagnostics) {
	resp, err := c.API.GetDataStream(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneSliceResponse(resp.Output.DataStreams, name, "datastream", "datastreams")
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) CreateDataStream(ctx context.Context, name string) (*types.DataStream, diag.Diagnostics) {
	resp, err := c.API.CreateDataStream(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetDataStream(ctx, name)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteDataStream(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteDataStream(ctx, name)
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

// ============================================================================

func (c *ElasticsearchClient) GetIngestPipeline(ctx context.Context, id string) (*elasticsearch.IngestPipeline, diag.Diagnostics) {
	resp, err := c.API.GetIngestPipeline(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneMapResponse(resp.Output, id, "ingest pipeline", "ingest pipelines")
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutIngestPipeline(ctx context.Context, id string, req elasticsearch.PutIngestPipelineRequest) (*elasticsearch.IngestPipeline, diag.Diagnostics) {
	resp, err := c.API.PutIngestPipeline(ctx, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetIngestPipeline(ctx, id)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteIngestPipeline(ctx context.Context, id string) diag.Diagnostics {
	resp, err := c.API.DeleteIngestPipeline(ctx, id)
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

// ============================================================================

func (c *ElasticsearchClient) GetIndex(ctx context.Context, name string) (*types.IndexState, diag.Diagnostics) {
	resp, err := c.API.GetIndex(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneMapResponse(resp.Output, name, "index", "indexes")
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) CreateIndex(ctx context.Context, name string, req elasticsearch.CreateIndexRequest) (*types.IndexState, diag.Diagnostics) {
	resp, err := c.API.CreateIndex(ctx, name, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetIndex(ctx, name)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteIndex(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteIndex(ctx, name)
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

// ============================================================================

func (c *ElasticsearchClient) GetIndexAlias(ctx context.Context, index string, alias string) (*types.AliasDefinition, diag.Diagnostics) {
	resp, err := c.API.GetIndexAlias(ctx, index, alias)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		index, diags := getOneMapResponse(resp.Output, index, "index alias", "index aliases")
		if diags.HasError() {
			return nil, diags
		}
		return getOneMapResponse(index.Aliases, alias, "index alias", "index aliases")
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutIndexAlias(ctx context.Context, index string, alias string, req elasticsearch.PutAliasRequest) (*types.AliasDefinition, diag.Diagnostics) {
	resp, err := c.API.PutIndexAlias(ctx, index, alias, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetIndexAlias(ctx, index, alias)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteIndexAlias(ctx context.Context, name string, alias string) diag.Diagnostics {
	resp, err := c.API.DeleteIndexAlias(ctx, name, alias)
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

// ============================================================================

func (c *ElasticsearchClient) GetIndexMappings(ctx context.Context, index string) (*types.IndexMappingRecord, diag.Diagnostics) {
	resp, err := c.API.GetIndexMapping(ctx, index)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneMapResponse(resp.Output, index, "index", "indexes")
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutIndexMapping(ctx context.Context, index string, req elasticsearch.PutIndexMappingRequest) (*types.IndexMappingRecord, diag.Diagnostics) {
	resp, err := c.API.PutIndexMapping(ctx, index, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetIndexMappings(ctx, index)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *ElasticsearchClient) GetIndexSettings(ctx context.Context, index string) (*types.IndexState, diag.Diagnostics) {
	resp, err := c.API.GetIndexSettings(ctx, index)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneMapResponse(resp.Output, index, "index", "indexes")
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutIndexSettings(ctx context.Context, index string, req elasticsearch.PutIndexSettingsRequest) (*types.IndexState, diag.Diagnostics) {
	resp, err := c.API.PutIndexSettings(ctx, index, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetIndexSettings(ctx, index)
	default:
		return nil, c.reportUnknownError(resp)
	}
}
