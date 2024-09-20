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

func TestAccResourceDataStream(t *testing.T) {
	name := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceDataStreamDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDataStreamCreate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_data_stream.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_data_stream.test", "indices.#", "1"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_data_stream.test", "template", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_data_stream.test", "ilm_policy", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_data_stream.test", "hidden", "false"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_data_stream.test", "system", "false"),
				),
			},
		},
	})
}

func testAccResourceDataStreamCreate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index_lifecycle" "test" {
	name = "%s"
	hot = {
    	min_age = "1h"
    	set_priority = {
      		priority = 10
    	}
    	rollover = {
      		max_age = "1d"
		}
		readonly = {}
  	}
	delete = {
		min_age = "2d"
		delete = {}
	}
}

resource "elasticstack_elasticsearch_index_template" "test" {
  name = "%s"
  index_patterns = ["%s*"]
  template = {
    settings = jsonencode({
    	index = {
	  		lifecycle = {
				name = elasticstack_elasticsearch_index_lifecycle.test.name
			}
		}
    })
  }
  data_stream = {}
}

resource "elasticstack_elasticsearch_data_stream" "test" {
	name = "%s"
	depends_on = [elasticstack_elasticsearch_index_template.test]
}
`, name, name, name, name)
}

func checkResourceDataStreamDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccElasticsearchClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_elasticsearch_data_stream" {
			continue
		}

		res, diags := client.GetDataStream(ctx, rs.Primary.ID)
		if diags.HasError() {
			return util.DiagsAsError(diags)
		}

		if res != nil {
			return fmt.Errorf("data stream id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
