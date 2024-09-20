package fleet_test

import (
	"fmt"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceEnrollmentTokensByPolicy(t *testing.T) {
	policyId := acctest.RandUUID()
	policyName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:             checkResourceAgentPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnrollmentTokensByPolicy(policyId, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.elasticstack_fleet_enrollment_tokens.test", "policy_id", policyId),
					resource.TestCheckResourceAttr("data.elasticstack_fleet_enrollment_tokens.test", "tokens.0.policy_id", policyId),
				),
			},
		},
	})
}

func testAccDataSourceEnrollmentTokensByPolicy(policyId string, name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_agent_policy" "test" {
	policy_id = "%s"
	name = "Policy %s"
	namespace = "default"
	description = "Agent Policy for testing Enrollment Tokens"
}

data "elasticstack_fleet_enrollment_tokens" "test" {
	policy_id = elasticstack_fleet_agent_policy.test.policy_id
}
`, policyId, name)
}
