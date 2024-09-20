package fleet

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &IntegrationResource{}
	_ resource.ResourceWithImportState = &IntegrationResource{}
)

func NewIntegrationResource(client *clients.FleetClient) *IntegrationResource {
	return &IntegrationResource{client: client}
}

type IntegrationResource struct {
	client *clients.FleetClient
}

func (r *IntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "fleet_integration")
}

func (r *IntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Manage installation of a Fleet integration package."
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of this resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The integration package name.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"version": schema.StringAttribute{
			Description: "The integration package version.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"force": schema.BoolAttribute{
			Description: "Set to true to force the requested action.",
			Optional:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"skip_destroy": schema.BoolAttribute{
			Description: "Set to true if you do not wish the integration package to be uninstalled at destroy time, and instead just remove the integration package from the Terraform state.",
			Optional:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
	}
}

func (r *IntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data integrationModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	version := data.Version.ValueString()
	_, diags = r.client.ReadPackage(ctx, name, version)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data integrationModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	version := data.Version.ValueString()
	force := data.Force.ValueBool()
	diags = r.client.InstallPackage(ctx, name, version, force)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("%s/%s", name, version))

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func (r *IntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data integrationModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	version := data.Version.ValueString()
	force := data.Force.ValueBool()
	diags = r.client.InstallPackage(ctx, name, version, force)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("%s_%s", name, version))

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func (r *IntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data integrationModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	version := data.Version.ValueString()
	force := data.Force.ValueBool()
	skipDestroy := data.SkipDestroy.ValueBool()
	if skipDestroy {
		tflog.Debug(ctx, "Skipping destroy of integration package", map[string]any{"name": name, "version": version})
	} else {
		diags = r.client.UninstallPackage(ctx, name, version, force)
		resp.Diagnostics.Append(diags...)
	}
}

type integrationModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Version     types.String `tfsdk:"version"`
	Force       types.Bool   `tfsdk:"force"`
	SkipDestroy types.Bool   `tfsdk:"skip_destroy"`
}
