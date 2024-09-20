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

func TestAccResourceIngestPipeline(t *testing.T) {
	name := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceIngestPipelineDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIngestPipelineCreate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("elasticstack_elasticsearch_ingest_pipeline.test", "id"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_ingest_pipeline.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_ingest_pipeline.test", "description", "Test Pipeline"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_ingest_pipeline.test", "processors.#", "2"),
				),
			},
			{
				Config: testAccResourceIngestPipelineUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("elasticstack_elasticsearch_ingest_pipeline.test", "id"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_ingest_pipeline.test", "name", name),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_ingest_pipeline.test", "description", "Test Pipeline"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_ingest_pipeline.test", "processors.#", "1"),
				),
			},
			{
				Config:            testAccResourceIngestPipelineUpdate(name),
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "elasticstack_elasticsearch_ingest_pipeline.test",
			},
		},
	})
}

func testAccResourceIngestPipelineCreate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_ingest_pipeline" "test" {
	name = "%s"
	description = "Test Pipeline"
	processors = [
		jsonencode({
			set = {
				description = "My set processor description"
				field = "_meta"
				value = "indexed"
			}
		}),
		jsonencode({
			json = {
				field = "data"
				target_field = "parsed_data"
			}
		}),
	]
}`, name)
}

func testAccResourceIngestPipelineUpdate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_ingest_pipeline" "test" {
	name = "%s"
	description = "Test Pipeline"
	processors = [
		jsonencode({
			set = {
				description = "My set processor description"
				field = "_meta"
				value = "indexed"
			}
		})
	]
}`, name)
}

func checkResourceIngestPipelineDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccElasticsearchClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_elasticsearch_ingest_pipeline" {
			continue
		}

		res, diags := client.GetIngestPipeline(ctx, rs.Primary.ID)
		if diags.HasError() {
			return util.DiagsAsError(diags)
		}

		if res != nil {
			return fmt.Errorf("Ingest pipeline id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
