package kibana

import (
	"context"
)

func (c *Client) ReadAlertingRule(ctx context.Context, space string, id string) (*ApiResponse[AlertingRule], error) {
	return doAPI[AlertingRule](
		c, ctx,
		"GET", "/s/{space}/api/alerting/rule/{id}",
		map[string]string{"space": space, "id": id},
		nil, nil,
	)
}

func (c *Client) CreateAlertingRule(ctx context.Context, space string, id string, body AlertingRule) (*ApiResponse[AlertingRule], error) {
	return doAPI[AlertingRule](
		c, ctx,
		"POST", "/s/{space}/api/alerting/rule/{id}",
		map[string]string{"space": space, "id": id},
		body, nil,
	)
}

func (c *Client) UpdateAlertingRule(ctx context.Context, space string, id string, body AlertingRule) (*ApiResponse[AlertingRule], error) {
	return doAPI[AlertingRule](
		c, ctx,
		"PUT", "/s/{space}/api/alerting/rule/{id}",
		map[string]string{"space": space, "id": id},
		body, nil,
	)
}

func (c *Client) DeleteAlertingRule(ctx context.Context, space string, id string) (*ApiResponse[EmptyResponse], error) {
	return doAPI[EmptyResponse](
		c, ctx,
		"DELETE", "/s/{space}/api/alerting/rule/{id}",
		map[string]string{"space": space, "id": id},
		nil, nil,
	)
}

func (c *Client) EnableAlertingRule(ctx context.Context, space string, id string) (*ApiResponse[EmptyResponse], error) {
	return doAPI[EmptyResponse](
		c, ctx,
		"POST", "/s/{space}/api/alerting/rule/{id}/_enable",
		map[string]string{"space": space, "id": id},
		nil, nil,
	)
}

func (c *Client) DisableAlertingRule(ctx context.Context, space string, id string) (*ApiResponse[EmptyResponse], error) {
	return doAPI[EmptyResponse](
		c, ctx,
		"POST", "/s/{space}/api/alerting/rule/{id}/_disable",
		map[string]string{"space": space, "id": id},
		nil, nil,
	)
}

// ============================================================================

type AlertingRule struct {
	Actions             []AlertingRuleAction         `json:"actions,omitempty"`
	ApiKeyCreatedByUser *bool                        `json:"api_key_created_by_user,omitempty"`
	ApiKeyOwner         *string                      `json:"api_key_owner,omitempty"`
	Consumer            *string                      `json:"consumer,omitempty"`
	CreatedAt           *string                      `json:"created_at,omitempty"`
	CreatedBy           *string                      `json:"created_by,omitempty"`
	Enabled             *bool                        `json:"enabled,omitempty"`
	ExecutionStatus     *AlertingRuleExecutionStatus `json:"execution_status,omitempty"`
	ID                  *string                      `json:"id,omitempty"`
	Name                *string                      `json:"name,omitempty"`
	NextRun             *string                      `json:"next_run,omitempty"`
	NotifyWhen          *string                      `json:"notify_when,omitempty"`
	Params              map[string]any               `json:"params,omitempty"`
	Revision            *int64                       `json:"revision,omitempty"`
	RuleTypeID          *string                      `json:"rule_type_id,omitempty"`
	Running             *bool                        `json:"running,omitempty"`
	Schedule            *AlertingRuleSchedule        `json:"schedule,omitempty"`
	ScheduledTaskID     *string                      `json:"scheduled_task_id,omitempty"`
	Tags                []string                     `json:"tags,omitempty"`
	UpdatedAt           *string                      `json:"updated_at,omitempty"`
	UpdatedBy           *string                      `json:"updated_by,omitempty"`
}

type AlertingRuleAction struct {
	ConnectorTypeId *string                      `json:"connector_type_id,omitempty"`
	Frequency       *AlertingRuleActionFrequency `json:"frequency,omitempty"`
	Group           *string                      `json:"group,omitempty"`
	ID              *string                      `json:"id,omitempty"`
	Params          map[string]any               `json:"params,omitempty"`
	UUID            *string                      `json:"uuid,omitempty"`
}

type AlertingRuleActionFrequency struct {
	Summary    *bool   `json:"summary,omitempty"`
	NotifyWhen *string `json:"notify_when,omitempty"`
	Throttle   *string `json:"throttle,omitempty"`
}

type AlertingRuleExecutionStatus struct {
	LastDuration      *int64  `json:"last_duration,omitempty"`
	LastExecutionDate *string `json:"last_execution_date,omitempty"`
	Status            *string `json:"status,omitempty"`
}

type AlertingRuleSchedule struct {
	Interval *string `json:"interval,omitempty"`
}
