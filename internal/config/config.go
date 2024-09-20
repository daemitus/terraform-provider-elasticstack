package config

import (
	"github.com/caarlos0/env/v10"
)

// ProviderConfig is the configuration for the Terraform provider.
type ProviderConfig struct {
	Elasticsearch ServiceConfig `envPrefix:"ELASTICSEARCH_"`
	Kibana        ServiceConfig `envPrefix:"KIBANA_"`
	Fleet         ServiceConfig `envPrefix:"FLEET_"`
}

// ServiceConfig is the common configuration for each Elasticstack service.
type ServiceConfig struct {
	Endpoint string `env:"ENDPOINT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	ApiKey   string `env:"API_KEY"`
	Insecure bool   `env:"INSECURE"`
}

func New() *ProviderConfig {
	return &ProviderConfig{}
}

func (cfg *ProviderConfig) FromSchema(s ProviderModel) *ProviderConfig {
	cfg.Elasticsearch.FromSchema(s.Elasticsearch)
	cfg.Kibana.FromSchema(s.Kibana)
	cfg.Fleet.FromSchema(s.Fleet)

	return cfg
}

func (cfg *ProviderConfig) WithEnv() (*ProviderConfig, error) {
	err := env.Parse(cfg)
	return cfg, err
}

func (cfg *ProviderConfig) WithDefaults() *ProviderConfig {
	if cfg.Elasticsearch.Endpoint == "" {
		cfg.Elasticsearch.Endpoint = "http://localhost:9200"
	}
	if cfg.Kibana.Endpoint == "" {
		cfg.Kibana.Endpoint = "http://localhost:5601"
	}
	if cfg.Fleet.Endpoint == "" {
		cfg.Fleet.Endpoint = "http://localhost:5601"
	}
	return cfg
}

func (cfg *ServiceConfig) FromSchema(s *ServiceModel) {
	if s == nil {
		return
	}

	cfg.Endpoint = s.Endpoint.ValueString()
	cfg.Username = s.Username.ValueString()
	cfg.Password = s.Password.ValueString()
	cfg.ApiKey = s.ApiKey.ValueString()
	cfg.Insecure = s.Insecure.ValueBool()
}
