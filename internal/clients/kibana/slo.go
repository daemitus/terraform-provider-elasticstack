package kibana

import (
	"context"
	"net/http"

	kbapi "github.com/elastic/terraform-provider-elasticstack/generated/kibana"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// GetSlo reads a specific SLO from the API.
func GetSlo(ctx context.Context, client *Client, spaceID string, sloID string) (*kbapi.SLOsSloWithSummaryResponse, diag.Diagnostics) {
	params := kbapi.GetSloOpParams{}

	resp, err := client.API.GetSloOpWithResponse(ctx, spaceID, sloID, &params)
	if err != nil {
		return nil, diag.FromErr(err)
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

// CreateSlo creates a new data view.
func CreateSlo(ctx context.Context, client *Client, spaceID string, req kbapi.SLOsCreateSloRequest) (*kbapi.SLOsSloWithSummaryResponse, diag.Diagnostics) {
	resp, err := client.API.CreateSloOpWithResponse(ctx, spaceID, req)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return GetSlo(ctx, client, spaceID, resp.JSON200.Id)
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// UpdateSlo updates an existing data view.
func UpdateSlo(ctx context.Context, client *Client, spaceID string, sloID string, req kbapi.SLOsUpdateSloRequest) (*kbapi.SLOsSloWithSummaryResponse, diag.Diagnostics) {
	resp, err := client.API.UpdateSloOpWithResponse(ctx, spaceID, sloID, req)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return GetSlo(ctx, client, spaceID, resp.JSON200.Id)
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// DeleteSlo deletes an existing data view.
func DeleteSlo(ctx context.Context, client *Client, spaceID string, sloID string) diag.Diagnostics {
	resp, err := client.API.DeleteSloOpWithResponse(ctx, spaceID, sloID)
	if err != nil {
		return diag.FromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return reportUnknownError(resp.StatusCode(), resp.Body)
	}
}
