package kibana_test

import (
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSpaces(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaces,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.elasticstack_kibana_spaces.test", "spaces.0.id"),
					resource.TestCheckResourceAttrSet("data.elasticstack_kibana_spaces.test", "spaces.0.name"),
				),
			},
		},
	})
}

var testAccDataSourceSpaces = `
provider "elasticstack" {}

data "elasticstack_kibana_spaces" "test" {}`
