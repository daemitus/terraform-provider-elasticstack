package fleet

import "context"

func (c *Client) HealthCheck(ctx context.Context) (*ApiResponse[HealthCheck], error) {
	return doAPI[HealthCheck](
		c, ctx,
		"GET", "/health_check",
		nil, nil, nil,
	)
}

type HealthCheck struct {
	Name   string `json:"name"`
	Host   string `json:"host"`
	Status string `json:"status"`
}
