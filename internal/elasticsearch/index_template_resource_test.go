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

func TestAccResourceIndexTemplate(t *testing.T) {
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceIndexTemplateDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIndexTemplateCreate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test1", "name", fmt.Sprintf("%s-indexes", name)),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test1", "index_patterns.0", fmt.Sprintf("%s-index1-*", name)),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test1", "priority", "42"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test1", "template.aliases.%", "1"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test2", "name", fmt.Sprintf("%s-streams", name)),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test2", "index_patterns.0", fmt.Sprintf("%s-stream1-*", name)),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test2", "data_stream.hidden", "true"),
				),
			},
			{
				Config: testAccResourceIndexTemplateUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test1", "name", fmt.Sprintf("%s-indexes", name)),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test1", "index_patterns.0", fmt.Sprintf("%s-index2-*", name)),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test1", "template.aliases.%", "2"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test2", "name", fmt.Sprintf("%s-streams", name)),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test2", "index_patterns.0", fmt.Sprintf("%s-stream2-*", name)),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_template.test2", "data_stream.hidden", "false"),
				),
			},
		},
	})
}

func testAccResourceIndexTemplateCreate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index_template" "test1" {
	name = "%s-indexes"
	priority = 42
	index_patterns = ["%s-index1-*"]
	template = {
		aliases = {
			my_template_test = {}
		}
		settings = jsonencode({
			index = {
				number_of_shards = "3"
			}
		})
	}
}

resource "elasticstack_elasticsearch_index_template" "test2" {
	name = "%s-streams"
	index_patterns = ["%s-stream1-*"]
	data_stream = {
		hidden = true
	}
}`, name, name, name, name)
}

func testAccResourceIndexTemplateUpdate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index_template" "test1" {
	name = "%s-indexes"
	index_patterns = ["%s-index2-*"]
	template = {
		aliases = {
			"my_template_test" = {}
			"alias2" = {}
		}
		settings = jsonencode({
			index = {
				number_of_shards = "3"
			}
		})
	}
}

resource "elasticstack_elasticsearch_index_template" "test2" {
	name = "%s-streams"
	index_patterns = ["%s-stream2-*"]
	data_stream = {
		hidden = false
	}
	template = {}
}`, name, name, name, name)
}

func checkResourceIndexTemplateDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccElasticsearchClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_elasticsearch_index_template" {
			continue
		}

		res, diags := client.GetIndexTemplate(ctx, rs.Primary.ID)
		if diags.HasError() {
			return util.DiagsAsError(diags)
		}

		if res != nil {
			return fmt.Errorf("index template id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
