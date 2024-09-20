package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &IngestPipelineResource{}
	_ resource.ResourceWithImportState = &IngestPipelineResource{}
)

func NewIngestPipelineResource(client *clients.ElasticsearchClient) *IngestPipelineResource {
	return &IngestPipelineResource{client: client}
}

type IngestPipelineResource struct {
	client *clients.ElasticsearchClient
}

func (r *IngestPipelineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "elasticsearch_ingest_pipeline")
}

func (r *IngestPipelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages tasks and resources related to ingest pipelines and processors. See: https://www.elastic.co/guide/en/elasticsearch/reference/current/ingest-apis.html",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Internal identifier of the resource",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the ingest pipeline.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the ingest pipeline.",
				Optional:    true,
			},
			"on_failure": schema.ListAttribute{
				Description: "Processors to run immediately after a processor failure. Each processor supports a processor-level `on_failure` value. If a processor without an `on_failure` value fails, Elasticsearch uses this pipeline-level parameter as a fallback. The processors in this parameter run sequentially in the order specified. Elasticsearch will not attempt to run the pipeline’s remaining processors. See: https://www.elastic.co/guide/en/elasticsearch/reference/current/processors.html. Each record must be a valid JSON document",
				Optional:    true,
				ElementType: jsontypes.NormalizedType{},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"processors": schema.ListAttribute{
				Description: "Processors used to perform transformations on documents before indexing. Processors run sequentially in the order specified. See: https://www.elastic.co/guide/en/elasticsearch/reference/current/processors.html. Each record must be a valid JSON document.",
				Required:    true,
				ElementType: jsontypes.NormalizedType{},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"metadata": schema.StringAttribute{
				Description: "Optional user metadata about the index template.",
				CustomType:  jsontypes.NormalizedType{},
				Optional:    true,
			},
		},
	}
}

func (r *IngestPipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IngestPipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ingestPipelineModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	var name string
	if data.ID.IsUnknown() {
		name = data.Name.ValueString()
	} else {
		name = data.ID.ValueString()
	}

	pipeline, diags := r.client.GetIngestPipeline(ctx, name)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = data.fromApi(pipeline)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IngestPipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ingestPipelineModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	putReq, diags := data.toApi(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	pipeline, diags := r.client.PutIngestPipeline(ctx, name, putReq)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(pipeline)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IngestPipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ingestPipelineModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	putReq, diags := data.toApi(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	pipeline, diags := r.client.PutIngestPipeline(ctx, name, putReq)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(pipeline)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IngestPipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ingestPipelineModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	diags = r.client.DeleteIngestPipeline(ctx, name)
	resp.Diagnostics.Append(diags...)
}

type ingestPipelineModel struct {
	ID          types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	Description types.String         `tfsdk:"description"`
	OnFailure   types.List           `tfsdk:"on_failure"`
	Processors  types.List           `tfsdk:"processors"`
	Metadata    jsontypes.Normalized `tfsdk:"metadata"`
}

func (m *ingestPipelineModel) toApi(ctx context.Context) (elasticsearch.PutIngestPipelineRequest, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	output := elasticsearch.PutIngestPipelineRequest{
		Description: m.Description.ValueStringPointer(),
		Meta_:       util.NormalizedTypeToMap[json.RawMessage](m.Metadata, path.AtName("metadata"), diags),
		OnFailure: util.ListTypeToSlice(ctx, m.OnFailure, path.AtName("on_failure"), diags, func(val jsontypes.Normalized, index int) map[string]any {
			return util.NormalizedTypeToMap[any](val, path.AtName("on_failure").AtListIndex(index), diags)
		}),
		Processors: util.ListTypeToSlice(ctx, m.Processors, path.AtName("processors"), diags, func(val jsontypes.Normalized, index int) map[string]any {
			return util.NormalizedTypeToMap[any](val, path.AtName("processors").AtListIndex(index), diags)
		}),
	}

	return output, diags
}

func (m *ingestPipelineModel) fromApi(item *elasticsearch.IngestPipeline) diag.Diagnostics {
	if item == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	if m.Name.IsUnknown() || m.Name.IsNull() {
		m.Name = m.ID
	} else {
		m.ID = m.Name
	}
	m.Description = types.StringPointerValue(item.Description)
	m.Metadata = util.MapToNormalizedType(item.Meta_, path.AtName("metadata"), diags)
	m.OnFailure = util.SliceToListType(item.OnFailure, jsontypes.NormalizedType{}, path.AtName("on_failure"), diags, func(item map[string]any, index int) attr.Value {
		path := path.AtName("on_failure").AtListIndex(index)
		out, err := util.JsonMarshalS(item)
		if err != nil {
			diags.AddAttributeError(path, "marshal failure", err.Error())
		}
		return jsontypes.NewNormalizedValue(out)
	})
	m.Processors = util.SliceToListType(item.Processors, jsontypes.NormalizedType{}, path.AtName("processors"), diags, func(item map[string]any, index int) attr.Value {
		path := path.AtName("processors").AtListIndex(index)
		out, err := util.JsonMarshalS(item)
		if err != nil {
			diags.AddAttributeError(path, "marshal failure", err.Error())
		}
		return jsontypes.NewNormalizedValue(out)
	})

	return diags
}
