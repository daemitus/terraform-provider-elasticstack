package fleet

import (
	"context"
	"time"
)

func (c *Client) ReadAgentPolicy(ctx context.Context, agentPolicyId string) (*ApiResponse[ReadAgentPolicyResponse], error) {
	return doAPI[ReadAgentPolicyResponse](
		c, ctx, "GET", "/agent_policies/{id}",
		map[string]string{"id": agentPolicyId},
		nil, nil,
	)
}

type ReadAgentPolicyResponse struct {
	Item AgentPolicy `json:"item"`
}

// ============================================================================

func (c *Client) CreateAgentPolicy(ctx context.Context, body CreateAgentPolicyRequest) (*ApiResponse[CreateAgentPolicyResponse], error) {
	return doAPI[CreateAgentPolicyResponse](
		c, ctx,
		"POST", "/agent_policies",
		nil, body, nil,
	)
}

type CreateAgentPolicyRequest struct {
	AgentFeatures      []CreateAgentPolicyRequestAgentFeatures `json:"agent_features,omitempty"`
	DataOutputId       *string                                 `json:"data_output_id"`
	Description        *string                                 `json:"description,omitempty"`
	DownloadSourceId   *string                                 `json:"download_source_id"`
	FleetServerHostId  *string                                 `json:"fleet_server_host_id"`
	Id                 *string                                 `json:"id,omitempty"`
	InactivityTimeout  *float32                                `json:"inactivity_timeout,omitempty"`
	IsProtected        *bool                                   `json:"is_protected,omitempty"`
	MonitoringEnabled  []AgentPolicyMonitoringEnabled          `json:"monitoring_enabled,omitempty"`
	MonitoringOutputId *string                                 `json:"monitoring_output_id"`
	Name               string                                  `json:"name"`
	Namespace          string                                  `json:"namespace"`
	UnenrollTimeout    *float32                                `json:"unenroll_timeout,omitempty"`
}

type CreateAgentPolicyRequestAgentFeatures struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
}

type CreateAgentPolicyResponse struct {
	Item *AgentPolicy `json:"item,omitempty"`
}

// ============================================================================

func (c *Client) UpdateAgentPolicy(ctx context.Context, agentPolicyId string, body UpdateAgentPolicyRequest) (*ApiResponse[UpdateAgentPolicyResponse], error) {
	return doAPI[UpdateAgentPolicyResponse](
		c, ctx,
		"PUT", "agent_policies/{id}",
		map[string]string{"id": agentPolicyId},
		body, nil,
	)
}

type UpdateAgentPolicyRequest struct {
	AgentFeatures      []UpdateAgentPolicyRequest_AgentFeatures `json:"agent_features,omitempty"`
	DataOutputId       *string                                  `json:"data_output_id"`
	Description        *string                                  `json:"description,omitempty"`
	DownloadSourceId   *string                                  `json:"download_source_id"`
	FleetServerHostId  *string                                  `json:"fleet_server_host_id"`
	InactivityTimeout  *float32                                 `json:"inactivity_timeout,omitempty"`
	IsProtected        *bool                                    `json:"is_protected,omitempty"`
	MonitoringEnabled  []AgentPolicyMonitoringEnabled           `json:"monitoring_enabled,omitempty"`
	MonitoringOutputId *string                                  `json:"monitoring_output_id"`
	Name               string                                   `json:"name"`
	Namespace          string                                   `json:"namespace"`
	UnenrollTimeout    *float32                                 `json:"unenroll_timeout,omitempty"`
}

type UpdateAgentPolicyRequest_AgentFeatures struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
}

type UpdateAgentPolicyResponse struct {
	Item AgentPolicy `json:"item"`
}

// ============================================================================

func (c *Client) DeleteAgentPolicy(ctx context.Context, body DeleteAgentPolicyRequest) (*ApiResponse[DeleteAgentPolicyResponse], error) {
	return doAPI[DeleteAgentPolicyResponse](
		c, ctx,
		"POST", "/agent_policies/delete",
		nil, body, nil,
	)
}

type DeleteAgentPolicyRequest struct {
	AgentPolicyId string `json:"agentPolicyId"`
}

type DeleteAgentPolicyResponse struct {
	Id      string `json:"id"`
	Success bool   `json:"success"`
}

// ============================================================================

type AgentPolicy struct {
	AgentFeatures      []AgentPolicyFeature           `json:"agent_features,omitempty"`
	Agents             *float32                       `json:"agents,omitempty"`
	DataOutputId       *string                        `json:"data_output_id"`
	Description        *string                        `json:"description,omitempty"`
	DownloadSourceId   *string                        `json:"download_source_id"`
	FleetServerHostId  *string                        `json:"fleet_server_host_id"`
	Id                 string                         `json:"id"`
	InactivityTimeout  *float32                       `json:"inactivity_timeout,omitempty"`
	IsProtected        *bool                          `json:"is_protected,omitempty"`
	MonitoringEnabled  []AgentPolicyMonitoringEnabled `json:"monitoring_enabled,omitempty"`
	MonitoringOutputId *string                        `json:"monitoring_output_id"`
	Name               string                         `json:"name"`
	Namespace          string                         `json:"namespace"`
	Overrides          map[string]any                 `json:"overrides"`
	Revision           *float32                       `json:"revision,omitempty"`
	UnenrollTimeout    *float32                       `json:"unenroll_timeout,omitempty"`
	UpdatedBy          *string                        `json:"updated_by,omitempty"`
	UpdatedOn          *time.Time                     `json:"updated_on,omitempty"`
}

type AgentPolicyFeature struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
}

type AgentPolicyMonitoringEnabled string

const (
	AgentPolicyMonitoringEnabledLogs    AgentPolicyMonitoringEnabled = "logs"
	AgentPolicyMonitoringEnabledMetrics AgentPolicyMonitoringEnabled = "metrics"
)
