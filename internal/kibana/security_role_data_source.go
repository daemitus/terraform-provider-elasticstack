package kibana

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &SecurityRoleDataSource{}
)

func NewSecurityRoleDataSource(client *clients.KibanaClient) *SecurityRoleDataSource {
	return &SecurityRoleDataSource{client: client}
}

type SecurityRoleDataSource struct {
	client *clients.KibanaClient
}

func (r *SecurityRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "kibana_security_role")
}

func (r *SecurityRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema.Description = "Retrieve a specific role. See https://www.elastic.co/guide/en/kibana/current/role-management-kibana.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of this resource.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name for the role.",
			Required:    true,
		},
		"elasticsearch": schema.SingleNestedAttribute{
			Description: "Elasticsearch cluster and index privileges.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"cluster": schema.ListAttribute{
					Description: "List of the cluster privileges.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"indices": schema.ListNestedAttribute{
					Description: "A list of indices permissions entries.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"allow_restricted_indices": schema.BoolAttribute{
								Description: "Allow access to restricted indices.",
								Computed:    true,
							},
							"field_security": schema.SingleNestedAttribute{
								Description: "The document fields that the owners of the role have read access to.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"grant": schema.ListAttribute{
										Description: "List of the fields to grant the access to.",
										Computed:    true,
										ElementType: types.StringType,
									},
									"except": schema.ListAttribute{
										Description: "List of the fields to which the grants will not be applied.",
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
							"query": schema.StringAttribute{
								Description: "A search query that defines the documents the owners of the role have read access to.",
								Computed:    true,
								CustomType:  jsontypes.NormalizedType{},
							},
							"names": schema.ListAttribute{
								Description: "A list of indices (or index name patterns) to which the permissions in this entry apply.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"privileges": schema.ListAttribute{
								Description: "The index level privileges that the owners of the role have on the specified indices.",
								Computed:    true,
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
					ElementType: types.StringType,
				},
			},
		},
		"kibana": schema.ListNestedAttribute{
			Description: "The list of objects that specify the Kibana privileges for the role.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"base": schema.ListAttribute{
						Description: "A base privilege. When specified, the base must be [\"all\"] or [\"read\"].",
						Computed:    true,
						ElementType: types.StringType,
					},
					"feature": schema.MapAttribute{
						Description: "Map of feature names to specific privileges. When the feature privileges are specified, you are unable to use the \"base\" section.",
						Computed:    true,
						ElementType: types.ListType{ElemType: types.StringType},
					},
					"spaces": schema.ListAttribute{
						Description: "The spaces to apply the privileges to. To grant access to all spaces, set to [\"*\"], or omit the value.",
						Computed:    true,
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

func (d *SecurityRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data roleModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	roleName := data.Name.ValueString()
	role, diags := d.client.ReadRole(ctx, roleName)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(role)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
