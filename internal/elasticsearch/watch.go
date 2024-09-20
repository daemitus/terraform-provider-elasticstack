package elasticsearch

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &WatchResource{}
	_ resource.ResourceWithImportState = &WatchResource{}
)

func NewWatchResource(client *clients.ElasticsearchClient) *WatchResource {
	return &WatchResource{client: client}
}

type WatchResource struct {
	client *clients.ElasticsearchClient
}

func (r *WatchResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "elasticsearch_watch")
}

func (r *WatchResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage Watches. See https://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api.html",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Internal identifier of the resource.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"active": schema.BoolAttribute{
				Description: "Defines whether the watch is active or inactive by default. The default value is true, which means the watch is active by default.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"trigger": schema.StringAttribute{
				Description: "The trigger that defines when the watch should run.",
				CustomType:  jsontypes.NormalizedType{},
				Required:    true,
			},
			"input": schema.StringAttribute{
				Description: "The input that defines the input that loads the data for the watch.",
				CustomType:  jsontypes.NormalizedType{},
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(`{"none":{}}`),
			},
			"condition": schema.StringAttribute{
				Description: "The condition that defines if the actions should be run.",
				CustomType:  jsontypes.NormalizedType{},
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(`{"always":{}}`),
			},
			"actions": schema.StringAttribute{
				Description: "The list of actions that will be run if the condition matches.",
				CustomType:  jsontypes.NormalizedType{},
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("{}"),
			},
			"metadata": schema.StringAttribute{
				Description: "Metadata json that will be copied into the history entries.",
				CustomType:  jsontypes.NormalizedType{},
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("{}"),
			},
			"transform": schema.StringAttribute{
				Description: "Processes the watch payload to prepare it for the watch actions.",
				CustomType:  jsontypes.NormalizedType{},
				Optional:    true,
			},
			"throttle_period_in_millis": schema.Int64Attribute{
				Description: "Minimum time in milliseconds between actions being run. Defaults to 5000.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(5000),
			},
		},
	}
}

func (r *WatchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WatchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data watchModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	id := data.ID.ValueString()
	watch, diags := r.client.GetWatch(ctx, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = data.fromApi(watch)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *WatchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data watchModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	putReq, diags := data.toApi()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	active := data.Active.ValueBool()
	watch, diags := r.client.PutWatch(ctx, id, active, putReq)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(watch)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *WatchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data watchModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	putReq, diags := data.toApi()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	active := data.Active.ValueBool()
	watch, diags := r.client.PutWatch(ctx, id, active, putReq)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(watch)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *WatchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data watchModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	id := data.ID.ValueString()
	diags = r.client.DeleteWatch(ctx, id)
	resp.Diagnostics.Append(diags...)
}

type watchModel struct {
	ID                     types.String         `tfsdk:"id"`
	Active                 types.Bool           `tfsdk:"active"`
	Trigger                jsontypes.Normalized `tfsdk:"trigger"`
	Input                  jsontypes.Normalized `tfsdk:"input"`
	Condition              jsontypes.Normalized `tfsdk:"condition"`
	Actions                jsontypes.Normalized `tfsdk:"actions"`
	Metadata               jsontypes.Normalized `tfsdk:"metadata"`
	Transform              jsontypes.Normalized `tfsdk:"transform"`
	ThrottlePeriodInMillis types.Int64          `tfsdk:"throttle_period_in_millis"`
}

func (m *watchModel) toApi() (elasticsearch.PutWatchRequest, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	output := elasticsearch.PutWatchRequest{
		Actions:        util.NormalizedTypeToMap[any](m.Actions, path.AtName("actions"), diags),
		Condition:      util.NormalizedTypeToMap[any](m.Condition, path.AtName("condition"), diags),
		Input:          util.NormalizedTypeToMap[any](m.Input, path.AtName("input"), diags),
		Metadata:       util.NormalizedTypeToMap[any](m.Metadata, path.AtName("metadata"), diags),
		Transform:      util.NormalizedTypeToMap[any](m.Transform, path.AtName("transform"), diags),
		Trigger:        util.NormalizedTypeToMap[any](m.Trigger, path.AtName("trigger"), diags),
		ThrottlePeriod: m.ThrottlePeriodInMillis.ValueInt64Pointer(),
	}

	return output, diags
}

func (m *watchModel) fromApi(item *elasticsearch.GetWatchResponse) diag.Diagnostics {
	if item == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.ID = types.StringValue(item.ID)
	m.Active = types.BoolValue(item.Status.State.Active)
	m.Trigger = util.MapToNormalizedType(item.Watch.Trigger, path.AtName("trigger"), diags)
	m.Input = util.MapToNormalizedType(item.Watch.Input, path.AtName("input"), diags)
	m.Condition = util.MapToNormalizedType(item.Watch.Condition, path.AtName("condition"), diags)
	m.Actions = util.MapToNormalizedType(item.Watch.Actions, path.AtName("actions"), diags)
	m.Metadata = util.MapToNormalizedType(item.Watch.Metadata, path.AtName("metadata"), diags)
	m.Transform = util.MapToNormalizedType(item.Watch.Transform, path.AtName("transform"), diags)
	m.ThrottlePeriodInMillis = types.Int64Value(*item.Watch.ThrottlePeriod)

	return diags
}
