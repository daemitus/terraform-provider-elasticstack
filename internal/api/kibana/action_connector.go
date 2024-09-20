package kibana

import (
	"context"
)

func (c *Client) ListConnectors(ctx context.Context, space string) (*ApiResponse[Connectors], error) {
	return doAPI[Connectors](
		c, ctx,
		"GET", "/s/{space}/api/actions/connectors",
		map[string]string{"space": space}, nil, nil,
	)
}

func (c *Client) ReadConnector(ctx context.Context, space string, id string) (*ApiResponse[Connector], error) {
	return doAPI[Connector](
		c, ctx,
		"GET", "/s/{space}/api/actions/connector/{id}",
		map[string]string{"id": id, "space": space},
		nil, nil,
	)
}

func (c *Client) CreateConnector(ctx context.Context, space string, body ConnectorRequest) (*ApiResponse[Connector], error) {
	return doAPI[Connector](
		c, ctx,
		"POST", "/s/{space}/api/actions/connector",
		map[string]string{"space": space},
		body, nil,
	)
}

func (c *Client) CreateConnectorWithID(ctx context.Context, space string, id string, body ConnectorRequest) (*ApiResponse[Connector], error) {
	return doAPI[Connector](
		c, ctx,
		"POST", "/s/{space}/api/actions/connector/{id}",
		map[string]string{"id": id, "space": space},
		body, nil,
	)
}

func (c *Client) UpdateConnector(ctx context.Context, space string, id string, body ConnectorRequest) (*ApiResponse[Connector], error) {
	return doAPI[Connector](
		c, ctx,
		"PUT", "/s/{space}/api/actions/connector/{id}",
		map[string]string{"id": id, "space": space},
		body, nil,
	)
}

func (c *Client) DeleteConnector(ctx context.Context, space string, id string) (*ApiResponse[EmptyResponse], error) {
	return doAPI[EmptyResponse](
		c, ctx,
		"DELETE", "/s/{space}/api/actions/connector/{id}",
		map[string]string{"id": id, "space": space},
		nil, nil,
	)
}

// ============================================================================

type ConnectorRequest struct {
	ConnectorTypeID string         `json:"connector_type_id"`
	Name            string         `json:"name"`
	Config          map[string]any `json:"config"`
	Secrets         map[string]any `json:"secrets"`
}

type Connectors []Connector

type Connector struct {
	ConnectorTypeID  string         `json:"connector_type_id"`
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	IsDeprecated     bool           `json:"is_deprecated"`
	IsMissingSecrets bool           `json:"is_missing_secrets,omitempty"`
	IsPreconfigured  bool           `json:"is_preconfigured"`
	IsSystemAction   bool           `json:"is_system_action,omitempty"`
	Config           map[string]any `json:"config"`
}
