package kibana

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/kibana"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &securityRoleResource{}
	_ resource.ResourceWithImportState = &securityRoleResource{}
)

func NewSecurityRoleResource(client *clients.KibanaClient) *securityRoleResource {
	return &securityRoleResource{client: client}
}

type securityRoleResource struct {
	client *clients.KibanaClient
}

func (r *securityRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "kibana_security_role")
}

func (r *securityRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Creates a Kibana role. See https://www.elastic.co/guide/en/kibana/master/role-management-kibana.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of this resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name for the role.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"elasticsearch": schema.SingleNestedAttribute{
			Description: "Elasticsearch cluster and index privileges.",
			Required:    true,
			Attributes: map[string]schema.Attribute{
				"cluster": schema.ListAttribute{
					Description: "List of the cluster privileges.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"indices": schema.ListNestedAttribute{
					Description: "A list of index permissions entries.",
					Optional:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"allow_restricted_indices": schema.BoolAttribute{
								Description: "Allow access to restricted indices.",
								Computed:    true,
								Optional:    true,
							},
							"field_security": schema.SingleNestedAttribute{
								Description: "The document fields that the owners of the role have read access to.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"grant": schema.ListAttribute{
										Description: "List of the fields to grant the access to.",
										Required:    true,
										ElementType: types.StringType,
									},
									"except": schema.ListAttribute{
										Description: "List of the fields to which the grants will not be applied.",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
							"query": schema.StringAttribute{
								Description: "A search query that defines the documents the owners of the role have read access to.",
								Optional:    true,
								CustomType:  jsontypes.NormalizedType{},
							},
							"names": schema.ListAttribute{
								Description: "A list of indices (or index name patterns) to which the permissions in this entry apply.",
								Required:    true,
								ElementType: types.StringType,
							},
							"privileges": schema.ListAttribute{
								Description: "The index level privileges that the owners of the role have on the specified indices.",
								Required:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
				"remote_indices": schema.ListNestedAttribute{
					Description: "A list of remote index permissions entries.",
					Optional:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"clusters": schema.ListAttribute{
								Description: "A list of remote clusters (or remote cluster name patterns) to which the permissions in this entry apply.",
								Required:    true,
								ElementType: types.StringType,
							},
							"allow_restricted_indices": schema.BoolAttribute{
								Description: "Allow access to restricted indices.",
								Computed:    true,
								Optional:    true,
							},
							"field_security": schema.SingleNestedAttribute{
								Description: "The document fields that the owners of the role have read access to.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"grant": schema.ListAttribute{
										Description: "List of the fields to grant the access to.",
										Required:    true,
										ElementType: types.StringType,
									},
									"except": schema.ListAttribute{
										Description: "List of the fields to which the grants will not be applied.",
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
							"query": schema.StringAttribute{
								Description: "A search query that defines the documents the owners of the role have read access to.",
								Optional:    true,
								CustomType:  jsontypes.NormalizedType{},
							},
							"names": schema.ListAttribute{
								Description: "A list of indices (or index name patterns) to which the permissions in this entry apply.",
								Required:    true,
								ElementType: types.StringType,
							},
							"privileges": schema.ListAttribute{
								Description: "The index level privileges that the owners of the role have on the specified indices.",
								Required:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
				"run_as": schema.ListAttribute{
					Description: "A list of usernames the owners of this role can impersonate.",
					Computed:    true,
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
		"kibana": schema.ListNestedAttribute{
			Description: "The list of objects that specify the Kibana privileges for the role.",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"base": schema.ListAttribute{
						Description: "A base privilege. When specified, the base must be [\"all\"] or [\"read\"].",
						Computed:    true,
						Optional:    true,
						ElementType: types.StringType,
					},
					"feature": schema.MapAttribute{
						Description: "Map of feature names to specific privileges. When the feature privileges are specified, you are unable to use the \"base\" section.",
						Computed:    true,
						Optional:    true,
						ElementType: types.ListType{ElemType: types.StringType},
					},
					"spaces": schema.ListAttribute{
						Description: "The spaces to apply the privileges to. To grant access to all spaces, set to [\"*\"], or omit the value.",
						Required:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
		"metadata": schema.StringAttribute{
			Description: "Optional meta-data.",
			Optional:    true,
			Computed:    true,
			CustomType:  jsontypes.NormalizedType{},
		},
	}
}

func (r *securityRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *securityRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data roleModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	roleName := data.Name.ValueString()
	role, diags := r.client.ReadRole(ctx, roleName)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(role)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *securityRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data roleModel

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

	body.Name = nil
	roleName := data.Name.ValueString()
	role, diags := r.client.PutRole(ctx, roleName, *body, true)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(role)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *securityRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data roleModel

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

	body.Name = nil
	roleName := data.Name.ValueString()
	role, diags := r.client.PutRole(ctx, roleName, *body, false)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(role)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *securityRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data roleModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	roleName := data.Name.ValueString()
	diags = r.client.DeleteRole(ctx, roleName)
	resp.Diagnostics.Append(diags...)
}

type roleModel struct {
	ID            types.String            `tfsdk:"id"`
	Name          types.String            `tfsdk:"name"`
	Elasticsearch *roleElasticsearchModel `tfsdk:"elasticsearch"`
	Kibana        []roleKibanaModel       `tfsdk:"kibana"`
	Metadata      jsontypes.Normalized    `tfsdk:"metadata"`
}

type roleElasticsearchModel struct {
	Clusters      types.List                          `tfsdk:"cluster"`
	Indices       []roleElasticsearchIndexModel       `tfsdk:"indices"`
	RemoteIndices []roleElasticsearchRemoteIndexModel `tfsdk:"remote_indices"`
	RunAs         types.List                          `tfsdk:"run_as"`
}

type roleElasticsearchIndexModel struct {
	AllowRestrictedIndices types.Bool                                `tfsdk:"allow_restricted_indices"`
	FieldSecurity          *roleElasticsearchIndexFieldSecurityModel `tfsdk:"field_security"`
	Names                  types.List                                `tfsdk:"names"`
	Query                  jsontypes.Normalized                      `tfsdk:"query"`
	Privileges             types.List                                `tfsdk:"privileges"`
}

type roleElasticsearchRemoteIndexModel struct {
	Clusters               types.List                                `tfsdk:"clusters"`
	AllowRestrictedIndices types.Bool                                `tfsdk:"allow_restricted_indices"`
	FieldSecurity          *roleElasticsearchIndexFieldSecurityModel `tfsdk:"field_security"`
	Names                  types.List                                `tfsdk:"names"`
	Query                  jsontypes.Normalized                      `tfsdk:"query"`
	Privileges             types.List                                `tfsdk:"privileges"`
}

type roleElasticsearchIndexFieldSecurityModel struct {
	Grants  types.List `tfsdk:"grant"`
	Excepts types.List `tfsdk:"except"`
}

type roleKibanaModel struct {
	Bases    types.List `tfsdk:"base"`
	Features types.Map  `tfsdk:"feature"`
	Spaces   types.List `tfsdk:"spaces"`
}

func (m *roleModel) toApi(ctx context.Context) (output *kibana.Role, diags diag.Diagnostics) {
	path := path.Empty()

	output = &kibana.Role{
		Name:     m.Name.ValueStringPointer(),
		Metadata: util.NormalizedTypeToMap[any](m.Metadata, path.AtName("metadata"), diags),
		Elasticsearch: util.TransformStruct(m.Elasticsearch, func(m roleElasticsearchModel) kibana.RoleElasticsearch {
			path := path.AtName("elasticsearch")
			return kibana.RoleElasticsearch{
				Indices: util.TransformSlice(m.Indices, func(m roleElasticsearchIndexModel, index int) kibana.RoleElasticsearchIndex {
					path := path.AtName("indices").AtListIndex(index)
					return kibana.RoleElasticsearchIndex{
						AllowRestrictedIndices: m.AllowRestrictedIndices.ValueBoolPointer(),
						Names:                  util.ListTypeToSliceBasic[string](ctx, m.Names, path.AtName("names"), diags),
						Privileges:             util.ListTypeToSliceBasic[string](ctx, m.Privileges, path.AtName("privileges"), diags),
						FieldSecurity: util.TransformStruct(m.FieldSecurity, func(m roleElasticsearchIndexFieldSecurityModel) kibana.RoleElasticsearchIndexFieldSecurity {
							path := path.AtName("field_security")
							return kibana.RoleElasticsearchIndexFieldSecurity{
								Grants:  util.ListTypeToSliceBasic[string](ctx, m.Grants, path.AtName("grants"), diags),
								Excepts: util.ListTypeToSliceBasic[string](ctx, m.Excepts, path.AtName("excepts"), diags),
							}
						}),
						Query: m.Query.ValueStringPointer(),
					}
				}),
				RemoteIndices: util.TransformSlice(m.RemoteIndices, func(m roleElasticsearchRemoteIndexModel, index int) kibana.RoleElasticsearchRemoteIndex {
					path := path.AtName("remote_indices").AtListIndex(index)
					return kibana.RoleElasticsearchRemoteIndex{
						Clusters:               util.ListTypeToSliceBasic[string](ctx, m.Clusters, path.AtName("clusters"), diags),
						AllowRestrictedIndices: m.AllowRestrictedIndices.ValueBoolPointer(),
						Names:                  util.ListTypeToSliceBasic[string](ctx, m.Names, path.AtName("names"), diags),
						Privileges:             util.ListTypeToSliceBasic[string](ctx, m.Privileges, path.AtName("privileges"), diags),
						FieldSecurity: util.TransformStruct(m.FieldSecurity, func(m roleElasticsearchIndexFieldSecurityModel) kibana.RoleElasticsearchIndexFieldSecurity {
							path := path.AtName("field_security")
							return kibana.RoleElasticsearchIndexFieldSecurity{
								Grants:  util.ListTypeToSliceBasic[string](ctx, m.Grants, path.AtName("grants"), diags),
								Excepts: util.ListTypeToSliceBasic[string](ctx, m.Excepts, path.AtName("excepts"), diags),
							}
						}),
						Query: m.Query.ValueStringPointer(),
					}
				}),
				Clusters: util.ListTypeToSliceBasic[string](ctx, m.Clusters, path.AtName("clusters"), diags),
				RunAs:    util.ListTypeToSliceBasic[string](ctx, m.RunAs, path.AtName("run_as"), diags),
			}
		}),
		Kibana: util.TransformSlice(m.Kibana, func(m roleKibanaModel, index int) kibana.RoleKibana {
			path := path.AtName("kibana").AtListIndex(index)
			return kibana.RoleKibana{
				Bases: util.ListTypeToSliceBasic[string](ctx, m.Bases, path.AtName("bases"), diags),
				Features: util.MapTypeToMap(ctx, m.Features, path.AtName("features"), diags, func(key string, value types.List) []string {
					path := path.AtName("features").AtMapKey(key)
					return util.ListTypeToSliceBasic[string](ctx, value, path, diags)
				}),
				Spaces: util.ListTypeToSliceBasic[string](ctx, m.Spaces, path.AtName("spaces"), diags),
			}
		}),
	}

	return
}

func (m *roleModel) fromApi(resp *kibana.Role) (diags diag.Diagnostics) {
	path := path.Empty()

	m.ID = types.StringPointerValue(resp.Name)
	m.Name = types.StringPointerValue(resp.Name)
	m.Elasticsearch = util.TransformStruct(resp.Elasticsearch, func(resp kibana.RoleElasticsearch) roleElasticsearchModel {
		path := path.AtName("elasticsearch")
		return roleElasticsearchModel{
			Clusters: util.SliceToListType_String(resp.Clusters, path.AtName("clusters"), diags),
			Indices: util.TransformSlice(resp.Indices, func(resp kibana.RoleElasticsearchIndex, index int) roleElasticsearchIndexModel {
				path := path.AtName("indices").AtListIndex(index)
				return roleElasticsearchIndexModel{
					AllowRestrictedIndices: types.BoolPointerValue(resp.AllowRestrictedIndices),
					FieldSecurity: util.TransformStruct(resp.FieldSecurity, func(resp kibana.RoleElasticsearchIndexFieldSecurity) roleElasticsearchIndexFieldSecurityModel {
						path := path.AtName("field_security")
						return roleElasticsearchIndexFieldSecurityModel{
							Grants:  util.SliceToListType_String(resp.Grants, path.AtName("grants"), diags),
							Excepts: util.SliceToListType_String(resp.Excepts, path.AtName("excepts"), diags),
						}
					}),
					Query:      jsontypes.NewNormalizedPointerValue(resp.Query),
					Names:      util.SliceToListType_String(resp.Names, path.AtName("names"), diags),
					Privileges: util.SliceToListType_String(resp.Privileges, path.AtName("privileges"), diags),
				}
			}),
			RemoteIndices: util.TransformSlice(resp.RemoteIndices, func(resp kibana.RoleElasticsearchRemoteIndex, index int) roleElasticsearchRemoteIndexModel {
				path := path.AtName("remote_indices").AtListIndex(index)
				return roleElasticsearchRemoteIndexModel{
					Clusters:               util.SliceToListType_String(resp.Clusters, path.AtName("clusters"), diags),
					AllowRestrictedIndices: types.BoolPointerValue(resp.AllowRestrictedIndices),
					FieldSecurity: util.TransformStruct(resp.FieldSecurity, func(resp kibana.RoleElasticsearchIndexFieldSecurity) roleElasticsearchIndexFieldSecurityModel {
						path := path.AtName("field_security")
						return roleElasticsearchIndexFieldSecurityModel{
							Grants:  util.SliceToListType_String(resp.Grants, path.AtName("grants"), diags),
							Excepts: util.SliceToListType_String(resp.Excepts, path.AtName("excepts"), diags),
						}
					}),
					Query:      jsontypes.NewNormalizedPointerValue(resp.Query),
					Names:      util.SliceToListType_String(resp.Names, path.AtName("names"), diags),
					Privileges: util.SliceToListType_String(resp.Privileges, path.AtName("privileges"), diags),
				}
			}),
			RunAs: util.SliceToListType_String(resp.RunAs, path.AtName("run_as"), diags),
		}
	})
	m.Kibana = util.TransformSlice(resp.Kibana, func(resp kibana.RoleKibana, index int) roleKibanaModel {
		path := path.AtName("kibana").AtListIndex(index)
		return roleKibanaModel{
			Bases: util.SliceToListType_String(resp.Bases, path.AtName("bases"), diags),
			Features: util.MapToMapType(resp.Features, types.ListType{ElemType: types.StringType}, path.AtName("features"), diags, func(key string, value []string) attr.Value {
				path := path.AtName("features").AtMapKey(key)
				return util.SliceToListType_String(value, path, diags)
			}),
			Spaces: util.SliceToListType_String(resp.Spaces, path.AtName("spaces"), diags),
		}
	})
	m.Metadata = util.MapToNormalizedType(resp.Metadata, path.AtName("metadata"), diags)

	return
}
