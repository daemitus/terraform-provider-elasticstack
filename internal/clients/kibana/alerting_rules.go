package kibana

import (
	"context"
	"net/http"

	kbapi "github.com/elastic/terraform-provider-elasticstack/generated/kibana"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// ReadAlertingRule reads a specific alerting rule from the API.
func ReadAlertingRule(ctx context.Context, client *Client, spaceID string, ruleID string) (*kbapi.RuleResponse, diag.Diagnostics) {
	resp, err := client.API.GetRuleWithResponse(ctx, spaceID, ruleID, nil)
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
func CreateAlertingRule(ctx context.Context, client *Client, spaceID string, ruleID string, req kbapi.CreateRuleJSONRequestBody) (*kbapi.RuleResponse, diag.Diagnostics) {
	resp, err := client.API.CreateRuleWithResponse(ctx, spaceID, ruleID, nil, req)
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
func UpdateAlertingRule(ctx context.Context, client *Client, spaceID string, ruleID string, req kbapi.UpdateRuleJSONRequestBody) (*kbapi.RuleResponse, diag.Diagnostics) {
	resp, err := client.API.UpdateRuleWithResponse(ctx, spaceID, ruleID, nil, req)
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
func DeleteAlertingRule(ctx context.Context, client *Client, spaceID string, ruleID string) diag.Diagnostics {
	resp, err := client.API.DeleteRuleWithResponse(ctx, spaceID, ruleID, nil)
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
func ReadAlertingRuleDefault(ctx context.Context, client *Client, ruleID string) (*kbapi.GetRuleResponseObject, diag.Diagnostics) {
	resp, err := client.API.GetRuleDefaultWithResponse(ctx, ruleID, nil)
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
func CreateAlertingRuleDefault(ctx context.Context, client *Client, ruleID string, req kbapi.CreateRuleDefaultJSONRequestBody) (*kbapi.CreateRuleResponseObject, diag.Diagnostics) {
	resp, err := client.API.CreateRuleDefaultWithResponse(ctx, ruleID, nil, req)
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
func UpdateAlertingRuleDefault(ctx context.Context, client *Client, spaceID string, ruleID string, req kbapi.UpdateRuleDefaultJSONRequestBody) (*kbapi.UpdateRuleResponseObject, diag.Diagnostics) {
	resp, err := client.API.UpdateRuleDefaultWithResponse(ctx, ruleID, nil, req)
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
