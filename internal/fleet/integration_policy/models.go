package integration_policy

import (
	"context"
	"encoding/json"
	"sort"

	fleetapi "github.com/elastic/terraform-provider-elasticstack/generated/fleet"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type integrationPolicyModel struct {
	ID                 types.String         `tfsdk:"id"`
	PolicyID           types.String         `tfsdk:"policy_id"`
	Name               types.String         `tfsdk:"name"`
	Namespace          types.String         `tfsdk:"namespace"`
	AgentPolicyID      types.String         `tfsdk:"agent_policy_id"`
	Description        types.String         `tfsdk:"description"`
	Enabled            types.Bool           `tfsdk:"enabled"`
	Force              types.Bool           `tfsdk:"force"`
	IntegrationName    types.String         `tfsdk:"integration_name"`
	IntegrationVersion types.String         `tfsdk:"integration_version"`
	Input              types.List           `tfsdk:"input"` //> integrationPolicyInputModel
	VarsJson           jsontypes.Normalized `tfsdk:"vars_json"`
}

type integrationPolicyInputModel struct {
	InputID     types.String         `tfsdk:"input_id"`
	Enabled     types.Bool           `tfsdk:"enabled"`
	StreamsJson jsontypes.Normalized `tfsdk:"streams_json"`
	VarsJson    jsontypes.Normalized `tfsdk:"vars_json"`
}

func (model *integrationPolicyModel) populateFromAPI(ctx context.Context, data *fleetapi.PackagePolicy) diag.Diagnostics {
	if data == nil {
		return nil
	}

	var diags diag.Diagnostics

	model.ID = types.StringValue(data.Id)
	model.PolicyID = types.StringValue(data.Id)
	model.Name = types.StringValue(data.Name)
	model.Namespace = types.StringPointerValue(data.Namespace)
	model.AgentPolicyID = types.StringPointerValue(data.PolicyId)
	model.Description = types.StringPointerValue(data.Description)
	model.Enabled = types.BoolPointerValue(data.Enabled)
	model.IntegrationName = types.StringValue(data.Package.Name)
	model.IntegrationVersion = types.StringValue(data.Package.Version)
	model.VarsJson = utils.MapToNormalizedType(data.Vars, path.Root("vars_json"), diags)

	{
		newInputs := utils.TransformMapToSlice(data.Inputs, path.Root("input"), diags,
			func(inputData fleetapi.PackagePolicyInput, meta utils.MapMeta) integrationPolicyInputModel {
				return integrationPolicyInputModel{
					InputID:     types.StringValue(meta.Key),
					Enabled:     types.BoolValue(inputData.Enabled),
					StreamsJson: utils.MapToNormalizedType(inputData.Streams, meta.Path.AtName("streams_json"), diags),
					VarsJson:    utils.MapToNormalizedType(inputData.Vars, meta.Path.AtName("vars_json"), diags),
				}
			})
		if newInputs == nil {
			model.Input = types.ListNull(getInputType())
		} else {
			oldInputs := utils.ListTypeAs[integrationPolicyInputModel](ctx, model.Input, path.Root("input"), diags)
			sortInputs(newInputs, oldInputs)

			inputList, d := types.ListValueFrom(ctx, getInputType(), newInputs)
			diags.Append(d...)

			model.Input = inputList
		}
	}

	return diags
}

func (model integrationPolicyModel) toAPIModel(ctx context.Context, isUpdate bool) (fleetapi.PackagePolicyRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := fleetapi.PackagePolicyRequest{
		Description: model.Description.ValueStringPointer(),
		Force:       model.Force.ValueBoolPointer(),
		Name:        model.Name.ValueString(),
		Namespace:   model.Namespace.ValueStringPointer(),
		Package: struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{
			Name:    model.IntegrationName.ValueString(),
			Version: model.IntegrationVersion.ValueString(),
		},
		PolicyId: model.AgentPolicyID.ValueString(),
		Vars:     utils.NormalizedTypeToMap[any](model.VarsJson, path.Root("vars_json"), diags),
	}

	if isUpdate {
		body.Id = model.ID.ValueStringPointer()
	}

	body.Inputs = utils.ListTypeToMap(ctx, model.Input, path.Root("input"), diags,
		func(inputModel integrationPolicyInputModel, meta utils.ListMeta) (string, fleetapi.PackagePolicyRequestInput) {
			return inputModel.InputID.ValueString(), fleetapi.PackagePolicyRequestInput{
				Enabled: inputModel.Enabled.ValueBoolPointer(),
				Streams: utils.NormalizedTypeToMap[fleetapi.PackagePolicyRequestInputStream](inputModel.StreamsJson, meta.Path.AtName("streams_json"), diags),
				Vars:    utils.NormalizedTypeToMap[any](inputModel.VarsJson, meta.Path.AtName("vars_json"), diags),
			}
		})

	return body, diags
}

// sortInputs will sort the 'incoming' list of input definitions based on
// the order of inputs defined in the 'existing' list. Inputs not present in
// 'existing' will be placed at the end of the list. Inputs are identified by
// their ID ('input_id'). The 'incoming' slice will be sorted in-place.
func sortInputs(incoming []integrationPolicyInputModel, existing []integrationPolicyInputModel) {
	idToIndex := make(map[string]int, len(existing))
	for index, inputData := range existing {
		inputID := inputData.InputID.ValueString()
		idToIndex[inputID] = index
	}

	sort.Slice(incoming, func(i, j int) bool {
		iID := incoming[i].InputID.ValueString()
		iIdx, ok := idToIndex[iID]
		if !ok {
			return false
		}

		jID := incoming[j].InputID.ValueString()
		jIdx, ok := idToIndex[jID]
		if !ok {
			return true
		}

		return iIdx < jIdx
	})
}

// The secret store is a map of policy secret reference IDs to the
// original value at time of creation. By replacing the ref when
// marshaling the state back to Terraform, we can prevent resource
// drift.
type secretStore map[string]any

// newSecretStore creates a new secretStore from the resource privateData.
func newSecretStore(ctx context.Context, private privateData) (store secretStore, diags diag.Diagnostics) {
	bytes, diags := private.GetKey(ctx, "secrets")
	if diags != nil {
		return
	}
	if bytes == nil {
		bytes = []byte("{}")
	}

	err := json.Unmarshal(bytes, &store)
	if err != nil {
		diags.AddError("could not unmarshal secret store", err.Error())
		return
	}

	return
}

// Save marshals the secretStore back to the provider.
func (s secretStore) Save(ctx context.Context, private privateData) (diags diag.Diagnostics) {
	bytes, err := json.Marshal(s)
	if err != nil {
		diags.AddError("could not marshal secret store", err.Error())
		return
	}

	return private.SetKey(ctx, "secrets", bytes)
}

// extractVars first extracts each wrapped var from the response,
// then replaces any secretRefs with the original value from the
// secret store if it exists.
func extractVars(resp *fleetapi.PackagePolicy, secrets secretStore) {
	// First unwrap the response from the {"value": ...} struct. Then for
	// any values that have `isSecretRef` set, replace with the original
	// value from reqVars. If req is empty, fetch from the store instead.
	handleVars := func(vars map[string]any) {
		for key, varAny := range vars {
			// All returned vars have the following structure:
			// {"type": "...", "value": <any>}
			if varMap, ok := varAny.(map[string]any); ok {
				if _, ok := varMap["type"]; ok {
					if wrappedVal, ok := varMap["value"]; ok {
						vars[key] = wrappedVal
						varAny = wrappedVal
					} else {
						// No need to keep null values
						delete(vars, key)
					}
				}
			}

			// Policy secrets have the following struture:
			// {"id": "...", "isSecretRef": true}
			if rmap, ok := varAny.(map[string]any); ok {
				if isRef, ok := rmap["isSecretRef"].(bool); ok && isRef {
					refID := rmap["id"].(string)
					if original, ok := secrets[refID]; ok {
						vars[key] = original
					}
				}
			}
		}
	}

	handleVars(resp.Vars)
	for _, input := range resp.Inputs {
		handleVars(input.Vars)
		for _, streamAny := range input.Streams {
			stream := streamAny.(map[string]any)
			streamVars := stream["vars"].(map[string]any)
			handleVars(streamVars)
		}
	}
}

// saveVars first extracts each wrapped var from the response,
// then replaces/saves the original value from the request with the ref from
// the response.
func saveVars(req fleetapi.PackagePolicyRequest, resp *fleetapi.PackagePolicy, secrets secretStore) {
	// Prune the store and only keep the currently used refs.
	if resp.SecretReferences != nil {
		for _, ref := range resp.SecretReferences {
			if ref.Id != nil {
				delete(secrets, *ref.Id)
			}
		}
	}

	// First unwrap the response from the {"value": ...} struct. Then for
	// any values that have `isSecretRef` set, replace with the original
	// value from reqVars. If req is empty, fetch from the store instead.
	handleVars := func(reqVars map[string]any, respVars map[string]any) {
		for key, respVar := range respVars {
			// All returned vars have the following structure:
			// {"type": "...", "value": <any>}
			if respVarMap, ok := respVar.(map[string]any); ok {
				if _, ok := respVarMap["type"]; ok {
					if wrappedVal, ok := respVarMap["value"]; ok {
						respVars[key] = wrappedVal
						respVar = wrappedVal
					} else {
						// No need to keep null values
						delete(respVars, key)
					}
				}
			}

			// Policy secrets have the following struture:
			// {"id": "...", "isSecretRef": true}
			if rmap, ok := respVar.(map[string]any); ok {
				if isRef, ok := rmap["isSecretRef"].(bool); ok && isRef {
					refID := rmap["id"].(string)
					originalVal := reqVars[key]
					secrets[refID] = originalVal
					respVars[key] = originalVal
				}
			}
		}
	}

	handleVars(req.Vars, resp.Vars)
	for inputID, inputReq := range req.Inputs {
		inputResp := resp.Inputs[inputID]
		handleVars(inputReq.Vars, inputResp.Vars)
		for streamID, streamReq := range inputReq.Streams {
			streamResp := inputResp.Streams[streamID].(map[string]any)
			streamRespVars := streamResp["vars"].(map[string]any)
			handleVars(streamReq.Vars, streamRespVars)
		}
	}
}
