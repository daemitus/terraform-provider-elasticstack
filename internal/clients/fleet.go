package clients

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/fleet"
	"github.com/daemitus/terraform-provider-elasticstack/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	p "go.openly.dev/pointy"
)

// FleetClient provides an API client for Fleet.
type FleetClient struct {
	baseClient
	API *fleet.Client
}

// NewFleetClient creates a new Fleet API client.
func NewFleetClient(config config.ServiceConfig) (*FleetClient, error) {
	api, err := fleet.NewClient(config)
	if err != nil {
		return nil, err
	}
	client := FleetClient{API: api}
	return &client, nil
}

// NewAccFleetClient creates a new Fleet API client for acceptance testing.
func NewAccFleetClient() *FleetClient {
	config, err := config.New().WithEnv()
	if err != nil {
		log.Fatal("Fleet config failure: %w", err)
	}

	client, err := NewFleetClient(config.Fleet)
	if err != nil {
		log.Fatal("Fleet client failure: %w", err)
	}

	return client
}

// ============================================================================

func (c *FleetClient) ListEnrollmentTokens(ctx context.Context) (fleet.EnrollmentApiKeys, diag.Diagnostics) {
	params := &fleet.ReadEnrollmentApiKeysParams{
		PerPage: p.Int(10000),
	}

	resp, err := c.API.ReadEnrollmentApiKeys(ctx, params)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Output.Items, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) ReadEnrollmentTokensByPolicy(ctx context.Context, policyID string) (fleet.EnrollmentApiKeys, diag.Diagnostics) {
	params := &fleet.ReadEnrollmentApiKeysParams{
		Kuery: p.String(fmt.Sprintf("policy_id:%s", policyID)),
	}

	resp, err := c.API.ReadEnrollmentApiKeys(ctx, params)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Output.Items, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *FleetClient) ReadAgentPolicy(ctx context.Context, id string) (*fleet.AgentPolicy, diag.Diagnostics) {
	resp, err := c.API.ReadAgentPolicy(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output.Item, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) CreateAgentPolicy(ctx context.Context, req fleet.CreateAgentPolicyRequest) (*fleet.AgentPolicy, diag.Diagnostics) {
	resp, err := c.API.CreateAgentPolicy(ctx, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Output.Item, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) UpdateAgentPolicy(ctx context.Context, id string, req fleet.UpdateAgentPolicyRequest) (*fleet.AgentPolicy, diag.Diagnostics) {
	resp, err := c.API.UpdateAgentPolicy(ctx, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output.Item, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) DeleteAgentPolicy(ctx context.Context, id string) diag.Diagnostics {
	body := fleet.DeleteAgentPolicyRequest{
		AgentPolicyId: id,
	}

	resp, err := c.API.DeleteAgentPolicy(ctx, body)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *FleetClient) ReadOutput(ctx context.Context, id string) (*fleet.OutputUnion, diag.Diagnostics) {
	resp, err := c.API.ReadOutput(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Output.Item, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) CreateOutput(ctx context.Context, req fleet.CreateOutputRequest) (*fleet.OutputUnion, diag.Diagnostics) {
	resp, err := c.API.CreateOutput(ctx, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Output.Item, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) UpdateOutput(ctx context.Context, id string, req fleet.UpdateOutputRequest) (*fleet.OutputUnion, diag.Diagnostics) {
	resp, err := c.API.UpdateOutput(ctx, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Output.Item, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) DeleteOutput(ctx context.Context, id string) diag.Diagnostics {
	resp, err := c.API.DeleteOutput(ctx, id)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *FleetClient) ReadFleetServerHost(ctx context.Context, id string) (*fleet.FleetServerHost, diag.Diagnostics) {
	resp, err := c.API.ReadFleetServerHosts(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output.Item, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) CreateFleetServerHost(ctx context.Context, req fleet.CreateFleetServerHostsRequest) (*fleet.FleetServerHost, diag.Diagnostics) {
	resp, err := c.API.CreateFleetServerHosts(ctx, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Output.Item, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) UpdateFleetServerHost(ctx context.Context, id string, req fleet.UpdateFleetServerHostsRequest) (*fleet.FleetServerHost, diag.Diagnostics) {
	resp, err := c.API.UpdateFleetServerHosts(ctx, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output.Item, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) DeleteFleetServerHost(ctx context.Context, id string) diag.Diagnostics {
	resp, err := c.API.DeleteFleetServerHosts(ctx, id)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *FleetClient) ReadPackagePolicy(ctx context.Context, id string) (*fleet.PackagePolicy, diag.Diagnostics) {
	resp, err := c.API.ReadPackagePolicy(ctx, id)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output.Item, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) CreatePackagePolicy(ctx context.Context, req fleet.PackagePolicyRequest) (*fleet.PackagePolicy, diag.Diagnostics) {
	resp, err := c.API.CreatePackagePolicy(ctx, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output.Item, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) UpdatePackagePolicy(ctx context.Context, id string, req fleet.PackagePolicyRequest) (*fleet.PackagePolicy, diag.Diagnostics) {
	resp, err := c.API.UpdatePackagePolicy(ctx, id, req)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return &resp.Output.Item, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) DeletePackagePolicy(ctx context.Context, id string, force bool) diag.Diagnostics {
	params := fleet.DeletePackagePolicyParams{Force: &force}
	resp, err := c.API.DeletePackagePolicy(ctx, id, &params)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================

func (c *FleetClient) ListPackages(ctx context.Context, prerelease bool) ([]fleet.SearchResult, diag.Diagnostics) {
	params := fleet.ListPackagesParams{
		Prerelease: &prerelease,
	}

	resp, err := c.API.ListPackages(ctx, &params)
	if err != nil {
		return nil, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return resp.Output.Items, nil
	default:
		return nil, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) ReadPackage(ctx context.Context, name, version string) (bool, diag.Diagnostics) {
	resp, err := c.API.ReadPackage(ctx, name, version, nil)
	if err != nil {
		return false, c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, c.reportFromErr(errors.New("package not found"))
	default:
		return false, c.reportUnknownError(resp)
	}
}

func (c *FleetClient) InstallPackage(ctx context.Context, name, version string, force bool) diag.Diagnostics {
	body := fleet.InstallPackageRequest{
		Force: &force,
	}

	resp, err := c.API.InstallPackage(ctx, name, version, body, nil)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

func (c *FleetClient) UninstallPackage(ctx context.Context, name, version string, force bool) diag.Diagnostics {
	body := fleet.DeletePackageRequest{
		Force: &force,
	}

	resp, err := c.API.DeletePackage(ctx, name, version, body, nil)
	if err != nil {
		return c.reportFromErr(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return c.reportUnknownError(resp)
	}
}

// ============================================================================
