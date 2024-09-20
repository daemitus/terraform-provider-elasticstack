package elasticsearch

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8/typedapi/security/createapikey"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/deleterole"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/deleterolemapping"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/deleteuser"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/getapikey"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/getrole"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/getrolemapping"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/getuser"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/invalidateapikey"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/putrole"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/putrolemapping"
	"github.com/elastic/go-elasticsearch/v8/typedapi/security/putuser"
)

func (c *Client) GetApiKey(ctx context.Context, id string) (*ApiResponse[GetApiKeyResponse], error) {
	return doApiPtr(c.http.Security.GetApiKey().Id(id).Do(ctx))
}

func (c *Client) CreateApiKey(ctx context.Context, req CreateApiKeyRequest) (*ApiResponse[CreateApiKeyResponse], error) {
	return doApiPtr(c.http.Security.CreateApiKey().Request(&req).Do(ctx))
}

func (c *Client) DeleteApiKey(ctx context.Context, id string) (*ApiResponse[DeleteApiKeyResponse], error) {
	return doApiPtr(c.http.Security.InvalidateApiKey().Id(id).Do(ctx))
}

type (
	GetApiKeyResponse    = getapikey.Response
	CreateApiKeyRequest  = createapikey.Request
	CreateApiKeyResponse = createapikey.Response
	DeleteApiKeyResponse = invalidateapikey.Response
)

// ============================================================================

func (c *Client) GetRole(ctx context.Context, name string) (*ApiResponse[GetRoleResponse], error) {
	return doApi(c.http.Security.GetRole().Name(name).Do(ctx))
}

func (c *Client) PutRole(ctx context.Context, name string, req PutRoleRequest) (*ApiResponse[PutRoleResponse], error) {
	return doApiPtr(c.http.Security.PutRole(name).Request(&req).Do(ctx))
}

func (c *Client) DeleteRole(ctx context.Context, name string) (*ApiResponse[DeleteRoleResponse], error) {
	return doApiPtr(c.http.Security.DeleteRole(name).Do(ctx))
}

type (
	GetRoleResponse    = getrole.Response
	PutRoleRequest     = putrole.Request
	PutRoleResponse    = putrole.Response
	DeleteRoleResponse = deleterole.Response
)

// ============================================================================

func (c *Client) GetRoleMapping(ctx context.Context, name string) (*ApiResponse[GetRoleMappingResponse], error) {
	return doApi(c.http.Security.GetRoleMapping().Name(name).Do(ctx))
}

func (c *Client) PutRoleMapping(ctx context.Context, name string, req PutRoleMappingRequest) (*ApiResponse[PutRoleMappingResponse], error) {
	return doApiPtr(c.http.Security.PutRoleMapping(name).Request(&req).Do(ctx))
}

func (c *Client) DeleteRoleMapping(ctx context.Context, name string) (*ApiResponse[DeleteRoleMappingResponse], error) {
	return doApiPtr(c.http.Security.DeleteRoleMapping(name).Do(ctx))
}

type (
	GetRoleMappingResponse    = getrolemapping.Response
	PutRoleMappingRequest     = putrolemapping.Request
	PutRoleMappingResponse    = putrolemapping.Response
	DeleteRoleMappingResponse = deleterolemapping.Response
)

// ============================================================================

func (c *Client) GetUser(ctx context.Context, username string) (*ApiResponse[GetUserResponse], error) {
	return doApi(c.http.Security.GetUser().Username(username).Do(ctx))
}

func (c *Client) PutUser(ctx context.Context, username string, req PutUserRequest) (*ApiResponse[PutUserResponse], error) {
	return doApiPtr(c.http.Security.PutUser(username).Request(&req).Do(ctx))
}

func (c *Client) DeleteUser(ctx context.Context, username string) (*ApiResponse[DeleteUserResponse], error) {
	return doApiPtr(c.http.Security.DeleteUser(username).Do(ctx))
}

type (
	GetUserResponse    = getuser.Response
	PutUserRequest     = putuser.Request
	PutUserResponse    = putuser.Response
	DeleteUserResponse = deleteuser.Response
)
