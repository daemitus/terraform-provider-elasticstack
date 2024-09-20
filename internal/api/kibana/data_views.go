package kibana

import (
	"context"
)

func (c *Client) ListDataViews(ctx context.Context, space string) (*ApiResponse[DataViews], error) {
	return doAPI[DataViews](
		c, ctx,
		"GET", "/s/{space}/api/data_views",
		map[string]string{"space": space},
		nil, nil,
	)
}

func (c *Client) CreateDataView(ctx context.Context, space string, body DataView) (*ApiResponse[DataView], error) {
	return doAPI[DataView](
		c, ctx,
		"POST", "/s/{space}/api/data_views/data_view",
		map[string]string{"space": space},
		body, nil,
	)
}

func (c *Client) ReadDataView(ctx context.Context, space string, id string) (*ApiResponse[DataView], error) {
	return doAPI[DataView](
		c, ctx,
		"GET", "/s/{space}/api/data_views/data_view/{id}",
		map[string]string{"id": id, "space": space},
		nil, nil,
	)
}

func (c *Client) UpdateDataView(ctx context.Context, space string, id string, body DataView) (*ApiResponse[DataView], error) {
	return doAPI[DataView](
		c, ctx,
		"PUT", "/s/{space}/api/data_views/data_view/{id}",
		map[string]string{"id": id, "space": space},
		body, nil,
	)
}

func (c *Client) DeleteDataView(ctx context.Context, space string, id string) (*ApiResponse[EmptyResponse], error) {
	return doAPI[EmptyResponse](
		c, ctx,
		"DELETE", "/s/{space}/api/data_views/data_view/{id}",
		map[string]string{"id": id, "space": space},
		nil, nil,
	)
}

// ============================================================================

type DataViews struct {
	DataViews []ListDataView `json:"data_view"`
}

type ListDataView struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Namespaces []string       `json:"namespaces"`
	Title      string         `json:"title"`
	TypeMeta   map[string]any `json:"typeMeta,omitempty"`
}

// ============================================================================

type DataView struct {
	DataView      DataViewInternal `json:"data_view"`
	Override      *bool            `json:"override,omitempty"`       // Create
	RefreshFields *bool            `json:"refresh_fields,omitempty"` // Update
}

type DataViewInternal struct {
	AllowNoIndex    *bool                           `json:"allowNoIndex,omitempty"`
	FieldAttrs      map[string]DataViewFieldAttr    `json:"fieldAttrs,omitempty"`
	FieldFormats    map[string]DataViewFieldFormat  `json:"fieldFormats,omitempty"`
	Fields          map[string]map[string]any       `json:"fields,omitempty"`
	ID              *string                         `json:"id,omitempty"`
	Name            *string                         `json:"name,omitempty"`
	Namespaces      []string                        `json:"namespaces,omitempty"`
	RuntimeFieldMap map[string]DataViewRuntimeField `json:"runtimeFieldMap,omitempty"`
	SourceFilters   []DataViewSourceFilter          `json:"sourceFilters,omitempty"`
	TimeFieldName   *string                         `json:"timeFieldName,omitempty"`
	Title           string                          `json:"title"`
	Type            *string                         `json:"type,omitempty"`
	TypeMeta        *DataViewTypeMeta               `json:"typeMeta,omitempty"`
	Version         *string                         `json:"version,omitempty"`
}

type DataViewFieldAttr struct {
	CustomDescription *string `json:"customDescription,omitempty"`
	CustomLabel       *string `json:"customLabel,omitempty"`
	Count             *int64  `json:"count,omitempty"`
}

type DataViewFieldFormat struct {
	ID     string                     `json:"id"`
	Params *DataViewFieldFormatParams `json:"params,omitempty"`
}

type DataViewFieldFormatParams struct {
	Pattern *string `json:"pattern,omitempty"`
}

type DataViewRuntimeField struct {
	Type   string                     `json:"type"`
	Script DataViewRuntimeFieldScript `json:"script"`
}

type DataViewRuntimeFieldScript struct {
	Source string `json:"source"`
}

type DataViewSourceFilter struct {
	Value string `json:"value"`
}

type DataViewTypeMeta struct {
	Aggs   map[string]any `json:"aggs,omitempty"`
	Params map[string]any `json:"params,omitempty"`
}
