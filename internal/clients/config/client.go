package config

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/fleet"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients/kibana"
)

type Client struct {
	UserAgent     string
	Elasticsearch *elasticsearch.Config
	Kibana        *kibana.Config
	Fleet         *fleet.Config
}
