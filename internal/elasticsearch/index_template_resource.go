package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &IndexTemplateResource{}
	_ resource.ResourceWithImportState = &IndexTemplateResource{}
)

func NewIndexTemplateResource(client *clients.ElasticsearchClient) *IndexTemplateResource {
	return &IndexTemplateResource{client: client}
}

type IndexTemplateResource struct {
	client *clients.ElasticsearchClient
}

func (r *IndexTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "elasticsearch_index_template")
}

func (r *IndexTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates or updates an index template. Index templates define settings, mappings, and aliases that can be applied automatically to new indices. See https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-put-template.html",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Internal identifier of the resource",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the index template to create.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"allow_auto_create": schema.BoolAttribute{
				Description: "If true, allows the automatic creation of indices.",
				Optional:    true,
			},
			"composed_of": schema.ListAttribute{
				Description: "An ordered list of component template names.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"data_stream": schema.SingleNestedAttribute{
				Description: "If this object is included, the template is used to create data streams and their backing indices. Supports an empty object.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"hidden": schema.BoolAttribute{
						Description: "If true, the data stream is hidden.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
			"index_patterns": schema.SetAttribute{
				Description: "Array of wildcard (*) expressions used to match the names of data streams and indices during creation.",
				Required:    true,
				ElementType: types.StringType,
			},
			"metadata": schema.StringAttribute{
				Description: "Optional user metadata about the index template.",
				Optional:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
			"priority": schema.Int64Attribute{
				Description: "Priority to determine index template precedence when a new data stream or index is created.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"template": schema.SingleNestedAttribute{
				Description: "Template to be applied. It may optionally include an aliases, mappings, or settings configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"aliases": schema.MapNestedAttribute{
						Description: "Aliases to add.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"filter": schema.StringAttribute{
									CustomType:  jsontypes.NormalizedType{},
									Description: "Query used to limit documents the alias can access.",
									Computed:    true,
									Optional:    true,
								},
								"index_routing": schema.StringAttribute{
									Description: "Value used to route indexing operations to a specific shard. If specified, this overwrites the `routing` value for indexing operations.",
									Computed:    true,
									Optional:    true,
								},
								"is_hidden": schema.BoolAttribute{
									Description: "If true, the alias is hidden.",
									Computed:    true,
									Optional:    true,
								},
								"is_write_index": schema.BoolAttribute{
									Description: "If true, the index is the write index for the alias.",
									Computed:    true,
									Optional:    true,
								},
								"routing": schema.StringAttribute{
									Description: "Value used to route indexing and search operations to a specific shard.",
									Computed:    true,
									Optional:    true,
								},
								"search_routing": schema.StringAttribute{
									Description: "Value used to route search operations to a specific shard. If specified, this overwrites the routing value for search operations.",
									Computed:    true,
									Optional:    true,
								},
							},
						},
					},
					"lifecycle": schema.MapNestedAttribute{
						Description: "Data lifecycle denotes that a data stream is managed by the data stream lifecycle and contains the configuration.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"data_retention": schema.StringAttribute{
									Description: "Retention period to keep each datastream.",
									Required:    true,
								},
							},
						},
					},
					"mappings": schema.StringAttribute{
						CustomType:  jsontypes.NormalizedType{},
						Description: "Mapping for fields in the index.",
						Optional:    true,
					},
					"settings": schema.StringAttribute{
						CustomType:  jsontypes.NormalizedType{},
						Description: "Configuration options for the index. See https://www.elastic.co/guide/en/elasticsearch/reference/current/index-modules.html#index-modules-settings",
						Optional:    true,
					},
				},
			},
			"version": schema.Int64Attribute{
				Description: "Version number used to manage index templates externally.",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func (r *IndexTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IndexTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data indexTemplateModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	template, diags := r.client.GetIndexTemplate(ctx, name)
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

func (r *IndexTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data indexTemplateModel

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
	template, diags := r.client.PutIndexTemplate(ctx, name, true, putReq)
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

func (r *IndexTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data indexTemplateModel

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
	template, diags := r.client.PutIndexTemplate(ctx, name, false, putReq)
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

func (r *IndexTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data indexTemplateModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	diags = r.client.DeleteIndexTemplate(ctx, name)
	resp.Diagnostics.Append(diags...)
}

type indexTemplateModel struct {
	ID              types.String                  `tfsdk:"id"`
	Name            types.String                  `tfsdk:"name"`
	AllowAutoCreate types.Bool                    `tfsdk:"allow_auto_create"`
	ComposedOf      types.List                    `tfsdk:"composed_of"`
	DataStream      *indexTemplateDataStreamModel `tfsdk:"data_stream"`
	IndexPatterns   types.Set                     `tfsdk:"index_patterns"`
	Metadata        jsontypes.Normalized          `tfsdk:"metadata"`
	Priority        types.Int64                   `tfsdk:"priority"`
	Template        *indexTemplateMappingModel    `tfsdk:"template"`
	Version         types.Int64                   `tfsdk:"version"`
}

type indexTemplateDataStreamModel struct {
	Hidden types.Bool `tfsdk:"hidden"`
}

type indexTemplateMappingModel struct {
	Aliases   map[string]indexTemplateAliasModel     `tfsdk:"aliases"`
	Lifecycle *indexTemplateDataStreamLifecycleModel `tfsdk:"lifecycle"`
	Mappings  jsontypes.Normalized                   `tfsdk:"mappings"`
	Settings  jsontypes.Normalized                   `tfsdk:"settings"`
}

type indexTemplateAliasModel struct {
	Filter        jsontypes.Normalized `tfsdk:"filter"`
	IndexRouting  types.String         `tfsdk:"index_routing"`
	IsHidden      types.Bool           `tfsdk:"is_hidden"`
	IsWriteIndex  types.Bool           `tfsdk:"is_write_index"`
	Routing       types.String         `tfsdk:"routing"`
	SearchRouting types.String         `tfsdk:"search_routing"`
}

type indexTemplateDataStreamLifecycleModel struct {
	DataRetention types.String `tfsdk:"data_retention"`
}

func (m *indexTemplateModel) toApi(ctx context.Context) (elasticsearch.PutIndexTemplateRequest, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	output := elasticsearch.PutIndexTemplateRequest{
		AllowAutoCreate: m.AllowAutoCreate.ValueBoolPointer(),
		ComposedOf:      util.ListTypeToSliceBasic[string](ctx, m.ComposedOf, path.AtName("composed_of"), diags),
		DataStream: util.TransformStruct(m.DataStream, func(m indexTemplateDataStreamModel) estypes.DataStreamVisibility {
			return estypes.DataStreamVisibility{
				Hidden: m.Hidden.ValueBoolPointer(),
			}
		}),
		IndexPatterns: util.SetTypeToSliceBasic[string](ctx, m.IndexPatterns, path.AtName("index_patterns"), diags),
		Meta_:         util.NormalizedTypeToMap[json.RawMessage](m.Metadata, path.AtName("metadata"), diags),
		Priority:      m.Priority.ValueInt64Pointer(),
		Template: util.TransformStruct(m.Template, func(m indexTemplateMappingModel) estypes.IndexTemplateMapping {
			return estypes.IndexTemplateMapping{
				Aliases: util.TransformMap(m.Aliases, func(aliasName string, alias indexTemplateAliasModel) estypes.Alias {
					path := path.AtMapKey(aliasName)
					return estypes.Alias{
						Filter:        util.NormalizedTypeToStruct[estypes.Query](alias.Filter, path.AtName("filter"), diags),
						IndexRouting:  alias.IndexRouting.ValueStringPointer(),
						IsHidden:      alias.IsHidden.ValueBoolPointer(),
						IsWriteIndex:  alias.IsWriteIndex.ValueBoolPointer(),
						Routing:       alias.Routing.ValueStringPointer(),
						SearchRouting: alias.SearchRouting.ValueStringPointer(),
					}
				}),
				Lifecycle: util.TransformStruct(m.Lifecycle, func(m indexTemplateDataStreamLifecycleModel) estypes.DataStreamLifecycle {
					return estypes.DataStreamLifecycle{
						DataRetention: m.DataRetention.ValueStringPointer(),
					}
				}),
				Mappings: util.NormalizedTypeToStruct[estypes.TypeMapping](m.Mappings, path.AtName("mappings"), diags),
				Settings: util.NormalizedTypeToStruct[estypes.IndexSettings](m.Settings, path.AtName("settings"), diags),
			}
		}),
		Version: m.Version.ValueInt64Pointer(),
	}

	return output, diags
}

func (m *indexTemplateModel) fromApi(item *estypes.IndexTemplateItem) diag.Diagnostics {
	if item == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.ID = types.StringValue(item.Name)
	m.Name = types.StringValue(item.Name)

	val := item.IndexTemplate
	m.AllowAutoCreate = types.BoolPointerValue(val.AllowAutoCreate)
	m.ComposedOf = util.SliceToListType_String(val.ComposedOf, path.AtName("composed_of"), diags)
	m.DataStream = util.TransformStruct(val.DataStream, func(val estypes.IndexTemplateDataStreamConfiguration) indexTemplateDataStreamModel {
		return indexTemplateDataStreamModel{
			Hidden: types.BoolPointerValue(val.Hidden),
		}
	})
	m.IndexPatterns = util.SliceToSetType_String(val.IndexPatterns, path.AtName("index_patterns"), diags)
	m.Metadata = util.MapToNormalizedType(val.Meta_, path.AtName("metadata"), diags)
	m.Priority = types.Int64PointerValue(val.Priority)
	m.Template = util.TransformStruct(val.Template, func(val estypes.IndexTemplateSummary) indexTemplateMappingModel {
		return indexTemplateMappingModel{
			Aliases: util.TransformMap(val.Aliases, func(aliasName string, alias estypes.Alias) indexTemplateAliasModel {
				path := path.AtMapKey(aliasName)
				return indexTemplateAliasModel{
					Filter:        util.StructToNormalizedType(alias.Filter, path.AtName("filter"), diags),
					IndexRouting:  types.StringPointerValue(alias.IndexRouting),
					IsHidden:      types.BoolPointerValue(alias.IsHidden),
					IsWriteIndex:  types.BoolPointerValue(alias.IsWriteIndex),
					Routing:       types.StringPointerValue(alias.Routing),
					SearchRouting: types.StringPointerValue(alias.SearchRouting),
				}
			}),
			Lifecycle: util.TransformStruct(val.Lifecycle, func(val estypes.DataStreamLifecycleWithRollover) indexTemplateDataStreamLifecycleModel {
				return indexTemplateDataStreamLifecycleModel{
					DataRetention: types.StringValue(val.DataRetention.(string)),
				}
			}),
			Mappings: util.StructToNormalizedType(val.Mappings, path.AtName("mappings"), diags),
			Settings: util.StructToNormalizedType(val.Settings, path.AtName("settings"), diags),
		}
	})
	m.Version = types.Int64PointerValue(val.Version)

	return diags
}
