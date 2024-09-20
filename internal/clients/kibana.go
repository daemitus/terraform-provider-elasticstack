package clients

import (
	"context"
	"log"
	"net/http"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/kibana"
	"github.com/daemitus/terraform-provider-elasticstack/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// KibanaClient provides an API client for Kibana.
type KibanaClient struct {
	baseClient
	API *kibana.Client
}

// NewKibanaClient creates a new Kibana API client.
func NewKibanaClient(config config.ServiceConfig) (*KibanaClient, error) {
	api, err := kibana.NewClient(config)
	if err != nil {
		return nil, err
	}
	client := KibanaClient{API: api}
	return &client, nil
}

// NewAccFleetClient creates a new Kibana API client for acceptance testing.
func NewAccKibanaClient() *KibanaClient {
	config, err := config.New().WithEnv()
	if err != nil {
		log.Fatal("Kibana config failure: %w", err)
	}

	client, err := NewKibanaClient(config.Kibana)
	if err != nil {
		log.Fatal("Kibana client failure: %w", err)
	}

	return client
}

// ============================================================================

func (c *KibanaClient) ListConnectors(ctx context.Context, space string) (kibana.Connectors, diag.Diagnostics) {
	resp, err := c.API.ListConnectors(ctx, space)
	if err != nil {
		return nil, c.reportFromErr(err)
	}
	if resp.Output == nil {
		return nil, c.reportUnknownError(resp)
	}
	return resp.Output, nil
}

func (c *KibanaClient) ReadConnector(ctx context.Context, space string, id string) (*kibana.Connector, diag.Diagnostics) {
	resp, err := c.API.ReadConnector(ctx, space, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) CreateConnector(ctx context.Context, space string, id string, req kibana.ConnectorRequest) (*kibana.Connector, diag.Diagnostics) {
	var resp *kibana.ApiResponse[kibana.Connector]
	var err error

	if id == "" {
		resp, err = c.API.CreateConnector(ctx, space, req)
	} else {
		resp, err = c.API.CreateConnectorWithID(ctx, space, id, req)
	}

	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) UpdateConnector(ctx context.Context, space string, id string, req kibana.ConnectorRequest) (*kibana.Connector, diag.Diagnostics) {
	resp, err := c.API.UpdateConnector(ctx, space, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) DeleteConnector(ctx context.Context, space string, id string) diag.Diagnostics {
	resp, err := c.API.DeleteConnector(ctx, space, id)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *KibanaClient) ListDataViews(ctx context.Context, space string) (*kibana.DataViews, diag.Diagnostics) {
	resp, err := c.API.ListDataViews(ctx, space)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) ReadDataView(ctx context.Context, space string, id string) (*kibana.DataView, diag.Diagnostics) {
	resp, err := c.API.ReadDataView(ctx, space, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) CreateDataView(ctx context.Context, space string, req kibana.DataView) (*kibana.DataView, diag.Diagnostics) {
	resp, err := c.API.CreateDataView(ctx, space, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) UpdateDataView(ctx context.Context, space string, id string, req kibana.DataView) (*kibana.DataView, diag.Diagnostics) {
	resp, err := c.API.UpdateDataView(ctx, space, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) DeleteDataView(ctx context.Context, space string, id string) diag.Diagnostics {
	resp, err := c.API.DeleteDataView(ctx, space, id)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *KibanaClient) ListSpaces(ctx context.Context) (kibana.Spaces, diag.Diagnostics) {
	resp, err := c.API.ListSpaces(ctx)
	if err != nil {
		return nil, c.reportFromErr(err)
	}
	if resp.Output == nil {
		return nil, c.reportUnknownError(resp)
	}
	return resp.Output, nil
}

func (c *KibanaClient) ReadSpace(ctx context.Context, id string) (*kibana.Space, diag.Diagnostics) {
	resp, err := c.API.ReadSpace(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) CreateSpace(ctx context.Context, req kibana.Space) (*kibana.Space, diag.Diagnostics) {
	resp, err := c.API.CreateSpace(ctx, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) UpdateSpace(ctx context.Context, id string, req kibana.Space) (*kibana.Space, diag.Diagnostics) {
	resp, err := c.API.UpdateSpace(ctx, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) DeleteSpace(ctx context.Context, id string) diag.Diagnostics {
	resp, err := c.API.DeleteSpace(ctx, id)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *KibanaClient) ReadRole(ctx context.Context, name string) (*kibana.Role, diag.Diagnostics) {
	resp, err := c.API.ReadRole(ctx, name)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) PutRole(ctx context.Context, name string, req kibana.Role, isCreate bool) (*kibana.Role, diag.Diagnostics) {
	var params *kibana.PutRoleParams
	if isCreate {
		params = &kibana.PutRoleParams{CreateOnly: true}
	} else {
		params = nil
	}

	resp, err := c.API.PutRole(ctx, name, req, params)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return c.ReadRole(ctx, name)
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) DeleteRole(ctx context.Context, name string) diag.Diagnostics {
	resp, err := c.API.DeleteRole(ctx, name)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *KibanaClient) ReadAlertingRule(ctx context.Context, space string, id string) (*kibana.AlertingRule, diag.Diagnostics) {
	resp, err := c.API.ReadAlertingRule(ctx, space, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) CreateAlertingRule(ctx context.Context, space string, id string, req kibana.AlertingRule) (*kibana.AlertingRule, diag.Diagnostics) {
	resp, err := c.API.CreateAlertingRule(ctx, space, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) UpdateAlertingRule(ctx context.Context, space string, id string, req kibana.AlertingRule) (*kibana.AlertingRule, diag.Diagnostics) {
	resp, err := c.API.UpdateAlertingRule(ctx, space, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) DeleteAlertingRule(ctx context.Context, space string, id string) diag.Diagnostics {
	resp, err := c.API.DeleteAlertingRule(ctx, space, id)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) EnableAlertingRule(ctx context.Context, space string, id string) diag.Diagnostics {
	resp, err := c.API.EnableAlertingRule(ctx, space, id)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

func (c *KibanaClient) DisableAlertingRule(ctx context.Context, space string, id string) diag.Diagnostics {
	resp, err := c.API.DisableAlertingRule(ctx, space, id)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *KibanaClient) ImportSavedObjects(ctx context.Context, space string, body []byte, params *kibana.ImportSavedObjectsParams) (*kibana.ImportSavedObjectsResponse, diag.Diagnostics) {
	resp, err := c.API.ImportSavedObjects(ctx, space, body, params)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}
