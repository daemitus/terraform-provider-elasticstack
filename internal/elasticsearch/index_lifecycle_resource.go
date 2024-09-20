package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &IndexLifecycleResource{}
	_ resource.ResourceWithImportState = &IndexLifecycleResource{}
)

func NewIndexLifecycleResource(client *clients.ElasticsearchClient) *IndexLifecycleResource {
	return &IndexLifecycleResource{client: client}
}

type IndexLifecycleResource struct {
	client *clients.ElasticsearchClient
}

func (r *IndexLifecycleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "elasticsearch_index_lifecycle")
}

func (r *IndexLifecycleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	requiresAtleastOneState := objectvalidator.AtLeastOneOf(
		path.MatchRoot("hot"),
		path.MatchRoot("warm"),
		path.MatchRoot("cold"),
		path.MatchRoot("frozen"),
		path.MatchRoot("delete"),
	)

	var supportedActions = map[string]schema.Attribute{
		"allocate": schema.SingleNestedAttribute{
			Description: "Updates the index settings to change which nodes are allowed to host the index shards and change the number of replicas.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"number_of_replicas": schema.Int64Attribute{
					Description: "Number of replicas to assign to the index. Default: `0`",
					Computed:    true,
					Optional:    true,
					Default:     int64default.StaticInt64(0),
				},
				"total_shards_per_node": schema.Int64Attribute{
					Description: "The maximum number of shards for the index on a single Elasticsearch node. Defaults to `-1` (unlimited).",
					Computed:    true,
					Optional:    true,
					Default:     int64default.StaticInt64(-1),
				},
				"include": schema.StringAttribute{
					Description: "Assigns an index to nodes that have at least one of the specified custom attributes. Must be valid JSON document.",
					CustomType:  jsontypes.NormalizedType{},
					Computed:    true,
					Optional:    true,
					Default:     stringdefault.StaticString("{}"),
				},
				"exclude": schema.StringAttribute{
					Description: "Assigns an index to nodes that have none of the specified custom attributes. Must be valid JSON document.",
					CustomType:  jsontypes.NormalizedType{},
					Computed:    true,
					Optional:    true,
					Default:     stringdefault.StaticString("{}"),
				},
				"require": schema.StringAttribute{
					Description: "Assigns an index to nodes that have all of the specified custom attributes. Must be valid JSON document.",
					CustomType:  jsontypes.NormalizedType{},
					Computed:    true,
					Optional:    true,
					Default:     stringdefault.StaticString("{}"),
				},
			},
		},
		"delete": schema.SingleNestedAttribute{
			Description: "Permanently removes the index.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"delete_searchable_snapshot": schema.BoolAttribute{
					Description: "Deletes the searchable snapshot created in a previous phase.",
					Computed:    true,
					Optional:    true,
					Default:     booldefault.StaticBool(true),
				},
			},
		},
		"forcemerge": schema.SingleNestedAttribute{
			Description: "Force merges the index into the specified maximum number of segments. This action makes the index read-only.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"max_num_segments": schema.Int64Attribute{
					Description: "Number of segments to merge to. To fully merge the index, set to 1.",
					Required:    true,
					Validators: []validator.Int64{
						int64validator.AtLeast(1),
					},
				},
				"index_codec": schema.StringAttribute{
					Description: "Codec used to compress the document store.",
					Optional:    true,
				},
			},
		},
		"freeze": schema.SingleNestedAttribute{
			Description: "Freeze the index to minimize its memory footprint.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Description: "Controls whether ILM freezes the index.",
					Computed:    true,
					Optional:    true,
					Default:     booldefault.StaticBool(true),
				},
			},
		},
		"migrate": schema.SingleNestedAttribute{
			Description: `Moves the index to the data tier that corresponds to the current phase by updating the "index.routing.allocation.include._tier_preference" index setting.`,
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Description: "Controls whether ILM automatically migrates the index during this phase.",
					Computed:    true,
					Optional:    true,
					Default:     booldefault.StaticBool(true),
				},
			},
		},
		"readonly": schema.SingleNestedAttribute{
			Description: "Makes the index read-only.",
			Optional:    true,
			Attributes:  map[string]schema.Attribute{},
		},
		"rollover": schema.SingleNestedAttribute{
			Description: "Rolls over a target to a new index when the existing index meets one or more of the rollover conditions.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"max_age": schema.StringAttribute{
					Description: "Triggers rollover after the maximum elapsed time from index creation is reached.",
					Optional:    true,
				},
				"max_docs": schema.Int64Attribute{
					Description: "Triggers rollover after the specified maximum number of documents is reached.",
					Optional:    true,
				},
				"max_size": schema.StringAttribute{
					Description: "Triggers rollover when the index reaches a certain size.",
					Optional:    true,
				},
				"max_primary_shard_size": schema.StringAttribute{
					Description: "Triggers rollover when the largest primary shard in the index reaches a certain size.",
					Optional:    true,
				},
				"max_primary_shard_docs": schema.Int64Attribute{
					Description: "Triggers rollover when the largest primary shard in the index reaches a specified maximum number of documents.",
					Optional:    true,
				},
				"min_age": schema.StringAttribute{
					Description: "Prevents rollover until after the minimum elapsed time from index creation is reached. Supported from Elasticsearch version **8.4**",
					Optional:    true,
				},
				"min_docs": schema.Int64Attribute{
					Description: "Prevents rollover until after the specified minimum number of documents is reached. Supported from Elasticsearch version **8.4**",
					Optional:    true,
				},
				"min_size": schema.StringAttribute{
					Description: "Prevents rollover until the index reaches a certain size.",
					Optional:    true,
				},
				"min_primary_shard_size": schema.StringAttribute{
					Description: "Prevents rollover until the largest primary shard in the index reaches a certain size. Supported from Elasticsearch version **8.4**",
					Optional:    true,
				},
				"min_primary_shard_docs": schema.Int64Attribute{
					Description: "Prevents rollover until the largest primary shard in the index reaches a certain number of documents. Supported from Elasticsearch version **8.4**",
					Optional:    true,
				},
			},
		},
		"searchable_snapshot": schema.SingleNestedAttribute{
			Description: "Takes a snapshot of the managed index in the configured repository and mounts it as a searchable snapshot.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"snapshot_repository": schema.StringAttribute{
					Description: "Repository used to store the snapshot.",
					Required:    true,
				},
				"force_merge_index": schema.BoolAttribute{
					Description: "Force merges the managed index to one segment.",
					Computed:    true,
					Optional:    true,
					Default:     booldefault.StaticBool(true),
				},
			},
		},
		"set_priority": schema.SingleNestedAttribute{
			Description: "Sets the priority of the index as soon as the policy enters the hot, warm, or cold phase. Higher priority indices are recovered before indices with lower priorities following a node restart. Default priority is 1.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"priority": schema.Int64Attribute{
					Description: "The priority for the index. Must be 0 or greater.",
					Required:    true,
					Validators: []validator.Int64{
						int64validator.AtLeast(0),
					},
				},
			},
		},
		"shrink": schema.SingleNestedAttribute{
			Description: "Sets a source index to read-only and shrinks it into a new index with fewer primary shards.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"number_of_shards": schema.Int64Attribute{
					Description: "Number of shards to shrink to.",
					Optional:    true,
				},
				"max_primary_shard_size": schema.StringAttribute{
					Description: "The max primary shard size for the target index.",
					Optional:    true,
				},
				"allow_write_after_shrink": schema.BoolAttribute{
					Description: "If true, the shrunken index is made writable by removing the write block.",
					Optional:    true,
				},
			},
		},
		"unfollow": schema.SingleNestedAttribute{
			Description: "Convert a follower index to a regular index. Performed automatically before a rollover, shrink, or searchable snapshot action.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Description: "Controls whether ILM makes the follower index a regular one.",
					Computed:    true,
					Optional:    true,
					Default:     booldefault.StaticBool(true),
				},
			},
		},
		"wait_for_snapshot": schema.SingleNestedAttribute{
			Description: "Waits for the specified SLM policy to be executed before removing the index. This ensures that a snapshot of the deleted index is available.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"policy": schema.StringAttribute{
					Description: "Name of the SLM policy that the delete action should wait for.",
					Required:    true,
				},
			},
		},
		"downsample": schema.SingleNestedAttribute{
			Description: "Roll up documents within a fixed interval to a single summary document. Reduces the index footprint by storing time series data at reduced granularity.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"fixed_interval": schema.StringAttribute{
					Description: "Downsampling interval",
					Required:    true,
				},
				"wait_timeout": schema.StringAttribute{
					Description: "Downsampling interval",
					Computed:    true,
					Optional:    true,
					Default:     stringdefault.StaticString("1d"),
				},
			},
		},
		"min_age": schema.StringAttribute{
			Description: "ILM moves indices through the lifecycle according to their age. To control the timing of these transitions, you set a minimum age for each phase.",
			Computed:    true,
			Optional:    true,
		},
	}

	getSchema := func(actions ...string) map[string]schema.Attribute {
		actions = append(actions, "min_age") // Always supported
		m := make(map[string]schema.Attribute)
		for _, key := range actions {
			if action, ok := supportedActions[key]; ok {
				m[key] = action
			} else {
				panic(fmt.Sprintf("%s phase action not found", key))
			}
		}
		return m
	}

	resp.Schema = schema.Schema{
		Description: "Creates or updates a lifecycle policy. See https://www.elastic.co/guide/en/elasticsearch/reference/current/ilm-put-lifecycle.html and https://www.elastic.co/guide/en/elasticsearch/reference/current/ilm-index-lifecycle.html",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Internal identifier of the resource",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Identifier for the policy.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata": schema.StringAttribute{
				Description: "Optional user metadata about the ilm policy. Must be valid JSON document.",
				CustomType:  jsontypes.NormalizedType{},
				Optional:    true,
			},
			"hot": schema.SingleNestedAttribute{
				Description: "The index is actively being updated and queried.",
				Optional:    true,
				Validators: []validator.Object{
					requiresAtleastOneState,
				},
				Attributes: getSchema(
					"downsample",
					"forcemerge",
					"readonly",
					"rollover",
					"searchable_snapshot",
					"set_priority",
					"shrink",
					"unfollow",
				),
			},
			"warm": schema.SingleNestedAttribute{
				Description: "The index is no longer being updated but is still being queried.",
				Optional:    true,
				Validators: []validator.Object{
					requiresAtleastOneState,
				},
				Attributes: getSchema(
					"allocate",
					"downsample",
					"forcemerge",
					"migrate",
					"readonly",
					"set_priority",
					"shrink",
					"unfollow",
				),
			},
			"cold": schema.SingleNestedAttribute{
				Description: "The index is no longer being updated and is queried infrequently. The information still needs to be searchable, but it's okay if those queries are slower.",
				Optional:    true,
				Validators: []validator.Object{
					requiresAtleastOneState,
				},
				Attributes: getSchema(
					"allocate",
					"downsample",
					"freeze",
					"migrate",
					"readonly",
					"searchable_snapshot",
					"set_priority",
					"unfollow",
				),
			},
			"frozen": schema.SingleNestedAttribute{
				Description: "The index is no longer being updated and is queried rarely. The information still needs to be searchable, but it's okay if those queries are extremely slow.",
				Optional:    true,
				Validators: []validator.Object{
					requiresAtleastOneState,
				},
				Attributes: getSchema(
					"searchable_snapshot",
					"unfollow",
				),
			},
			"delete": schema.SingleNestedAttribute{
				Description: "The index is no longer needed and can safely be removed.",
				Optional:    true,
				Attributes: getSchema(
					"delete",
					"wait_for_snapshot",
				),
			},
			"modified_date": schema.StringAttribute{
				Description: "The DateTime of the last modification.",
				Computed:    true,
			},
		},
	}
}

func (r *IndexLifecycleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IndexLifecycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ilmPolicyModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	policy, diags := r.client.GetIlmPolicy(ctx, name)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = data.fromApi(name, policy)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IndexLifecycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ilmPolicyModel

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

	name := data.Name.ValueString()
	policy, diags := r.client.PutIlmPolicy(ctx, name, putReq)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(name, policy)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IndexLifecycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ilmPolicyModel

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

	name := data.Name.ValueString()
	policy, diags := r.client.PutIlmPolicy(ctx, name, putReq)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(name, policy)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IndexLifecycleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ilmPolicyModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	diags = r.client.DeleteIlmPolicy(ctx, name)
	resp.Diagnostics.Append(diags...)
}

type ilmPolicyModel struct {
	ID           types.String         `tfsdk:"id"`
	Name         types.String         `tfsdk:"name"`
	Metadata     jsontypes.Normalized `tfsdk:"metadata"`
	Hot          *ilmPolicyHot        `tfsdk:"hot"`
	Warm         *ilmPolicyWarm       `tfsdk:"warm"`
	Cold         *ilmPolicyCold       `tfsdk:"cold"`
	Frozen       *ilmPolicyFrozen     `tfsdk:"frozen"`
	Delete       *ilmPolicyDelete     `tfsdk:"delete"`
	ModifiedDate types.String         `tfsdk:"modified_date"`
}

type ilmPolicyHot struct {
	Downsample         *ilmPolicyPhaseDownsample         `tfsdk:"downsample"`
	ForceMerge         *ilmPolicyPhaseForceMerge         `tfsdk:"forcemerge"`
	ReadOnly           *ilmPolicyPhaseReadOnly           `tfsdk:"readonly"`
	Rollover           *ilmPolicyPhaseRollover           `tfsdk:"rollover"`
	SearchableSnapshot *ilmPolicyPhaseSearchableSnapshot `tfsdk:"searchable_snapshot"`
	SetPriority        *ilmPolicyPhaseSetPriority        `tfsdk:"set_priority"`
	Shrink             *ilmPolicyPhaseShrink             `tfsdk:"shrink"`
	Unfollow           *ilmPolicyPhaseUnfollow           `tfsdk:"unfollow"`
	MinAge             types.String                      `tfsdk:"min_age"`
}

type ilmPolicyWarm struct {
	Allocate    *ilmPolicyPhaseAllocate    `tfsdk:"allocate"`
	Downsample  *ilmPolicyPhaseDownsample  `tfsdk:"downsample"`
	ForceMerge  *ilmPolicyPhaseForceMerge  `tfsdk:"forcemerge"`
	Migrate     *ilmPolicyPhaseMigrate     `tfsdk:"migrate"`
	ReadOnly    *ilmPolicyPhaseReadOnly    `tfsdk:"readonly"`
	SetPriority *ilmPolicyPhaseSetPriority `tfsdk:"set_priority"`
	Shrink      *ilmPolicyPhaseShrink      `tfsdk:"shrink"`
	Unfollow    *ilmPolicyPhaseUnfollow    `tfsdk:"unfollow"`
	MinAge      types.String               `tfsdk:"min_age"`
}

type ilmPolicyCold struct {
	Allocate           *ilmPolicyPhaseAllocate           `tfsdk:"allocate"`
	Downsample         *ilmPolicyPhaseDownsample         `tfsdk:"downsample"`
	Freeze             *ilmPolicyPhaseFreeze             `tfsdk:"freeze"`
	Migrate            *ilmPolicyPhaseMigrate            `tfsdk:"migrate"`
	ReadOnly           *ilmPolicyPhaseReadOnly           `tfsdk:"readonly"`
	SearchableSnapshot *ilmPolicyPhaseSearchableSnapshot `tfsdk:"searchable_snapshot"`
	SetPriority        *ilmPolicyPhaseSetPriority        `tfsdk:"set_priority"`
	Unfollow           *ilmPolicyPhaseUnfollow           `tfsdk:"unfollow"`
	MinAge             types.String                      `tfsdk:"min_age"`
}

type ilmPolicyFrozen struct {
	SearchableSnapshot *ilmPolicyPhaseSearchableSnapshot `tfsdk:"searchable_snapshot"`
	Unfollow           *ilmPolicyPhaseUnfollow           `tfsdk:"unfollow"`
	MinAge             types.String                      `tfsdk:"min_age"`
}

type ilmPolicyDelete struct {
	Delete          *ilmPolicyPhaseDelete          `tfsdk:"delete"`
	WaitForSnapshot *ilmPolicyPhaseWaitForSnapshot `tfsdk:"wait_for_snapshot"`
	MinAge          types.String                   `tfsdk:"min_age"`
}

type ilmPolicyPhaseRollover struct {
	MaxAge              types.String `tfsdk:"max_age"`
	MaxDocs             types.Int64  `tfsdk:"max_docs"`
	MaxSize             types.String `tfsdk:"max_size"`
	MaxPrimaryShardSize types.String `tfsdk:"max_primary_shard_size"`
	MaxPrimaryShardDocs types.Int64  `tfsdk:"max_primary_shard_docs"`
	MinAge              types.String `tfsdk:"min_age"`
	MinDocs             types.Int64  `tfsdk:"min_docs"`
	MinSize             types.String `tfsdk:"min_size"`
	MinPrimaryShardSize types.String `tfsdk:"min_primary_shard_size"`
	MinPrimaryShardDocs types.Int64  `tfsdk:"min_primary_shard_docs"`
}

type ilmPolicyPhaseForceMerge struct {
	MaxNumSegments types.Int64  `tfsdk:"max_num_segments"`
	IndexCodec     types.String `tfsdk:"index_codec"`
}

type ilmPolicyPhaseShrink struct {
	NumberOfShards        types.Int64  `tfsdk:"number_of_shards"`
	MaxPrimaryShardSize   types.String `tfsdk:"max_primary_shard_size"`
	AllowWriteAfterShrink types.Bool   `tfsdk:"allow_write_after_shrink"`
}

type ilmPolicyPhaseAllocate struct {
	NumberOfReplicas   types.Int64          `tfsdk:"number_of_replicas"`
	TotalShardsPerNode types.Int64          `tfsdk:"total_shards_per_node"`
	Include            jsontypes.Normalized `tfsdk:"include"`
	Exclude            jsontypes.Normalized `tfsdk:"exclude"`
	Require            jsontypes.Normalized `tfsdk:"require"`
}

type ilmPolicyPhaseDelete struct {
	DeleteSearchableSnapshot types.Bool `tfsdk:"delete_searchable_snapshot"`
}

type ilmPolicyPhaseFreeze struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type ilmPolicyPhaseMigrate struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type ilmPolicyPhaseReadOnly struct {
}

type ilmPolicyPhaseSearchableSnapshot struct {
	SnapshotRepository types.String `tfsdk:"snapshot_repository"`
	ForceMergeIndex    types.Bool   `tfsdk:"force_merge_index"`
}

type ilmPolicyPhaseSetPriority struct {
	Priority types.Int64 `tfsdk:"priority"`
}

type ilmPolicyPhaseWaitForSnapshot struct {
	Policy types.String `tfsdk:"policy"`
}

type ilmPolicyPhaseDownsample struct {
	FixedInterval types.String `tfsdk:"fixed_interval"`
	WaitTimeout   types.String `tfsdk:"wait_timeout"`
}

type ilmPolicyPhaseUnfollow struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

func (m *ilmPolicyPhaseAllocate) toApi(phase string, diags diag.Diagnostics) *phaseActionAllocate {
	if m == nil {
		return nil
	}

	path := path.Root(phase)
	return &phaseActionAllocate{
		NumberOfReplicas:   m.NumberOfReplicas.ValueInt64Pointer(),
		TotalShardsPerNode: m.TotalShardsPerNode.ValueInt64Pointer(),
		Include:            util.NormalizedTypeToMap[any](m.Include, path.AtName("include"), diags),
		Exclude:            util.NormalizedTypeToMap[any](m.Exclude, path.AtName("exclude"), diags),
		Require:            util.NormalizedTypeToMap[any](m.Require, path.AtName("require"), diags),
	}
}

func (m *ilmPolicyPhaseDelete) toApi() *phaseActionDelete {
	if m == nil {
		return nil
	}

	return &phaseActionDelete{
		DeleteSearchableSnapshot: m.DeleteSearchableSnapshot.ValueBoolPointer(),
	}
}

func (m *ilmPolicyPhaseFreeze) toApi() *phaseActionFreeze {
	if m == nil {
		return nil
	}

	return &phaseActionFreeze{
		Enabled: m.Enabled.ValueBoolPointer(),
	}
}

func (m *ilmPolicyPhaseMigrate) toApi() *phaseActionMigrate {
	if m == nil {
		return nil
	}

	return &phaseActionMigrate{
		Enabled: m.Enabled.ValueBoolPointer(),
	}
}

func (m *ilmPolicyPhaseReadOnly) toApi() *phaseActionReadonly {
	if m == nil {
		return nil
	}

	return &phaseActionReadonly{}
}

func (m *ilmPolicyPhaseSearchableSnapshot) toApi() *phaseActionSearchableSnapshot {
	if m == nil {
		return nil
	}

	return &phaseActionSearchableSnapshot{
		SnapshotRepository: m.SnapshotRepository.ValueStringPointer(),
		ForceMergeIndex:    m.ForceMergeIndex.ValueBoolPointer(),
	}
}

func (m *ilmPolicyPhaseSetPriority) toApi() *phaseActionSetPriority {
	if m == nil {
		return nil
	}

	return &phaseActionSetPriority{
		Priority: m.Priority.ValueInt64Pointer(),
	}
}

func (m *ilmPolicyPhaseWaitForSnapshot) toApi() *phaseActionWaitForSnapshot {
	if m == nil {
		return nil
	}

	return &phaseActionWaitForSnapshot{
		Policy: m.Policy.ValueStringPointer(),
	}
}

func (m *ilmPolicyPhaseDownsample) toApi() *phaseActionDownsample {
	if m == nil {
		return nil
	}

	return &phaseActionDownsample{
		FixedInterval: m.FixedInterval.ValueStringPointer(),
		WaitTimeout:   m.WaitTimeout.ValueStringPointer(),
	}
}

func (m *ilmPolicyPhaseUnfollow) toApi() *phaseActionUnfollow {
	if m == nil {
		return nil
	}

	return &phaseActionUnfollow{
		Enabled: m.Enabled.ValueBoolPointer(),
	}
}

func (m *ilmPolicyPhaseForceMerge) toApi() *phaseActionForceMerge {
	if m == nil {
		return nil
	}

	return &phaseActionForceMerge{
		MaxNumSegments: m.MaxNumSegments.ValueInt64Pointer(),
		IndexCodec:     m.IndexCodec.ValueStringPointer(),
	}
}

func (m *ilmPolicyPhaseRollover) toApi() *phaseActionRollover {
	if m == nil {
		return nil
	}

	return &phaseActionRollover{
		MaxAge:              m.MaxAge.ValueStringPointer(),
		MaxDocs:             m.MaxDocs.ValueInt64Pointer(),
		MaxSize:             m.MaxSize.ValueStringPointer(),
		MaxPrimaryShardSize: m.MaxPrimaryShardSize.ValueStringPointer(),
		MaxPrimaryShardDocs: m.MaxPrimaryShardDocs.ValueInt64Pointer(),
		MinAge:              m.MinAge.ValueStringPointer(),
		MinDocs:             m.MinDocs.ValueInt64Pointer(),
		MinSize:             m.MinSize.ValueStringPointer(),
		MinPrimaryShardSize: m.MinPrimaryShardSize.ValueStringPointer(),
		MinPrimaryShardDocs: m.MinPrimaryShardDocs.ValueInt64Pointer(),
	}
}

func (m *ilmPolicyPhaseShrink) toApi() *phaseActionShrink {
	if m == nil {
		return nil
	}

	return &phaseActionShrink{
		NumberOfShards:        m.NumberOfShards.ValueInt64Pointer(),
		MaxPrimaryShardSize:   m.MaxPrimaryShardSize.ValueStringPointer(),
		AllowWriteAfterShrink: m.AllowWriteAfterShrink.ValueBoolPointer(),
	}
}

type phaseActions struct {
	Allocate           *phaseActionAllocate           `json:"allocate,omitempty"`
	Delete             *phaseActionDelete             `json:"delete,omitempty"`
	Downsample         *phaseActionDownsample         `json:"downsample,omitempty"`
	ForceMerge         *phaseActionForceMerge         `json:"forcemerge,omitempty"`
	Freeze             *phaseActionFreeze             `json:"freeze,omitempty"`
	Migrate            *phaseActionMigrate            `json:"migrate,omitempty"`
	Readonly           *phaseActionReadonly           `json:"readonly,omitempty"`
	Rollover           *phaseActionRollover           `json:"rollover,omitempty"`
	SearchableSnapshot *phaseActionSearchableSnapshot `json:"searchable_snapshot,omitempty"`
	SetPriority        *phaseActionSetPriority        `json:"set_priority,omitempty"`
	Shrink             *phaseActionShrink             `json:"shrink,omitempty"`
	WaitForSnapshot    *phaseActionWaitForSnapshot    `json:"wait_for_snapshot,omitempty"`
	Unfollow           *phaseActionUnfollow           `json:"unfollow,omitempty"`
}

type phaseActionAllocate struct {
	NumberOfReplicas   *int64         `json:"number_of_replicas,omitempty"`
	TotalShardsPerNode *int64         `json:"total_shards_per_node,omitempty"`
	Include            map[string]any `json:"include,omitempty"`
	Exclude            map[string]any `json:"exclude,omitempty"`
	Require            map[string]any `json:"require,omitempty"`
}

type phaseActionDelete struct {
	DeleteSearchableSnapshot *bool `json:"delete_searchable_snapshot,omitempty"`
}

type phaseActionDownsample struct {
	FixedInterval *string `json:"fixed_interval,omitempty"`
	WaitTimeout   *string `json:"wait_timeout,omitempty"`
}

type phaseActionForceMerge struct {
	MaxNumSegments *int64  `json:"max_num_segments,omitempty"`
	IndexCodec     *string `json:"index_codec,omitempty"`
}

type phaseActionFreeze struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type phaseActionMigrate struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type phaseActionReadonly struct {
}

type phaseActionRollover struct {
	MaxAge              *string `json:"max_age,omitempty"`
	MaxDocs             *int64  `json:"max_docs,omitempty"`
	MaxSize             *string `json:"max_size,omitempty"`
	MaxPrimaryShardSize *string `json:"max_primary_shard_size,omitempty"`
	MaxPrimaryShardDocs *int64  `json:"max_primary_shard_docs,omitempty"`
	MinAge              *string `json:"min_age,omitempty"`
	MinDocs             *int64  `json:"min_docs,omitempty"`
	MinSize             *string `json:"min_size,omitempty"`
	MinPrimaryShardSize *string `json:"min_primary_shard_size,omitempty"`
	MinPrimaryShardDocs *int64  `json:"min_primary_shard_docs,omitempty"`
}

type phaseActionSearchableSnapshot struct {
	SnapshotRepository *string `json:"snapshot_repository,omitempty"`
	ForceMergeIndex    *bool   `json:"force_merge_index,omitempty"`
}

type phaseActionSetPriority struct {
	Priority *int64 `json:"priority,omitempty"`
}

type phaseActionShrink struct {
	NumberOfShards        *int64  `json:"number_of_shards,omitempty"`
	MaxPrimaryShardSize   *string `json:"max_primary_shard_size,omitempty"`
	AllowWriteAfterShrink *bool   `json:"allow_write_after_shrink,omitempty"`
}

type phaseActionWaitForSnapshot struct {
	Policy *string `json:"policy,omitempty"`
}

type phaseActionUnfollow struct {
	Enabled *bool `json:"enabled,omitempty"`
}

func (m *ilmPolicyModel) toApi() (elasticsearch.PutIlmPolicyRequest, diag.Diagnostics) {
	path := path.Empty()
	var diags diag.Diagnostics

	asDuration := func(val types.String) *estypes.Duration {
		var out *estypes.Duration
		p := val.ValueStringPointer()
		if p != nil && *p != "" {
			v := estypes.Duration(*p)
			out = &v
		}
		return out
	}

	marshalActions := func(name string, actions *phaseActions) json.RawMessage {
		bytes, err := util.JsonMarshal(actions)
		if err != nil {
			diags.AddAttributeError(path.AtName(name), "marshal failure", err.Error())
		}
		return bytes
	}

	output := elasticsearch.PutIlmPolicyRequest{
		Policy: &estypes.IlmPolicy{
			Meta_: util.NormalizedTypeToMap[json.RawMessage](m.Metadata, path.AtName("metadata"), diags),
			Phases: estypes.Phases{
				Hot: util.TransformStruct(m.Hot, func(m ilmPolicyHot) estypes.Phase {
					return estypes.Phase{
						Actions: marshalActions("hot", &phaseActions{
							Downsample:         m.Downsample.toApi(),
							ForceMerge:         m.ForceMerge.toApi(),
							Readonly:           m.ReadOnly.toApi(),
							Rollover:           m.Rollover.toApi(),
							SearchableSnapshot: m.SearchableSnapshot.toApi(),
							SetPriority:        m.SetPriority.toApi(),
							Shrink:             m.Shrink.toApi(),
							Unfollow:           m.Unfollow.toApi(),
						}),
						MinAge: asDuration(m.MinAge),
					}
				}),
				Warm: util.TransformStruct(m.Warm, func(m ilmPolicyWarm) estypes.Phase {
					return estypes.Phase{
						Actions: marshalActions("warm", &phaseActions{
							Allocate:    m.Allocate.toApi("warm", diags),
							Downsample:  m.Downsample.toApi(),
							ForceMerge:  m.ForceMerge.toApi(),
							Migrate:     m.Migrate.toApi(),
							Readonly:    m.ReadOnly.toApi(),
							SetPriority: m.SetPriority.toApi(),
							Unfollow:    m.Unfollow.toApi(),
							Shrink:      m.Shrink.toApi(),
						}),
						MinAge: asDuration(m.MinAge),
					}
				}),
				Cold: util.TransformStruct(m.Cold, func(m ilmPolicyCold) estypes.Phase {
					return estypes.Phase{
						Actions: marshalActions("cold", &phaseActions{
							Allocate:           m.Allocate.toApi("cold", diags),
							Downsample:         m.Downsample.toApi(),
							Freeze:             m.Freeze.toApi(),
							Migrate:            m.Migrate.toApi(),
							Readonly:           m.ReadOnly.toApi(),
							SearchableSnapshot: m.SearchableSnapshot.toApi(),
							SetPriority:        m.SetPriority.toApi(),
							Unfollow:           m.Unfollow.toApi(),
						}),
						MinAge: asDuration(m.MinAge),
					}
				}),
				Frozen: util.TransformStruct(m.Frozen, func(m ilmPolicyFrozen) estypes.Phase {
					return estypes.Phase{
						Actions: marshalActions("frozen", &phaseActions{
							SearchableSnapshot: m.SearchableSnapshot.toApi(),
							Unfollow:           m.Unfollow.toApi(),
						}),
						MinAge: asDuration(m.MinAge),
					}
				}),
				Delete: util.TransformStruct(m.Delete, func(m ilmPolicyDelete) estypes.Phase {
					return estypes.Phase{
						Actions: marshalActions("delete", &phaseActions{
							Delete:          m.Delete.toApi(),
							WaitForSnapshot: m.WaitForSnapshot.toApi(),
						}),
						MinAge: asDuration(m.MinAge),
					}
				}),
			},
		},
	}

	return output, diags
}

func (m *ilmPolicyModel) fromApi(name string, resp *estypes.Lifecycle) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	var diags diag.Diagnostics

	minAge := func(val *estypes.Duration) types.String {
		if val == nil {
			return types.StringNull()
		}
		v := *val
		return types.StringValue(v.(string))
	}

	allocateValue := func(resp phaseActions, path path.Path) (out *ilmPolicyPhaseAllocate) {
		if resp.Allocate != nil {
			out = &ilmPolicyPhaseAllocate{
				NumberOfReplicas:   types.Int64PointerValue(resp.Allocate.NumberOfReplicas),
				TotalShardsPerNode: types.Int64PointerValue(resp.Allocate.TotalShardsPerNode),
				Include:            util.MapToNormalizedType(resp.Allocate.Include, path.AtName("include"), diags),
				Exclude:            util.MapToNormalizedType(resp.Allocate.Exclude, path.AtName("exlude"), diags),
				Require:            util.MapToNormalizedType(resp.Allocate.Require, path.AtName("require"), diags),
			}
		}
		return
	}

	deleteValue := func(resp phaseActions) (out *ilmPolicyPhaseDelete) {
		if resp.Delete != nil {
			out = &ilmPolicyPhaseDelete{
				DeleteSearchableSnapshot: types.BoolPointerValue(resp.Delete.DeleteSearchableSnapshot),
			}
		}
		return
	}

	downsampleValue := func(resp phaseActions) (out *ilmPolicyPhaseDownsample) {
		if resp.Downsample != nil {
			out = &ilmPolicyPhaseDownsample{
				FixedInterval: types.StringPointerValue(resp.Downsample.FixedInterval),
				WaitTimeout:   types.StringPointerValue(resp.Downsample.WaitTimeout),
			}
		}
		return
	}

	freezeValue := func(resp phaseActions) (out *ilmPolicyPhaseFreeze) {
		if resp.Freeze != nil {
			out = &ilmPolicyPhaseFreeze{
				Enabled: types.BoolPointerValue(resp.Freeze.Enabled),
			}
		}
		return
	}

	migrateValue := func(resp phaseActions) (out *ilmPolicyPhaseMigrate) {
		if resp.Migrate != nil {
			out = &ilmPolicyPhaseMigrate{
				Enabled: types.BoolPointerValue(resp.Migrate.Enabled),
			}
		}
		return
	}

	readonlyValue := func(resp phaseActions) (out *ilmPolicyPhaseReadOnly) {
		if resp.Readonly != nil {
			out = &ilmPolicyPhaseReadOnly{}
		}
		return
	}

	searchableSnapshotValue := func(resp phaseActions) (out *ilmPolicyPhaseSearchableSnapshot) {
		if resp.SearchableSnapshot != nil {
			out = &ilmPolicyPhaseSearchableSnapshot{
				SnapshotRepository: types.StringPointerValue(resp.SearchableSnapshot.SnapshotRepository),
				ForceMergeIndex:    types.BoolPointerValue(resp.SearchableSnapshot.ForceMergeIndex),
			}
		}
		return
	}

	setPriorityValue := func(resp phaseActions) (out *ilmPolicyPhaseSetPriority) {
		if resp.SetPriority != nil {
			out = &ilmPolicyPhaseSetPriority{
				Priority: types.Int64PointerValue(resp.SetPriority.Priority),
			}
		}
		return
	}

	unfollowValue := func(resp phaseActions) (out *ilmPolicyPhaseUnfollow) {
		if resp.Unfollow != nil {
			out = &ilmPolicyPhaseUnfollow{
				Enabled: types.BoolPointerValue(resp.Unfollow.Enabled),
			}
		}
		return
	}

	waitForSnapshotValue := func(resp phaseActions) (out *ilmPolicyPhaseWaitForSnapshot) {
		if resp.WaitForSnapshot != nil {
			out = &ilmPolicyPhaseWaitForSnapshot{
				Policy: types.StringPointerValue(resp.WaitForSnapshot.Policy),
			}
		}
		return
	}

	forceMergeValue := func(resp phaseActions) (out *ilmPolicyPhaseForceMerge) {
		if resp.ForceMerge != nil {
			out = &ilmPolicyPhaseForceMerge{
				MaxNumSegments: types.Int64PointerValue(resp.ForceMerge.MaxNumSegments),
				IndexCodec:     types.StringPointerValue(resp.ForceMerge.IndexCodec),
			}
		}
		return
	}

	rolloverValue := func(resp phaseActions) (out *ilmPolicyPhaseRollover) {
		if resp.Rollover != nil {
			out = &ilmPolicyPhaseRollover{
				MaxAge:              types.StringPointerValue(resp.Rollover.MaxAge),
				MaxDocs:             types.Int64PointerValue(resp.Rollover.MaxDocs),
				MaxSize:             types.StringPointerValue(resp.Rollover.MaxSize),
				MaxPrimaryShardSize: types.StringPointerValue(resp.Rollover.MaxPrimaryShardSize),
				MaxPrimaryShardDocs: types.Int64PointerValue(resp.Rollover.MaxPrimaryShardDocs),
				MinAge:              types.StringPointerValue(resp.Rollover.MinAge),
				MinDocs:             types.Int64PointerValue(resp.Rollover.MinDocs),
				MinSize:             types.StringPointerValue(resp.Rollover.MinSize),
				MinPrimaryShardSize: types.StringPointerValue(resp.Rollover.MinPrimaryShardSize),
				MinPrimaryShardDocs: types.Int64PointerValue(resp.Rollover.MinPrimaryShardDocs),
			}
		}
		return
	}

	shrinkValue := func(resp phaseActions) (out *ilmPolicyPhaseShrink) {
		if resp.Shrink != nil {
			out = &ilmPolicyPhaseShrink{
				NumberOfShards:        types.Int64PointerValue(resp.Shrink.NumberOfShards),
				MaxPrimaryShardSize:   types.StringPointerValue(resp.Shrink.MaxPrimaryShardSize),
				AllowWriteAfterShrink: types.BoolPointerValue(resp.Shrink.AllowWriteAfterShrink),
			}
		}
		return
	}

	path := path.Empty()

	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)
	m.Metadata = util.MapToNormalizedType(resp.Policy.Meta_, path.AtName("metadata"), diags)
	m.Hot = util.TransformStruct(resp.Policy.Phases.Hot, func(resp estypes.Phase) ilmPolicyHot {
		actions, err := util.JsonUnmarshal[phaseActions](resp.Actions)
		if err != nil {
			diags.AddAttributeError(path.AtName("hot"), "unmarshal failure", err.Error())
			return ilmPolicyHot{}
		}

		return ilmPolicyHot{
			Downsample:         downsampleValue(actions),
			ReadOnly:           readonlyValue(actions),
			SearchableSnapshot: searchableSnapshotValue(actions),
			SetPriority:        setPriorityValue(actions),
			Unfollow:           unfollowValue(actions),
			ForceMerge:         forceMergeValue(actions),
			Rollover:           rolloverValue(actions),
			Shrink:             shrinkValue(actions),
			MinAge:             minAge(resp.MinAge),
		}
	})
	m.Warm = util.TransformStruct(resp.Policy.Phases.Warm, func(resp estypes.Phase) ilmPolicyWarm {
		path := path.AtName("warm")
		actions, err := util.JsonUnmarshal[phaseActions](resp.Actions)
		if err != nil {
			diags.AddAttributeError(path, "unmarshal failure", err.Error())
			return ilmPolicyWarm{}
		}

		return ilmPolicyWarm{
			Allocate:    allocateValue(actions, path),
			Downsample:  downsampleValue(actions),
			Migrate:     migrateValue(actions),
			ReadOnly:    readonlyValue(actions),
			SetPriority: setPriorityValue(actions),
			Unfollow:    unfollowValue(actions),
			ForceMerge:  forceMergeValue(actions),
			Shrink:      shrinkValue(actions),
			MinAge:      minAge(resp.MinAge),
		}
	})
	m.Cold = util.TransformStruct(resp.Policy.Phases.Cold, func(resp estypes.Phase) ilmPolicyCold {
		path := path.AtName("cold")
		actions, err := util.JsonUnmarshal[phaseActions](resp.Actions)
		if err != nil {
			diags.AddAttributeError(path, "unmarshal failure", err.Error())
			return ilmPolicyCold{}
		}

		return ilmPolicyCold{
			Allocate:           allocateValue(actions, path),
			Downsample:         downsampleValue(actions),
			Freeze:             freezeValue(actions),
			Migrate:            migrateValue(actions),
			ReadOnly:           readonlyValue(actions),
			SearchableSnapshot: searchableSnapshotValue(actions),
			SetPriority:        setPriorityValue(actions),
			Unfollow:           unfollowValue(actions),
			MinAge:             minAge(resp.MinAge),
		}
	})
	m.Frozen = util.TransformStruct(resp.Policy.Phases.Frozen, func(resp estypes.Phase) ilmPolicyFrozen {
		actions, err := util.JsonUnmarshal[phaseActions](resp.Actions)
		if err != nil {
			diags.AddAttributeError(path.AtName("frozen"), "unmarshal failure", err.Error())
			return ilmPolicyFrozen{}
		}

		return ilmPolicyFrozen{
			SearchableSnapshot: searchableSnapshotValue(actions),
			Unfollow:           unfollowValue(actions),
			MinAge:             minAge(resp.MinAge),
		}
	})
	m.Delete = util.TransformStruct(resp.Policy.Phases.Delete, func(resp estypes.Phase) ilmPolicyDelete {
		actions, err := util.JsonUnmarshal[phaseActions](resp.Actions)
		if err != nil {
			diags.AddAttributeError(path.AtName("delete"), "unmarshal failure", err.Error())
			return ilmPolicyDelete{}
		}

		return ilmPolicyDelete{
			Delete:          deleteValue(actions),
			WaitForSnapshot: waitForSnapshotValue(actions),
			MinAge:          minAge(resp.MinAge),
		}
	})
	m.ModifiedDate = types.StringValue(resp.ModifiedDate.(string))

	return diags
}
