package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/config"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/fleet"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana"
	"github.com/elastic/terraform-provider-elasticstack/internal/models"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/go-version"
	fwdiags "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CompositeId struct {
	ClusterId  string
	ResourceId string
}

func CompositeIdFromStr(id string) (*CompositeId, diag.Diagnostics) {
	var diags diag.Diagnostics
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Wrong resource ID.",
			Detail:   "Resource ID must have following format: <cluster_uuid>/<resource identifier>",
		})
		return nil, diags
	}
	return &CompositeId{
			ClusterId:  idParts[0],
			ResourceId: idParts[1],
		},
		diags
}

func ResourceIDFromStr(id string) (string, diag.Diagnostics) {
	compID, diags := CompositeIdFromStr(id)
	if diags.HasError() {
		return "", diags
	}
	return compID.ResourceId, nil
}

func (c *CompositeId) String() string {
	return fmt.Sprintf("%s/%s", c.ClusterId, c.ResourceId)
}

type ApiClient struct {
	elasticsearch            *elasticsearch.Client
	elasticsearchClusterInfo *models.ClusterInfo
	kibana                   *kibana.Client
	fleet                    *fleet.Client
	version                  string
}

func NewApiClientFuncFromSDK(version string) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return newApiClientFromSDK(d, version)
	}
}

func NewAcceptanceTestingClient() (*ApiClient, error) {
	version := "tf-acceptance-testing"
	cfg := config.NewFromEnv(version)

	es, err := elasticsearch.NewClient(*cfg.Elasticsearch)
	if err != nil {
		return nil, err
	}

	kibanaClient, err := kibana.NewClient(*cfg.Kibana)
	if err != nil {
		return nil, err
	}

	fleetClient, err := fleet.NewClient(*cfg.Fleet)
	if err != nil {
		return nil, err
	}

	return &ApiClient{
			elasticsearch: es,
			kibana:        kibanaClient,
			fleet:         fleetClient,
			version:       version,
		},
		nil
}

func NewApiClientFromFramework(ctx context.Context, cfg config.ProviderConfiguration, version string) (*ApiClient, fwdiags.Diagnostics) {
	clientCfg, diags := config.NewFromFramework(ctx, cfg, version)
	if diags.HasError() {
		return nil, diags
	}

	client, err := newApiClientFromConfig(clientCfg, version)
	if err != nil {
		return nil, fwdiags.Diagnostics{
			fwdiags.NewErrorDiagnostic("Failed to create API client", err.Error()),
		}
	}

	return client, nil
}

func ConvertProviderData(providerData any) (*ApiClient, fwdiags.Diagnostics) {
	var diags fwdiags.Diagnostics

	if providerData == nil {
		return nil, diags
	}

	client, ok := providerData.(*ApiClient)
	if !ok {
		diags.AddError(
			"Unexpected Provider Data",
			fmt.Sprintf("Expected *ApiClient, got: %T. Please report this issue to the provider developers.", providerData),
		)

		return nil, diags
	}
	if client == nil {
		diags.AddError(
			"Unconfigured Client",
			"Expected configured client. Please report this issue to the provider developers.",
		)
	}
	return client, diags
}

func MaybeNewApiClientFromFrameworkResource(ctx context.Context, esConnList types.List, defaultClient *ApiClient) (*ApiClient, fwdiags.Diagnostics) {
	var esConns []config.ElasticsearchConnection
	if diags := esConnList.ElementsAs(ctx, &esConns, true); diags.HasError() {
		return nil, diags
	}

	if len(esConns) == 0 {
		return defaultClient, nil
	}

	cfg, diags := config.NewFromFramework(ctx, config.ProviderConfiguration{Elasticsearch: esConns}, defaultClient.version)
	if diags.HasError() {
		return nil, diags
	}

	esClient, err := buildEsClient(cfg)
	if err != nil {
		return nil, fwdiags.Diagnostics{fwdiags.NewErrorDiagnostic(err.Error(), err.Error())}
	}

	return &ApiClient{
		elasticsearch:            esClient,
		elasticsearchClusterInfo: defaultClient.elasticsearchClusterInfo,
		kibana:                   defaultClient.kibana,
		fleet:                    defaultClient.fleet,
		version:                  defaultClient.version,
	}, diags
}

func NewApiClientFromSDKResource(d *schema.ResourceData, meta interface{}) (*ApiClient, diag.Diagnostics) {
	defaultClient := meta.(*ApiClient)
	version := defaultClient.version
	resourceConfig, diags := config.NewFromSDKResource(d, version)
	if diags.HasError() {
		return nil, diags
	}

	if resourceConfig == nil {
		return defaultClient, nil
	}

	esClient, err := buildEsClient(*resourceConfig)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return &ApiClient{
		elasticsearch:            esClient,
		elasticsearchClusterInfo: defaultClient.elasticsearchClusterInfo,
		kibana:                   defaultClient.kibana,
		fleet:                    defaultClient.fleet,
		version:                  version,
	}, diags
}

func (a *ApiClient) GetESClient() (*elasticsearch.Client, error) {
	if a.elasticsearch == nil {
		return nil, errors.New("elasticsearch client not found")
	}

	return a.elasticsearch, nil
}

func (a *ApiClient) GetKibanaClient() (*kibana.Client, error) {
	if a.kibana == nil {
		return nil, errors.New("kibana client not found")
	}

	return a.kibana, nil
}

func (a *ApiClient) GetFleetClient() (*fleet.Client, error) {
	if a.fleet == nil {
		return nil, errors.New("fleet client not found")
	}

	return a.fleet, nil
}

func (a *ApiClient) ID(ctx context.Context, resourceId string) (*CompositeId, diag.Diagnostics) {
	var diags diag.Diagnostics
	clusterId, diags := a.ClusterID(ctx)
	if diags.HasError() {
		return nil, diags
	}
	return &CompositeId{*clusterId, resourceId}, diags
}

func (a *ApiClient) serverInfo(ctx context.Context) (*models.ClusterInfo, diag.Diagnostics) {
	if a.elasticsearchClusterInfo != nil {
		return a.elasticsearchClusterInfo, nil
	}

	var diags diag.Diagnostics
	esClient, err := a.GetESClient()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	res, err := esClient.Info(esClient.Info.WithContext(ctx))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	defer res.Body.Close()
	if diags := utils.CheckError(res, "Unable to connect to the Elasticsearch cluster"); diags.HasError() {
		return nil, diags
	}

	info := models.ClusterInfo{}
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, diag.FromErr(err)
	}
	// cache info
	a.elasticsearchClusterInfo = &info

	return &info, diags
}

func (a *ApiClient) ServerVersion(ctx context.Context) (*version.Version, diag.Diagnostics) {
	info, diags := a.serverInfo(ctx)
	if diags.HasError() {
		return nil, diags
	}

	rawVersion := info.Version.Number
	serverVersion, err := version.NewVersion(rawVersion)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return serverVersion, nil
}

func (a *ApiClient) ServerFlavor(ctx context.Context) (string, diag.Diagnostics) {
	info, diags := a.serverInfo(ctx)
	if diags.HasError() {
		return "", diags
	}

	return info.Version.BuildFlavor, nil
}

func (a *ApiClient) ClusterID(ctx context.Context) (*string, diag.Diagnostics) {
	info, diags := a.serverInfo(ctx)
	if diags.HasError() {
		return nil, diags
	}

	if uuid := info.ClusterUUID; uuid != "" && uuid != "_na_" {
		tflog.Trace(ctx, fmt.Sprintf("cluster UUID: %s", uuid))
		return &uuid, diags
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to get cluster UUID",
		Detail: `Unable to get cluster UUID.
		There might be a problem with permissions or cluster is still starting up and UUID has not been populated yet.`,
	})
	return nil, diags
}

func buildEsClient(cfg config.Client) (*elasticsearch.Client, error) {
	if cfg.Elasticsearch == nil {
		return nil, nil
	}

	es, err := elasticsearch.NewClient(*cfg.Elasticsearch)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Elasticsearch client: %w", err)
	}

	return es, nil
}

func buildKibanaClient(cfg config.Client) (*kibana.Client, error) {
	client, err := kibana.NewClient(*cfg.Kibana)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Kibana client: %w", err)
	}

	return client, nil
}

func buildFleetClient(cfg config.Client) (*fleet.Client, error) {
	client, err := fleet.NewClient(*cfg.Fleet)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Fleet client: %w", err)
	}

	return client, nil
}

func newApiClientFromSDK(d *schema.ResourceData, version string) (*ApiClient, diag.Diagnostics) {
	cfg, diags := config.NewFromSDK(d, version)
	if diags.HasError() {
		return nil, diags
	}

	client, err := newApiClientFromConfig(cfg, version)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}

func newApiClientFromConfig(cfg config.Client, version string) (*ApiClient, error) {
	client := &ApiClient{version: version}

	if cfg.Elasticsearch != nil {
		esClient, err := buildEsClient(cfg)
		if err != nil {
			return nil, err
		}
		client.elasticsearch = esClient
	}

	if cfg.Kibana != nil {
		kibanaClient, err := buildKibanaClient(cfg)
		if err != nil {
			return nil, err
		}
		client.kibana = kibanaClient
	}

	if cfg.Fleet != nil {
		fleetClient, err := buildFleetClient(cfg)
		if err != nil {
			return nil, err
		}
		client.fleet = fleetClient
	}

	return client, nil
}
