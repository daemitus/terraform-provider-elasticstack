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

func TestAccResourceIndexLifecycle(t *testing.T) {
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceIndexLifecycleDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIndexLifecycleCreate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.min_age", "1h"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.set_priority.priority", "10"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.rollover.max_age", "1d"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.readonly.%", "0"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "delete.min_age", "0ms"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "cold"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "frozen"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "delete.%", "3"),
				),
			},
			{
				Config: testAccResourceIndexLifecycleRemoveActions(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.min_age", "1h"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.set_priority.priority", "0"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.rollover.max_age", "2d"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.readonly.%", "0"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm.min_age", "0ms"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm.set_priority.priority", "60"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm.readonly.%", "0"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm.allocate.number_of_replicas", "1"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "cold"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "frozen"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "delete"),
				),
			},
			{
				Config: testAccResourceIndexLifecycleTotalShardsPerNode(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.min_age", "1h"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm.min_age", "0ms"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm.set_priority.priority", "60"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm.readonly.%", "0"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm.allocate.number_of_replicas", "1"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm.allocate.total_shards_per_node", "200"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "cold"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "frozen"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "delete"),
				),
			},
			{
				Config: testAccResourceIndexLifecycleDownsampleNoTimeout(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.min_age", "1h"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.set_priority.priority", "10"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.downsample.fixed_interval", "1d"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.rollover.max_age", "1d"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.readonly.%", "0"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "delete.min_age", "0ms"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "cold"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "frozen"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "delete.%", "3"),
				),
			},
			{
				Config: testAccResourceIndexLifecycleDownsample(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.min_age", "1h"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.set_priority.priority", "10"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.downsample.fixed_interval", "1d"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.downsample.wait_timeout", "1d"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.rollover.max_age", "1d"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "hot.readonly.%", "0"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "delete.min_age", "0ms"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "warm"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "cold"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "frozen"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test", "delete.%", "3"),
				),
			},
		},
	})
}

func TestAccResourceIndexLifecycleRolloverConditions(t *testing.T) {
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		CheckDestroy:             checkResourceIndexLifecycleDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIndexLifecycleCreateWithRolloverConditions(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "hot.rollover.max_age", "7d"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "hot.rollover.max_docs", "10000"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "hot.rollover.max_size", "100gb"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "hot.rollover.max_primary_shard_size", "50gb"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "hot.rollover.min_age", "3d"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "hot.rollover.min_docs", "1000"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "hot.rollover.min_size", "50gb"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "hot.rollover.min_primary_shard_size", "25gb"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index_lifecycle.test_rollover", "hot.rollover.min_primary_shard_docs", "500"),
				),
			},
		},
	})
}

func testAccResourceIndexLifecycleCreate(name string) string {
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
		delete = {}
	}
}
`, name)
}

func testAccResourceIndexLifecycleRemoveActions(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index_lifecycle" "test" {
	name = "%s"
	hot = {
		min_age = "1h"
		set_priority = {
			priority = 0
		}
		rollover = {
			max_age = "2d"
		}
	}
	warm = {
		min_age = "0ms"
		set_priority = {
			priority = 60
		}
		readonly = {}
		allocate = {
			exclude = jsonencode({
				box_type = "hot"
			})
			number_of_replicas = 1
		}
		shrink = {
			max_primary_shard_size = "50gb"
		}
	}
}
 `, name)
}

func testAccResourceIndexLifecycleTotalShardsPerNode(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index_lifecycle" "test" {
	name = "%s"
	hot = {
		min_age = "1h"
		set_priority = {
			priority = 0
		}
		rollover = {
			max_age = "2d"
		}
	}
	warm = {
		min_age = "0ms"
		set_priority = {
			priority = 60
		}
		readonly = {}
		allocate = {
			exclude = jsonencode({
				box_type = "hot"
			})
			number_of_replicas = 1
			total_shards_per_node = 200
		}
	}
}
`, name)
}

func testAccResourceIndexLifecycleCreateWithRolloverConditions(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index_lifecycle" "test_rollover" {
	name = "%s"
	hot = {
		rollover = {
			max_age = "7d"
			max_docs = 10000
			max_size = "100gb"
			max_primary_shard_size = "50gb"
			min_age = "3d"
			min_docs = 1000
			min_size = "50gb"
			min_primary_shard_size = "25gb"
			min_primary_shard_docs = 500
		}
		readonly = {}
	}
	delete = {
		delete = {}
	}
}
`, name)
}

func testAccResourceIndexLifecycleDownsampleNoTimeout(name string) string {
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
		downsample = {
			fixed_interval = "1d"
		}
		readonly = {}
	}
	delete = {
		delete = {}
	}
}
`, name)
}

func testAccResourceIndexLifecycleDownsample(name string) string {
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
		downsample = {
			fixed_interval = "1d"
			wait_timeout = "1d"
		}
		readonly = {}
	}
	delete = {
		delete = {}
	}
}
`, name)
}

func checkResourceIndexLifecycleDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccElasticsearchClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_elasticsearch_index_lifecycle" {
			continue
		}

		res, diags := client.GetIlmPolicy(ctx, rs.Primary.ID)
		if diags.HasError() {
			return util.DiagsAsError(diags)
		}

		if res != nil {
			return fmt.Errorf("ILM policy id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
