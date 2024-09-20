package elasticsearch_test

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

func TestAccResourceComponentTemplate(t *testing.T) {
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceComponentTemplateDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceComponentTemplateCreate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_component_template.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_component_template.test", "template.aliases.%", "1"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_component_template.test", "template.settings", `{"index":{"number_of_shards":"3"}}`),
				),
			},
			{
				Config: testAccResourceComponentTemplateUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_component_template.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_component_template.test", "template.aliases.%", "2"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_component_template.test", "template.settings", `{"index":{"number_of_shards":"4"}}`),
				),
			},
		},
	})
}

func testAccResourceComponentTemplateCreate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_component_template" "test" {
	name = "%s"
	template = {
		aliases = {
			test1 = {}
		}
		settings = jsonencode({
			index = {
				number_of_shards = "3"
			}
		})
	}
}`, name)
}

func testAccResourceComponentTemplateUpdate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_component_template" "test" {
	name = "%s"
	template = {
		aliases = {
			test1 = {}
			test2 = {}
		}
		settings = jsonencode({
			index = {
				number_of_shards = "4"
			}
		})
	}
}`, name)
}

func checkResourceComponentTemplateDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccElasticsearchClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_elasticsearch_component_template" {
			continue
		}

		res, diags := client.GetComponentTemplate(ctx, rs.Primary.ID)
		if diags.HasError() {
			return util.DiagsAsError(diags)
		}

		if res != nil {
			return fmt.Errorf("index template id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
