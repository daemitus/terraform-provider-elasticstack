package fleet

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/fleet"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &IntegrationPolicyResource{}
	_ resource.ResourceWithImportState = &IntegrationPolicyResource{}
)

func NewIntegrationPolicyResource(client *clients.FleetClient) *IntegrationPolicyResource {
	return &IntegrationPolicyResource{client: client}
}

type IntegrationPolicyResource struct {
	client *clients.FleetClient
}

func (r *IntegrationPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "fleet_integration_policy")
}

func (r *IntegrationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Creates a new Fleet Integration Policy. See https://www.elastic.co/guide/en/fleet/current/add-integration-to-policy.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of this resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"policy_id": schema.StringAttribute{
			Description: "Unique identifier of the integration policy.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the integration policy.",
			Required:    true,
		},
		"namespace": schema.StringAttribute{
			Description: "The namespace of the integration policy.",
			Required:    true,
		},
		"agent_policy_id": schema.StringAttribute{
			Description: "ID of the agent policy.",
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description of the integration policy.",
			Optional:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Enable the integration policy.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(true),
		},
		"force": schema.BoolAttribute{
			Description: "Force operations, such as creation and deletion, to occur.",
			Optional:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"integration_name": schema.StringAttribute{
			Description: "The name of the integration package.",
			Required:    true,
		},
		"integration_version": schema.StringAttribute{
			Description: "The version of the integration package.",
			Required:    true,
		},
		"inputs": schema.MapNestedAttribute{
			Description: "A mapping of the input identifier to its config.",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable the input.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(true),
					},
					"vars_json": schema.StringAttribute{
						Description: "Input-level variables in JSON format.",
						CustomType:  jsontypes.NormalizedType{},
						Computed:    true,
						Optional:    true,
					},
					"streams": schema.MapNestedAttribute{
						Description: "Input streams.",
						Computed:    true,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "Enable the stream.",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(true),
								},
								"vars_json": schema.StringAttribute{
									Description: "Stream-level variables in JSON format.",
									CustomType:  jsontypes.NormalizedType{},
									Computed:    true,
									Optional:    true,
								},
							},
						},
					},
				},
			},
		},
		"vars_json": schema.StringAttribute{
			Description: "Integration-level variables in JSON format.",
			CustomType:  jsontypes.NormalizedType{},
			Computed:    true,
			Optional:    true,
		},
	}
}

func (r *IntegrationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IntegrationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data integrationPolicyModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	policyId := data.PolicyId.ValueString()
	policy, diags := r.client.ReadPackagePolicy(ctx, policyId)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(ctx, policy, resp.Private)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IntegrationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data integrationPolicyModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi(false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, diags := r.client.CreatePackagePolicy(ctx, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(ctx, policy, resp.Private)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IntegrationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data integrationPolicyModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi(true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := data.PolicyId.ValueString()
	policy, diags := r.client.UpdatePackagePolicy(ctx, policyId, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(ctx, policy, resp.Private)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IntegrationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data integrationPolicyModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	policyId := data.PolicyId.ValueString()
	force := data.Force.ValueBool()
	diags = r.client.DeletePackagePolicy(ctx, policyId, force)
	resp.Diagnostics.Append(diags...)
}

type integrationPolicyModel struct {
	Id                 types.String                           `tfsdk:"id"`
	PolicyId           types.String                           `tfsdk:"policy_id"`
	Name               types.String                           `tfsdk:"name"`
	Namespace          types.String                           `tfsdk:"namespace"`
	AgentPolicyId      types.String                           `tfsdk:"agent_policy_id"`
	Description        types.String                           `tfsdk:"description"`
	Enabled            types.Bool                             `tfsdk:"enabled"`
	Force              types.Bool                             `tfsdk:"force"`
	IntegrationName    types.String                           `tfsdk:"integration_name"`
	IntegrationVersion types.String                           `tfsdk:"integration_version"`
	Inputs             map[string]integrationPolicyInputModel `tfsdk:"inputs"`
	VarsJson           jsontypes.Normalized                   `tfsdk:"vars_json"`
}

type integrationPolicyInputModel struct {
	Enabled  types.Bool                              `tfsdk:"enabled"`
	Streams  map[string]integrationPolicyStreamModel `tfsdk:"streams"`
	VarsJson jsontypes.Normalized                    `tfsdk:"vars_json"`
}

type integrationPolicyStreamModel struct {
	Enabled  types.Bool           `tfsdk:"enabled"`
	VarsJson jsontypes.Normalized `tfsdk:"vars_json"`
}

func (m *integrationPolicyModel) toApi(isUpdate bool) (*fleet.PackagePolicyRequest, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	body := &fleet.PackagePolicyRequest{
		Description: m.Description.ValueStringPointer(),
		Force:       m.Force.ValueBoolPointer(),
		Id:          nil,
		Inputs:      map[string]fleet.PackagePolicyRequestInput{},
		Name:        m.Name.ValueString(),
		Namespace:   m.Namespace.ValueStringPointer(),
		Package: fleet.PackagePolicyRequestPackage{
			Name:    m.IntegrationName.ValueString(),
			Version: m.IntegrationVersion.ValueString(),
		},
		PolicyId: m.AgentPolicyId.ValueString(),
		Vars:     nil,
	}
	if isUpdate {
		body.Id = m.Id.ValueStringPointer()
	}
	body.Vars = util.NormalizedTypeToMap[any](m.VarsJson, path.AtName("vars_json"), diags)

	for inputName, inputData := range m.Inputs {
		path := path.AtName("inputs").AtMapKey(inputName)
		input := fleet.PackagePolicyRequestInput{
			Enabled: inputData.Enabled.ValueBoolPointer(),
			Streams: map[string]fleet.PackagePolicyRequestInputStream{},
			Vars:    util.NormalizedTypeToMap[any](inputData.VarsJson, path.AtName("vars_json"), diags),
		}
		body.Inputs[inputName] = input

		for streamName, streamData := range inputData.Streams {
			path := path.AtName("streams").AtMapKey(streamName)
			stream := fleet.PackagePolicyRequestInputStream{
				Enabled: streamData.Enabled.ValueBoolPointer(),
				Vars:    util.NormalizedTypeToMap[any](streamData.VarsJson, path.AtName("vars_json"), diags),
			}
			input.Streams[streamName] = stream
		}
	}
	return body, diags
}

func (m *integrationPolicyModel) fromApi(ctx context.Context, resp *fleet.PackagePolicy, private privateData) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	secretStore := newSecretStore(ctx, private)
	secretStore.RemoveAllExcept(resp.SecretReferences)

	var diags diag.Diagnostics
	originalInputs := m.Inputs

	m.Id = types.StringValue(resp.Id)
	m.PolicyId = types.StringValue(resp.Id)
	m.Name = types.StringValue(resp.Name)
	m.Namespace = types.StringPointerValue(resp.Namespace)
	m.AgentPolicyId = types.StringPointerValue(resp.PolicyId)
	m.Description = types.StringPointerValue(resp.Description)
	m.Enabled = types.BoolPointerValue(resp.Enabled)
	m.IntegrationName = types.StringValue(resp.Package.Name)
	m.IntegrationVersion = types.StringValue(resp.Package.Version)
	m.Inputs = map[string]integrationPolicyInputModel{}
	m.VarsJson = fromApiIntegrationPolicyVars(m.VarsJson, resp.Vars, secretStore)

	for _, inputResp := range resp.Inputs {
		inputName := fmt.Sprintf("%s-%s", inputResp.PolicyTemplate, inputResp.Type)
		m.Inputs[inputName] = integrationPolicyInputModel{
			Enabled: types.BoolValue(inputResp.Enabled),
			Streams: map[string]integrationPolicyStreamModel{},
			VarsJson: fromApiIntegrationPolicyVars(
				originalInputs[inputName].VarsJson,
				inputResp.Vars,
				secretStore,
			),
		}

		for _, streamResp := range inputResp.Streams {
			streamName := streamResp.DataStream.Dataset
			m.Inputs[inputName].Streams[streamName] = integrationPolicyStreamModel{
				Enabled: types.BoolValue(streamResp.Enabled),
				VarsJson: fromApiIntegrationPolicyVars(
					originalInputs[inputName].Streams[streamName].VarsJson,
					streamResp.Vars,
					secretStore,
				),
			}
		}
	}

	secretStore.Save(ctx, private)

	return diags
}

func fromApiIntegrationPolicyVars(
	dataVarsJson jsontypes.Normalized,
	respVars map[string]fleet.PackagePolicyVar,
	secretStore secretStore,
) jsontypes.Normalized {
	if dataVarsJson.IsUnknown() || dataVarsJson.IsNull() {
		return jsontypes.NewNormalizedNull()
	}

	vars := make(map[string]any)

	// Unmarshal dataVarsJson
	var dataVars map[string]any
	if diags := dataVarsJson.Unmarshal(&dataVars); diags != nil {
		panic(fmt.Errorf("%s: %s", diags[0].Summary(), diags[0].Detail()))
	}

	// Check each response var for a secret ref
	// Fetch the actual value from the store.
	for varName, respVar := range respVars {
		switch v := respVar.Value.(type) {
		case map[string]any:
			// Does this match {"isSecretRef": true}?
			isSecretRef, ok := v["isSecretRef"].(bool)
			if !ok || !isSecretRef {
				vars[varName] = v
				continue
			}

			// Assumed to exist if isSecretRef is true
			secretId := v["id"].(string)

			// If this is a create/Update operation, dataVars will have the
			// original secret value. Save it to the store for read operations.
			if varValue, ok := dataVars[varName]; ok && varValue != nil {
				secretStore[varName] = varValue
			}

			// If the secret => real value exists, assign it.
			// Otherwise use the entire map value to trigger an update.
			if secret, ok := secretStore[secretId]; ok {
				vars[varName] = secret
			} else {
				vars[varName] = v
			}
		default:
			// Many optional vars will default to null, ignore them.
			if v != nil {
				vars[varName] = v
			}
		}
	}

	// Marshal and return
	varsJson := lo.Must(json.Marshal(vars))
	return jsontypes.NewNormalizedValue(string(varsJson))
}

func newSecretStore(ctx context.Context, private privateData) secretStore {
	bytes, diags := private.GetKey(ctx, "secrets")
	if diags != nil {
		panic(fmt.Errorf("%s: %s", diags[0].Summary(), diags[0].Detail()))
	}
	if bytes == nil {
		bytes = []byte("{}")
	}

	var data secretStore
	lo.Must0(json.Unmarshal(bytes, &data))
	return data
}

type secretStore map[string]any

// RemoveAllExcept only keeps the given refIds in the secretStore.
func (s secretStore) RemoveAllExcept(refs []fleet.PackagePolicySecretRef) {
	ids := make([]string, 0, len(refs))
	for idx, secretRef := range refs {
		ids[idx] = secretRef.Id
	}
	for id := range s {
		if !slices.Contains(ids, id) {
			delete(s, id)
		}
	}
}

// Save marshals the secretStore back to the provider.
func (s secretStore) Save(ctx context.Context, private privateData) {
	bytes := lo.Must(json.Marshal(s))
	private.SetKey(ctx, "secrets", bytes)
}

// Equivalent to privatestate.ProviderData
type privateData interface {
	// GetKey returns the private state data associated with the given key.
	//
	// If the key is reserved for framework usage, an error diagnostic
	// is returned. If the key is valid, but private state data is not found,
	// nil is returned.
	//
	// The naming of keys only matters in context of a single resource,
	// however care should be taken that any historical keys are not reused
	// without accounting for older resource instances that may still have
	// older data at the key.
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)

	// SetKey sets the private state data at the given key.
	//
	// If the key is reserved for framework usage, an error diagnostic
	// is returned. The data must be valid JSON and UTF-8 safe or an error
	// diagnostic is returned.
	//
	// The naming of keys only matters in context of a single resource,
	// however care should be taken that any historical keys are not reused
	// without accounting for older resource instances that may still have
	// older data at the key.
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}
