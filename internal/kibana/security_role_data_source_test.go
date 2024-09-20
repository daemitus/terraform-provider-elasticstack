package kibana_test

import (
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceKibanaSecurityRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecurityRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.elasticstack_kibana_security_role.test", "name", "data_source_test"),
					resource.TestCheckNoResourceAttr("data.elasticstack_kibana_security_role.test", "elasticsearch.indices.0.field_security.#"),
					util.TestCheckResourceListAttr("data.elasticstack_kibana_security_role.test", "elasticsearch.run_as", []string{"elastic", "kibana"}),
					resource.TestCheckNoResourceAttr("data.elasticstack_kibana_security_role.test", "kibana.0.feature.#"),
					util.TestCheckResourceListAttr("data.elasticstack_kibana_security_role.test", "kibana.0.base", []string{"all"}),
					util.TestCheckResourceListAttr("data.elasticstack_kibana_security_role.test", "kibana.0.spaces", []string{"default"}),
				),
			},
		},
	})
}

const testAccDataSourceSecurityRole = `
provider "elasticstack" {}

resource "elasticstack_kibana_security_role" "test" {
	name = "data_source_test"
	elasticsearch = {
		cluster = ["create_snapshot"]
		indices = [
			{
				field_security = {
					grant = ["test*"]
					exclude = ["test123"]
				}
				names = ["sample"]
				privileges = ["create", "read", "write"]
				query = jsonencode({
					match_all = {}
				})
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
}

data "elasticstack_kibana_security_role" "test" {
	name = elasticstack_kibana_security_role.test.name
}
`
