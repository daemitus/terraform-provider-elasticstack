package kibana2

import (
	"context"
	"net/http"

	kbapi "github.com/elastic/terraform-provider-elasticstack/generated/kibana"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// GetSpaces reads all spaces from the API.
func GetSpaces(ctx context.Context, client *Client) ([]kbapi.KibanaSpace, diag.Diagnostics) {
	resp, err := client.API.GetSpacesSpaceWithResponse(ctx, nil)
	if err != nil {
		return nil, utils.FrameworkDiagFromError(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return *resp.JSON200, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// GetSpace reads a specific space from the API.
func GetSpace(ctx context.Context, client *Client, spaceID string) (*kbapi.KibanaSpace, diag.Diagnostics) {
	resp, err := client.API.GetSpacesSpaceIdWithResponse(ctx, spaceID)
	if err != nil {
		return nil, utils.FrameworkDiagFromError(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// CreateSpace creates a new space.
func CreateSpace(ctx context.Context, client *Client, req kbapi.PostSpacesSpaceJSONRequestBody) (*kbapi.KibanaSpace, diag.Diagnostics) {
	resp, err := client.API.PostSpacesSpaceWithResponse(ctx, req)
	if err != nil {
		return nil, utils.FrameworkDiagFromError(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// UpdateSpace updates an existing space.
func UpdateSpace(ctx context.Context, client *Client, spaceID string, req kbapi.PutSpacesSpaceIdJSONRequestBody) (*kbapi.KibanaSpace, diag.Diagnostics) {
	resp, err := client.API.PutSpacesSpaceIdWithResponse(ctx, spaceID, req)
	if err != nil {
		return nil, utils.FrameworkDiagFromError(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// DeleteSpace deletes an existing space.
func DeleteSpace(ctx context.Context, client *Client, spaceID string) diag.Diagnostics {
	resp, err := client.API.DeleteSpacesSpaceIdWithResponse(ctx, spaceID)
	if err != nil {
		return utils.FrameworkDiagFromError(err)
	}

	switch resp.StatusCode() {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return reportUnknownError(resp.StatusCode(), resp.Body)
	}
}
