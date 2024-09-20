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

func TestResourceWatch(t *testing.T) {
	id := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceWatchDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWatchCreate(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "id", id),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "active", "false"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "trigger", `{"schedule":{"cron":"0 0/1 * * * ?"}}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "input", `{"none":{}}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "condition", `{"always":{}}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "actions", `{}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "metadata", `{}`),
				),
			},
			{
				Config: testAccResourceWatchUpdate(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "id", id),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "active", "true"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "trigger", `{"schedule":{"cron":"0 0/2 * * * ?"}}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "input", `{"simple":{"name":"example"}}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "condition", `{"never":{}}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "actions", `{"log":{"logging":{"level":"info","text":"example logging text"}}}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "metadata", `{"example_key":"example_value"}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "transform", `{"search":{"request":{"body":{"query":{"match_all":{}}},"indices":[],"rest_total_hits_as_int":true,"search_type":"query_then_fetch"}}}`),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_watch.test", "throttle_period_in_millis", "10000"),
				),
			},
			{
				Config:            testAccResourceWatchUpdate(id),
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "elasticstack_elasticsearch_watch.test",
			},
		},
	})
}

func testAccResourceWatchCreate(watchID string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_watch" "test" {
	id = "%s"
	active = false
	trigger = jsonencode({
		schedule = {
			cron = "0 0/1 * * * ?"
		}
	})

}`, watchID)
}

func testAccResourceWatchUpdate(watchID string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_watch" "test" {
	id = "%s"
	active = true
	trigger = jsonencode({
		schedule = {
			cron = "0 0/2 * * * ?"
		}
	})
	input = jsonencode({
		simple = {
			name = "example"
		}
	})
	condition = jsonencode({
		never = {}
	})
	actions = jsonencode({
		log = {
			logging = {
				level = "info"
				text = "example logging text"
			}
		}
	})
	metadata = jsonencode({
		example_key = "example_value"
	})
	transform = jsonencode({
		search = {
			request = {
				body = {
					query = {
						match_all = {}
					}
				}
				indices = []
				rest_total_hits_as_int = true
				search_type = "query_then_fetch"
			}
		}
	})
	throttle_period_in_millis = 10000
}`, watchID)
}

func checkResourceWatchDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccElasticsearchClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_elasticsearch_watch" {
			continue
		}

		res, diags := client.GetWatch(ctx, rs.Primary.ID)
		if diags.HasError() {
			return util.DiagsAsError(diags)
		}

		if res != nil {
			return fmt.Errorf("watch id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
