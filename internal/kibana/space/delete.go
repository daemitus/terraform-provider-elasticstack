package space

import (
	"context"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *spaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateModel spaceModel

	diags := req.State.Get(ctx, &stateModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.client.GetKibana2Client()
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	spaceId := stateModel.SpaceID.ValueString()
	diags = kibana2.DeleteSpace(ctx, client, spaceId)
	resp.Diagnostics.Append(diags...)
}
