package kibana

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/kibana"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"

	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

var (
	_ datasource.DataSource = &ActionConnectorDataSource{}
)

func NewActionConnectorDataSource(client *clients.KibanaClient) *ActionConnectorDataSource {
	return &ActionConnectorDataSource{client: client}
}

type ActionConnectorDataSource struct {
	client *clients.KibanaClient
}

func (r *ActionConnectorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "kibana_action_connector")
}

func (r *ActionConnectorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema.Description = "Retrieve a specific role. See https://www.elastic.co/guide/en/kibana/current/role-management-api.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"connector_id": schema.StringAttribute{
			Description: "A UUID v1 or v4 to use instead of a randomly generated ID.",
			Computed:    true,
		},
		"space_id": schema.StringAttribute{
			Description: "An identifier for the space. If space_id is not provided, the default space is used.",
			Computed:    true,
			Optional:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the connector. While this name does not have to be unique, a distinctive name can help you identify a connector.",
			Computed:    true,
			Optional:    true,
		},
		"connector_type_id": schema.StringAttribute{
			Description: "The ID of the connector type, e.g. `.index`.",
			Computed:    true,
			Optional:    true,
		},
		"config": schema.StringAttribute{
			Description: "The configuration for the connector. Configuration properties vary depending on the connector type.",
			CustomType:  jsontypes.NormalizedType{},
			Computed:    true,
		},
		"is_deprecated": schema.BoolAttribute{
			Description: "Indicates whether the connector type is deprecated.",
			Computed:    true,
		},
		"is_missing_secrets": schema.BoolAttribute{
			Description: "Indicates whether secrets are missing for the connector.",
			Computed:    true,
		},
		"is_preconfigured": schema.BoolAttribute{
			Description: "Indicates whether it is a preconfigured connector.",
			Computed:    true,
		},
	}
}

func (d *ActionConnectorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data connectorDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	space := data.SpaceID.ValueString()
	if space == "" {
		space = "default"
		data.SpaceID = types.StringValue(space)
	}

	conns, diags := d.client.ListConnectors(ctx, space)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	connType := data.ConnectorTypeID.ValueString()
	conns = lo.Filter(conns, func(conn kibana.Connector, index int) bool {
		if name != "" && name != conn.Name {
			return false
		}
		if connType != "" && connType != conn.ConnectorTypeID {
			return false
		}
		return true
	})

	if len(conns) == 0 {
		diags.AddError(
			"connector not found",
			fmt.Sprintf("connector with name [%s/%s] and type [%s] not found", space, name, connType),
		)
		return
	}

	if len(conns) > 1 {
		diags.AddError(
			"more than one connector found",
			fmt.Sprintf("multiple connectors found with name [%s/%s] and type [%s]", space, name, connType),
		)
		return
	}

	diags = data.fromApi(&conns[0])
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type connectorDataSourceModel struct {
	ConnectorID      types.String         `tfsdk:"connector_id"`
	SpaceID          types.String         `tfsdk:"space_id"`
	Name             types.String         `tfsdk:"name"`
	ConnectorTypeID  types.String         `tfsdk:"connector_type_id"`
	Config           jsontypes.Normalized `tfsdk:"config"`
	IsDeprecated     types.Bool           `tfsdk:"is_deprecated"`
	IsMissingSecrets types.Bool           `tfsdk:"is_missing_secrets"`
	IsPreconfigured  types.Bool           `tfsdk:"is_preconfigured"`
}

func (m *connectorDataSourceModel) fromApi(resp *kibana.Connector) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.ConnectorID = types.StringValue(resp.ID)
	m.Name = types.StringValue(resp.Name)
	m.ConnectorTypeID = types.StringValue(resp.ConnectorTypeID)
	m.IsDeprecated = types.BoolValue(resp.IsDeprecated)
	m.IsMissingSecrets = types.BoolValue(resp.IsMissingSecrets)
	m.IsPreconfigured = types.BoolValue(resp.IsPreconfigured)
	m.Config = util.MapToNormalizedType(lo.OmitByValues(resp.Config, []any{nil}), path.AtName("config"), diags)

	return diags
}
