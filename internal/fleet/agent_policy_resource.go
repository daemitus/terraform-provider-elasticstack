package fleet

import (
	"context"
	"fmt"
	"slices"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/fleet"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &AgentPolicyResource{}
	_ resource.ResourceWithImportState = &AgentPolicyResource{}
)

func NewAgentPolicyResource(client *clients.FleetClient) *AgentPolicyResource {
	return &AgentPolicyResource{client: client}
}

type AgentPolicyResource struct {
	client *clients.FleetClient
}

func (r *AgentPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "fleet_agent_policy")
}

func (r *AgentPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Creates a new Fleet Agent Policy. See https://www.elastic.co/guide/en/fleet/current/agent-policy.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of this resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"policy_id": schema.StringAttribute{
			Description: "Unique identifier of the agent policy.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the agent policy.",
			Required:    true,
		},
		"namespace": schema.StringAttribute{
			Description: "The namespace of the agent policy.",
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description of the agent policy.",
			Optional:    true,
		},
		"data_output_id": schema.StringAttribute{
			Description: "The identifier for the data output.",
			Optional:    true,
		},
		"monitoring_output_id": schema.StringAttribute{
			Description: "The identifier for monitoring output.",
			Optional:    true,
		},
		"fleet_server_host_id": schema.StringAttribute{
			Description: "The identifier for the Fleet server host.",
			Optional:    true,
		},
		"download_source_id": schema.StringAttribute{
			Description: "The identifier for the Elastic Agent binary download server.",
			Optional:    true,
		},
		"monitor_logs": schema.BoolAttribute{
			Description: "Enable collection of agent logs.",
			Computed:    true,
			Optional:    true,
		},
		"monitor_metrics": schema.BoolAttribute{
			Description: "Enable collection of agent metrics.",
			Computed:    true,
			Optional:    true,
		},
		"skip_destroy": schema.BoolAttribute{
			Description: "Set to true if you do not wish the agent policy to be deleted at destroy time, and instead just remove the agent policy from the Terraform state.",
			Optional:    true,
		},
	}
}

func (r *AgentPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AgentPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data agentPolicyModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	policyId := data.PolicyId.ValueString()
	policy, diags := r.client.ReadAgentPolicy(ctx, policyId)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.fromApi(policy)
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *AgentPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data agentPolicyModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	monitoring := []fleet.AgentPolicyMonitoringEnabled{}
	if data.MonitorLogs.ValueBool() {
		monitoring = append(monitoring, fleet.AgentPolicyMonitoringEnabledLogs)
	}
	if data.MonitorMetrics.ValueBool() {
		monitoring = append(monitoring, fleet.AgentPolicyMonitoringEnabledMetrics)
	}

	body := fleet.CreateAgentPolicyRequest{
		AgentFeatures:      nil,
		DataOutputId:       data.DataOutputId.ValueStringPointer(),
		Description:        data.Description.ValueStringPointer(),
		DownloadSourceId:   data.DownloadSourceId.ValueStringPointer(),
		FleetServerHostId:  data.FleetServerHostId.ValueStringPointer(),
		Id:                 data.PolicyId.ValueStringPointer(),
		InactivityTimeout:  nil,
		IsProtected:        nil,
		MonitoringEnabled:  monitoring,
		MonitoringOutputId: data.MonitoringOutputId.ValueStringPointer(),
		Name:               data.Name.ValueString(),
		Namespace:          data.Namespace.ValueString(),
		UnenrollTimeout:    nil,
	}
	policy, diags := r.client.CreateAgentPolicy(ctx, body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.fromApi(policy)
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *AgentPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data agentPolicyModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	monitoring := []fleet.AgentPolicyMonitoringEnabled{}
	if data.MonitorLogs.ValueBool() {
		monitoring = append(monitoring, fleet.AgentPolicyMonitoringEnabledLogs)
	}
	if data.MonitorMetrics.ValueBool() {
		monitoring = append(monitoring, fleet.AgentPolicyMonitoringEnabledMetrics)
	}

	policyId := data.PolicyId.ValueString()
	body := fleet.UpdateAgentPolicyRequest{
		AgentFeatures:      nil,
		DataOutputId:       data.DataOutputId.ValueStringPointer(),
		Description:        data.Description.ValueStringPointer(),
		DownloadSourceId:   data.DownloadSourceId.ValueStringPointer(),
		FleetServerHostId:  data.FleetServerHostId.ValueStringPointer(),
		InactivityTimeout:  nil,
		IsProtected:        nil,
		MonitoringEnabled:  monitoring,
		MonitoringOutputId: data.MonitoringOutputId.ValueStringPointer(),
		Name:               data.Name.ValueString(),
		Namespace:          data.Namespace.ValueString(),
		UnenrollTimeout:    nil,
	}
	policy, diags := r.client.UpdateAgentPolicy(ctx, policyId, body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.fromApi(policy)
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *AgentPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data agentPolicyModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	policyId := data.PolicyId.ValueString()
	skipDestroy := data.SkipDestroy.ValueBool()
	if skipDestroy {
		tflog.Debug(ctx, "Skipping destroy of Agent Policy", map[string]any{"policy_id": policyId})
	} else {
		diags = r.client.DeleteAgentPolicy(ctx, policyId)
		resp.Diagnostics.Append(diags...)
	}
}

type agentPolicyModel struct {
	Id                 types.String `tfsdk:"id"`
	PolicyId           types.String `tfsdk:"policy_id"`
	Name               types.String `tfsdk:"name"`
	Namespace          types.String `tfsdk:"namespace"`
	Description        types.String `tfsdk:"description"`
	DataOutputId       types.String `tfsdk:"data_output_id"`
	MonitoringOutputId types.String `tfsdk:"monitoring_output_id"`
	FleetServerHostId  types.String `tfsdk:"fleet_server_host_id"`
	DownloadSourceId   types.String `tfsdk:"download_source_id"`
	MonitorLogs        types.Bool   `tfsdk:"monitor_logs"`
	MonitorMetrics     types.Bool   `tfsdk:"monitor_metrics"`
	SkipDestroy        types.Bool   `tfsdk:"skip_destroy"`
}

func (m *agentPolicyModel) fromApi(data *fleet.AgentPolicy) {
	if data == nil {
		return
	}
	m.Id = types.StringValue(data.Id)
	m.PolicyId = types.StringValue(data.Id)
	m.DataOutputId = types.StringPointerValue(data.DataOutputId)
	m.Description = types.StringPointerValue(data.Description)
	m.DownloadSourceId = types.StringPointerValue(data.DownloadSourceId)
	m.FleetServerHostId = types.StringPointerValue(data.FleetServerHostId)
	m.MonitorLogs = types.BoolValue(slices.Contains(data.MonitoringEnabled, fleet.AgentPolicyMonitoringEnabledLogs))
	m.MonitorMetrics = types.BoolValue(slices.Contains(data.MonitoringEnabled, fleet.AgentPolicyMonitoringEnabledMetrics))
	m.MonitoringOutputId = types.StringPointerValue(data.MonitoringOutputId)
	m.Name = types.StringValue(data.Name)
	m.Namespace = types.StringValue(data.Namespace)
}
