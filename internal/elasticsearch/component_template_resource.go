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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ComponentTemplateResource{}
	_ resource.ResourceWithImportState = &ComponentTemplateResource{}
)

func NewComponentTemplateResource(client *clients.ElasticsearchClient) *ComponentTemplateResource {
	return &ComponentTemplateResource{client: client}
}

type ComponentTemplateResource struct {
	client *clients.ElasticsearchClient
}

func (r *ComponentTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "elasticsearch_component_template")
}

func (r *ComponentTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates or updates a component template. Component templates are building blocks for constructing index templates that specify index mappings, settings, and aliases. See https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-component-template.html",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Internal identifier of the resource",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the component template to create.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata": schema.StringAttribute{
				Description: "Optional user metadata about the index template.",
				Optional:    true,
				CustomType:  jsontypes.NormalizedType{},
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
				Description: "Version number used to manage component templates externally.",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func (r *ComponentTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ComponentTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data componentTemplateModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	template, diags := r.client.GetComponentTemplate(ctx, name)
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

func (r *ComponentTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data componentTemplateModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	putReq, diags := data.toApi()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	template, diags := r.client.PutComponentTemplate(ctx, name, putReq)
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

func (r *ComponentTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data componentTemplateModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	putReq, diags := data.toApi()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	template, diags := r.client.PutComponentTemplate(ctx, name, putReq)
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

func (r *ComponentTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data componentTemplateModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	diags = r.client.DeleteComponentTemplate(ctx, name)
	resp.Diagnostics.Append(diags...)
}

type componentTemplateModel struct {
	ID       types.String                   `tfsdk:"id"`
	Name     types.String                   `tfsdk:"name"`
	Metadata jsontypes.Normalized           `tfsdk:"metadata"`
	Template *componentTemplateMappingModel `tfsdk:"template"`
	Version  types.Int64                    `tfsdk:"version"`
}

type componentTemplateMappingModel struct {
	Aliases   map[string]componentTemplateAliasModel     `tfsdk:"aliases"`
	Lifecycle *componentTemplateDataStreamLifecycleModel `tfsdk:"lifecycle"`
	Mappings  jsontypes.Normalized                       `tfsdk:"mappings"`
	Settings  jsontypes.Normalized                       `tfsdk:"settings"`
}

type componentTemplateAliasModel struct {
	Filter        jsontypes.Normalized `tfsdk:"filter"`
	IndexRouting  types.String         `tfsdk:"index_routing"`
	IsHidden      types.Bool           `tfsdk:"is_hidden"`
	IsWriteIndex  types.Bool           `tfsdk:"is_write_index"`
	Routing       types.String         `tfsdk:"routing"`
	SearchRouting types.String         `tfsdk:"search_routing"`
}

type componentTemplateDataStreamLifecycleModel struct {
	DataRetention types.String `tfsdk:"data_retention"`
}

func (m *componentTemplateModel) toApi() (elasticsearch.PutComponentTemplateRequest, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	output := elasticsearch.PutComponentTemplateRequest{
		Meta_: util.NormalizedTypeToMap[json.RawMessage](m.Metadata, path.AtName("metadata"), diags),
		Template: *util.TransformStruct(m.Template, func(m componentTemplateMappingModel) estypes.IndexState {
			return estypes.IndexState{
				Aliases: util.TransformMap(m.Aliases, func(aliasName string, alias componentTemplateAliasModel) estypes.Alias {
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
				Lifecycle: util.TransformStruct(m.Lifecycle, func(m componentTemplateDataStreamLifecycleModel) estypes.DataStreamLifecycle {
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

func (m *componentTemplateModel) fromApi(item *estypes.ClusterComponentTemplate) diag.Diagnostics {
	if item == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.ID = types.StringValue(item.Name)
	m.Name = types.StringValue(item.Name)

	val := item.ComponentTemplate
	m.Metadata = util.MapToNormalizedType(val.Meta_, path.AtName("metadata"), diags)
	m.Template = util.TransformStruct(&val.Template, func(val estypes.ComponentTemplateSummary) componentTemplateMappingModel {
		return componentTemplateMappingModel{
			Aliases: util.TransformMap(val.Aliases, func(aliasName string, alias estypes.AliasDefinition) componentTemplateAliasModel {
				path := path.AtMapKey(aliasName)
				return componentTemplateAliasModel{
					Filter:        util.StructToNormalizedType(alias.Filter, path.AtName("filter"), diags),
					IndexRouting:  types.StringPointerValue(alias.IndexRouting),
					IsHidden:      types.BoolPointerValue(alias.IsHidden),
					IsWriteIndex:  types.BoolPointerValue(alias.IsWriteIndex),
					Routing:       types.StringPointerValue(alias.Routing),
					SearchRouting: types.StringPointerValue(alias.SearchRouting),
				}
			}),
			Lifecycle: util.TransformStruct(val.Lifecycle, func(val estypes.DataStreamLifecycleWithRollover) componentTemplateDataStreamLifecycleModel {
				return componentTemplateDataStreamLifecycleModel{
					DataRetention: types.StringValue(val.DataRetention.(string)),
				}
			}),
			Mappings: util.StructToNormalizedType(val.Mappings, path.AtName("mappings"), diags),
			Settings: util.MapToNormalizedType(val.Settings, path.AtName("settings"), diags),
		}
	})
	m.Version = types.Int64PointerValue(val.Version)

	return diags
}
