package fleet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceOutputElasticsearch(t *testing.T) {
	policyName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceOutputDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOutputCreateElasticsearch(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "name", fmt.Sprintf("Elasticsearch Output %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "id", "elasticsearch-output"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "type", "elasticsearch"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "config_yaml", "\"ssl.verification_mode\": \"none\"\n"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "default_integrations", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "default_monitoring", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "hosts.0", "http://elasticsearch:9200"),
				),
			},
			{
				Config: testAccResourceOutputUpdateElasticsearch(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "name", fmt.Sprintf("Updated Elasticsearch Output %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "id", "elasticsearch-output"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "type", "elasticsearch"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "config_yaml", "\"ssl.verification_mode\": \"none\"\n"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "default_integrations", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "default_monitoring", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "hosts.0", "http://elasticsearch:9200"),
				),
			},
		},
	})
}

func TestAccResourceOutputRemoteElasticsearch(t *testing.T) {
	policyName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceOutputDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOutputCreateRemoteElasticsearch(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "name", fmt.Sprintf("Remote Elasticsearch Output %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "id", "remote-elasticsearch-output"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "type", "remote_elasticsearch"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "config_yaml", "\"ssl.verification_mode\": \"none\"\n"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "default_integrations", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "default_monitoring", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "hosts.0", "http://elasticsearch:9200"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "service_token", "2682f20b-f3bb-4d60-9a45-e3b13809fbf5"),
				),
			},
			{
				Config: testAccResourceOutputUpdateRemoteElasticsearch(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "name", fmt.Sprintf("Updated Remote Elasticsearch Output %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "id", "remote-elasticsearch-output"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "type", "remote_elasticsearch"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "config_yaml", "\"ssl.verification_mode\": \"none\"\n"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "default_integrations", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "default_monitoring", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "hosts.0", "http://elasticsearch:9200"),
					resource.TestCheckResourceAttr("elasticstack_fleet_output.test_output", "service_token", "2682f20b-f3bb-4d60-9a45-e3b13809fbf5"),
				),
			},
		},
	})
}

func testAccResourceOutputCreateElasticsearch(id string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_output" "test_output" {
	name = "Elasticsearch Output %s"
	output_id = "elasticsearch-output"
	type = "elasticsearch"
	config_yaml = yamlencode({"ssl.verification_mode" : "none"})
	default_integrations = false
	default_monitoring = false
	hosts = ["http://elasticsearch:9200"]
}`, id)
}

func testAccResourceOutputUpdateElasticsearch(id string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_output" "test_output" {
	name = "Updated Elasticsearch Output %s"
	output_id = "elasticsearch-output"
	type = "elasticsearch"
	config_yaml = yamlencode({"ssl.verification_mode" : "none"})
	default_integrations = false
	default_monitoring = false
	hosts = ["http://elasticsearch:9200"]
}`, id)
}

func testAccResourceOutputCreateRemoteElasticsearch(id string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_output" "test_output" {
	name = "Remote Elasticsearch Output %s"
	output_id = "remote-elasticsearch-output"
	type = "remote_elasticsearch"
	config_yaml = yamlencode({"ssl.verification_mode" : "none"})
	default_integrations = false
	default_monitoring = false
	hosts = ["http://elasticsearch:9200"]
	service_token = "2682f20b-f3bb-4d60-9a45-e3b13809fbf5"
}`, id)
}

func testAccResourceOutputUpdateRemoteElasticsearch(id string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_output" "test_output" {
	name = "Updated Remote Elasticsearch Output %s"
	output_id = "remote-elasticsearch-output"
	type = "remote_elasticsearch"
	config_yaml = yamlencode({"ssl.verification_mode" : "none"})
	default_integrations = false
	default_monitoring = false
	hosts = ["http://elasticsearch:9200"]
	service_token = "2682f20b-f3bb-4d60-9a45-e3b13809fbf5"
}`, id)
}

func checkResourceOutputDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccFleetClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_fleet_output" {
			continue
		}

		output, diags := client.ReadOutput(ctx, rs.Primary.ID)
		if diags.HasError() {
			return fmt.Errorf(diags.Errors()[0].Summary())
		}
		if output != nil {
			return fmt.Errorf("output id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
