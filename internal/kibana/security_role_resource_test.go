package kibana_test

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

func TestAccResourceKibanaSecurityRole(t *testing.T) {
	roleName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceSecurityRoleDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceKibanaSecurityRoleCreate(roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_kibana_security_role.test", "name", roleName),
					resource.TestCheckNoResourceAttr("elasticstack_kibana_security_role.test", "elasticsearch.run_as.#"),
					util.TestCheckResourceListAttr("elasticstack_kibana_security_role.test", "elasticsearch.indices.0.names", []string{"sample"}),
					util.TestCheckResourceListAttr("elasticstack_kibana_security_role.test", "elasticsearch.indices.0.field_security.grant", []string{"sample"}),
					resource.TestCheckNoResourceAttr("elasticstack_kibana_security_role.test", "kibana.0.base.#"),
					util.TestCheckResourceListAttr("elasticstack_kibana_security_role.test", "kibana.0.feature.discover", []string{"minimal_read", "url_create", "store_search_session"}),
					util.TestCheckResourceListAttr("elasticstack_kibana_security_role.test", "kibana.0.spaces", []string{"default"}),
				),
			},
			{
				Config: testAccResourceKibanaSecurityRoleUpdate(roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_kibana_security_role.test", "name", roleName),
					resource.TestCheckNoResourceAttr("elasticstack_kibana_security_role.test", "elasticsearch.indices.0.field_security.#"),
					util.TestCheckResourceListAttr("elasticstack_kibana_security_role.test", "elasticsearch.run_as", []string{"elastic", "kibana"}),
					util.TestCheckResourceListAttr("elasticstack_kibana_security_role.test", "kibana.0.base", []string{"all"}),
					resource.TestCheckNoResourceAttr("elasticstack_kibana_security_role.test", "kibana.0.feature.#"),
					util.TestCheckResourceListAttr("elasticstack_kibana_security_role.test", "kibana.0.spaces", []string{"default"}),
				),
			},
		},
	})
}

func testAccResourceKibanaSecurityRoleCreate(roleName string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_kibana_security_role" "test" {
	name = "%s"
	elasticsearch = {
		cluster = ["create_snapshot"]
		indices = [
		{
			field_security = {
				grant = ["sample"]
			}
			names = ["sample"]
			privileges = ["create", "read", "write"]
		}
	]
	remote_indices = [
		{
			clusters = ["*"]
			field_security = {
				grant = ["sample"]
			}
			names = ["sample"]
			privileges = ["create", "read", "write"]
		}
	]
	}
	kibana = [
		{
			feature = {
				actions = ["read"]
				advancedSettings = ["read"]
				discover = ["minimal_read", "url_create", "store_search_session"]
				generalCases = ["minimal_read", "cases_delete"]
				observabilityCases = ["minimal_read", "cases_delete"]
				osquery = ["minimal_read", "live_queries_all", "run_saved_queries", "saved_queries_read", "packs_all"]
				rulesSettings = ["minimal_read", "readFlappingSettings"]
				securitySolutionCases = ["minimal_read", "cases_delete"]
				visualize = ["minimal_read", "url_create"]
			}
			spaces = ["default"]
		}
	]
}`, roleName)
}

func testAccResourceKibanaSecurityRoleUpdate(roleName string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_kibana_security_role" "test" {
	name = "%s"
	elasticsearch = {
		cluster = ["create_snapshot"]
		indices = [
			{
				names = ["sample"]
				privileges = ["create", "read", "write"]
			}
		]
		run_as = ["elastic", "kibana"]
	}
	kibana = [
		{
			base = ["all"]
			spaces = ["default"]
		}
	]
}`, roleName)
}

func checkResourceSecurityRoleDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccKibanaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_kibana_security_role" {
			continue
		}

		role, diags := client.ReadRole(ctx, rs.Primary.ID)
		if diags.HasError() {
			return fmt.Errorf(diags.Errors()[0].Summary())
		}
		if role != nil {
			return fmt.Errorf("role id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
