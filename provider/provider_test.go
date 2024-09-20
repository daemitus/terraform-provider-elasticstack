package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/daemitus/terraform-provider-elasticstack/internal/config"
	"github.com/daemitus/terraform-provider-elasticstack/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/samber/lo"
)

func TestProvider(t *testing.T) {
	provider.NewProvider("test")
}

func TestProviderElasticsearchConfiguration(t *testing.T) {
	cfg := lo.Must(config.New().WithEnv())
	os.Unsetenv("ELASTICSEARCH_ENDPOINT")
	os.Unsetenv("ELASTICSEARCH_USERNAME")
	os.Unsetenv("ELASTICSEARCH_PASSWORD")
	os.Unsetenv("ELASTICSEARCH_API_KEY")
	os.Unsetenv("ELASTICSEARCH_INSECURE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderElasticsearchBasicConfiguration(cfg.Elasticsearch),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.elasticstack_elasticsearch_cluster_health.test", "status", "green"),
				),
			},
			{
				Config: testProviderElasticsearchApiKeyConfiguration(cfg.Elasticsearch),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.elasticstack_elasticsearch_cluster_health.test", "status", "green"),
				),
			},
		},
	})
}

func TestProviderKibanaConfiguration(t *testing.T) {
	cfg := lo.Must(config.New().WithEnv())
	os.Unsetenv("KIBANA_ENDPOINT")
	os.Unsetenv("KIBANA_USERNAME")
	os.Unsetenv("KIBANA_PASSWORD")
	os.Unsetenv("KIBANA_API_KEY")
	os.Unsetenv("KIBANA_INSECURE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderKibanaBasicConfiguration(cfg.Kibana),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("elasticstack_kibana_space.acc_test", "name"),
				),
			},
			{
				Config: testProviderKibanaApiKeyConfiguration(cfg.Kibana),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("elasticstack_kibana_space.acc_test", "name"),
				),
			},
		},
	})
}

func TestProviderFleetConfiguration(t *testing.T) {
	cfg := lo.Must(config.New().WithEnv())
	os.Unsetenv("FLEET_ENDPOINT")
	os.Unsetenv("FLEET_USERNAME")
	os.Unsetenv("FLEET_PASSWORD")
	os.Unsetenv("FLEET_API_KEY")
	os.Unsetenv("FLEET_INSECURE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderFleetBasicConfiguration(cfg.Fleet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.elasticstack_fleet_enrollment_tokens.test", "tokens.#"),
				),
			},
			{
				Config: testProviderFleetApiKeyConfiguration(cfg.Fleet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.elasticstack_fleet_enrollment_tokens.test", "tokens.#"),
				),
			},
		},
	})
}

func testProviderElasticsearchBasicConfiguration(cfg config.ServiceConfig) string {
	return fmt.Sprintf(`
provider "elasticstack" {
	elasticsearch {
		endpoint = "%s"
		username = "%s"
		password = "%s"
	}
}

data "elasticstack_elasticsearch_cluster_health" "test" {}
`, cfg.Endpoint, cfg.Username, cfg.Password)
}

func testProviderElasticsearchApiKeyConfiguration(cfg config.ServiceConfig) string {
	return fmt.Sprintf(`
provider "elasticstack" {
	elasticsearch {
		endpoint = "%s"
		api_key	= "%s"
	}
}

data "elasticstack_elasticsearch_cluster_health" "acc_test" {}
`, cfg.Endpoint, cfg.ApiKey)
}

func testProviderKibanaBasicConfiguration(cfg config.ServiceConfig) string {
	return fmt.Sprintf(`
provider "elasticstack" {
	kibana {
		endpoint = "%s"
		username = "%s"
		password = "%s"
	}
}

resource "elasticstack_kibana_space" "acc_test" {
	space_id = "acc_test_space"
	name = "Acceptance Test Space"
}
`, cfg.Endpoint, cfg.Username, cfg.Password)
}

func testProviderKibanaApiKeyConfiguration(cfg config.ServiceConfig) string {
	return fmt.Sprintf(`
provider "elasticstack" {
	kibana {
		endpoint = "%s"
		api_key	= "%s"
	}
}

resource "elasticstack_kibana_space" "acc_test" {
	space_id = "acc_test_space"
	name = "Acceptance Test Space"
}
`, cfg.Endpoint, cfg.ApiKey)
}

func testProviderFleetBasicConfiguration(cfg config.ServiceConfig) string {
	return fmt.Sprintf(`
provider "elasticstack" {
	fleet {
		endpoint = "%s"
		username = "%s"
		password = "%s"
	}
}

data "elasticstack_fleet_enrollment_tokens" "test" {}
`, cfg.Endpoint, cfg.Username, cfg.Password)
}

func testProviderFleetApiKeyConfiguration(cfg config.ServiceConfig) string {
	return fmt.Sprintf(`
provider "elasticstack" {
	fleet {
		endpoint = "%s"
		api_key	= "%s"
	}
}

data "elasticstack_fleet_enrollment_tokens" "test" {}
`, cfg.Endpoint, cfg.ApiKey)
}
