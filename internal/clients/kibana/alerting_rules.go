package kibana

import (
	"context"
	"net/http"

	kbapi "github.com/elastic/terraform-provider-elasticstack/generated/kibana"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// ReadAlertingRule reads a specific alerting rule from the API.
func ReadAlertingRule(ctx context.Context, client *Client, spaceID string, ruleID string) (*kbapi.AlertingRuleResponseProperties, diag.Diagnostics) {
	resp, err := client.API.GetRuleWithResponse(ctx, spaceID, ruleID)
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

// CreateAlertingRule creates a new alerting rule.
func CreateAlertingRule(ctx context.Context, client *Client, spaceID string, req kbapi.AlertingCreateRuleRequest) (*kbapi.AlertingRuleResponseProperties, diag.Diagnostics) {
	resp, err := client.API.CreateRuleWithResponse(ctx, spaceID, req)
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

// UpdateAlertingRule updates an existing alerting rule.
func UpdateAlertingRule(ctx context.Context, client *Client, spaceID string, viewId string, req kbapi.AlertingUpdateRuleRequest) (*kbapi.AlertingRuleResponseProperties, diag.Diagnostics) {
	resp, err := client.API.UpdateRuleWithResponse(ctx, spaceID, viewId, req)
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

// DeleteAlertingRule deletes an existing alerting rule.
func DeleteAlertingRule(ctx context.Context, client *Client, spaceID string, viewId string) diag.Diagnostics {
	resp, err := client.API.DeleteRuleWithResponse(ctx, spaceID, viewId, nil)
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

// ReadAlertingRuleDefault reads a specific alerting rule from the API.
func ReadAlertingRuleDefault(ctx context.Context, client *Client, viewId string) (*kbapi.AlertingRuleResponseProperties, diag.Diagnostics) {
	resp, err := client.API.GetRuleDefaultWithResponse(ctx, viewId)
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

// CreateAlertingRuleDefault creates a new alerting rule.
func CreateAlertingRuleDefault(ctx context.Context, client *Client, spaceID string, req kbapi.AlertingCreateRuleRequest) (*kbapi.AlertingRuleResponseProperties, diag.Diagnostics) {
	resp, err := client.API.CreateRuleDefaultWithResponse(ctx, req)
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

// UpdateAlertingRuleDefault updates an existing alerting rule.
func UpdateAlertingRuleDefault(ctx context.Context, client *Client, spaceID string, viewId string, req kbapi.AlertingUpdateRuleRequest) (*kbapi.AlertingRuleResponseProperties, diag.Diagnostics) {
	resp, err := client.API.UpdateRuleDefaultWithResponse(ctx, viewId, req)
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

// DeleteAlertingRuleDefault deletes an existing alerting rule.
func DeleteAlertingRuleDefault(ctx context.Context, client *Client, ruleID string) diag.Diagnostics {
	resp, err := client.API.DeleteRuleDefaultWithResponse(ctx, ruleID, nil)
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
