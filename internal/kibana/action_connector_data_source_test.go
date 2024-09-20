package kibana_test

import (
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceActionConnector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceActionConnector,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.elasticstack_kibana_action_connector.test", "connector_id"),
					resource.TestCheckResourceAttr("data.elasticstack_kibana_action_connector.test", "name", "test_slack"),
					resource.TestCheckResourceAttr("data.elasticstack_kibana_action_connector.test", "space_id", "default"),
					resource.TestCheckResourceAttr("data.elasticstack_kibana_action_connector.test", "connector_type_id", ".slack"),
				),
			},
		},
	})
}

const testAccDataSourceActionConnector = `
provider "elasticstack" {}

resource "elasticstack_kibana_action_connector" "test" {
	name = "test_slack"
	connector_type_id = ".slack"
	secrets = jsonencode({
		webhookUrl = "https://internet.com"
	})
}

data "elasticstack_kibana_action_connector" "test" {
	name = elasticstack_kibana_action_connector.test.name
}`
