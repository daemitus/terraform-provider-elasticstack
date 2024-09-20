package clients

import (
	"context"
	"net/http"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (c *ElasticsearchClient) GetClusterInfo(ctx context.Context, target string) (*elasticsearch.GetClusterInfoResponse, diag.Diagnostics) {
	resp, err := c.API.GetClusterInfo(ctx, target)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) GetClusterHealth(ctx context.Context) (*elasticsearch.GetClusterHealthResponse, diag.Diagnostics) {
	resp, err := c.API.GetClusterHealth(ctx)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *ElasticsearchClient) GetSlmPolicy(ctx context.Context, id string) (*types.SnapshotLifecycle, diag.Diagnostics) {
	resp, err := c.API.GetSlmPolicy(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneMapResponse(resp.Output, id, "SLM lifecycle", "SLM lifecycles")
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutSlmPolicy(ctx context.Context, id string, req elasticsearch.PutSlmPolicyRequest) (*types.SnapshotLifecycle, diag.Diagnostics) {
	resp, err := c.API.PutSlmPolicy(ctx, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetSlmPolicy(ctx, id)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteSlmPolicy(ctx context.Context, id string) diag.Diagnostics {
	resp, err := c.API.DeleteSlmPolicy(ctx, id)
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

func (c *ElasticsearchClient) GetSnapshot(ctx context.Context, repository string, snapshot string) (*elasticsearch.GetSnapshotResponse, diag.Diagnostics) {
	resp, err := c.API.GetSnapshot(ctx, repository, snapshot)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) CreateSnapshot(ctx context.Context, repository string, snapshot string, req elasticsearch.CreateSnapshotRequest) (*elasticsearch.CreateSnapshotResponse, diag.Diagnostics) {
	resp, err := c.API.CreateSnapshot(ctx, repository, snapshot, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteSnapshot(ctx context.Context, repository string, snapshot string) diag.Diagnostics {
	resp, err := c.API.DeleteSnapshot(ctx, repository, snapshot)
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

func (c *ElasticsearchClient) GetRepository(ctx context.Context, repository string) (*elasticsearch.GetRepositoryResponse, diag.Diagnostics) {
	resp, err := c.API.GetRepository(ctx, repository)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) CreateRepository(ctx context.Context, repository string, req elasticsearch.CreateRepositoryRequest) (*elasticsearch.CreateRepositoryResponse, diag.Diagnostics) {
	resp, err := c.API.CreateRepository(ctx, repository, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteRepository(ctx context.Context, repository string) diag.Diagnostics {
	resp, err := c.API.DeleteRepository(ctx, repository)
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

func (c *ElasticsearchClient) GetClusterSettings(ctx context.Context) (*elasticsearch.GetClusterSettingsResponse, diag.Diagnostics) {
	resp, err := c.API.GetClusterSettings(ctx)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutClusterSettings(ctx context.Context, req elasticsearch.PutClusterSettingsRequest) (*elasticsearch.GetClusterSettingsResponse, diag.Diagnostics) {
	resp, err := c.API.UpdateClusterSettings(ctx, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetClusterSettings(ctx)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *ElasticsearchClient) GetScript(ctx context.Context, id string) (*types.StoredScript, diag.Diagnostics) {
	resp, err := c.API.ReadScript(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Output.Script, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutScript(ctx context.Context, id string, req elasticsearch.PutScriptRequest) (*types.StoredScript, diag.Diagnostics) {
	resp, err := c.API.PutScript(ctx, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetScript(ctx, id)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteScript(ctx context.Context, id string) diag.Diagnostics {
	resp, err := c.API.DeleteScript(ctx, id)
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
