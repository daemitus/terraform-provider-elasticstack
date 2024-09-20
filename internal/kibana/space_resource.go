package kibana

import (
	"context"
	"fmt"
	"regexp"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/kibana"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource                = &SpaceResource{}
	_ resource.ResourceWithImportState = &SpaceResource{}
)

func NewSpaceResource(client *clients.KibanaClient) *SpaceResource {
	return &SpaceResource{client: client}
}

type SpaceResource struct {
	client *clients.KibanaClient
}

func (r *SpaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "kibana_space")
}

func (r *SpaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Creates a Kibana space. See https://www.elastic.co/guide/en/kibana/current/spaces-kibana.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The space ID that is part of the Kibana URL when inside the space.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The display name for the space.",
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description for the space.",
			Optional:    true,
		},
		"disabled_features": schema.ListAttribute{
			Description: "The list of disabled features for the space. To get a list of available feature IDs, use the Features API (https://www.elastic.co/guide/en/kibana/master/features-api-get.html).",
			Computed:    true,
			Optional:    true,
			ElementType: types.StringType,
		},
		"initials": schema.StringAttribute{
			Description: "The initials shown in the space avatar. By default, the initials are automatically generated from the space name. Initials must be 1 or 2 characters.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 2),
			},
		},
		"color": schema.StringAttribute{
			Description: "The hexadecimal color code used in the space avatar. By default, the color is automatically generated from the space name.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile("^#[0-9a-fA-F]{6}$"), "",
				),
			},
		},
		"image_url": schema.StringAttribute{
			Description: "The data-URL encoded image to display in the space avatar.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile("^data:image/.+$"), "",
				),
			},
		},
	}
}

func (r *SpaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *SpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data spaceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	spaceId := data.ID.ValueString()
	space, diags := r.client.ReadSpace(ctx, spaceId)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(space)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *SpaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data spaceModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, diags := r.client.CreateSpace(ctx, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(space)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *SpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data spaceModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.ID.ValueString()
	space, diags := r.client.UpdateSpace(ctx, spaceId, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(space)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *SpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data spaceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	spaceId := data.ID.ValueString()
	diags = r.client.DeleteSpace(ctx, spaceId)
	resp.Diagnostics.Append(diags...)
}

type spaceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	DisabledFeatures types.List   `tfsdk:"disabled_features"`
	Initials         types.String `tfsdk:"initials"`
	Color            types.String `tfsdk:"color"`
	ImageURL         types.String `tfsdk:"image_url"`
}

func (m *spaceModel) toApi(ctx context.Context) (*kibana.Space, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	output := &kibana.Space{
		ID:               m.ID.ValueString(),
		Name:             m.Name.ValueString(),
		Description:      m.Description.ValueString(),
		DisabledFeatures: util.ListTypeToSliceBasic[string](ctx, m.DisabledFeatures, path.AtName("disabled_features"), diags),
		Initials:         m.Initials.ValueString(),
		Color:            m.Color.ValueString(),
		ImageURL:         m.ImageURL.ValueString(),
	}

	return output, diags
}

func (m *spaceModel) fromApi(resp *kibana.Space) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.ID = types.StringValue(resp.ID)
	m.Name = types.StringValue(resp.Name)
	m.Description = types.StringValue(resp.Description)
	m.DisabledFeatures = util.SliceToListType_String(resp.DisabledFeatures, path.AtName("disabled_features"), diags)
	m.Initials = types.StringValue(resp.Initials)
	m.Color = types.StringValue(resp.Color)
	m.ImageURL = types.StringValue(resp.ImageURL)

	return diags
}
