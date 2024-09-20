package elasticsearch

import (
	"bytes"
	"context"

	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/elastic/go-elasticsearch/v8/typedapi/watcher/deletewatch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/watcher/putwatch"
)

func (c *Client) GetWatch(ctx context.Context, id string) (*ApiResponse[GetWatchResponse], error) {
	return doApiResp[GetWatchResponse](c.http.Watcher.GetWatch(id).Perform(ctx))
}

func (c *Client) PutWatch(ctx context.Context, id string, active bool, req PutWatchRequest) (*ApiResponse[PutWatchResponse], error) {
	val, err := util.JsonMarshal(req)
	if err != nil {
		return nil, err
	}
	return doApiPtr(c.http.Watcher.PutWatch(id).Raw(bytes.NewBuffer(val)).Active(active).Do(ctx))
}

func (c *Client) DeleteWatch(ctx context.Context, id string) (*ApiResponse[DeleteWatchResponse], error) {
	return doApiPtr(c.http.Watcher.DeleteWatch(id).Do(ctx))
}

type GetWatchResponse struct {
	ID     string       `json:"_id"`
	Status *WatchStatus `json:"status,omitempty"`
	Watch  *Watch       `json:"watch,omitempty"`
}

type WatchStatus struct {
	State WatchState `json:"state"`
}

type WatchState struct {
	Active bool `json:"active"`
}

type Watch struct {
	Trigger        map[string]any `json:"trigger"`
	Input          map[string]any `json:"input"`
	Condition      map[string]any `json:"condition"`
	Actions        map[string]any `json:"actions"`
	Metadata       map[string]any `json:"metadata"`
	Transform      map[string]any `json:"transform,omitempty"`
	ThrottlePeriod *int64         `json:"throttle_period_in_millis,omitempty"`
}

type (
	PutWatchRequest     = Watch
	PutWatchResponse    = putwatch.Response
	DeleteWatchResponse = deletewatch.Response
)
