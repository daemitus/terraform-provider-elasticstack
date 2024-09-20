package elasticsearch_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceIndexAliases(t *testing.T) {
	indexName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceIndexDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIndexAliasesCreate(indexName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test", "name", indexName),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test", "aliases.%", "2"),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test", "aliases.test_alias_2.filter", `{"term":{"user.id":{"value":"developer"}}}`),
				),
			},
			{
				Config: testAccResourceIndexAliasesUpdate(indexName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test", "name", indexName),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test", "aliases.%", "1"),
					resource.TestCheckNoResourceAttr("elasticstack_elasticsearch_index.test", "aliases.test_alias_2.filter"),
				),
			},
		},
	})
}

func TestAccResourceIndexSettings(t *testing.T) {
	indexName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceIndexDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIndexSettingsCreate(indexName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_settings", "name", indexName),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_settings", "settings", `{"index":{"codec":"best_compression","number_of_replicas":"2"}}`),
				),
			},
			{
				Config: testAccResourceIndexSettingsUpdate(indexName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_settings", "name", indexName),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_settings", "settings", `{"index":{"number_of_replicas":"2","search":{"idle":{"after":"20s"}}}}`),
				),
			},
		},
	})
}

func TestAccResourceIndexMappings(t *testing.T) {
	indexName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceIndexDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIndexMappingsCreate(indexName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_mappings", "name", indexName),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_mappings", "mappings", `{"properties":{"field1":{"type":"text"},"field2":{"type":"text"}}}`),
				),
			},
			{
				Config: testAccResourceIndexMappingsUpdate(indexName, true),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccResourceIndexMappingsPostUpdate(indexName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_mappings", "name", indexName),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_mappings", "mappings", `{"properties":{"field1":{"type":"text"}}}`),
				),
				ExpectError: regexp.MustCompile("cannot destroy index without setting deletion_protection=false and running `terraform apply`"),
			},
			{
				Config: testAccResourceIndexMappingsUpdate(indexName, false),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccResourceIndexMappingsPostUpdate(indexName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_mappings", "name", indexName),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test_mappings", "mappings", `{"properties":{"field1":{"type":"boolean"}}}`),
				),
			},
		},
	})
}

func testAccResourceIndexAliasesCreate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index" "test" {
	name = "%s"
	aliases = {
		test_alias_1 = {}
		test_alias_2 = {
			filter = jsonencode({
				term = { "user.id" = { value = "developer" } }
			})
		}
	}
	deletion_protection = false
}`, name)
}

func testAccResourceIndexAliasesUpdate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index" "test" {
	name = "%s"
	aliases = {
		test_alias_1 = {}
	}
	deletion_protection = false
}`, name)
}

func testAccResourceIndexSettingsCreate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index" "test_settings" {
	name = "%s"
	settings = jsonencode({
		index = {
		codec = "best_compression"
		number_of_replicas = "2"
		}
	})
	deletion_protection = false
}`, name)
}

func testAccResourceIndexSettingsUpdate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index" "test_settings" {
	name = "%s"
	settings = jsonencode({
		index = {
			number_of_replicas = "2"
			search = { idle = { after = "20s" } }
		}
	})
	deletion_protection = false
}`, name)
}

func testAccResourceIndexMappingsCreate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index" "test_mappings" {
	name = "%s"
	mappings = jsonencode({
		properties = {
			field1 = { type = "text" }
			field2 = { type = "text" }
		}
	})
	deletion_protection = true
}`, name)
}

func testAccResourceIndexMappingsUpdate(name string, preventDestroy bool) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index" "test_mappings" {
	name = "%s"
	mappings = jsonencode({
		properties = {
			field1 = { type = "text" }
			field3 = { type = "text" }
		}
	})
	deletion_protection = %t
}`, name, preventDestroy)
}

func testAccResourceIndexMappingsPostUpdate(name string, preventDestroy bool) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index" "test_mappings" {
	name = "%s"
	mappings = jsonencode({
		properties = {
			field1 = { type = "boolean" }
		}
	})
	deletion_protection = %t
}`, name, preventDestroy)
}

func checkResourceIndexDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccElasticsearchClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_elasticsearch_index" {
			continue
		}

		res, diags := client.GetIndex(ctx, rs.Primary.ID)
		if diags.HasError() {
			return util.DiagsAsError(diags)
		}

		if res != nil {
			return fmt.Errorf("index id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
