package kibana

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/kibana"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	p "go.openly.dev/pointy"
)

var (
	_ resource.Resource                = &alertingRuleResource{}
	_ resource.ResourceWithImportState = &alertingRuleResource{}
)

func NewAlertingRuleResource(client *clients.KibanaClient) *alertingRuleResource {
	return &alertingRuleResource{client: client}
}

type alertingRuleResource struct {
	client *clients.KibanaClient
}

func (r *alertingRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "kibana_alerting_rule")
}

func (r *alertingRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema.Description = "Creates a Kibana rule. See https://www.elastic.co/guide/en/kibana/master/create-rule-api.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of this resource.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"space_id": schema.StringAttribute{
			Description: "An identifier for the space. If space_id is not provided, the default space is used.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("default"),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the rule. While this name does not have to be unique, a distinctive name can help you identify a rule.",
			Required:    true,
		},
		"consumer": schema.StringAttribute{
			Description: "The name of the application or feature that owns the rule.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"params": schema.StringAttribute{
			Description: "The rule parameters, which differ for each rule type.",
			CustomType:  jsontypes.NormalizedType{},
			Required:    true,
		},
		"rule_type_id": schema.StringAttribute{
			Description: "The ID of the rule type that you want to call when the rule is scheduled to run. For more information about the valid values, list the rule types using [Get rule types API](https://www.elastic.co/guide/en/kibana/master/list-rule-types-api.html) or refer to the [Rule types documentation](https://www.elastic.co/guide/en/kibana/master/rule-types.html).",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"interval": schema.StringAttribute{
			Description: "The check interval, which specifies how frequently the rule conditions are checked. The interval must be specified in seconds, minutes, hours or days.",
			Required:    true,
			Validators: []validator.String{
				util.ElasticDurationValidator(),
			},
		},
		"actions": schema.ListNestedAttribute{
			Description: "An action that runs under defined conditions.",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"group": schema.StringAttribute{
						Description: "The group name, which affects when the action runs (for example, when the threshold is met or when the alert is recovered). Each rule type has a list of valid action group names.",
						Computed:    true,
						Optional:    true,
						Default:     stringdefault.StaticString("default"),
					},
					"frequency": schema.SingleNestedAttribute{
						Description: "The properties that affect how often actions are generated.",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"summary": schema.BoolAttribute{
								Description: "Indicates whether the action is a summary.",
								Computed:    true,
								Optional:    true,
								Default:     booldefault.StaticBool(false),
							},
							"notify_when": schema.StringAttribute{
								Description: "Defines how often alerts generate actions. Valid values include: `onActionGroupChange`: Actions run when the alert status changes; `onActiveAlert`: Actions run when the alert becomes active and at each check interval while the rule conditions are met; `onThrottleInterval`: Actions run when the alert becomes active and at the interval specified in the throttle property while the rule conditions are met.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("onActionGroupChange", "onActiveAlert", "onThrottleInterval"),
								},
							},
							"throttle": schema.StringAttribute{
								Description: "Defines how often an alert generates repeated actions. This custom action interval must be specified in seconds, minutes, hours, or days. For example, 10m or 1h. This property is applicable only if `notify_when` is `onThrottleInterval`.",
								Optional:    true,
								Validators: []validator.String{
									util.DurationValidator(),
								},
							},
						},
					},
					"id": schema.StringAttribute{
						Description: "The identifier for the connector saved object.",
						Required:    true,
					},
					"params": schema.StringAttribute{
						Description: "The parameters for the action, which are sent to the connector.",
						CustomType:  jsontypes.NormalizedType{},
						Required:    true,
					},
				},
			},
		},
		"enabled": schema.BoolAttribute{
			Description: "Indicates if you want to run the rule on an interval basis.",
			Optional:    true,
		},
		"tags": schema.ListAttribute{
			Description: "A list of tag names that are applied to the rule.",
			Computed:    true,
			Optional:    true,
			ElementType: types.StringType,
		},
		"scheduled_task_id": schema.StringAttribute{
			Description: "ID of the scheduled task that will execute the alert.",
			Computed:    true,
		},
		"last_execution_status": schema.StringAttribute{
			Description: "Status of the last execution of this rule.",
			Computed:    true,
		},
		"last_execution_date": schema.StringAttribute{
			Description: "Date of the last execution of this rule.",
			Computed:    true,
		},
	}
}

func (r *alertingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *alertingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data alertingRuleModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	space := data.SpaceID.ValueString()
	ruleID := data.ID.ValueString()
	rule, diags := r.client.ReadAlertingRule(ctx, space, ruleID)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(rule)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *alertingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data alertingRuleModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	space := data.SpaceID.ValueString()
	ruleID := data.ID.ValueString()
	rule, diags := r.client.CreateAlertingRule(ctx, space, ruleID, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(rule)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *alertingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data alertingRuleModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	body, diags := data.toApi(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body.Consumer = nil
	body.Enabled = nil
	body.RuleTypeID = nil
	space := data.SpaceID.ValueString()
	ruleID := data.ID.ValueString()
	rule, diags := r.client.UpdateAlertingRule(ctx, space, ruleID, *body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	shouldBeEnabled := data.Enabled.ValueBool()
	if shouldBeEnabled && !*rule.Enabled {
		diags = r.client.EnableAlertingRule(ctx, space, ruleID)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		rule.Enabled = p.Bool(true)
	}

	if !shouldBeEnabled && *rule.Enabled {
		diags = r.client.DisableAlertingRule(ctx, space, ruleID)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		rule.Enabled = p.Bool(false)
	}

	diags = data.fromApi(rule)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *alertingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data alertingRuleModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	space := data.SpaceID.ValueString()
	ruleID := data.ID.ValueString()
	diags = r.client.DeleteAlertingRule(ctx, space, ruleID)
	resp.Diagnostics.Append(diags...)
}

type alertingRuleModel struct {
	ID                  types.String              `tfsdk:"id"`
	SpaceID             types.String              `tfsdk:"space_id"`
	Name                types.String              `tfsdk:"name"`
	Consumer            types.String              `tfsdk:"consumer"`
	Params              jsontypes.Normalized      `tfsdk:"params"`
	RuleTypeID          types.String              `tfsdk:"rule_type_id"`
	Interval            types.String              `tfsdk:"interval"`
	Actions             []alertingRuleActionModel `tfsdk:"actions"`
	Enabled             types.Bool                `tfsdk:"enabled"`
	Tags                types.List                `tfsdk:"tags"`
	ScheduledTaskID     types.String              `tfsdk:"scheduled_task_id"`
	LastExecutionStatus types.String              `tfsdk:"last_execution_status"`
	LastExecutionDate   types.String              `tfsdk:"last_execution_date"`
}

type alertingRuleActionModel struct {
	Group     types.String                      `tfsdk:"group"`
	Frequency *alertingRuleActionFrequencyModel `tfsdk:"frequency"`
	ID        types.String                      `tfsdk:"id"`
	Params    jsontypes.Normalized              `tfsdk:"params"`
}

type alertingRuleActionFrequencyModel struct {
	Summary    types.Bool   `tfsdk:"summary"`
	NotifyWhen types.String `tfsdk:"notify_when"`
	Throttle   types.String `tfsdk:"throttle"`
}

func (m *alertingRuleModel) toApi(ctx context.Context) (*kibana.AlertingRule, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	req := &kibana.AlertingRule{
		Actions: util.TransformSlice(m.Actions, func(m alertingRuleActionModel, index int) kibana.AlertingRuleAction {
			path := path.AtName("actions").AtListIndex(index)
			return kibana.AlertingRuleAction{
				Frequency: util.TransformStruct(m.Frequency, func(m alertingRuleActionFrequencyModel) kibana.AlertingRuleActionFrequency {
					return kibana.AlertingRuleActionFrequency{
						Summary:    m.Summary.ValueBoolPointer(),
						NotifyWhen: m.NotifyWhen.ValueStringPointer(),
						Throttle:   m.Throttle.ValueStringPointer(),
					}
				}),
				Group:  m.Group.ValueStringPointer(),
				ID:     m.ID.ValueStringPointer(),
				Params: util.NormalizedTypeToMap[any](m.Params, path.AtName("params"), diags),
			}
		}),
		Consumer:   m.Consumer.ValueStringPointer(),
		Enabled:    m.Enabled.ValueBoolPointer(),
		Name:       m.Name.ValueStringPointer(),
		Params:     util.NormalizedTypeToMap[any](m.Params, path.AtName("params"), diags),
		RuleTypeID: m.RuleTypeID.ValueStringPointer(),
		Schedule: &kibana.AlertingRuleSchedule{
			Interval: m.Interval.ValueStringPointer(),
		},
		Tags: util.ListTypeToSliceBasic[string](ctx, m.Tags, path.AtName("tags"), diags),
	}

	return req, diags
}

func (m *alertingRuleModel) fromApi(resp *kibana.AlertingRule) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.ID = types.StringPointerValue(resp.ID)
	m.Name = types.StringPointerValue(resp.Name)
	m.Consumer = types.StringPointerValue(resp.Consumer)
	m.Params = util.MapToNormalizedType(resp.Params, path.AtName("params"), diags)
	m.RuleTypeID = types.StringPointerValue(resp.RuleTypeID)
	m.Interval = types.StringPointerValue(resp.Schedule.Interval)
	m.Actions = util.TransformSlice(resp.Actions, func(resp kibana.AlertingRuleAction, index int) alertingRuleActionModel {
		path := path.AtName("actions").AtListIndex(index)
		return alertingRuleActionModel{
			Group: types.StringPointerValue(resp.Group),
			Frequency: util.TransformStruct(resp.Frequency, func(resp kibana.AlertingRuleActionFrequency) alertingRuleActionFrequencyModel {
				return alertingRuleActionFrequencyModel{
					Summary:    types.BoolPointerValue(resp.Summary),
					NotifyWhen: types.StringPointerValue(resp.NotifyWhen),
					Throttle:   types.StringPointerValue(resp.Throttle),
				}
			}),
			ID:     types.StringPointerValue(resp.ID),
			Params: util.MapToNormalizedType(resp.Params, path.AtName("params"), diags),
		}
	})
	m.Enabled = types.BoolPointerValue(resp.Enabled)
	m.Tags = util.SliceToListType_String(resp.Tags, path.AtName("tags"), diags)
	m.ScheduledTaskID = types.StringPointerValue(resp.ScheduledTaskID)
	m.LastExecutionStatus = types.StringPointerValue(resp.ExecutionStatus.Status)
	m.LastExecutionDate = types.StringPointerValue(resp.ExecutionStatus.LastExecutionDate)

	return diags
}
