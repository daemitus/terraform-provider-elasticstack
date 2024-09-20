package fleet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceFleetServerHost(t *testing.T) {
	hostId := acctest.RandUUID()
	hostName := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceFleetServerHostDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFleetServerHostCreate(hostId, hostName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_server_host.test_host", "id", hostId),
					resource.TestCheckResourceAttr("elasticstack_fleet_server_host.test_host", "name", fmt.Sprintf("Created FleetServerHost %s", hostName)),
					resource.TestCheckResourceAttr("elasticstack_fleet_server_host.test_host", "default", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_server_host.test_host", "hosts.0", "https://fleet-server:8220"),
				),
			},
			{
				Config: testAccResourceFleetServerHostUpdate(hostId, hostName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_fleet_server_host.test_host", "id", hostId),
					resource.TestCheckResourceAttr("elasticstack_fleet_server_host.test_host", "name", fmt.Sprintf("Updated FleetServerHost %s", hostName)),
					resource.TestCheckResourceAttr("elasticstack_fleet_server_host.test_host", "default", "false"),
					resource.TestCheckResourceAttr("elasticstack_fleet_server_host.test_host", "hosts.0", "https://fleet-server:8220"),
				),
			},
		},
	})
}

func testAccResourceFleetServerHostCreate(hostId string, hostName string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_server_host" "test_host" {
	host_id = "%s"
	name = "Created FleetServerHost %s"
	default =	false
	hosts = ["https://fleet-server:8220"]
}`, hostId, hostName)
}

func testAccResourceFleetServerHostUpdate(hostId string, hostName string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_fleet_server_host" "test_host" {
	host_id = "%s"
	name = "Updated FleetServerHost %s"
	default =	false
	hosts = ["https://fleet-server:8220"]
}`, hostId, hostName)
}

func checkResourceFleetServerHostDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccFleetClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_fleet_server_host" {
			continue
		}

		host, diags := client.ReadFleetServerHost(ctx, rs.Primary.ID)
		if diags.HasError() {
			return fmt.Errorf(diags.Errors()[0].Summary())
		}
		if host != nil {
			return fmt.Errorf("fleet server host id=%v still exists, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
