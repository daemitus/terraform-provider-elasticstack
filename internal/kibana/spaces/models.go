package spaces

import (
	"context"

	kbapi "github.com/elastic/terraform-provider-elasticstack/generated/kibana"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// spacesModel maps the data source schema data.
type spacesModel struct {
	ID     types.String `tfsdk:"id"`
	Spaces types.List   `tfsdk:"spaces"` //> spaceModel
}

// spaceModel maps spaces schema data.
type spaceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	DisabledFeatures types.List   `tfsdk:"disabled_features"`
	Initials         types.String `tfsdk:"initials"`
	Color            types.String `tfsdk:"color"`
	ImageUrl         types.String `tfsdk:"image_url"`
}

func (m *spacesModel) populateFromAPI(ctx context.Context, data kbapi.KibanaSpaces) diag.Diagnostics {
	if data == nil {
		return nil
	}

	var diags diag.Diagnostics

	m.ID = types.StringValue("spaces")
	m.Spaces = utils.SliceToListType(ctx, data, getSpaceType(), path.Root("spaces"), &diags,
		func(item kbapi.KibanaSpace, meta utils.ListMeta) spaceModel {
			return spaceModel{
				ID:               types.StringValue(item.Id),
				Name:             types.StringValue(item.Name),
				Description:      types.StringPointerValue(item.Description),
				Initials:         types.StringPointerValue(item.Initials),
				Color:            types.StringPointerValue(item.Color),
				ImageUrl:         types.StringPointerValue(item.ImageUrl),
				DisabledFeatures: utils.SliceToListType_String(ctx, utils.Deref(item.DisabledFeatures), meta.Path.AtName("disabled_features"), &diags),
			}
		})

	return diags
}
