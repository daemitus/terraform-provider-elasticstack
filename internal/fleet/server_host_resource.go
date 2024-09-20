package fleet

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/fleet"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                = &ServerHostResource{}
	_ resource.ResourceWithImportState = &ServerHostResource{}
)

func NewServerHostResource(client *clients.FleetClient) *ServerHostResource {
	return &ServerHostResource{client: client}
}

type ServerHostResource struct {
	client *clients.FleetClient
}

func (r *ServerHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "fleet_server_host")
}

func (r *ServerHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Creates a new Fleet Server Host."
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of this resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"host_id": schema.StringAttribute{
			Description: "Unique identifier of the Fleet server host.",
			Computed:    true,
			Optional:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the Fleet server host.",
			Required:    true,
		},
		"hosts": schema.ListAttribute{
			Description: "A list of hosts.",
			Required:    true,
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
		},
		"default": schema.BoolAttribute{
			Description: "Set as default.",
			Optional:    true,
		},
	}
}

func (r *ServerHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ServerHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data serverHostModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	hostId := data.HostId.ValueString()
	host, diags := r.client.ReadFleetServerHost(ctx, hostId)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(host)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ServerHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data serverHostModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	var hostUrls []string
	diags = data.Hosts.ElementsAs(ctx, &hostUrls, false)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body := fleet.CreateFleetServerHostsRequest{
		HostUrls:  hostUrls,
		Id:        data.HostId.ValueStringPointer(),
		IsDefault: data.Default.ValueBoolPointer(),
		Name:      data.Name.ValueString(),
	}
	host, diags := r.client.CreateFleetServerHost(ctx, body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(host)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ServerHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data serverHostModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	var hostUrls []string
	diags = data.Hosts.ElementsAs(ctx, &hostUrls, false)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	hostId := data.HostId.ValueString()
	body := fleet.UpdateFleetServerHostsRequest{
		HostUrls:  &hostUrls,
		IsDefault: data.Default.ValueBoolPointer(),
		Name:      data.Name.ValueStringPointer(),
	}
	host, diags := r.client.UpdateFleetServerHost(ctx, hostId, body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(host)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ServerHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data serverHostModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	hostId := data.HostId.ValueString()
	diags = r.client.DeleteFleetServerHost(ctx, hostId)
	resp.Diagnostics.Append(diags...)
}

type serverHostModel struct {
	Id      types.String `tfsdk:"id"`
	HostId  types.String `tfsdk:"host_id"`
	Name    types.String `tfsdk:"name"`
	Hosts   types.List   `tfsdk:"hosts"`
	Default types.Bool   `tfsdk:"default"`
}

func (m *serverHostModel) fromApi(data *fleet.FleetServerHost) diag.Diagnostics {
	if data == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.Id = types.StringValue(data.Id)
	m.HostId = types.StringValue(data.Id)
	m.Name = types.StringPointerValue(data.Name)
	m.Hosts = util.SliceToListType_String(data.HostUrls, path.AtName("hosts"), diags)
	m.Default = types.BoolValue(data.IsDefault)

	return diags
}
