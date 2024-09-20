package kibana

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/kibana"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &ActionConnectorResource{}
	_ resource.ResourceWithImportState = &ActionConnectorResource{}
)

func NewActionConnectorResource(client *clients.KibanaClient) *ActionConnectorResource {
	return &ActionConnectorResource{client: client}
}

type ActionConnectorResource struct {
	client *clients.KibanaClient
}

func (r *ActionConnectorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "kibana_action_connector")
}

func (r *ActionConnectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Creates a Kibana action connector. See https://www.elastic.co/guide/en/kibana/current/action-types.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "A UUID v1 or v4 to use instead of a randomly generated ID.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"space_id": schema.StringAttribute{
			Description: "An identifier for the space. If space_id is not provided, the default space is used.",
			Computed:    true,
			Optional:    true,
			Default:     stringdefault.StaticString("default"),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the connector. While this name does not have to be unique, a distinctive name can help you identify a connector.",
			Required:    true,
		},
		"connector_type_id": schema.StringAttribute{
			Description: "The ID of the connector type, e.g. `.index`.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"config": schema.StringAttribute{
			Description: "The configuration for the connector. Configuration properties vary depending on the connector type.",
			CustomType:  jsontypes.NormalizedType{},
			Computed:    true,
			Optional:    true,
			Default:     stringdefault.StaticString("{}"),
		},
		"secrets": schema.StringAttribute{
			Description: "The secrets configuration for the connector. Secrets configuration properties vary depending on the connector type.",
			CustomType:  jsontypes.NormalizedType{},
			Computed:    true,
			Optional:    true,
			Default:     stringdefault.StaticString("{}"),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
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

func (r *ActionConnectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ActionConnectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data connectorModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	space := data.SpaceID.ValueString()
	connId := data.ID.ValueString()
	conn, diags := r.client.ReadConnector(ctx, space, connId)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(conn)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ActionConnectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data connectorModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi()
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	space := data.SpaceID.ValueString()
	connId := data.ID.ValueString()
	conn, diags := r.client.CreateConnector(ctx, space, connId, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(conn)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ActionConnectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data connectorModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi()
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	space := data.SpaceID.ValueString()
	connId := data.ID.ValueString()
	output, diags := r.client.UpdateConnector(ctx, space, connId, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(output)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ActionConnectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data connectorModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	space := data.SpaceID.ValueString()
	connId := data.ID.ValueString()
	diags = r.client.DeleteConnector(ctx, space, connId)
	resp.Diagnostics.Append(diags...)
}

type connectorModel struct {
	ID               types.String         `tfsdk:"id"`
	SpaceID          types.String         `tfsdk:"space_id"`
	Name             types.String         `tfsdk:"name"`
	ConnectorTypeID  types.String         `tfsdk:"connector_type_id"`
	Config           jsontypes.Normalized `tfsdk:"config"`
	Secrets          jsontypes.Normalized `tfsdk:"secrets"`
	IsDeprecated     types.Bool           `tfsdk:"is_deprecated"`
	IsMissingSecrets types.Bool           `tfsdk:"is_missing_secrets"`
	IsPreconfigured  types.Bool           `tfsdk:"is_preconfigured"`
}

func (m *connectorModel) toApi() (*kibana.ConnectorRequest, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	req := &kibana.ConnectorRequest{
		ConnectorTypeID: m.ConnectorTypeID.ValueString(),
		Name:            m.Name.ValueString(),
		Config:          util.NormalizedTypeToMap[any](m.Config, path.AtName("config"), diags),
		Secrets:         util.NormalizedTypeToMap[any](m.Secrets, path.AtName("secrets"), diags),
	}

	return req, diags
}

func (m *connectorModel) fromApi(resp *kibana.Connector) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.ID = types.StringValue(resp.ID)
	m.Name = types.StringValue(resp.Name)
	m.ConnectorTypeID = types.StringValue(resp.ConnectorTypeID)
	m.IsDeprecated = types.BoolValue(resp.IsDeprecated)
	m.IsMissingSecrets = types.BoolValue(resp.IsMissingSecrets)
	m.IsPreconfigured = types.BoolValue(resp.IsPreconfigured)
	m.Config = util.MapToNormalizedType(lo.OmitByValues(resp.Config, []any{nil}), path.AtName("config"), diags)

	return diags
}
