package config_test

import (
	"os"
	"testing"

	"github.com/daemitus/terraform-provider-elasticstack/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestNewFromSchema(t *testing.T) {
	schema := config.ProviderModel{
		Elasticsearch: &config.ServiceModel{
			Endpoint: types.StringValue("endpoint"),
			Username: types.StringValue("username"),
			Password: types.StringValue("password"),
			ApiKey:   types.StringValue("apikey"),
			Insecure: types.BoolValue(true),
		},
		Kibana: nil,
		Fleet:  nil,
	}

	actual := config.New().FromSchema(schema)
	expected := &config.ProviderConfig{
		Elasticsearch: config.ServiceConfig{
			Endpoint: "endpoint",
			Username: "username",
			Password: "password",
			ApiKey:   "apikey",
			Insecure: true,
		},
		Kibana: config.ServiceConfig{},
		Fleet:  config.ServiceConfig{},
	}

	assert.Equal(t, expected, actual)
}

func TestNewFromSchemaWithEnv(t *testing.T) {
	os.Setenv("ELASTICSEARCH_ENDPOINT", "new-endpoint")
	os.Setenv("ELASTICSEARCH_USERNAME", "new-username")
	os.Setenv("ELASTICSEARCH_PASSWORD", "")
	os.Setenv("ELASTICSEARCH_API_KEY", "")
	os.Setenv("ELASTICSEARCH_INSECURE", "false")
	os.Unsetenv("KIBANA_ENDPOINT")
	os.Unsetenv("KIBANA_USERNAME")
	os.Unsetenv("KIBANA_PASSWORD")
	os.Unsetenv("KIBANA_API_KEY")
	os.Unsetenv("KIBANA_INSECURE")
	os.Unsetenv("FLEET_ENDPOINT")
	os.Unsetenv("FLEET_USERNAME")
	os.Unsetenv("FLEET_PASSWORD")
	os.Unsetenv("FLEET_API_KEY")
	os.Unsetenv("FLEET_INSECURE")

	schema := config.ProviderModel{
		Elasticsearch: &config.ServiceModel{
			Endpoint: types.StringValue("endpoint"),
			Username: types.StringValue("username"),
			Password: types.StringValue("password"),
			ApiKey:   types.StringValue("apikey"),
			Insecure: types.BoolValue(true),
		},
		Kibana: nil,
		Fleet:  nil,
	}

	actual, err := config.New().FromSchema(schema).WithEnv()
	if err != nil {
		t.Error(err.Error())
	}
	expected := &config.ProviderConfig{
		Elasticsearch: config.ServiceConfig{
			Endpoint: "new-endpoint",
			Username: "new-username",
			Password: "password",
			ApiKey:   "apikey",
			Insecure: false,
		},
		Kibana: config.ServiceConfig{},
		Fleet:  config.ServiceConfig{},
	}

	assert.Equal(t, expected, actual)
}
