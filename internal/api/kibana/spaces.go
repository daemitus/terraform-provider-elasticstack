package kibana

import (
	"context"
)

func (c *Client) ListSpaces(ctx context.Context) (*ApiResponse[Spaces], error) {
	return doAPI[Spaces](
		c, ctx,
		"GET", "/api/spaces/space",
		nil, nil, nil,
	)
}

func (c *Client) CreateSpace(ctx context.Context, body Space) (*ApiResponse[Space], error) {
	return doAPI[Space](
		c, ctx,
		"POST", "/api/spaces/space",
		nil, body, nil,
	)
}

func (c *Client) ReadSpace(ctx context.Context, id string) (*ApiResponse[Space], error) {
	return doAPI[Space](
		c, ctx,
		"GET", "/api/spaces/space/{id}",
		map[string]string{"id": id},
		nil, nil,
	)
}

func (c *Client) UpdateSpace(ctx context.Context, id string, body Space) (*ApiResponse[Space], error) {
	return doAPI[Space](
		c, ctx,
		"PUT", "/api/spaces/space/{id}",
		map[string]string{"id": id},
		body, nil,
	)
}

func (c *Client) DeleteSpace(ctx context.Context, id string) (*ApiResponse[EmptyResponse], error) {
	return doAPI[EmptyResponse](
		c, ctx,
		"DELETE", "/api/spaces/space/{id}",
		map[string]string{"id": id},
		nil, nil,
	)
}

// ============================================================================

type Spaces []Space

type Space struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	DisabledFeatures []string `json:"disabledFeatures,omitempty"`
	Reserved         bool     `json:"_reserved,omitempty"`
	Initials         string   `json:"initials,omitempty"`
	Color            string   `json:"color,omitempty"`
	ImageURL         string   `json:"imageUrl,omitempty"`
}
