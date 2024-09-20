package kibana

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/kibana"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                     = &SavedObjectResource{}
	_ resource.ResourceWithConfigValidators = &SavedObjectResource{}
)

func NewImportSavedObjectsResource(client *clients.KibanaClient) *SavedObjectResource {
	return &SavedObjectResource{client: client}
}

type SavedObjectResource struct {
	client *clients.KibanaClient
}

func (r *SavedObjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "kibana_import_saved_objects")
}

func (r *SavedObjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Imports saved objects from the referenced file."
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of this resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"space_id": schema.StringAttribute{
			Description: "An identifier for the space. If space_id is not provided, the default space is used.",
			Computed:    true,
			Optional:    true,
			Default:     stringdefault.StaticString("default"),
		},
		"ignore_import_errors": schema.BoolAttribute{
			Description: "If set to true, errors during the import process will not fail the configuration application",
			Optional:    true,
		},
		"create_new_copies": schema.BoolAttribute{
			Description: "Creates copies of saved objects, regenerates each object ID, and resets the origin. When used, potential conflict errors are avoided.",
			Optional:    true,
		},
		"overwrite": schema.BoolAttribute{
			Description: "Overwrites saved objects when they already exist. When used, potential conflict errors are automatically resolved by overwriting the destination object.",
			Optional:    true,
		},
		"compatibility_mode": schema.BoolAttribute{
			Description: "Applies various adjustments to the saved objects that are being imported to maintain compatibility between different Kibana versions. Use this option only if you encounter issues with imported saved objects.",
			Optional:    true,
		},
		"file_contents": schema.StringAttribute{
			Description: "The contents of the exported saved objects file.",
			Required:    true,
		},
		"success": schema.BoolAttribute{
			Description: "Indicates when the import was successfully completed. When set to false, some objects may not have been created. For additional information, refer to the errors and success_results properties.",
			Computed:    true,
		},
		"success_count": schema.Int64Attribute{
			Description: "Indicates the number of successfully imported records.",
			Computed:    true,
		},
		"success_results": schema.ListAttribute{
			Description: "Details for each successfully imported record.",
			Computed:    true,
			ElementType: importSavedObjectSuccessType,
		},
		"errors": schema.ListAttribute{
			Description: "Details for each imported record failure.",
			Computed:    true,
			ElementType: importSavedObjectErrorType,
		},
	}
}

func (r *SavedObjectResource) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("create_new_copies"),
			path.MatchRoot("overwrite"),
			path.MatchRoot("compatibility_mode"),
		),
	}
}

func (r *SavedObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Read is not supported for elasticstack_kibana_import_saved_objects")
}

func (r *SavedObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data importSavedObjectModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	r.importSavedObjects(ctx, &data, resp.Diagnostics)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *SavedObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data importSavedObjectModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	r.importSavedObjects(ctx, &data, resp.Diagnostics)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *SavedObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Delete is not supported for elasticstack_kibana_import_saved_objects")
}

func (r *SavedObjectResource) importSavedObjects(ctx context.Context, data *importSavedObjectModel, diags diag.Diagnostics) {
	space := data.SpaceID.ValueString()
	params := &kibana.ImportSavedObjectsParams{
		CompatibilityMode: data.CompatibilityMode.ValueBoolPointer(),
		CreateNewCopies:   data.CreateNewCopies.ValueBoolPointer(),
		Overwrite:         data.Overwrite.ValueBoolPointer(),
	}
	body := []byte(data.FileContents.ValueString())

	resp, d := r.client.ImportSavedObjects(ctx, space, body, params)
	diags.Append(d...)
	if d.HasError() {
		return
	}

	d = data.fromApi(resp)
	diags.Append(d...)
	if d.HasError() {
		return
	}

	if !resp.Success && !data.IgnoreImportErrors.ValueBool() {
		diags.AddError("not all objects were imported successfully", "see errors attribute for more details")
	}
}

type importSavedObjectModel struct {
	ID                 types.String `tfsdk:"id"`
	SpaceID            types.String `tfsdk:"space_id"`
	IgnoreImportErrors types.Bool   `tfsdk:"ignore_import_errors"`
	CreateNewCopies    types.Bool   `tfsdk:"create_new_copies"`
	Overwrite          types.Bool   `tfsdk:"overwrite"`
	CompatibilityMode  types.Bool   `tfsdk:"compatibility_mode"`
	FileContents       types.String `tfsdk:"file_contents"`
	Success            types.Bool   `tfsdk:"success"`
	SuccessCount       types.Int64  `tfsdk:"success_count"`
	SuccessResults     types.List   `tfsdk:"success_results"`
	Errors             types.List   `tfsdk:"errors"`
}

var importSavedObjectSuccessType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":             types.StringType,
		"type":           types.StringType,
		"destination_id": types.StringType,
		"meta": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"icon":  types.StringType,
				"title": types.StringType,
			},
		},
	},
}
var importSavedObjectErrorType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":    types.StringType,
		"type":  types.StringType,
		"title": types.StringType,
		"error": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type": types.StringType,
			},
		},
		"meta": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"icon":  types.StringType,
				"title": types.StringType,
			},
		},
	},
}

func (m *importSavedObjectModel) fromApi(resp *kibana.ImportSavedObjectsResponse) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	if m.ID.IsUnknown() || m.ID.IsNull() {
		m.ID = types.StringValue(uuid.New().String())
	}

	m.Success = types.BoolValue(resp.Success)
	m.SuccessCount = types.Int64Value(resp.SuccessCount)
	m.SuccessResults = util.SliceToListType(resp.SuccessResults, importSavedObjectSuccessType, path.AtName("success_results"), diags,
		func(model kibana.ImportSavedObjectsSuccess, index int) attr.Value {
			path := path.AtName("success_results").AtListIndex(index)
			return util.MapToObjectType(importSavedObjectSuccessType.AttrTypes, map[string]attr.Value{
				"id":             types.StringValue(model.ID),
				"type":           types.StringValue(model.Type),
				"destination_id": types.StringValue(model.DestinationID),
				"meta": util.MapToObjectType(importSavedObjectSuccessType.AttrTypes["meta"].(types.ObjectType).AttrTypes, map[string]attr.Value{
					"icon":  types.StringValue(model.Meta.Icon),
					"title": types.StringValue(model.Meta.Title),
				}, path.AtMapKey("meta"), diags),
			}, path, diags)
		},
	)
	m.Errors = util.SliceToListType(
		resp.Errors, importSavedObjectErrorType, path.AtName("errors"), diags,
		func(model kibana.ImportSavedObjectsError, index int) attr.Value {
			path := path.AtName("errors").AtListIndex(index)
			return util.MapToObjectType(importSavedObjectSuccessType.AttrTypes, map[string]attr.Value{
				"id":   types.StringValue(model.ID),
				"type": types.StringValue(model.Type),
				"error": util.MapToObjectType(importSavedObjectSuccessType.AttrTypes["error"].(types.ObjectType).AttrTypes, map[string]attr.Value{
					"type": types.StringValue(model.Error.Type),
				}, path.AtMapKey("error"), diags),
				"meta": util.MapToObjectType(importSavedObjectSuccessType.AttrTypes["meta"].(types.ObjectType).AttrTypes, map[string]attr.Value{
					"icon":  types.StringValue(model.Meta.Icon),
					"title": types.StringValue(model.Meta.Title),
				}, path.AtMapKey("meta"), diags),
			}, path, diags)
		},
	)

	return diags
}
