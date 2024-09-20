package fleet

import (
	"context"
)

func (c *Client) ReadEnrollmentApiKeys(ctx context.Context, params *ReadEnrollmentApiKeysParams) (*ApiResponse[ReadEnrollmentApiKeysResponse], error) {
	return doAPI[ReadEnrollmentApiKeysResponse](
		c, ctx,
		"GET", "/enrollment_api_keys",
		nil, nil, params,
	)
}

type ReadEnrollmentApiKeysParams struct {
	PerPage *int    `url:"perPage,omitempty"`
	Page    *int    `url:"page,omitempty"`
	Kuery   *string `url:"kuery,omitempty"`
}

type ReadEnrollmentApiKeysResponse struct {
	Items   EnrollmentApiKeys `json:"items"`
	Page    float32           `json:"page"`
	PerPage float32           `json:"perPage"`
	Total   float32           `json:"total"`
}

type EnrollmentApiKeys []EnrollmentApiKey

type EnrollmentApiKey struct {
	Active    bool    `json:"active"`
	ApiKey    string  `json:"api_key"`
	ApiKeyId  string  `json:"api_key_id"`
	CreatedAt string  `json:"created_at"`
	Id        string  `json:"id"`
	Name      *string `json:"name,omitempty"`
	PolicyId  *string `json:"policy_id,omitempty"`
}
