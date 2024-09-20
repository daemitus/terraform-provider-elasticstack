package kibana_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceAlertingRule(t *testing.T) {
	name := acctest.RandString(22)

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceAlertingRuleDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertingRuleCreate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "id", "af22bd1c-8fb3-4020-9249-a4ac5471624b"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "name", fmt.Sprintf("Created %s", name)),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "consumer", "alerts"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "actions.0.frequency.notify_when", "onActiveAlert"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "rule_type_id", ".index-threshold"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "interval", "1m"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceAlertingRuleUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "id", "af22bd1c-8fb3-4020-9249-a4ac5471624b"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "name", fmt.Sprintf("Updated %s", name)),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "consumer", "alerts"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "actions.0.frequency.notify_when", "onActiveAlert"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "rule_type_id", ".index-threshold"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "interval", "10m"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "enabled", "false"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "tags.0", "first"),
					resource.TestCheckResourceAttr("elasticstack_kibana_alerting_rule.test_rule", "tags.1", "second"),
				),
			},
		},
	})
}

func testAccResourceAlertingRuleCreate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_kibana_action_connector" "log" {
	connector_type_id = ".server-log"
	name = "Kibana Log"
}

resource "elasticstack_kibana_alerting_rule" "test_rule" {
	name = "Created %s"
	id = "af22bd1c-8fb3-4020-9249-a4ac5471624b"
	consumer = "alerts"
	actions = [{
	group = "threshold met"
		frequency = {
		notify_when = "onActiveAlert"
	}
	id = elasticstack_kibana_action_connector.log.id
	params = jsonencode({
		level = "info"
		message = ""
		})
	}]
	params = jsonencode({
		aggType = "avg"
		groupBy = "top"
		termSize = 10
		timeWindowSize = 10
		timeWindowUnit = "s"
		threshold = [10]
		thresholdComparator = ">"
		index = ["test-index"]
		timeField = "@timestamp"
		aggField = "version"
		termField = "name"
	})
	rule_type_id = ".index-threshold"
	interval = "1m"
	enabled = true
}`, name)
}

func testAccResourceAlertingRuleUpdate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_kibana_action_connector" "log" {
	connector_type_id = ".server-log"
	name = "Kibana Log"
}

resource "elasticstack_kibana_alerting_rule" "test_rule" {
	name = "Updated %s"
	id = "af22bd1c-8fb3-4020-9249-a4ac5471624b"
	consumer = "alerts"
	actions = [{
	group = "threshold met"
		frequency = {
		notify_when = "onActiveAlert"
	}
	id = elasticstack_kibana_action_connector.log.id
	params = jsonencode({
		level = "info"
		message = ""
		})
	}]
	params = jsonencode({
		aggType = "avg"
		groupBy = "top"
		termSize = 10
		timeWindowSize = 10
		timeWindowUnit = "s"
		threshold = [10]
		thresholdComparator = ">"
		index = ["test-index"]
		timeField = "@timestamp"
		aggField = "version"
		termField = "name"
	})
	rule_type_id = ".index-threshold"
	interval = "10m"
	enabled = false
	tags = ["first", "second"]
}`, name)
}

func checkResourceAlertingRuleDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccKibanaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_kibana_alerting_rule" {
			continue
		}

		rule, diags := client.ReadAlertingRule(ctx, "default", rs.Primary.ID)
		if diags.HasError() {
			return fmt.Errorf(diags.Errors()[0].Summary())
		}

		if rule != nil {
			return fmt.Errorf("alerting rule id=%v, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
