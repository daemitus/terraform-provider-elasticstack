package provider

import (
	"context"

	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/config"
	"github.com/daemitus/terraform-provider-elasticstack/internal/elasticsearch"
	"github.com/daemitus/terraform-provider-elasticstack/internal/fleet"
	"github.com/daemitus/terraform-provider-elasticstack/internal/kibana"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var _ provider.Provider = &elasticstackProvider{} // interface check

// ProtoV6ProviderServerFactory wraps the provider in the ProviderServer interface.
func ProtoV6ProviderServerFactory(ctx context.Context, version string) (func() tfprotov6.ProviderServer, error) {
	provider := NewProvider(version)
	providerServer := providerserver.NewProtocol6(provider)
	return providerServer, nil
}

// NewProvider wraps the provider in the Provider interface.
func NewProvider(version string) provider.Provider {
	return &elasticstackProvider{
		version: version,
	}
}

// elasticstackProvider is the provider implementation.
type elasticstackProvider struct {
	version       string
	Elasticsearch *clients.ElasticsearchClient
	Kibana        *clients.KibanaClient
	Fleet         *clients.FleetClient
}

func (p *elasticstackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "elasticstack"
	resp.Version = p.version
}

func (p *elasticstackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	serviceBlock := func(name string) schema.SingleNestedBlock {
		expr := path.Root(name).Expression()
		return schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"api_key": schema.StringAttribute{
					Description: "API Key to use for authentication.",
					Optional:    true,
					Sensitive:   true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(
							expr.AtName("username"),
							expr.AtName("password"),
						),
					},
				},
				"endpoint": schema.StringAttribute{
					Description: "An endpoint that the Terraform provider will point to, this must include the http(s) schema and port number.",
					Optional:    true,
				},
				"insecure": schema.BoolAttribute{
					Description: "Disable TLS certificate validation.",
					Optional:    true,
				},
				"username": schema.StringAttribute{
					Description: "Username to use for authentication.",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.AlsoRequires(expr.AtName("password")),
						stringvalidator.ConflictsWith(expr.AtName("api_key")),
					},
				},
				"password": schema.StringAttribute{
					Description: "Password to use for authentication.",
					Optional:    true,
					Sensitive:   true,
					Validators: []validator.String{
						stringvalidator.AlsoRequires(expr.AtName("username")),
						stringvalidator.ConflictsWith(expr.AtName("api_key")),
					},
				},
			},
		}
	}

	resp.Schema.Blocks = map[string]schema.Block{
		"elasticsearch": serviceBlock("elasticsearch"),
		"kibana":        serviceBlock("kibana"),
		"fleet":         serviceBlock("fleet"),
	}
}

func (p *elasticstackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg config.ProviderModel
	var err error

	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	config, err := config.New().FromSchema(cfg).WithEnv()
	if err != nil {
		diag := diag.NewErrorDiagnostic("Error parsing provider config", err.Error())
		resp.Diagnostics.Append(diag)
	}

	if p.Elasticsearch, err = clients.NewElasticsearchClient(config.Elasticsearch); err != nil {
		diag := diag.NewErrorDiagnostic("Error creating Elasticsearch client", err.Error())
		resp.Diagnostics.Append(diag)
	}
	if p.Kibana, err = clients.NewKibanaClient(config.Kibana); err != nil {
		diag := diag.NewErrorDiagnostic("Error creating Kibana client", err.Error())
		resp.Diagnostics.Append(diag)
	}
	if p.Fleet, err = clients.NewFleetClient(config.Fleet); err != nil {
		diag := diag.NewErrorDiagnostic("Error creating Fleet client", err.Error())
		resp.Diagnostics.Append(diag)
	}
}

func (p *elasticstackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return elasticsearch.NewComponentTemplateResource(p.Elasticsearch) },
		func() resource.Resource { return elasticsearch.NewDataStreamResource(p.Elasticsearch) },
		func() resource.Resource { return elasticsearch.NewIndexLifecycleResource(p.Elasticsearch) },
		func() resource.Resource { return elasticsearch.NewIndexResource(p.Elasticsearch) },
		func() resource.Resource { return elasticsearch.NewIndexTemplateResource(p.Elasticsearch) },
		func() resource.Resource { return elasticsearch.NewIngestPipelineResource(p.Elasticsearch) },
		func() resource.Resource { return elasticsearch.NewWatchResource(p.Elasticsearch) },
		func() resource.Resource { return kibana.NewActionConnectorResource(p.Kibana) },
		func() resource.Resource { return kibana.NewAlertingRuleResource(p.Kibana) },
		func() resource.Resource { return kibana.NewDataViewResourceResource(p.Kibana) },
		func() resource.Resource { return kibana.NewImportSavedObjectsResource(p.Kibana) },
		func() resource.Resource { return kibana.NewSecurityRoleResource(p.Kibana) },
		func() resource.Resource { return kibana.NewSpaceResource(p.Kibana) },
		func() resource.Resource { return fleet.NewAgentPolicyResource(p.Fleet) },
		func() resource.Resource { return fleet.NewIntegrationResource(p.Fleet) },
		func() resource.Resource { return fleet.NewIntegrationPolicyResource(p.Fleet) },
		func() resource.Resource { return fleet.NewOutputResource(p.Fleet) },
		func() resource.Resource { return fleet.NewServerHostResource(p.Fleet) },
	}
}

func (p *elasticstackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return kibana.NewActionConnectorDataSource(p.Kibana) },
		func() datasource.DataSource { return kibana.NewSecurityRoleDataSource(p.Kibana) },
		func() datasource.DataSource { return kibana.NewSpacesDataSource(p.Kibana) },
		func() datasource.DataSource { return fleet.NewEnrollmentTokensDataSource(p.Fleet) },
		func() datasource.DataSource { return fleet.NewIntegrationDataSource(p.Fleet) },
	}
}
