package fleet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceIntegrationPolicy(t *testing.T) {
	policyName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceIntegrationPolicyDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIntegrationPolicy(policyName, "Created", `[]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "name", policyName),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "description", "Created Integration Policy"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "integration_name", "tcp"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "integration_version", "1.16.0"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.tcp-tcp.enabled`, "true"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.tcp-tcp.streams.tcp.generic.enabled`, "true"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.tcp-tcp.streams.tcp.generic.vars_json`, `{"custom":"","data_stream.dataset":"tcp.generic","listen_address":"localhost","listen_port":8080,"ssl":"","syslog_options":"field: message","tags":[]}`),
				),
			},
			{
				Config: testAccResourceIntegrationPolicy(policyName, "Updated", `["redfish","bluefish"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "name", policyName),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "description", "Updated Integration Policy"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "integration_name", "tcp"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "integration_version", "1.16.0"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.tcp-tcp.enabled`, "true"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.tcp-tcp.streams.tcp.generic.enabled`, "true"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.tcp-tcp.streams.tcp.generic.vars_json`, `{"custom":"","data_stream.dataset":"tcp.generic","listen_address":"localhost","listen_port":8080,"ssl":"","syslog_options":"field: message","tags":["redfish","bluefish"]}`),
				),
			},
		},
	})
}

func TestAccResourceIntegrationPolicyWithSecrets(t *testing.T) {
	policyName := acctest.RandString(22)

	checkVarsJson := func(key string, expected string) resource.CheckResourceAttrWithFunc {
		return func(value string) error {
			vars, err := util.JsonUnmarshalS[map[string]any](value)
			if err != nil {
				return err
			}
			connString, ok := vars[key]
			if !ok {
				return fmt.Errorf(`vars_json missing key "%s"`, key)
			}
			if connString != expected {
				return fmt.Errorf(`vars_json key "%s" expected: "%s", actual: "%s"`, key, expected, value)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceIntegrationPolicyDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIntegrationPolicyWithSecrets(policyName, "Created"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "name", policyName),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "description", "Created Integration Policy"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "integration_name", "m365_defender"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "integration_version", "2.14.3"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.m365_defender-azure-eventhub.enabled`, "true"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.m365_defender-azure-eventhub.streams.m365_defender.event.enabled`, "true"),
					resource.TestCheckResourceAttrWith("elasticstack_fleet_integration_policy.test_policy", `inputs.m365_defender-azure-eventhub.streams.m365_defender.event.vars_json`, checkVarsJson("connection_string", "Endpoint=sb://placeholder-connection-string")),
				),
			},
			{
				Config: testAccResourceIntegrationPolicyWithSecrets(policyName, "Updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "name", policyName),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "description", "Updated Integration Policy"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "integration_name", "m365_defender"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", "integration_version", "2.14.3"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.m365_defender-azure-eventhub.enabled`, "true"),
					resource.TestCheckResourceAttr("elasticstack_fleet_integration_policy.test_policy", `inputs.m365_defender-azure-eventhub.streams.m365_defender.event.enabled`, "true"),
					resource.TestCheckResourceAttrWith("elasticstack_fleet_integration_policy.test_policy", `inputs.m365_defender-azure-eventhub.streams.m365_defender.event.vars_json`, checkVarsJson("connection_string", "Endpoint=sb://placeholder-connection-string")),
				),
			},
		},
	})
}

func testAccResourceIntegrationPolicy(id string, prefix string, tags string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_integration" "test_policy" {
	name = "tcp"
	version = "1.16.0"
	force = true
}

resource "elasticstack_fleet_agent_policy" "test_policy" {
	name = "Agent Policy %s"
	namespace = "default"
	description = "IntegrationPolicyTest Agent Policy"
	monitor_logs = true
	monitor_metrics = true
	skip_destroy = false
}

data "elasticstack_fleet_enrollment_tokens" "test_policy" {
	policy_id = elasticstack_fleet_agent_policy.test_policy.policy_id
}

resource "elasticstack_fleet_integration_policy" "test_policy" {
	name = "%s"
	namespace = "default"
	description = "%s Integration Policy"
	agent_policy_id = elasticstack_fleet_agent_policy.test_policy.policy_id
	integration_name = elasticstack_fleet_integration.test_policy.name
	integration_version = elasticstack_fleet_integration.test_policy.version
	inputs = {
		tcp-tcp = {
			streams = {
				"tcp.generic" = {
					enabled = true
					vars_json = jsonencode({
						listen_address = "localhost"
						listen_port = 8080
						"data_stream.dataset" = "tcp.generic"
						tags = %s
						syslog_options = "field: message"
						ssl = ""
						custom = ""
					})
				}
			}
		}
	}
}`, id, id, prefix, tags)
}

func testAccResourceIntegrationPolicyWithSecrets(id string, prefix string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_integration" "test_policy" {
	name = "m365_defender"
	version = "2.14.3"
	force = true
}

resource "elasticstack_fleet_agent_policy" "test_policy" {
	name = "Agent Policy %s"
	namespace = "default"
	description = "IntegrationPolicyTest Agent Policy"
	monitor_logs = true
	monitor_metrics = true
	skip_destroy = false
}

data "elasticstack_fleet_enrollment_tokens" "test_policy" {
	policy_id = elasticstack_fleet_agent_policy.test_policy.policy_id
}

resource "elasticstack_fleet_integration_policy" "test_policy" {
	name = "%s"
	namespace = "default"
	description = "%s Integration Policy"
	agent_policy_id = elasticstack_fleet_agent_policy.test_policy.policy_id
	integration_name = elasticstack_fleet_integration.test_policy.name
	integration_version = elasticstack_fleet_integration.test_policy.version
	inputs = {
		"m365_defender-azure-eventhub" = {
			enabled = true
			streams = {
				"m365_defender.event" = {
					enabled = true
					vars_json = jsonencode({
						eventhub = "placeholder-eventhub"
						consumer_group = "placeholder-consumer-group"
						connection_string = "Endpoint=sb://placeholder-connection-string"
						storage_account = "placeholder-storage-account"
						storage_account_key = "placeholder-storage-account-key"
						tags = ["forwarded"]
						processors = ""
						preserve_original_event = false
						preserve_duplicate_custom_fields = false
					})
				}
			}
		}
		"m365_defender-httpjson" = {
			enabled = true
			vars_json = jsonencode({
				login_url = "https://login.microsoftonline.com"
				token_endpoint = "oauth2/v2.0/token"
				client_id = "placeholder-client-id"
				client_secret = "placeholder-client-secret"
				tenant_id = "placeholder-tenant-id"
				ssl = ""
			})
			streams = {
				"m365_defender.alert" = {
					enabled = false
					vars_json = jsonencode({
						request_url = "https://graph.microsoft.com"
						initial_interval = "24h"
						interval = "5m"
						batch_size = 2000
						http_client_timeout = "30s"
						tags = ["forwarded"]
						preserve_original_event = false
						preserve_duplicate_custom_fields = false
						processors = ""
					})
				}
				"m365_defender.incident" = {
					enabled = true
					vars_json = jsonencode({
						request_url = "https://graph.microsoft.com"
						initial_interval = "168h"
						interval = "1m"
						batch_size = 50
						http_client_timeout = "30s"
						tags = ["forwarded"]
						processors = ""
						preserve_original_event = false
						preserve_duplicate_custom_fields = false
					})
				}
				"m365_defender.log" = {
					enabled = false
					vars_json = jsonencode({
						interval = "5m"
						initial_interval = "168h"
						request_url = "https://api.security.microsoft.com"
						tags = ["forwarded"]
						processors = ""
						preserve_original_event = false
					})
				}
			}
		}
	}
}`, id, id, prefix)
}

func checkResourceIntegrationPolicyDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccFleetClient()

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "elasticstack_fleet_integration_policy":
			integrationPolicy, diags := client.ReadPackagePolicy(ctx, rs.Primary.ID)
			if diags.HasError() {
				return fmt.Errorf(diags.Errors()[0].Summary())
			}
			if integrationPolicy != nil {
				return fmt.Errorf("integration policy id=%v still exists, but it should have been removed", rs.Primary.ID)
			}
		case "elasticstack_fleet_agent_policy":
			agentPolicy, diags := client.ReadAgentPolicy(ctx, rs.Primary.ID)
			if diags.HasError() {
				return fmt.Errorf(diags.Errors()[0].Summary())
			}
			if agentPolicy != nil {
				return fmt.Errorf("agent policy id=%v still exists, but it should have been removed", rs.Primary.ID)
			}
		default:
			continue
		}
	}
	return nil
}
