package elasticsearch

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &IndexResource{}
	_ resource.ResourceWithImportState = &IndexResource{}
)

func NewIndexResource(client *clients.ElasticsearchClient) *IndexResource {
	return &IndexResource{client: client}
}

type IndexResource struct {
	client *clients.ElasticsearchClient
}

func (r *IndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "elasticsearch_index")
}

func (r *IndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates Elasticsearch indices. See https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Internal identifier of the resource",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the index you wish to create.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					stringvalidator.NoneOf(".", ".."),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^-_+]`), "cannot start with -, _, +"),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9!$%&'()+.;=@[\]^{}~_-]+$`), "must contain lower case alphanumeric characters and selected punctuation, see: https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html#indices-create-api-path-params"),
				},
			},
			"aliases": schema.MapNestedAttribute{
				Description: "Aliases for the index.",
				Computed:    true,
				Optional:    true,
				Default: mapdefault.StaticValue(
					types.MapValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"filter":         jsontypes.NormalizedType{},
								"index_routing":  types.StringType,
								"is_hidden":      types.BoolType,
								"is_write_index": types.BoolType,
								"routing":        types.StringType,
								"search_routing": types.StringType,
							},
						},
						map[string]attr.Value{},
					)),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"filter": schema.StringAttribute{
							Description: "Query used to limit documents the alias can access.",
							CustomType:  jsontypes.NormalizedType{},
							Optional:    true,
						},
						"index_routing": schema.StringAttribute{
							Description: "Value used to route indexing operations to a specific shard. If specified, this overwrites the `routing` value for indexing operations.",
							Optional:    true,
						},
						"is_hidden": schema.BoolAttribute{
							Description: "If true, the alias is hidden.",
							Optional:    true,
						},
						"is_write_index": schema.BoolAttribute{
							Description: "If true, the index is the write index for the alias.",
							Optional:    true,
						},
						"routing": schema.StringAttribute{
							Description: "Value used to route indexing and search operations to a specific shard.",
							Optional:    true,
						},
						"search_routing": schema.StringAttribute{
							Description: "Value used to route search operations to a specific shard. If specified, this overwrites the routing value for search operations.",
							Optional:    true,
						},
					},
				},
			},
			"mappings": schema.StringAttribute{
				Description: "Mapping for fields in the index.\nIf specified, this mapping can include: field names, [field data types](https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-types.html), [mapping parameters](https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-params.html).\n**NOTE:**\n- Changing datatypes in the existing _mappings_ will force index to be re-created.\n- Removing field will be ignored by default same as elasticsearch. You need to recreate the index to remove field completely.",
				CustomType:  jsontypes.NormalizedType{},
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					mappingRequiresReplace(),
				},
			},
			"settings": schema.StringAttribute{
				Description: "Configuration options for the index.",
				CustomType:  jsontypes.NormalizedType{},
				Computed:    true,
				Optional:    true,
			},
			"deletion_protection": schema.BoolAttribute{
				Description: "Whether to allow Terraform to destroy the index. Unless this field is set to false in Terraform state, a terraform destroy or terraform apply command that deletes the instance will fail.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *IndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data indexModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	index, diags := r.client.GetIndex(ctx, name)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = data.fromApi(index)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data indexModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	index, diags := data.toApi()
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	name := data.Name.ValueString()
	body := elasticsearch.CreateIndexRequest{
		Aliases:  index.Aliases,
		Mappings: index.Mappings,
		Settings: index.Settings,
	}
	indexResp, diags := r.client.CreateIndex(ctx, name, body)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = data.fromApi(indexResp)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan indexModel
	var state indexModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	indexName := state.Name.ValueString()
	path := path.Empty()

	{
		planAliases := lo.Keys(plan.Aliases)
		stateAliases := lo.Keys(state.Aliases)

		for _, aliasName := range lo.Intersect(planAliases, stateAliases) {
			planAlias := plan.Aliases[aliasName]
			stateAlias := state.Aliases[aliasName]

			// Update aliases that aren't equal
			if !reflect.DeepEqual(planAlias, stateAlias) {
				diags = r.putAlias(ctx, indexName, aliasName, &planAlias)
				resp.Diagnostics.Append(diags...)
				if diags.HasError() {
					return
				}

				state.Aliases[aliasName] = planAlias
			}
		}

		planDiff, stateDiff := lo.Difference(planAliases, stateAliases)

		// Delete aliases missing from the plan
		for _, aliasName := range stateDiff {
			diags = r.client.DeleteIndexAlias(ctx, indexName, aliasName)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			delete(state.Aliases, aliasName)
		}

		// Add aliases missing from the state
		for _, aliasName := range planDiff {
			alias := plan.Aliases[aliasName]
			diags = r.putAlias(ctx, indexName, aliasName, &alias)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			state.Aliases[aliasName] = alias
		}
	}

	{
		planMappings := util.NormalizedTypeToStruct[estypes.TypeMapping](plan.Mappings, path.AtName("mappings"), resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		stateMappings := util.NormalizedTypeToStruct[estypes.TypeMapping](state.Mappings, path.AtName("mappings"), resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if planMappings != nil && !reflect.DeepEqual(planMappings, stateMappings) {
			mappingsResp, diags := r.client.PutIndexMapping(ctx, indexName, elasticsearch.PutIndexMappingRequest{
				DateDetection:      planMappings.DateDetection,
				Dynamic:            planMappings.Dynamic,
				DynamicDateFormats: planMappings.DynamicDateFormats,
				DynamicTemplates:   planMappings.DynamicTemplates,
				FieldNames_:        planMappings.FieldNames_,
				Meta_:              planMappings.Meta_,
				NumericDetection:   planMappings.NumericDetection,
				Properties:         planMappings.Properties,
				Routing_:           planMappings.Routing_,
				Runtime:            planMappings.Runtime,
				Source_:            planMappings.Source_,
			})
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			stateMappings = &mappingsResp.Mappings
			diags = formatMappings(planMappings, stateMappings)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			state.Mappings = util.StructToNormalizedType(stateMappings, path.AtName("mappings"), resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	{
		planSettings := util.NormalizedTypeToStruct[estypes.IndexSettings](plan.Settings, path.AtName("settings"), resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		stateSettings := util.NormalizedTypeToStruct[estypes.IndexSettings](state.Settings, path.AtName("settings"), resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if planSettings != nil && !reflect.DeepEqual(planSettings, stateSettings) {
			settingsResp, diags := r.client.PutIndexSettings(ctx, indexName, *planSettings)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			stateSettings = settingsResp.Settings
			diags = formatSettings(planSettings, stateSettings)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			state.Settings = util.StructToNormalizedType(stateSettings, path.AtName("settings"), resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	state.DeletionProtection = plan.DeletionProtection

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func (r *IndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data indexModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if data.DeletionProtection.ValueBool() {
		resp.Diagnostics.AddError("cannot destroy index without setting deletion_protection=false and running `terraform apply`", "")
		return
	}

	name := data.Name.ValueString()
	diags = r.client.DeleteIndex(ctx, name)
	resp.Diagnostics.Append(diags...)
}

func (r *IndexResource) putAlias(ctx context.Context, indexName string, aliasName string, alias *indexAliasModel) diag.Diagnostics {
	var diags diag.Diagnostics

	path := path.Root("aliases").AtMapKey(aliasName)
	req, d := alias.toApi(path)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	aliasResp, d := r.client.PutIndexAlias(ctx, indexName, aliasName, req)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	d = alias.fromApi(aliasResp, path)
	diags.Append(d...)

	return diags
}

func formatMappings(plan *estypes.TypeMapping, state *estypes.TypeMapping) diag.Diagnostics {
	if plan == nil {
		return nil
	}

	planMap, err := util.StructToMap(plan)
	if err != nil {
		return diag.Diagnostics{diag.NewAttributeErrorDiagnostic(path.Root("mappings"), "struct to map failure", err.Error())}
	}
	planMappings := util.Flatten(planMap)

	stateMap, err := util.StructToMap(state)
	if err != nil {
		return diag.Diagnostics{diag.NewAttributeErrorDiagnostic(path.Root("mappings"), "struct to map failure", err.Error())}
	}
	stateMappings := util.Flatten(stateMap)

	// delete from the state any key not in the plan
	// only track mappings in the configration
	for key := range stateMappings {
		if _, ok := planMappings[key]; !ok {
			delete(stateMappings, key)
		}
	}

	newState, err := util.MapToStruct[estypes.TypeMapping](util.Unflatten(stateMappings))
	if err != nil {
		return diag.Diagnostics{diag.NewAttributeErrorDiagnostic(path.Root("mappings"), "map to struct failure", err.Error())}
	}

	*state = *newState
	return nil
}

func formatSettings(plan *estypes.IndexSettings, state *estypes.IndexSettings) diag.Diagnostics {
	if plan == nil {
		return nil
	}

	planMap, err := util.StructToMap(plan)
	if err != nil {
		return diag.Diagnostics{diag.NewAttributeErrorDiagnostic(path.Root("settings"), "struct to map failure", err.Error())}
	}
	planSettings := util.Flatten(planMap)

	stateMap, err := util.StructToMap(state)
	if err != nil {
		return diag.Diagnostics{diag.NewAttributeErrorDiagnostic(path.Root("settings"), "struct to map failure", err.Error())}
	}
	stateSettings := util.Flatten(stateMap)

	// delete from the state any key not in the plan
	// only track settings in the configration
	for key := range stateSettings {
		if _, ok := planSettings[key]; !ok {
			delete(stateSettings, key)
		}
	}

	newState, err := util.MapToStruct[estypes.IndexSettings](util.Unflatten(stateSettings))
	if err != nil {
		return diag.Diagnostics{diag.NewAttributeErrorDiagnostic(path.Root("settings"), "map to struct failure", err.Error())}
	}

	*state = *newState

	return nil
}

type indexModel struct {
	ID                 types.String               `tfsdk:"id"`
	Name               types.String               `tfsdk:"name"`
	Aliases            map[string]indexAliasModel `tfsdk:"aliases"`
	Mappings           jsontypes.Normalized       `tfsdk:"mappings"`
	Settings           jsontypes.Normalized       `tfsdk:"settings"`
	DeletionProtection types.Bool                 `tfsdk:"deletion_protection"`
}

type indexAliasModel struct {
	Filter        jsontypes.Normalized `tfsdk:"filter"`
	IndexRouting  types.String         `tfsdk:"index_routing"`
	IsHidden      types.Bool           `tfsdk:"is_hidden"`
	IsWriteIndex  types.Bool           `tfsdk:"is_write_index"`
	Routing       types.String         `tfsdk:"routing"`
	SearchRouting types.String         `tfsdk:"search_routing"`
}

func (m *indexModel) toApi() (estypes.IndexState, diag.Diagnostics) {
	var diags diag.Diagnostics

	path := path.Empty()
	output := estypes.IndexState{
		Aliases: util.TransformMap(m.Aliases, func(key string, m indexAliasModel) estypes.Alias {
			return estypes.Alias{
				Filter:        util.NormalizedTypeToStruct[estypes.Query](m.Filter, path.AtName("filter"), diags),
				IndexRouting:  m.IndexRouting.ValueStringPointer(),
				IsHidden:      m.IsHidden.ValueBoolPointer(),
				IsWriteIndex:  m.IsWriteIndex.ValueBoolPointer(),
				Routing:       m.Routing.ValueStringPointer(),
				SearchRouting: m.SearchRouting.ValueStringPointer(),
			}
		}),
		Mappings: util.NormalizedTypeToStruct[estypes.TypeMapping](m.Mappings, path.AtName("mappings"), diags),
		Settings: util.NormalizedTypeToStruct[estypes.IndexSettings](m.Settings, path.AtName("settings"), diags),
	}

	return output, diags
}

func (m *indexModel) fromApi(resp *estypes.IndexState) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	path := path.Empty()
	var diags diag.Diagnostics

	m.ID = m.Name
	m.Aliases = util.TransformMap(resp.Aliases, func(key string, value estypes.Alias) indexAliasModel {
		path := path.AtMapKey("key")
		return indexAliasModel{
			Filter:        util.StructToNormalizedType(value.Filter, path.AtName("filter"), diags),
			IndexRouting:  types.StringPointerValue(value.IndexRouting),
			IsHidden:      types.BoolPointerValue(value.IsHidden),
			IsWriteIndex:  types.BoolPointerValue(value.IsWriteIndex),
			Routing:       types.StringPointerValue(value.Routing),
			SearchRouting: types.StringPointerValue(value.SearchRouting),
		}
	})

	planMappings := util.NormalizedTypeToStruct[estypes.TypeMapping](m.Mappings, path.AtName("mappings"), diags)
	planSettings := util.NormalizedTypeToStruct[estypes.IndexSettings](m.Settings, path.AtName("settings"), diags)
	diags.Append(formatMappings(planMappings, resp.Mappings)...)
	diags.Append(formatSettings(planSettings, resp.Settings)...)

	m.Mappings = util.StructToNormalizedType(resp.Mappings, path.AtName("mappings"), diags)
	m.Settings = util.StructToNormalizedType(resp.Settings, path.AtName("settings"), diags)

	return diags
}

func (m *indexAliasModel) toApi(path path.Path) (elasticsearch.PutAliasRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	output := elasticsearch.PutAliasRequest{
		Filter:        util.NormalizedTypeToStruct[estypes.Query](m.Filter, path.AtName("filter"), diags),
		IndexRouting:  m.IndexRouting.ValueStringPointer(),
		IsWriteIndex:  m.IsWriteIndex.ValueBoolPointer(),
		Routing:       m.Routing.ValueStringPointer(),
		SearchRouting: m.SearchRouting.ValueStringPointer(),
	}

	return output, diags
}

func (m *indexAliasModel) fromApi(resp *estypes.AliasDefinition, path path.Path) diag.Diagnostics {
	if resp == nil {
		return nil
	}

	var diags diag.Diagnostics

	m.Filter = util.StructToNormalizedType(resp.Filter, path.AtName("filter"), diags)
	m.IndexRouting = types.StringPointerValue(resp.IndexRouting)
	m.IsHidden = types.BoolPointerValue(resp.IsHidden)
	m.IsWriteIndex = types.BoolPointerValue(resp.IsWriteIndex)
	m.Routing = types.StringPointerValue(resp.Routing)
	m.SearchRouting = types.StringPointerValue(resp.SearchRouting)

	return diags
}

func mappingRequiresReplace() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			var plan indexModel
			var state indexModel
			path := path.Empty().AtName("mappings")

			diags := req.Plan.Get(ctx, &plan)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			diags = req.State.Get(ctx, &state)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			if plan.Mappings.IsUnknown() || plan.Mappings.IsNull() {
				return
			}

			planMappings := util.NormalizedTypeToStruct[estypes.TypeMapping](plan.Mappings, path, diags)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			stateMappings := util.NormalizedTypeToStruct[estypes.TypeMapping](state.Mappings, path, diags)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			flattenedPlan, err := flattenMapping(planMappings)
			if err != nil {
				resp.Diagnostics.AddAttributeError(path, "failure flattening plan mappings", err.Error())
				return
			}

			flattenedState, err := flattenMapping(stateMappings)
			if err != nil {
				resp.Diagnostics.AddAttributeError(path, "failure flattening state mappings", err.Error())
				return
			}

			planKeys := lo.Keys(flattenedPlan)
			stateKeys := lo.Keys(flattenedState)

			for _, key := range lo.Intersect(planKeys, stateKeys) {
				t1 := flattenedPlan[key]
				t2 := flattenedState[key]
				if t1 != t2 {
					resp.RequiresReplace = true
				}
			}

			_, stateDiff := lo.Difference(planKeys, stateKeys)
			for _, key := range stateDiff {
				diags.AddAttributeWarning(path, fmt.Sprintf("removing %s property in mappings is ignored, if you neeed to remove it completely, please recreate the index", key), "")
			}
		},
		"If the type of a mapping changes, Terraform will destroy and recreate the resource.",
		"If the type of a mapping changes, Terraform will destroy and recreate the resource.",
	)
}

func flattenMapping(mappings *estypes.TypeMapping) (map[string]string, error) {
	flattened := make(map[string]string)

	getPropertyInfo := func(property estypes.Property) (string, map[string]estypes.Property, error) {
		switch prop := property.(type) {
		case *estypes.BinaryProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.BooleanProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.DynamicProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.JoinProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.KeywordProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.MatchOnlyTextProperty:
			return prop.Type, prop.Fields, nil
		case *estypes.PercolatorProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.RankFeatureProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.RankFeaturesProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.SearchAsYouTypeProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.TextProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.VersionProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.WildcardProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.DateNanosProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.DateProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.AggregateMetricDoubleProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.DenseVectorProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.SparseVectorProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.FlattenedProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.NestedProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.ObjectProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.CompletionProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.ConstantKeywordProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.FieldAliasProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.HistogramProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.IpProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.Murmur3HashProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.TokenCountProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.GeoPointProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.GeoShapeProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.PointProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.ShapeProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.ByteNumberProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.DoubleNumberProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.FloatNumberProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.HalfFloatNumberProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.IntegerNumberProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.LongNumberProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.ScaledFloatNumberProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.ShortNumberProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.UnsignedLongNumberProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.DateRangeProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.DoubleRangeProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.FloatRangeProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.IntegerRangeProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.IpRangeProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.LongRangeProperty:
			return prop.Type, prop.Properties, nil
		case *estypes.IcuCollationProperty:
			return prop.Type, prop.Properties, nil
		default:
			return "", nil, fmt.Errorf("unhandled mapping type: %T", prop)
		}
	}

	var flattener func(string, map[string]estypes.Property, map[string]string) error
	flattener = func(path string, src map[string]estypes.Property, dst map[string]string) error {
		if len(path) > 0 {
			path += "."
		}

		for key, props := range src {
			key = path + key

			typ, children, err := getPropertyInfo(props)
			if err != nil {
				return fmt.Errorf("key: %s: %w", key, err)
			}

			dst[key] = typ
			err = flattener(key, children, dst)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := flattener("", mappings.Properties, flattened)
	return flattened, err
}
