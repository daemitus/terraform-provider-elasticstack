package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/typedapi/cluster/deletecomponenttemplate"
	"github.com/elastic/go-elasticsearch/v8/typedapi/cluster/getcomponenttemplate"
	"github.com/elastic/go-elasticsearch/v8/typedapi/cluster/putcomponenttemplate"
	deleteilmpolicy "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/deletelifecycle"
	getilmpolicy "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/getlifecycle"
	putilmpolicy "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/putlifecycle"
	createindex "github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/createdatastream"
	deleteindex "github.com/elastic/go-elasticsearch/v8/typedapi/indices/delete"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/deletealias"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/deletedatastream"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/deleteindextemplate"
	getindex "github.com/elastic/go-elasticsearch/v8/typedapi/indices/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/getalias"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/getdatastream"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/getindextemplate"
	getindexmapping "github.com/elastic/go-elasticsearch/v8/typedapi/indices/getmapping"
	getindexsettings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/getsettings"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/putalias"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/putindextemplate"
	putindexmapping "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putmapping"
	putindexsettings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putsettings"
)

func (c *Client) GetIndexTemplate(ctx context.Context, name string) (*ApiResponse[GetIndexTemplateResponse], error) {
	return doApiPtr(c.http.Indices.GetIndexTemplate().Name(name).Do(ctx))
}

func (c *Client) PutIndexTemplate(ctx context.Context, name string, isCreate bool, req PutIndexTemplateRequest) (*ApiResponse[PutIndexTemplateResponse], error) {
	return doApiPtr(c.http.Indices.PutIndexTemplate(name).Create(isCreate).Request(&req).Do(ctx))
}

func (c *Client) DeleteIndexTemplate(ctx context.Context, name string) (*ApiResponse[DeleteIndexTemplateResponse], error) {
	return doApiPtr(c.http.Indices.DeleteIndexTemplate(name).Do(ctx))
}

type (
	GetIndexTemplateResponse    = getindextemplate.Response
	PutIndexTemplateRequest     = putindextemplate.Request
	PutIndexTemplateResponse    = putindextemplate.Response
	DeleteIndexTemplateResponse = deleteindextemplate.Response
)

// ============================================================================

func (c *Client) GetComponentTemplate(ctx context.Context, name string) (*ApiResponse[GetComponentTemplateResponse], error) {
	return doApiPtr(c.http.Cluster.GetComponentTemplate().Name(name).Do(ctx))
}

func (c *Client) PutComponentTemplate(ctx context.Context, name string, req PutComponentTemplateRequest) (*ApiResponse[PutComponentTemplateResponse], error) {
	return doApiPtr(c.http.Cluster.PutComponentTemplate(name).Request(&req).Do(ctx))
}

func (c *Client) DeleteComponentTemplate(ctx context.Context, name string) (*ApiResponse[DeleteComponentTemplateResponse], error) {
	return doApiPtr(c.http.Cluster.DeleteComponentTemplate(name).Do(ctx))
}

type (
	GetComponentTemplateResponse    = getcomponenttemplate.Response
	PutComponentTemplateRequest     = putcomponenttemplate.Request
	PutComponentTemplateResponse    = putcomponenttemplate.Response
	DeleteComponentTemplateResponse = deletecomponenttemplate.Response
)

// ============================================================================

func (c *Client) GetIlmPolicy(ctx context.Context, name string) (*ApiResponse[GetIlmPolicyResponse], error) {
	return doApi(c.http.Ilm.GetLifecycle().Policy(name).Do(ctx))
}

func (c *Client) PutIlmPolicy(ctx context.Context, name string, req PutIlmPolicyRequest) (*ApiResponse[PutIlmPolicyResponse], error) {
	return doApiPtr(c.http.Ilm.PutLifecycle(name).Request(&req).Do(ctx))
}

func (c *Client) DeleteIlmPolicy(ctx context.Context, name string) (*ApiResponse[DeleteIlmPolicyResponse], error) {
	return doApiPtr(c.http.Ilm.DeleteLifecycle(name).Do(ctx))
}

type (
	GetIlmPolicyResponse    = getilmpolicy.Response
	PutIlmPolicyRequest     = putilmpolicy.Request
	PutIlmPolicyResponse    = putilmpolicy.Response
	DeleteIlmPolicyResponse = deleteilmpolicy.Response
)

// ============================================================================

func (c *Client) GetDataStream(ctx context.Context, name string) (*ApiResponse[GetDataStreamResponse], error) {
	return doApiPtr(c.http.Indices.GetDataStream().Name(name).Do(ctx))
}

func (c *Client) CreateDataStream(ctx context.Context, name string) (*ApiResponse[CreateDataStreamResponse], error) {
	return doApiPtr(c.http.Indices.CreateDataStream(name).Do(ctx))
}

func (c *Client) DeleteDataStream(ctx context.Context, name string) (*ApiResponse[DeleteDataStreamResponse], error) {
	return doApiPtr(c.http.Indices.DeleteDataStream(name).Do(ctx))
}

type (
	GetDataStreamResponse    = getdatastream.Response
	CreateDataStreamResponse = createdatastream.Response
	DeleteDataStreamResponse = deletedatastream.Response
)

// ============================================================================

func (c *Client) GetIngestPipeline(ctx context.Context, id string) (*ApiResponse[GetIngestPipelineResponse], error) {
	return doApiResp[GetIngestPipelineResponse](c.http.Ingest.GetPipeline().Id(id).Perform(ctx))
}

func (c *Client) PutIngestPipeline(ctx context.Context, id string, req PutIngestPipelineRequest) (*ApiResponse[PutIngestPipelineResponse], error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return doApiResp[PutIngestPipelineResponse](c.http.Ingest.PutPipeline(id).Raw(bytes.NewReader(body)).Perform(ctx))
}

func (c *Client) DeleteIngestPipeline(ctx context.Context, id string) (*ApiResponse[DeleteIngestPipelineResponse], error) {
	return doApiResp[DeleteIngestPipelineResponse](c.http.Ingest.DeletePipeline(id).Perform(ctx))
}

type IngestPipeline struct {
	Description *string                    `json:"description,omitempty"`
	Meta_       map[string]json.RawMessage `json:"_meta,omitempty"`
	OnFailure   []map[string]any           `json:"on_failure,omitempty"`
	Processors  []map[string]any           `json:"processors,omitempty"`
	Version     *int64                     `json:"version,omitempty"`
}

type GetIngestPipelineResponse map[string]IngestPipeline

type PutIngestPipelineRequest IngestPipeline

type PutIngestPipelineResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

type DeleteIngestPipelineResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

// ============================================================================

func (c *Client) GetIndex(ctx context.Context, name string) (*ApiResponse[GetIndexResponse], error) {
	return doApi(c.http.Indices.Get(name).Do(ctx))
}

func (c *Client) CreateIndex(ctx context.Context, name string, req CreateIndexRequest) (*ApiResponse[CreateIndexResponse], error) {
	return doApiPtr(c.http.Indices.Create(name).Request(&req).Do(ctx))
}

func (c *Client) DeleteIndex(ctx context.Context, name string) (*ApiResponse[DeleteIndexResponse], error) {
	return doApiPtr(c.http.Indices.Delete(name).Do(ctx))
}

type (
	GetIndexResponse    = getindex.Response
	CreateIndexRequest  = createindex.Request
	CreateIndexResponse = createindex.Response
	DeleteIndexResponse = deleteindex.Response
)

// ============================================================================

func (c *Client) GetIndexAlias(ctx context.Context, index string, alias string) (*ApiResponse[GetAliasResponse], error) {
	return doApi(c.http.Indices.GetAlias().Index(index).Name(alias).Do(ctx))
}

func (c *Client) PutIndexAlias(ctx context.Context, index string, alias string, req PutAliasRequest) (*ApiResponse[PutAliasResponse], error) {
	return doApiPtr(c.http.Indices.PutAlias(index, alias).Request(&req).Do(ctx))
}

func (c *Client) DeleteIndexAlias(ctx context.Context, index string, alias string) (*ApiResponse[DeleteAliasResponse], error) {
	return doApiPtr(c.http.Indices.DeleteAlias(index, alias).Do(ctx))
}

type (
	GetAliasResponse    = getalias.Response
	PutAliasRequest     = putalias.Request
	PutAliasResponse    = putalias.Response
	DeleteAliasResponse = deletealias.Response
)

// ============================================================================

func (c *Client) GetIndexMapping(ctx context.Context, index string) (*ApiResponse[GetIndexMappingResponse], error) {
	return doApi(c.http.Indices.GetMapping().Index(index).Do(ctx))
}

func (c *Client) PutIndexMapping(ctx context.Context, index string, req PutIndexMappingRequest) (*ApiResponse[PutIndexMappingResponse], error) {
	return doApiPtr(c.http.Indices.PutMapping(index).Request(&req).Do(ctx))
}

type (
	GetIndexMappingResponse = getindexmapping.Response
	PutIndexMappingRequest  = putindexmapping.Request
	PutIndexMappingResponse = putindexmapping.Response
)

// ============================================================================

func (c *Client) GetIndexSettings(ctx context.Context, index string) (*ApiResponse[GetIndexSettingsResponse], error) {
	return doApi(c.http.Indices.GetSettings().Index(index).Do(ctx))
}

func (c *Client) PutIndexSettings(ctx context.Context, index string, req PutIndexSettingsRequest) (*ApiResponse[PutIndexSettingsResponse], error) {
	return doApiPtr(c.http.Indices.PutSettings().Indices(index).Request(&req).Do(ctx))
}

type (
	GetIndexSettingsResponse = getindexsettings.Response
	PutIndexSettingsRequest  = putindexsettings.Request
	PutIndexSettingsResponse = putindexsettings.Response
)
