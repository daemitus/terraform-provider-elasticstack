package kibana_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceDataView(t *testing.T) {
	indexName := fmt.Sprintf("test-%s", acctest.RandString(4))

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceDataViewDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDataViewCreate(indexName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("elasticstack_kibana_data_view.test", "id"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "override", "false"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.name", indexName),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.source_filters.#", "2"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.field_formats.event_time.id", "date_nanos"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.field_formats.machine.ram.params.pattern", "0,0.[000] b"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.runtime_field_map.runtime_shape_name.script_source", "emit(doc['shape_name'].value)"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.field_attrs.ingest_failure.custom_label", "error.ingest_failure"),
				),
			},
			{
				Config: testAccResourceDataViewUpdate(indexName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("elasticstack_kibana_data_view.test", "id"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "override", "true"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.name", indexName),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.source_filters.#", "0"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.field_formats.#", "0"),
					resource.TestCheckResourceAttr("elasticstack_kibana_data_view.test", "data_view.runtime_field_map.#", "0"),
				),
			},
			{
				Config:                  testAccResourceDataViewUpdate(indexName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_id", "override"},
				ResourceName:            "elasticstack_kibana_data_view.test",
			},
		},
	})
}

func testAccResourceDataViewCreate(indexName string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index" "my_index" {
	name = "%s"
	deletion_protection = false
}

resource "elasticstack_kibana_data_view" "test" {
	data_view = {
		name = "%s"
		title = "%s*"
		time_field_name = "@timestamp"
		source_filters = ["event_time", "machine.ram"]
		allow_no_index = true
		namespaces = ["bar", "default", "foo"]
		field_formats = {
			event_time = {
				id = "date_nanos"
			}
			"machine.ram" = {
				id = "number"
				params = {
					pattern = "0,0.[000] b"
				}
			}
		}
		runtime_field_map = {
			runtime_shape_name = {
				type = "keyword"
				script_source = "emit(doc['shape_name'].value)"
			}
		}
		field_attrs = {
			ingest_failure = { 
				custom_label = "error.ingest_failure"
				count = 6
			}
		}
	}
	override = false
}`, indexName, indexName, indexName)
}

func testAccResourceDataViewUpdate(indexName string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_elasticsearch_index" "my_index" {
	name = "%s"
	deletion_protection = false
}

resource "elasticstack_kibana_data_view" "test" {
	data_view = {
		name = "%s"
		title = "%s*"
		time_field_name = "@timestamp"
		allow_no_index = true
	}
	override = true
}`, indexName, indexName, indexName)
}

func checkResourceDataViewDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccKibanaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_kibana_data_view" {
			continue
		}

		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		dataView, diags := client.ReadDataView(ctx, parts[0], parts[1])
		if diags.HasError() {
			return fmt.Errorf(diags.Errors()[0].Summary())
		}
		if dataView != nil {
			return fmt.Errorf("data view id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
