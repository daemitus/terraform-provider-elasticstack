package kibana

import (
	"context"
	"encoding/json"
	"net/http"

	kbapi "github.com/elastic/terraform-provider-elasticstack/generated/kibana"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func SearchConnectors(ctx context.Context, client *Client, spaceID string, connectorName string, connectorTypeID string) ([]kbapi.ConnectorResponseProperties, diag.Diagnostics) {
	resp, err := client.API.GetAllConnectorsWithResponse(ctx, spaceID, nil)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return filterSearchResults(resp.JSON200, connectorTypeID, connectorName)
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// ReadConnector reads a specific connector from the API.
func ReadConnector(ctx context.Context, client Client, spaceID string, connectorID string) (*kbapi.ConnectorResponseProperties, diag.Diagnostics) {
	resp, err := client.API.GetConnectorWithResponse(ctx, spaceID, connectorID, nil)
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

// CreateConnector creates a new connector.
func CreateConnector(ctx context.Context, client *Client, spaceID string, connectorID string, req kbapi.CreateConnectorJSONRequestBody) (*kbapi.ConnectorResponseProperties, diag.Diagnostics) {
	resp, err := client.API.CreateConnectorWithResponse(ctx, spaceID, connectorID, nil, req)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// UpdateConnector updates an existing connector.
func UpdateConnector(ctx context.Context, client *Client, spaceID string, connectorID string, req kbapi.UpdateConnectorJSONRequestBody) (*kbapi.ConnectorResponseProperties, diag.Diagnostics) {
	resp, err := client.API.UpdateConnectorWithResponse(ctx, spaceID, connectorID, nil, req)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// DeleteConnector deletes an existing connector.
func DeleteConnector(ctx context.Context, client *Client, spaceID string, connectorID string) diag.Diagnostics {
	resp, err := client.API.DeleteConnectorWithResponse(ctx, spaceID, connectorID, nil)
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

func SearchConnectorsDefault(ctx context.Context, client *Client, connectorName string, connectorTypeID string) ([]kbapi.ConnectorResponseProperties, diag.Diagnostics) {
	resp, err := client.API.GetAllConnectorsDefaultWithResponse(ctx, nil)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return filterSearchResults(resp.JSON200, connectorTypeID, connectorName)
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// ReadConnectorDefault reads a specific connector from the API.
func ReadConnectorDefault(ctx context.Context, client Client, connectorID string) (*kbapi.ConnectorResponseProperties, diag.Diagnostics) {
	resp, err := client.API.GetConnectorDefaultWithResponse(ctx, connectorID, nil)
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

// CreateConnectorDefault creates a new connector.
func CreateConnectorDefault(ctx context.Context, client *Client, connectorID string, req kbapi.CreateConnectorDefaultJSONRequestBody) (*kbapi.ConnectorResponseProperties, diag.Diagnostics) {
	resp, err := client.API.CreateConnectorDefaultWithResponse(ctx, connectorID, nil, req)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// UpdateConnectorDefault updates an existing connector.
func UpdateConnectorDefault(ctx context.Context, client *Client, connectorID string, req kbapi.UpdateConnectorDefaultJSONRequestBody) (*kbapi.ConnectorResponseProperties, diag.Diagnostics) {
	resp, err := client.API.UpdateConnectorDefaultWithResponse(ctx, connectorID, nil, req)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.JSON200, nil
	default:
		return nil, reportUnknownError(resp.StatusCode(), resp.Body)
	}
}

// DeleteConnectorDefault deletes an existing connector.
func DeleteConnectorDefault(ctx context.Context, client *Client, connectorID string) diag.Diagnostics {
	resp, err := client.API.DeleteConnectorDefaultWithResponse(ctx, connectorID, nil)
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

func filterSearchResults(results *[]kbapi.ConnectorResponseProperties, connectorTypeID string, connectorName string) ([]kbapi.ConnectorResponseProperties, diag.Diagnostics) {
	if results == nil {
		return nil, nil
	}

	type Connector struct {
		Name            string `json:"name"`
		ConnectorTypeID string `json:"connector_type_id"`
	}

	matches := make([]kbapi.ConnectorResponseProperties, 0)

	for _, union := range *results {
		bytes, err := union.MarshalJSON()
		if err != nil {
			return nil, diag.FromErr(err)
		}

		var base Connector
		err = json.Unmarshal(bytes, &base)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		if connectorTypeID != "" && connectorTypeID != base.ConnectorTypeID {
			continue
		}

		if connectorName != "" && connectorName != base.Name {
			continue
		}

		matches = append(matches, union)
	}
	return matches, nil
}
