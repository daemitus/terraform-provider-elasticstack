package kibana

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/kibana"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

var (
	_ datasource.DataSource = &SpacesDataSource{}
)

func NewSpacesDataSource(client *clients.KibanaClient) *SpacesDataSource {
	return &SpacesDataSource{client: client}
}

type SpacesDataSource struct {
	client *clients.KibanaClient
}

func (r *SpacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "kibana_spaces")
}

func (r *SpacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema.Description = "Retrieve all Kibana spaces. See https://www.elastic.co/guide/en/kibana/current/spaces-api-get-all.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"spaces": schema.ListNestedAttribute{
			Computed:    true,
			Description: "A list of all available spaces.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The space ID that is part of the Kibana URL when inside the space.",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "The display name for the space.",
						Computed:    true,
					},
					"description": schema.StringAttribute{
						Description: "The description for the space.",
						Computed:    true,
					},
					"disabled_features": schema.ListAttribute{
						Description: "The list of disabled features for the space. To get a list of available feature IDs, use the Features API (https://www.elastic.co/guide/en/kibana/master/features-api-get.html).",
						Computed:    true,
						ElementType: types.StringType,
					},
					"initials": schema.StringAttribute{
						Description: "The initials shown in the space avatar. By default, the initials are automatically generated from the space name. Initials must be 1 or 2 characters.",
						Computed:    true,
					},
					"color": schema.StringAttribute{
						Description: "The hexadecimal color code used in the space avatar. By default, the color is automatically generated from the space name.",
						Computed:    true,
					},
					"image_url": schema.StringAttribute{
						Description: "The data-URL encoded image to display in the space avatar.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *SpacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data spacesModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	spaces, diags := d.client.ListSpaces(ctx)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(spaces)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

type spacesModel struct {
	Spaces []spaceModel `tfsdk:"spaces"`
}

func (m *spacesModel) fromApi(resp kibana.Spaces) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Spaces = lo.Map(resp, func(space kibana.Space, index int) spaceModel {
		m := spaceModel{}
		d := m.fromApi(&space)
		diags.Append(d...)
		return m
	})

	return diags
}
