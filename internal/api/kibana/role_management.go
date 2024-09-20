package kibana

import (
	"context"
)

func (c *Client) ListRoles(ctx context.Context) (*ApiResponse[Roles], error) {
	return doAPI[Roles](
		c, ctx,
		"GET", "/api/security/role",
		nil, nil, nil,
	)
}

// ============================================================================

func (c *Client) PutRole(ctx context.Context, name string, body Role, params *PutRoleParams) (*ApiResponse[PutRoleResponse], error) {
	return doAPI[PutRoleResponse](
		c, ctx,
		"PUT", "/api/security/role/{name}",
		map[string]string{"name": name},
		body, params,
	)
}

type PutRoleParams struct {
	CreateOnly bool `url:"createOnly,omitempty"`
}

type PutRoleResponse struct {
	Role struct {
		Created bool `json:"created"`
	} `json:"role"`
}

// ============================================================================

func (c *Client) ReadRole(ctx context.Context, name string) (*ApiResponse[Role], error) {
	return doAPI[Role](
		c, ctx,
		"GET", "/api/security/role/{name}",
		map[string]string{"name": name},
		nil, nil,
	)
}

// ============================================================================

func (c *Client) DeleteRole(ctx context.Context, name string) (*ApiResponse[DeleteRoleResponse], error) {
	return doAPI[DeleteRoleResponse](
		c, ctx,
		"DELETE", "/api/security/role/{name}",
		map[string]string{"name": name},
		nil, nil,
	)
}

type DeleteRoleResponse struct {
	Found bool `json:"found"`
}

// ============================================================================

type Roles []Role

type Role struct {
	Name            *string                `json:"name,omitempty"`
	Metadata        map[string]any         `json:"metadata,omitempty"`
	TransientMedata *RoleTransientMetadata `json:"transient_metadata,omitempty"`
	Elasticsearch   *RoleElasticsearch     `json:"elasticsearch,omitempty"`
	Kibana          []RoleKibana           `json:"kibana,omitempty"`
}

type RoleTransientMetadata struct {
	Enabled bool `json:"enabled,omitempty"`
}

type RoleElasticsearch struct {
	Clusters      []string                       `json:"cluster,omitempty"`
	Indices       []RoleElasticsearchIndex       `json:"indices,omitempty"`
	RemoteIndices []RoleElasticsearchRemoteIndex `json:"remote_indices,omitempty"`
	RunAs         []string                       `json:"run_as,omitempty"`
}

type RoleKibana struct {
	Bases    []string            `json:"base,omitempty"`
	Features map[string][]string `json:"feature,omitempty"`
	Spaces   []string            `json:"spaces,omitempty"`
}

type RoleElasticsearchIndex struct {
	AllowRestrictedIndices *bool                                `json:"allow_restricted_indices,omitempty"`
	Names                  []string                             `json:"names,omitempty"`
	Privileges             []string                             `json:"privileges,omitempty"`
	FieldSecurity          *RoleElasticsearchIndexFieldSecurity `json:"field_security,omitempty"`
	Query                  *string                              `json:"query,omitempty"`
}

type RoleElasticsearchRemoteIndex struct {
	Clusters               []string                             `json:"clusters,omitempty"`
	AllowRestrictedIndices *bool                                `json:"allow_restricted_indices,omitempty"`
	Names                  []string                             `json:"names,omitempty"`
	Privileges             []string                             `json:"privileges,omitempty"`
	FieldSecurity          *RoleElasticsearchIndexFieldSecurity `json:"field_security,omitempty"`
	Query                  *string                              `json:"query,omitempty"`
}

type RoleElasticsearchIndexFieldSecurity struct {
	Grants  []string `json:"grant,omitempty"`
	Excepts []string `json:"excepts,omitempty"`
}
