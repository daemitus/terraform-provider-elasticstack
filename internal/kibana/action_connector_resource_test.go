package kibana_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/acctest"
	"github.com/daemitus/terraform-provider-elasticstack/internal/clients"
	"github.com/daemitus/terraform-provider-elasticstack/internal/util"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/samber/lo"
)

func TestAccResourceConnectorCasesWebhook(t *testing.T) {
	connType := ".cases-webhook"
	name := acctest.RandString(22)

	createName := fmt.Sprintf("Created %s", name)
	createConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"createIncidentJson":                  "{}",
		"createIncidentMethod":                "post",
		"createIncidentResponseKey":           "key",
		"createIncidentUrl":                   "https://www.elastic.co/",
		"getIncidentResponseExternalTitleKey": "title",
		"getIncidentUrl":                      "https://www.elastic.co/",
		"hasAuth":                             false,
		"updateIncidentJson":                  "{}",
		"updateIncidentMethod":                "put",
		"updateIncidentUrl":                   "https://www.elastic.co/",
		"viewIncidentUrl":                     "https://www.elastic.co/",
	}))
	createSecrets := lo.Must(util.JsonMarshalS(map[string]any{
		"user":     "user1",
		"password": "password1",
	}))

	updateName := fmt.Sprintf("Updated %s", name)
	updateConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"createIncidentJson":                  "{}",
		"createIncidentMethod":                "put",
		"createIncidentResponseKey":           "key",
		"createIncidentUrl":                   "https://www.elastic.co/",
		"getIncidentResponseExternalTitleKey": "title",
		"getIncidentUrl":                      "https://www.elastic.co/",
		"hasAuth":                             true,
		"updateIncidentJson":                  "{}",
		"updateIncidentMethod":                "post",
		"updateIncidentUrl":                   "https://elasticsearch.com/",
		"viewIncidentUrl":                     "https://www.elastic.co/",
	}))
	updateSecrets := lo.Must(util.JsonMarshalS(map[string]any{
		"user":     "user2",
		"password": "password2",
	}))

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceConnectorDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorCreate(connType, createName, createConfig, createSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, createName, createConfig, createSecrets),
				),
			},
			{
				Config: testAccResourceConnectorUpdate(connType, updateName, updateConfig, updateSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, updateName, updateConfig, updateSecrets),
				),
			},
		},
	})
}

func TestAccResourceConnectorEmail(t *testing.T) {
	connType := ".email"
	name := acctest.RandString(22)

	createName := fmt.Sprintf("Created %s", name)
	createConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"from":    "test@elastic.co",
		"port":    111,
		"host":    "localhost",
		"hasAuth": false,
		"service": "other",
	}))
	createSecrets := lo.Must(util.JsonMarshalS(map[string]any{}))

	updateName := fmt.Sprintf("Updated %s", name)
	updateConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"from":    "test2@elastic.co",
		"port":    222,
		"host":    "localhost",
		"hasAuth": false,
		"service": "other",
	}))
	updateSecrets := lo.Must(util.JsonMarshalS(map[string]any{
		"user":     "user1",
		"password": "password1",
	}))

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceConnectorDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorCreate(connType, createName, createConfig, createSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, createName, createConfig, createSecrets),
				),
			},
			{
				Config: testAccResourceConnectorUpdate(connType, updateName, updateConfig, updateSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, updateName, updateConfig, updateSecrets),
				),
			},
		},
	})
}

func TestAccResourceConnectorIndex(t *testing.T) {
	connType := ".index"
	name := acctest.RandString(22)

	createName := fmt.Sprintf("Created %s", name)
	createConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"index":   ".kibana",
		"refresh": true,
	}))
	createSecrets := lo.Must(util.JsonMarshalS(map[string]any{}))

	updateName := fmt.Sprintf("Updated %s", name)
	updateConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"index":   ".kibana",
		"refresh": false,
	}))
	updateSecrets := lo.Must(util.JsonMarshalS(map[string]any{}))

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceConnectorDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorCreate(connType, createName, createConfig, createSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, createName, createConfig, createSecrets),
				),
			},
			{
				Config: testAccResourceConnectorUpdate(connType, updateName, updateConfig, updateSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, updateName, updateConfig, updateSecrets),
				),
			},
		},
	})
}

func TestAccResourceConnectorJira(t *testing.T) {
	connType := ".jira"
	name := acctest.RandString(22)

	createName := fmt.Sprintf("Created %s", name)
	createConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"apiUrl":     "url1",
		"projectKey": "project1",
	}))
	createSecrets := lo.Must(util.JsonMarshalS(map[string]any{
		"apiToken": "secret1",
		"email":    "email1",
	}))

	updateName := fmt.Sprintf("Updated %s", name)
	updateConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"apiUrl":     "url2",
		"projectKey": "project2",
	}))
	updateSecrets := lo.Must(util.JsonMarshalS(map[string]any{
		"apiToken": "secret2",
		"email":    "email2",
	}))

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceConnectorDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorCreate(connType, createName, createConfig, createSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, createName, createConfig, createSecrets),
				),
			},
			{
				Config: testAccResourceConnectorUpdate(connType, updateName, updateConfig, updateSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, updateName, updateConfig, updateSecrets),
				),
			},
		},
	})
}

func TestAccResourceConnectorServerLog(t *testing.T) {
	connType := ".server-log"
	name := acctest.RandString(22)

	createName := fmt.Sprintf("Created %s", name)
	createConfig := lo.Must(util.JsonMarshalS(map[string]any{}))
	createSecrets := lo.Must(util.JsonMarshalS(map[string]any{}))

	updateName := fmt.Sprintf("Updated %s", name)
	updateConfig := lo.Must(util.JsonMarshalS(map[string]any{}))
	updateSecrets := lo.Must(util.JsonMarshalS(map[string]any{}))

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceConnectorDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorCreate(connType, createName, createConfig, createSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, createName, createConfig, createSecrets),
				),
			},
			{
				Config: testAccResourceConnectorUpdate(connType, updateName, updateConfig, updateSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, updateName, updateConfig, updateSecrets),
				),
			},
		},
	})
}

func TestAccResourceConnectorSlack(t *testing.T) {
	connType := ".slack"
	name := acctest.RandString(22)

	createName := fmt.Sprintf("Created %s", name)
	createConfig := lo.Must(util.JsonMarshalS(map[string]any{}))
	createSecrets := lo.Must(util.JsonMarshalS(map[string]any{
		"webhookUrl": "https://elastic.co",
	}))

	updateName := fmt.Sprintf("Updated %s", name)
	updateConfig := lo.Must(util.JsonMarshalS(map[string]any{}))
	updateSecrets := lo.Must(util.JsonMarshalS(map[string]any{
		"webhookUrl": "https://elasticsearch.com",
	}))

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceConnectorDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorCreate(connType, createName, createConfig, createSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, createName, createConfig, createSecrets),
				),
			},
			{
				Config: testAccResourceConnectorUpdate(connType, updateName, updateConfig, updateSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, updateName, updateConfig, updateSecrets),
				),
			},
		},
	})
}

func TestAccResourceConnectorSlackApi(t *testing.T) {
	connType := ".slack_api"
	name := acctest.RandString(22)

	createName := fmt.Sprintf("Created %s", name)
	createConfig := lo.Must(util.JsonMarshalS(map[string]any{}))
	createSecrets := lo.Must(util.JsonMarshalS(map[string]any{
		"token": "my-token",
	}))

	updateName := fmt.Sprintf("Updated %s", name)
	updateConfig := lo.Must(util.JsonMarshalS(map[string]any{}))
	updateSecrets := lo.Must(util.JsonMarshalS(map[string]any{
		"token": "my-updated-token",
	}))

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceConnectorDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorCreate(connType, createName, createConfig, createSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, createName, createConfig, createSecrets),
				),
			},
			{
				Config: testAccResourceConnectorUpdate(connType, updateName, updateConfig, updateSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, updateName, updateConfig, updateSecrets),
				),
			},
		},
	})
}

func TestAccResourceConnectorWebhook(t *testing.T) {
	connType := ".webhook"
	name := acctest.RandString(22)

	createName := fmt.Sprintf("Created %s", name)
	createConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"url":     "https://elastic.co",
		"hasAuth": true,
		"method":  "post",
	}))
	createSecrets := lo.Must(util.JsonMarshalS(map[string]any{}))

	updateName := fmt.Sprintf("Updated %s", name)
	updateConfig := lo.Must(util.JsonMarshalS(map[string]any{
		"url":     "https://elasticsearch.com",
		"hasAuth": true,
		"method":  "post",
	}))
	updateSecrets := lo.Must(util.JsonMarshalS(map[string]any{}))

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		CheckDestroy:             checkResourceConnectorDestroy,
		ProtoV6ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnectorCreate(connType, createName, createConfig, createSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, createName, createConfig, createSecrets),
				),
			},
			{
				Config: testAccResourceConnectorUpdate(connType, updateName, updateConfig, updateSecrets),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceConnectorCommon(connType, updateName, updateConfig, updateSecrets),
				),
			},
		},
	})
}

func testAccResourceConnectorCreate(connType string, name string, config string, secrets string) string {
	config = lo.Must(util.JsonMarshalS(config))
	secrets = lo.Must(util.JsonMarshalS(secrets))
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_kibana_action_connector" "test" {
	connector_type_id = "%s"
	name = "%s"
	config = %s
	secrets = %s
}`, connType, name, config, secrets)
}

func testAccResourceConnectorUpdate(connType string, name string, config string, secrets string) string {
	config = lo.Must(util.JsonMarshalS(config))
	secrets = lo.Must(util.JsonMarshalS(secrets))
	return fmt.Sprintf(`
provider "elasticstack" {}

resource "elasticstack_kibana_action_connector" "test" {
	connector_type_id = "%s"
	name = "%s"
	config = %s
	secrets = %s
}`, connType, name, config, secrets)
}

func testAccResourceConnectorCommon(connType string, name string, config string, secrets string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("elasticstack_kibana_action_connector.test", "name", name),
		resource.TestCheckResourceAttr("elasticstack_kibana_action_connector.test", "connector_type_id", connType),
		resource.TestCheckResourceAttr("elasticstack_kibana_action_connector.test", "is_deprecated", "false"),
		resource.TestCheckResourceAttr("elasticstack_kibana_action_connector.test", "is_missing_secrets", "false"),
		resource.TestCheckResourceAttr("elasticstack_kibana_action_connector.test", "is_preconfigured", "false"),
		resource.TestCheckResourceAttr("elasticstack_kibana_action_connector.test", "config", config),
		resource.TestCheckResourceAttr("elasticstack_kibana_action_connector.test", "secrets", secrets),
	)
}

func checkResourceConnectorDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := clients.NewAccKibanaClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_kibana_action_connector" {
			continue
		}

		space := rs.Primary.Attributes["space_id"]
		res, diags := client.ReadConnector(ctx, space, rs.Primary.ID)
		if diags.HasError() {
			return fmt.Errorf(diags.Errors()[0].Summary())
		}

		if res != nil {
			return fmt.Errorf("connector id=%v, but it should have been removed", rs.Primary.ID)
		}
	}
	return nil
}
