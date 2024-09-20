package fleet

import (
	"context"
	"fmt"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/fleet"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &EnrollmentTokensDataSource{}
)

func NewEnrollmentTokensDataSource(client *clients.FleetClient) *EnrollmentTokensDataSource {
	return &EnrollmentTokensDataSource{client: client}
}

type EnrollmentTokensDataSource struct {
	client *clients.FleetClient
}

func (d *EnrollmentTokensDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, "fleet_enrollment_tokens")
}

func (d *EnrollmentTokensDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema.Description = "Retrieves Elasticsearch API keys used to enroll Elastic Agents in Fleet. See: https://www.elastic.co/guide/en/fleet/current/fleet-enrollment-tokens.html"
	resp.Schema.Attributes = map[string]schema.Attribute{
		"policy_id": schema.StringAttribute{
			Description: "The identifier of the target agent policy. Only the enrollment tokens associated with this agent policy will be selected.",
			Required:    true,
		},
		"tokens": schema.ListNestedAttribute{
			Description: "A list of enrollment tokens.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"key_id": schema.StringAttribute{
						Description: "The unique identifier of the enrollment token.",
						Computed:    true,
					},
					"api_key": schema.StringAttribute{
						Description: "The API key.",
						Computed:    true,
						Sensitive:   true,
					},
					"api_key_id": schema.StringAttribute{
						Description: "The API key identifier.",
						Computed:    true,
					},
					"created_at": schema.StringAttribute{
						Description: "The time at which the enrollment token was created.",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the enrollment token.",
						Computed:    true,
					},
					"active": schema.BoolAttribute{
						Description: "Indicates if the enrollment token is active.",
						Computed:    true,
					},
					"policy_id": schema.StringAttribute{
						Description: "The identifier of the associated agent policy.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *EnrollmentTokensDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data enrollmentTokensModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	var tokens fleet.EnrollmentApiKeys
	policyId := data.PolicyId.ValueString()
	if policyId == "" {
		tokens, diags = d.client.ListEnrollmentTokens(ctx)
	} else {
		tokens, diags = d.client.ReadEnrollmentTokensByPolicy(ctx, policyId)
	}
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.fromApi(tokens)
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

type enrollmentTokensModel struct {
	PolicyId types.String           `tfsdk:"policy_id"`
	Tokens   []enrollmentTokenModel `tfsdk:"tokens"`
}

type enrollmentTokenModel struct {
	KeyId     types.String `tfsdk:"key_id"`
	ApiKey    types.String `tfsdk:"api_key"`
	ApiKeyId  types.String `tfsdk:"api_key_id"`
	CreatedAt types.String `tfsdk:"created_at"`
	Name      types.String `tfsdk:"name"`
	Active    types.Bool   `tfsdk:"active"`
	PolicyId  types.String `tfsdk:"policy_id"`
}

func (m *enrollmentTokensModel) fromApi(data fleet.EnrollmentApiKeys) {
	m.Tokens = []enrollmentTokenModel{}
	for _, token := range data {
		t := enrollmentTokenModel{}
		t.fromApi(token)
		m.Tokens = append(m.Tokens, t)
	}
}

func (m *enrollmentTokenModel) fromApi(data fleet.EnrollmentApiKey) {
	m.KeyId = types.StringValue(data.Id)
	m.Active = types.BoolValue(data.Active)
	m.ApiKey = types.StringValue(data.ApiKey)
	m.ApiKeyId = types.StringValue(data.ApiKeyId)
	m.CreatedAt = types.StringValue(data.CreatedAt)
	m.Name = types.StringPointerValue(data.Name)
	m.PolicyId = types.StringPointerValue(data.PolicyId)
}
