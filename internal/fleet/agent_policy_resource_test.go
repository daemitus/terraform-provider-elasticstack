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

func TestAccResourceAgentPolicy(t *testing.T) {
	policyName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:             checkResourceAgentPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAgentPolicyCreate(policyName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "name", fmt.Sprintf("Policy %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "namespace", "default"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "description", "Test Agent Policy"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "monitor_logs", "true"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "monitor_metrics", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "skip_destroy", "false"),
				),
			},
			{
				Config: testAccResourceAgentPolicyUpdate(policyName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "name", fmt.Sprintf("Updated Policy %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "namespace", "default"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "description", "This policy was updated"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "monitor_logs", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "monitor_metrics", "true"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "skip_destroy", "false"),
				),
			},
		},
	})
}

func TestAccResourceAgentPolicySkipDestroy(t *testing.T) {
	policyName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceAgentPolicySkipDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAgentPolicyCreate(policyName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "name", fmt.Sprintf("Policy %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "namespace", "default"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "description", "Test Agent Policy"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "monitor_logs", "true"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "monitor_metrics", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_agent_policy.test_policy", "skip_destroy", "true"),
				),
			},
		},
	})
}

func testAccResourceAgentPolicyCreate(id string, skipDestroy bool) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_agent_policy" "test_policy" {
	name = "Policy %s"
	namespace = "default"
	description = "Test Agent Policy"
	monitor_logs = true
	monitor_metrics = false
	skip_destroy = %t
}`, id, skipDestroy)
}

func testAccResourceAgentPolicyUpdate(id string, skipDestroy bool) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_agent_policy" "test_policy" {
	name = "Updated Policy %s"
	namespace = "default"
	description = "This policy was updated"
	monitor_logs = false
	monitor_metrics = true
	skip_destroy = %t
}`, id, skipDestroy)
}

func checkResourceAgentPolicyDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccFleetClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_fleet_agent_policy" {
			continue
		}

		agentPolicy, diags := client.ReadAgentPolicy(ctx, rs.Primary.ID)
		if diags.HasError() {
			return fmt.Errorf(diags.Errors()[0].Summary())
		}
		if agentPolicy != nil {
			return fmt.Errorf("agent policy id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}

	return nil
}

func checkResourceAgentPolicySkipDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccFleetClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_fleet_agent_policy" {
			continue
		}

		agentPolicy, diags := client.ReadAgentPolicy(ctx, rs.Primary.ID)
		if diags.HasError() {
			return fmt.Errorf(diags.Errors()[0].Summary())
		}
		if agentPolicy == nil {
			return fmt.Errorf("agent policy id=%v does not exist, but should still exist when skip_destroy is true", rs.Primary.ID)
		}

		diags = client.DeleteAgentPolicy(ctx, rs.Primary.ID)
		if diags.HasError() {
			return fmt.Errorf(diags.Errors()[0].Summary())
		}
	}

	return nil
}
