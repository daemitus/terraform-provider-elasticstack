package fleet

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/fleet"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &IntegrationDataSource{}
)

func NewIntegrationDataSource(client *clients.FleetClient) *IntegrationDataSource {
	return &IntegrationDataSource{client: client}
}

type IntegrationDataSource struct {
	client *clients.FleetClient
}

func (d *IntegrationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "fleet_integration")
}

func (d *IntegrationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema.Description = "Retrieves the latest version of an integration package in Fleet."
	resp.Schema.Attributes = map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "The integration package name.",
			Required:    true,
		},
		"prerelease": schema.BoolAttribute{
			Description: "Include prerelease packages.",
			Optional:    true,
		},
		"version": schema.StringAttribute{
			Description: "The integration package version.",
			Computed:    true,
		},
	}
}

func (d *IntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data integrationDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	prerelease := data.Prerelease.ValueBool()
	packages, diags := d.client.ListPackages(ctx, prerelease)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	pkgName := data.Name.ValueString()
	data.fromApi(pkgName, packages)
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

type integrationDataSourceModel struct {
	Name       types.String `tfsdk:"name"`
	Prerelease types.Bool   `tfsdk:"prerelease"`
	Version    types.String `tfsdk:"version"`
}

func (m *integrationDataSourceModel) fromApi(pkgName string, data fleet.SearchResults) {
	m.Version = types.StringNull()
	for _, pkg := range data {
		if pkg.Name == pkgName {
			m.Version = types.StringValue(pkg.Version)
			break
		}
	}
}
