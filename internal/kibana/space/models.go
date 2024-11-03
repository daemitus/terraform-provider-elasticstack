package space

import (
	"context"

	kbapi "github.com/elastic/terraform-provider-elasticstack/generated/kibana"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type spaceModel struct {
	ID               types.String `tfsdk:"id"`
	SpaceID          types.String `tfsdk:"space_id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	DisabledFeatures types.Set    `tfsdk:"disabled_features"`
	Initials         types.String `tfsdk:"initials"`
	Color            types.String `tfsdk:"color"`
	ImageUrl         types.String `tfsdk:"image_url"`
}

func (m *spaceModel) populateFromAPI(ctx context.Context, data *kbapi.KibanaSpace) diag.Diagnostics {
	if data == nil {
		return nil
	}

	var diags diag.Diagnostics

	m.ID = types.StringValue(data.Id)
	m.SpaceID = types.StringValue(data.Id)
	m.Name = types.StringValue(data.Name)
	m.Description = types.StringPointerValue(data.Description)
	m.Initials = types.StringPointerValue(data.Initials)
	m.Color = types.StringPointerValue(data.Color)
	m.ImageUrl = types.StringPointerValue(data.ImageUrl)
	m.DisabledFeatures = utils.SemanticEqualEmptySet(ctx, m.DisabledFeatures,
		utils.SliceToSetType_String(ctx, utils.Deref(data.DisabledFeatures), path.Root("disabled_features"), &diags))

	return diags
}

func (m spaceModel) toAPIModel(ctx context.Context) (kbapi.KibanaSpace, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := kbapi.KibanaSpace{
		Color:            utils.ValueStringPointer(m.Color),
		Description:      m.Description.ValueStringPointer(),
		DisabledFeatures: utils.SliceRef(utils.SetTypeToSlice_String(ctx, m.DisabledFeatures, path.Root("disabled_features"), &diags)),
		Id:               m.SpaceID.ValueString(),
		ImageUrl:         utils.ValueStringPointer(m.ImageUrl),
		Initials:         utils.ValueStringPointer(m.Initials),
		Name:             m.Name.ValueString(),
	}

	return body, diags
}
