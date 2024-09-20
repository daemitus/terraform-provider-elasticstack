package elasticsearch

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8/typedapi/transform/deletetransform"
	"github.com/elastic/go-elasticsearch/v8/typedapi/transform/gettransform"
	"github.com/elastic/go-elasticsearch/v8/typedapi/transform/puttransform"
)

func (c *Client) GetTransform(ctx context.Context, id string) (*ApiResponse[GetTransformResponse], error) {
	return doApiPtr(c.http.Transform.GetTransform().TransformId(id).Do(ctx))
}

func (c *Client) PutTransform(ctx context.Context, id string, req PutTransformRequest) (*ApiResponse[PutTransformResponse], error) {
	return doApiPtr(c.http.Transform.PutTransform(id).Request(&req).Do(ctx))
}

func (c *Client) DeleteTransform(ctx context.Context, id string) (*ApiResponse[DeleteTransformResponse], error) {
	return doApiPtr(c.http.Transform.DeleteTransform(id).Do(ctx))
}

type (
	GetTransformResponse    = gettransform.Response
	PutTransformRequest     = puttransform.Request
	PutTransformResponse    = puttransform.Response
	DeleteTransformResponse = deletetransform.Response
)
