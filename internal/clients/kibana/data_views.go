package kibana

import (
	"context"
	"net/http"

	kbapi "github.com/elastic/terraform-provider-elasticstack/generated/kibana"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ReadDataView reads a specific data view from the API.
func ReadDataView(ctx context.Context, client *Client, spaceID string, viewID string) (*kbapi.DataViewResponseObject, diag.Diagnostics) {
	resp, err := client.API.GetDataViewWithResponse(ctx, spaceID, viewID)
	if err != nil {
		return nil, fromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, reportUnknownErrorFw(resp.StatusCode(), resp.Body)
	}
}

// CreateDataView creates a new data view.
func CreateDataView(ctx context.Context, client *Client, spaceID string, req kbapi.CreateDataViewRequestObject) (*kbapi.DataViewResponseObject, diag.Diagnostics) {
	resp, err := client.API.CreateDataViewWithResponse(ctx, spaceID, req)
	if err != nil {
		return nil, fromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownErrorFw(resp.StatusCode(), resp.Body)
	}
}

// UpdateDataView updates an existing data view.
func UpdateDataView(ctx context.Context, client *Client, spaceID string, viewID string, req kbapi.UpdateDataViewRequestObject) (*kbapi.DataViewResponseObject, diag.Diagnostics) {
	resp, err := client.API.UpdateDataViewWithResponse(ctx, spaceID, viewID, req)
	if err != nil {
		return nil, fromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownErrorFw(resp.StatusCode(), resp.Body)
	}
}

// DeleteDataView deletes an existing data view.
func DeleteDataView(ctx context.Context, client *Client, spaceID string, viewID string) diag.Diagnostics {
	resp, err := client.API.DeleteDataViewWithResponse(ctx, spaceID, viewID, nil)
	if err != nil {
		return fromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return reportUnknownErrorFw(resp.StatusCode(), resp.Body)
	}
}

// ReadDataViewDefault reads a specific data view from the API.
func ReadDataViewDefault(ctx context.Context, client *Client, viewID string) (*kbapi.DataViewResponseObject, diag.Diagnostics) {
	resp, err := client.API.GetDataViewDefaultWithResponse(ctx, viewID)
	if err != nil {
		return nil, fromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, reportUnknownErrorFw(resp.StatusCode(), resp.Body)
	}
}

// CreateDataViewDefault creates a new data view.
func CreateDataViewDefault(ctx context.Context, client *Client, req kbapi.CreateDataViewRequestObject) (*kbapi.DataViewResponseObject, diag.Diagnostics) {
	resp, err := client.API.CreateDataViewDefaultWithResponse(ctx, req)
	if err != nil {
		return nil, fromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownErrorFw(resp.StatusCode(), resp.Body)
	}
}

// UpdateDataViewDefault updates an existing data view.
func UpdateDataViewDefault(ctx context.Context, client *Client, viewID string, req kbapi.UpdateDataViewRequestObject) (*kbapi.DataViewResponseObject, diag.Diagnostics) {
	resp, err := client.API.UpdateDataViewDefaultWithResponse(ctx, viewID, req)
	if err != nil {
		return nil, fromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownErrorFw(resp.StatusCode(), resp.Body)
	}
}

// DeleteDataViewDefault deletes an existing data view.
func DeleteDataViewDefault(ctx context.Context, client *Client, viewID string) diag.Diagnostics {
	resp, err := client.API.DeleteDataViewDefaultWithResponse(ctx, viewID, nil)
	if err != nil {
		return fromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return reportUnknownErrorFw(resp.StatusCode(), resp.Body)
	}
}
