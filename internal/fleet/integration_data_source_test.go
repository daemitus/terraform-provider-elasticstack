package fleet_test

import (
	"errors"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceIntegration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIntegration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.elasticstack_fleet_integration.test", "name", "tcp"),
					resource.TestCheckResourceAttrWith("data.elasticstack_fleet_integration.test", "version", func(value string) error {
						if value == "" {
							return errors.New("attribute was empty when expected")
						}
						return nil
					}),
				),
			},
		},
	})
}

const testAccDataSourceIntegration = `
provider "elasticstack" {}

data "elasticstack_fleet_integration" "test" {
	name = "tcp"
}
`
