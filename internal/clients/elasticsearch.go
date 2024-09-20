package clients

import (
	"fmt"
	"log"

	"github.com/daemitus/terraform-provider-elasticstack/internal/api/elasticsearch"
	"github.com/daemitus/terraform-provider-elasticstack/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ElasticsearchClient provides an API client for Elastic Fleet.
type ElasticsearchClient struct {
	baseClient
	API *elasticsearch.Client
}

// NewElasticsearchClient creates a new Elasticsearch API client.
func NewElasticsearchClient(config config.ServiceConfig) (*ElasticsearchClient, error) {
	api, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}
	client := ElasticsearchClient{API: api}
	return &client, nil
}

func NewAccElasticsearchClient() *ElasticsearchClient {
	config, err := config.New().WithEnv()
	if err != nil {
		log.Fatal("Elasticsearch config failure: %w", err)
	}

	client, err := NewElasticsearchClient(config.Elasticsearch)
	if err != nil {
		log.Fatal("Elasticsearch client failure: %w", err)
	}

	return client
}

func getOneSliceResponse[T any](items []T, name string, singular string, plural string) (*T, diag.Diagnostics) {
	var diags diag.Diagnostics
	if count := len(items); count > 1 {
		diags.AddError(
			fmt.Sprintf("Multiple %s returned", plural),
			fmt.Sprintf("Elasticsearch API returned %d when requested %s '%s'.", count, singular, name),
		)
		return nil, diags
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &items[0], nil
}

func getOneMapResponse[T any](items map[string]T, name string, singular string, plural string) (*T, diag.Diagnostics) {
	var diags diag.Diagnostics
	if count := len(items); count > 1 {
		diags.AddError(
			fmt.Sprintf("Multiple %s returned", plural),
			fmt.Sprintf("Elasticsearch API returned %d when requested %s '%s'.", count, singular, name),
		)
		return nil, diags
	}
	if len(items) == 0 {
		return nil, nil
	}
	item := items[name]
	return &item, nil
}
