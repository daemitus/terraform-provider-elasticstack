package config

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderModel struct {
	Elasticsearch *ServiceModel `tfsdk:"elasticsearch"`
	Kibana        *ServiceModel `tfsdk:"kibana"`
	Fleet         *ServiceModel `tfsdk:"fleet"`
}

type ServiceModel struct {
	ApiKey   types.String `tfsdk:"api_key"`
	Endpoint types.String `tfsdk:"endpoint"`
	Insecure types.Bool   `tfsdk:"insecure"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}
