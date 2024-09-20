package fleet

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/fleet"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &OutputResource{}
	_ resource.ResourceWithImportState = &OutputResource{}
)

func NewOutputResource(client *clients.FleetClient) *OutputResource {
	return &OutputResource{client: client}
}

type OutputResource struct {
	client *clients.FleetClient
}

func (r *OutputResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "fleet_output")
}

func (r *OutputResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Creates a new Fleet Output."
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of this resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"output_id": schema.StringAttribute{
			Description: "Unique identifier of the output.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the output.",
			Required:    true,
		},
		"type": schema.StringAttribute{
			Description: "The output type.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(fleet.OutputTypeElasticsearch),
					string(fleet.OutputTypeRemoteElasticsearch),
				),
			},
		},
		"hosts": schema.ListAttribute{
			Description: "A list of hosts.",
			Optional:    true,
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
		},
		"default_integrations": schema.BoolAttribute{
			Description: "Make this output the default for agent integrations.",
			Optional:    true,
		},
		"default_monitoring": schema.BoolAttribute{
			Description: "Make this output the default for agent monitoring.",
			Optional:    true,
		},
		"preset": schema.StringAttribute{
			Description: "Performance tuning presets are curated output settings for common use cases.",
			Computed:    true,
			Optional:    true,
			Default:     stringdefault.StaticString("balanced"),
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(fleet.OutputPresetBalanced),
					string(fleet.OutputPresetCustom),
					string(fleet.OutputPresetThroughput),
					string(fleet.OutputPresetScale),
					string(fleet.OutputPresetLatency),
				),
			},
		},
		"service_token": schema.StringAttribute{
			Description: "A service token for authentication when using remote_elasticsearch.",
			Computed:    true,
			Optional:    true,
			Sensitive:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"service_token_secret_id": schema.StringAttribute{
			Description: "The service token secret storage ID.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"config_yaml": schema.StringAttribute{
			Description: "Advanced YAML configuration. YAML settings here will be added to the output section of each agent policy.",
			Optional:    true,
			Sensitive:   true,
		},
	}
}

func (r *OutputResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *OutputResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data outputModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	outputId := data.OutputId.ValueString()
	output, diags := r.client.ReadOutput(ctx, outputId)
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

func (r *OutputResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data outputModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	union, diags := data.toApi(ctx, false)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body := fleet.CreateOutputRequest(*union)
	output, diags := r.client.CreateOutput(ctx, body)
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

func (r *OutputResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data outputModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	union, diags := data.toApi(ctx, true)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	outputId := data.OutputId.ValueString()
	body := fleet.UpdateOutputRequest(*union)
	output, diags := r.client.UpdateOutput(ctx, outputId, body)
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

func (r *OutputResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data outputModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	outputId := data.OutputId.ValueString()
	diags = r.client.DeleteOutput(ctx, outputId)
	resp.Diagnostics.Append(diags...)
}

type outputModel struct {
	Id                   types.String `tfsdk:"id"`
	OutputId             types.String `tfsdk:"output_id"`
	Name                 types.String `tfsdk:"name"`
	Type                 types.String `tfsdk:"type"`
	Hosts                types.List   `tfsdk:"hosts"`
	DefaultIntegrations  types.Bool   `tfsdk:"default_integrations"`
	DefaultMonitoring    types.Bool   `tfsdk:"default_monitoring"`
	Preset               types.String `tfsdk:"preset"`
	ServiceToken         types.String `tfsdk:"service_token"`
	ServiceTokenSecretId types.String `tfsdk:"service_token_secret_id"`
	ConfigYaml           types.String `tfsdk:"config_yaml"`
}

func (m *outputModel) toApi(ctx context.Context, isUpdate bool) (*fleet.OutputUnion, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	switch discriminator := m.Type.ValueString(); discriminator {
	case string(fleet.OutputTypeElasticsearch):
		output := fleet.OutputElasticsearch{
			Id:                  lo.Ternary(!isUpdate, m.OutputId.ValueStringPointer(), nil),
			Name:                m.Name.ValueString(),
			Type:                fleet.OutputType(m.Type.ValueString()),
			Hosts:               util.ListTypeToSliceBasic[string](ctx, m.Hosts, path.AtName("hosts"), diags),
			IsDefault:           m.DefaultIntegrations.ValueBoolPointer(),
			IsDefaultMonitoring: m.DefaultMonitoring.ValueBoolPointer(),
			Preset:              (*fleet.OutputPreset)(m.Preset.ValueStringPointer()),
			ConfigYaml:          m.ConfigYaml.ValueStringPointer(),
		}
		union, err := output.AsUnion()
		if err != nil {
			diags.AddError("output marshal failure", err.Error())
			return nil, diags
		}
		return union, diags

	case string(fleet.OutputTypeRemoteElasticsearch):
		output := fleet.OutputRemoteElasticsearch{
			Id:                  lo.Ternary(!isUpdate, m.OutputId.ValueStringPointer(), nil),
			Name:                m.Name.ValueString(),
			Type:                fleet.OutputType(m.Type.ValueString()),
			Hosts:               util.ListTypeToSliceBasic[string](ctx, m.Hosts, path.AtName("hosts"), diags),
			IsDefault:           m.DefaultIntegrations.ValueBoolPointer(),
			IsDefaultMonitoring: m.DefaultMonitoring.ValueBoolPointer(),
			Preset:              (*fleet.OutputPreset)(m.Preset.ValueStringPointer()),
			ConfigYaml:          m.ConfigYaml.ValueStringPointer(),
			Secrets:             &fleet.OutputRemoteElasticsearchSecrets{},
		}
		serviceToken, err := fleet.OutputRemoteElasticsearchSecretString(m.ServiceToken.ValueString()).AsUnion()
		if err != nil {
			diags.AddError("service_token marshal failure", err.Error())
			return nil, diags
		}
		output.Secrets.ServiceToken = *serviceToken

		union, err := output.AsUnion()
		if err != nil {
			diags.AddError("output marshal failure", err.Error())
			return nil, diags
		}
		return union, diags

	default:
		diags.AddError("unsupported output type", fmt.Sprintf("type: %s", discriminator))
		return nil, diags
	}
}

func (m *outputModel) fromApi(resp *fleet.OutputUnion) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	discriminator, err := resp.Discriminator()
	if err != nil {
		diags.AddError("discriminator unmarshal failure", err.Error())
		return diags
	}

	switch discriminator {
	case string(fleet.OutputTypeElasticsearch):
		output, err := resp.AsOutputElasticsearch()
		if err != nil {
			diags.AddError("elasticsearch unmarshal failure", err.Error())
			return diags
		}

		m.Id = types.StringPointerValue(output.Id)
		m.OutputId = types.StringPointerValue(output.Id)
		m.Name = types.StringValue(output.Name)
		m.Type = types.StringValue(string(output.Type))
		m.Hosts = util.SliceToListType_String(output.Hosts, path.AtName("hosts"), diags)
		m.DefaultIntegrations = types.BoolPointerValue(output.IsDefault)
		m.DefaultMonitoring = types.BoolPointerValue(output.IsDefaultMonitoring)
		m.Preset = types.StringPointerValue((*string)(output.Preset))
		m.ConfigYaml = types.StringPointerValue(output.ConfigYaml)
		m.ServiceToken = types.StringNull()
		m.ServiceTokenSecretId = types.StringNull()

	case string(fleet.OutputTypeRemoteElasticsearch):
		output, err := resp.AsOutputRemoteElasticsearch()
		if err != nil {
			diags.AddError("remote_elasticsearch unmarshal failure", err.Error())
			return diags
		}

		m.Id = types.StringPointerValue(output.Id)
		m.OutputId = types.StringPointerValue(output.Id)
		m.Name = types.StringValue(output.Name)
		m.Type = types.StringValue(string(output.Type))
		m.Hosts = util.SliceToListType_String(output.Hosts, path.AtName("hosts"), diags)
		m.DefaultIntegrations = types.BoolPointerValue(output.IsDefault)
		m.DefaultMonitoring = types.BoolPointerValue(output.IsDefaultMonitoring)
		m.Preset = types.StringPointerValue((*string)(output.Preset))
		m.ConfigYaml = types.StringPointerValue(output.ConfigYaml)

		if output.ServiceToken != nil {
			m.ServiceToken = types.StringPointerValue(output.ServiceToken)
			m.ServiceTokenSecretId = types.StringNull()
		} else {
			if output.Secrets == nil {
				diags.AddError("expected value was null", "path: output.secrets")
			} else {
				token := output.Secrets.ServiceToken
				secret, err := token.AsSecret()
				if err != nil {
					m.ServiceTokenSecretId = types.StringValue(secret.Id)
				} else {
					value, err := token.AsString()
					if err != nil {
						diags.AddError("service_token unmarshal failure", err.Error())
					} else {
						m.ServiceTokenSecretId = types.StringValue(string(value))
					}

				}
			}
		}

	default:
		diags.AddError("unsupported output type", fmt.Sprintf("type: %s", discriminator))
	}

	return diags
}
