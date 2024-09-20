package kibana

import (
	"context"
	"fmt"
	"strings"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/kibana"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &DataViewResource{}
	_ resource.ResourceWithImportState = &DataViewResource{}
)

func NewDataViewResourceResource(client *clients.KibanaClient) *DataViewResource {
	return &DataViewResource{client: client}
}

type DataViewResource struct {
	client *clients.KibanaClient
}

func (r *DataViewResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "kibana_data_view")
}

func (r *DataViewResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Kibana data views",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Space ID and Data View ID combined.",
			},
			"space_id": schema.StringAttribute{
				Description: "An identifier for the space. If space_id is not provided, the default space is used.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("default"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_view": schema.SingleNestedAttribute{
				Description: "Data view details.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "Saved object ID.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
							stringplanmodifier.RequiresReplace(),
						},
					},
					"allow_no_index": schema.BoolAttribute{
						Description: "Allows the Data view saved object to exist before the data is available.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"field_attrs": schema.MapNestedAttribute{
						Description: "Map of field attributes by field name.",
						Computed:    true,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"custom_description": schema.StringAttribute{
									Description: "Custom description for the field.",
									Optional:    true,
								},
								"custom_label": schema.StringAttribute{
									Description: "Custom label for the field.",
									Optional:    true,
								},
								"count": schema.Int64Attribute{
									Description: "Popularity count for the field.",
									Optional:    true,
								},
							},
						},
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.RequiresReplace(),
						},
						Default: mapdefault.StaticValue(types.MapValueMust(
							types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"custom_description": types.StringType,
									"custom_label":       types.StringType,
									"count":              types.Int64Type,
								},
							},
							map[string]attr.Value{},
						)),
					},
					"field_formats": schema.MapNestedAttribute{
						Description: "Map of field formats by field name.",
						Computed:    true,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required: true,
								},
								"params": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"pattern": schema.StringAttribute{
											Required: true,
										},
									},
								},
							},
						},
						Default: mapdefault.StaticValue(types.MapValueMust(
							types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"id": types.StringType,
									"params": types.ObjectType{
										AttrTypes: map[string]attr.Type{
											"pattern": types.StringType,
										},
									},
								},
							},
							map[string]attr.Value{},
						)),
					},
					"name": schema.StringAttribute{
						Description: "The Data view name.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"namespaces": schema.ListAttribute{
						Description: "Array of space IDs for sharing the Data view between multiple spaces.",
						ElementType: types.StringType,
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.RequiresReplace(),
						},
					},
					"runtime_field_map": schema.MapNestedAttribute{
						Description: "Map of runtime field definitions by field name.",
						Computed:    true,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Mapping type of the runtime field. For more information, check [Field data types](https://www.elastic.co/guide/en/elasticsearch/reference/8.11/mapping-types.html).",
									Required:            true,
								},
								"script_source": schema.StringAttribute{
									Description: "Script of the runtime field.",
									Required:    true,
								},
							},
						},
						Default: mapdefault.StaticValue(types.MapValueMust(
							types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"type":          types.StringType,
									"script_source": types.StringType,
								},
							},
							map[string]attr.Value{},
						)),
					},
					"source_filters": schema.ListAttribute{
						Description: "List of field names you want to filter out in Discover.",
						ElementType: types.StringType,
						Computed:    true,
						Optional:    true,
					},
					"time_field_name": schema.StringAttribute{
						Description: "Timestamp field name, which you use for time-based Data views.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"title": schema.StringAttribute{
						Description: "Comma-separated list of data streams, indices, and aliases that you want to search. Supports wildcards (*).",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
			"override": schema.BoolAttribute{
				Description: "Overrides an existing data view if a data view with the provided title already exists.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *DataViewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DataViewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dataViewModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	parts := strings.SplitN(data.ID.ValueString(), "/", 2)
	dataView, diags := r.client.ReadDataView(ctx, parts[0], parts[1])
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(dataView)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *DataViewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data dataViewModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	dataView, diags := r.client.CreateDataView(ctx, spaceId, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(dataView)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *DataViewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data dataViewModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.DataView.ID.ValueString()
	spaceId := data.SpaceID.ValueString()
	dataView, diags := r.client.UpdateDataView(ctx, spaceId, id, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(dataView)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *DataViewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data dataViewModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	id := data.DataView.ID.ValueString()
	spaceId := data.SpaceID.ValueString()
	diags = r.client.DeleteDataView(ctx, spaceId, id)
	resp.Diagnostics.Append(diags...)
}

type dataViewModel struct {
	ID       types.String        `tfsdk:"id"`
	DataView *dataViewInnerModel `tfsdk:"data_view"`
	SpaceID  types.String        `tfsdk:"space_id"`
	Override types.Bool          `tfsdk:"override"`
}

type dataViewInnerModel struct {
	ID              types.String                        `tfsdk:"id"`
	AllowNoIndex    types.Bool                          `tfsdk:"allow_no_index"`
	FieldAttributes map[string]dataViewFieldAttrModel   `tfsdk:"field_attrs"`
	FieldFormats    map[string]dataViewFieldFormatModel `tfsdk:"field_formats"`
	Name            types.String                        `tfsdk:"name"`
	Namespaces      types.List                          `tfsdk:"namespaces"`
	RuntimeFieldMap map[string]dataViewRuntimeField     `tfsdk:"runtime_field_map"`
	SourceFilters   types.List                          `tfsdk:"source_filters"`
	TimeFieldName   types.String                        `tfsdk:"time_field_name"`
	Title           types.String                        `tfsdk:"title"`
}

type dataViewFieldAttrModel struct {
	CustomDescription types.String `tfsdk:"custom_description"`
	CustomLabel       types.String `tfsdk:"custom_label"`
	Count             types.Int64  `tfsdk:"count"`
}

type dataViewFieldFormatModel struct {
	ID     types.String                    `tfsdk:"id"`
	Params *dataViewFieldFormatParamsModel `tfsdk:"params"`
}

type dataViewFieldFormatParamsModel struct {
	Pattern types.String `tfsdk:"pattern"`
}

type dataViewRuntimeField struct {
	Type         types.String `tfsdk:"type"`
	ScriptSource types.String `tfsdk:"script_source"`
}

func (m *dataViewModel) toApi(ctx context.Context) (*kibana.DataView, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	output := &kibana.DataView{
		DataView: kibana.DataViewInternal{
			AllowNoIndex: m.DataView.AllowNoIndex.ValueBoolPointer(),
			FieldAttrs: util.TransformMap(m.DataView.FieldAttributes, func(key string, model dataViewFieldAttrModel) kibana.DataViewFieldAttr {
				return kibana.DataViewFieldAttr{
					CustomDescription: model.CustomDescription.ValueStringPointer(),
					CustomLabel:       model.CustomLabel.ValueStringPointer(),
					Count:             model.Count.ValueInt64Pointer(),
				}
			}),
			FieldFormats: util.TransformMap(m.DataView.FieldFormats, func(key string, model dataViewFieldFormatModel) kibana.DataViewFieldFormat {
				return kibana.DataViewFieldFormat{
					ID: model.ID.ValueString(),
					Params: util.TransformStruct(model.Params, func(m dataViewFieldFormatParamsModel) kibana.DataViewFieldFormatParams {
						return kibana.DataViewFieldFormatParams{
							Pattern: m.Pattern.ValueStringPointer(),
						}
					}),
				}
			}),
			ID:         m.DataView.ID.ValueStringPointer(),
			Name:       m.DataView.Name.ValueStringPointer(),
			Namespaces: util.ListTypeToSliceBasic[string](ctx, m.DataView.Namespaces, path.AtName("namespaces"), diags),
			RuntimeFieldMap: util.TransformMap(m.DataView.RuntimeFieldMap, func(key string, model dataViewRuntimeField) kibana.DataViewRuntimeField {
				return kibana.DataViewRuntimeField{
					Type: model.Type.ValueString(),
					Script: kibana.DataViewRuntimeFieldScript{
						Source: model.ScriptSource.ValueString(),
					},
				}
			}),
			SourceFilters: util.ListTypeToSlice(ctx, m.DataView.SourceFilters, path.AtName("source_filters"), diags, func(val string, index int) kibana.DataViewSourceFilter {
				return kibana.DataViewSourceFilter{Value: val}
			}),
			TimeFieldName: m.DataView.TimeFieldName.ValueStringPointer(),
			Title:         m.DataView.Title.ValueString(),
		},
		Override: m.Override.ValueBoolPointer(),
	}

	return output, diags
}

func (m *dataViewModel) fromApi(resp *kibana.DataView) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	if m.ID.IsUnknown() {
		m.ID = types.StringValue(fmt.Sprintf("%s/%s", m.SpaceID.ValueString(), *resp.DataView.ID))
	}

	m.DataView = &dataViewInnerModel{
		ID:           types.StringPointerValue(resp.DataView.ID),
		Name:         types.StringPointerValue(resp.DataView.Name),
		AllowNoIndex: types.BoolPointerValue(resp.DataView.AllowNoIndex),
		FieldAttributes: util.TransformMap(resp.DataView.FieldAttrs, func(key string, resp kibana.DataViewFieldAttr) dataViewFieldAttrModel {
			return dataViewFieldAttrModel{
				CustomDescription: types.StringPointerValue(resp.CustomDescription),
				CustomLabel:       types.StringPointerValue(resp.CustomLabel),
				Count:             types.Int64PointerValue(resp.Count),
			}
		}),
		FieldFormats: util.TransformMap(resp.DataView.FieldFormats, func(key string, resp kibana.DataViewFieldFormat) dataViewFieldFormatModel {
			return dataViewFieldFormatModel{
				ID: types.StringValue(resp.ID),
				Params: util.TransformStruct(resp.Params, func(resp kibana.DataViewFieldFormatParams) dataViewFieldFormatParamsModel {
					return dataViewFieldFormatParamsModel{
						Pattern: types.StringPointerValue(resp.Pattern),
					}
				}),
			}
		}),
		Namespaces: util.SliceToListType_String(resp.DataView.Namespaces, path.AtName("namespaces"), diags),
		RuntimeFieldMap: util.TransformMap(resp.DataView.RuntimeFieldMap, func(key string, resp kibana.DataViewRuntimeField) dataViewRuntimeField {
			return dataViewRuntimeField{
				Type:         types.StringValue(resp.Type),
				ScriptSource: types.StringValue(resp.Script.Source),
			}
		}),
		SourceFilters: util.SliceToListType(resp.DataView.SourceFilters, types.StringType, path.AtName("source_filters"), diags, func(resp kibana.DataViewSourceFilter, index int) attr.Value {
			return types.StringValue(resp.Value)
		}),
		TimeFieldName: types.StringPointerValue(resp.DataView.TimeFieldName),
		Title:         types.StringValue(resp.DataView.Title),
	}

	return diags
}
