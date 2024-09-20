package fleet

import (
	"context"
)

func (c *Client) ReadFleetServerHosts(ctx context.Context, itemId string) (*ApiResponse[ReadFleetServerHostsResponse], error) {
	return doAPI[ReadFleetServerHostsResponse](
		c, ctx,
		"GET", "/fleet_server_hosts/{id}",
		map[string]string{"id": itemId},
		nil, nil,
	)
}

type ReadFleetServerHostsResponse struct {
	Item FleetServerHost `json:"item"`
}

// ============================================================================

func (c *Client) CreateFleetServerHosts(ctx context.Context, body CreateFleetServerHostsRequest) (*ApiResponse[CreateFleetServerHostsResponse], error) {
	return doAPI[CreateFleetServerHostsResponse](
		c, ctx,
		"POST", "/fleet_server_hosts",
		nil, body, nil,
	)
}

type CreateFleetServerHostsRequest struct {
	HostUrls  []string `json:"host_urls"`
	Id        *string  `json:"id,omitempty"`
	IsDefault *bool    `json:"is_default,omitempty"`
	Name      string   `json:"name"`
}

type CreateFleetServerHostsResponse struct {
	Item *FleetServerHost `json:"item,omitempty"`
}

// ============================================================================

func (c *Client) UpdateFleetServerHosts(ctx context.Context, itemId string, body UpdateFleetServerHostsRequest) (*ApiResponse[UpdateFleetServerHostsResponse], error) {
	return doAPI[UpdateFleetServerHostsResponse](
		c, ctx,
		"PUT", "/fleet_server_hosts/{id}",
		map[string]string{"id": itemId},
		body, nil,
	)
}

type UpdateFleetServerHostsRequest struct {
	HostUrls  *[]string `json:"host_urls,omitempty"`
	IsDefault *bool     `json:"is_default,omitempty"`
	Name      *string   `json:"name,omitempty"`
}

type UpdateFleetServerHostsResponse struct {
	Item FleetServerHost `json:"item"`
}

// ============================================================================

func (c *Client) DeleteFleetServerHosts(ctx context.Context, itemId string) (*ApiResponse[DeleteFleetServerHostsResponse], error) {
	return doAPI[DeleteFleetServerHostsResponse](
		c, ctx,
		"DELETE", "/fleet_server_hosts/{id}",
		map[string]string{"id": itemId},
		nil, nil,
	)
}

type DeleteFleetServerHostsResponse struct {
	Id string `json:"id"`
}

// ============================================================================

type FleetServerHost struct {
	HostUrls        []string `json:"host_urls"`
	Id              string   `json:"id"`
	IsDefault       bool     `json:"is_default"`
	IsPreconfigured bool     `json:"is_preconfigured"`
	Name            *string  `json:"name,omitempty"`
}
