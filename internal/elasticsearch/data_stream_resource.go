package elasticsearch

import (
	"context"
	"fmt"
	"regexp"

	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util/xstringvalidator"
	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource                = &DataStreamResource{}
	_ resource.ResourceWithImportState = &DataStreamResource{}
)

func NewDataStreamResource(client *clients.ElasticsearchClient) *DataStreamResource {
	return &DataStreamResource{client: client}
}

type DataStreamResource struct {
	client *clients.ElasticsearchClient
}

func (r *DataStreamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "elasticsearch_data_stream")
}

func (r *DataStreamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages data streams. See https://www.elastic.co/guide/en/elasticsearch/reference/current/data-stream-apis.html",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Internal identifier of the resource",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the data stream to create.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					stringvalidator.NoneOf(".", ".."),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9!$%&'()+.;=@[\]^{}~_-]+$`), "must contain lower case alphanumeric characters and selected punctuation. See https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-data-stream.html#indices-create-data-stream-api-path-params"),
					xstringvalidator.RegexDoesNotMatch(regexp.MustCompile(`^[-_+]`), "cannot start with -, _, +"),
				},
			},
			"timestamp_field": schema.StringAttribute{
				Description: "Contains information about the data stream's @timestamp field.",
				Computed:    true,
			},
			"indices": schema.ListNestedAttribute{
				Description: "Array of objects containing information about the data stream's backing indices. The last item in this array contains information about the stream's current write index.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"index_name": schema.StringAttribute{
							Description: "Name of the backing index.",
							Computed:    true,
						},
						"index_uuid": schema.StringAttribute{
							Description: "Universally unique identifier (UUID) for the index.",
							Computed:    true,
						},
					},
				},
			},
			"generation": schema.Int64Attribute{
				Description: "Current generation for the data stream.",
				Computed:    true,
			},
			"metadata": schema.StringAttribute{
				CustomType:  jsontypes.NormalizedType{},
				Description: "Custom metadata for the stream, copied from the _meta object of the stream's matching index template.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Health status of the data stream.",
				Computed:    true,
			},
			"template": schema.StringAttribute{
				Description: "Name of the index template used to create the data stream's backing indices.",
				Computed:    true,
			},
			"ilm_policy": schema.StringAttribute{
				Description: "Name of the current ILM lifecycle policy in the stream's matching index template.",
				Computed:    true,
			},
			"hidden": schema.BoolAttribute{
				Description: "If `true`, the data stream is hidden.",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				Description: "If `true`, the data stream is created and managed by an Elastic stack component and cannot be modified through normal user interaction.",
				Computed:    true,
			},
			"replicated": schema.BoolAttribute{
				Description: "If `true`, the data stream is created and managed by cross-cluster replication and the local cluster can not write into this data stream or change its mappings.",
				Computed:    true,
			},
		},
	}
}

func (r *DataStreamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DataStreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dataStreamModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	template, diags := r.client.GetDataStream(ctx, name)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = data.fromApi(template)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *DataStreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data dataStreamModel

	// Fixme: `the target type cannot handle unknown values`
	req.Plan.SetAttribute(ctx, path.Root("indices"), types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"index_name": types.StringType,
			"index_uuid": types.StringType,
		},
	}))

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	template, diags := r.client.CreateDataStream(ctx, name)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(template)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *DataStreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Cannot update data stream", "update not supported")
}

func (r *DataStreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data dataStreamModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	diags = r.client.DeleteDataStream(ctx, name)
	resp.Diagnostics.Append(diags...)
}

type dataStreamModel struct {
	ID             types.String         `tfsdk:"id"`
	Name           types.String         `tfsdk:"name"`
	TimestampField types.String         `tfsdk:"timestamp_field"`
	Indices        []dataStreamIndex    `tfsdk:"indices"`
	Generation     types.Int64          `tfsdk:"generation"`
	Metadata       jsontypes.Normalized `tfsdk:"metadata"`
	Status         types.String         `tfsdk:"status"`
	Template       types.String         `tfsdk:"template"`
	IlmPolicy      types.String         `tfsdk:"ilm_policy"`
	Hidden         types.Bool           `tfsdk:"hidden"`
	System         types.Bool           `tfsdk:"system"`
	Replicated     types.Bool           `tfsdk:"replicated"`
}

type dataStreamIndex struct {
	IndexName types.String `tfsdk:"index_name"`
	IndexUUID types.String `tfsdk:"index_uuid"`
}

func (m *dataStreamModel) fromApi(item *estypes.DataStream) diag.Diagnostics {
	if item == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.ID = types.StringValue(item.Name)
	m.Name = types.StringValue(item.Name)
	m.TimestampField = types.StringValue(item.TimestampField.Name)
	m.Indices = util.TransformSlice(item.Indices, func(item estypes.DataStreamIndex, index int) dataStreamIndex {
		return dataStreamIndex{
			IndexName: types.StringValue(item.IndexName),
			IndexUUID: types.StringValue(item.IndexUuid),
		}
	})
	m.Generation = util.IntToInt64Type(item.Generation)
	m.Metadata = util.MapToNormalizedType(item.Meta_, path.AtName("metadata"), diags)
	m.Status = types.StringValue(item.Status.Name)
	m.Template = types.StringValue(item.Template)
	m.IlmPolicy = types.StringPointerValue(item.IlmPolicy)
	m.Hidden = types.BoolValue(item.Hidden)
	m.System = types.BoolPointerValue(item.System)
	m.Replicated = types.BoolPointerValue(item.Replicated)

	return diags
}
