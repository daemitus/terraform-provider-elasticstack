package clients

import (
	"context"
	"net/http"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (c *ElasticsearchClient) GetApiKey(ctx context.Context, id string) (*types.ApiKey, diag.Diagnostics) {
	resp, err := c.API.GetApiKey(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneSliceResponse(resp.Output.ApiKeys, id, "API key", "API keys")
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) CreateApiKey(ctx context.Context, req elasticsearch.CreateApiKeyRequest) (*types.ApiKey, diag.Diagnostics) {
	resp, err := c.API.CreateApiKey(ctx, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetApiKey(ctx, resp.Output.Id)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteApiKey(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteApiKey(ctx, name)
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

func (c *ElasticsearchClient) GetRole(ctx context.Context, name string) (*types.Role, diag.Diagnostics) {
	resp, err := c.API.GetRole(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneMapResponse(resp.Output, name, "role", "roles")
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutRole(ctx context.Context, name string, req elasticsearch.PutRoleRequest) (*types.Role, diag.Diagnostics) {
	resp, err := c.API.PutRole(ctx, name, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetRole(ctx, name)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteRole(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteRole(ctx, name)
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

func (c *ElasticsearchClient) GetRoleMapping(ctx context.Context, name string) (*types.SecurityRoleMapping, diag.Diagnostics) {
	resp, err := c.API.GetRoleMapping(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneMapResponse(resp.Output, name, "role mapping", "role mappings")
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutRoleMapping(ctx context.Context, name string, req elasticsearch.PutRoleMappingRequest) (*types.SecurityRoleMapping, diag.Diagnostics) {
	resp, err := c.API.PutRoleMapping(ctx, name, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetRoleMapping(ctx, name)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteRoleMapping(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteRoleMapping(ctx, name)
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

func (c *ElasticsearchClient) GetUser(ctx context.Context, username string) (*types.User, diag.Diagnostics) {
	resp, err := c.API.GetUser(ctx, username)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return getOneMapResponse(resp.Output, username, "user", "users")
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) PutUser(ctx context.Context, username string, req elasticsearch.PutUserRequest) (*types.User, diag.Diagnostics) {
	resp, err := c.API.PutUser(ctx, username, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return c.GetUser(ctx, username)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *ElasticsearchClient) DeleteUser(ctx context.Context, username string) diag.Diagnostics {
	resp, err := c.API.DeleteUser(ctx, username)
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
