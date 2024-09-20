package fleet

import (
	"context"
	"encoding/json"
)

func (c *Client) ReadOutput(ctx context.Context, outputId string) (*ApiResponse[ReadOutputResponse], error) {
	return doAPI[ReadOutputResponse](
		c, ctx,
		"GET", "/outputs/{id}",
		map[string]string{"id": outputId},
		nil, nil,
	)
}

type ReadOutputResponse struct {
	Item *OutputUnion `json:"item,omitempty"`
}

// ============================================================================

func (c *Client) CreateOutput(ctx context.Context, body CreateOutputRequest) (*ApiResponse[CreateOutputResponse], error) {
	return doAPI[CreateOutputResponse](
		c, ctx,
		"POST", "/outputs",
		nil, body.raw, nil,
	)
}

type CreateOutputRequest OutputUnion

type CreateOutputResponse struct {
	Item *OutputUnion `json:"item,omitempty"`
}

// ============================================================================

func (c *Client) UpdateOutput(ctx context.Context, outputId string, body UpdateOutputRequest) (*ApiResponse[UpdateOutputResponse], error) {
	return doAPI[UpdateOutputResponse](
		c, ctx,
		"PUT", "/outputs/{id}",
		map[string]string{"id": outputId},
		body.raw, nil,
	)
}

type UpdateOutputRequest OutputUnion

type UpdateOutputResponse struct {
	Item *OutputUnion `json:"item,omitempty"`
}

// ============================================================================

func (c *Client) DeleteOutput(ctx context.Context, outputId string) (*ApiResponse[DeleteOutputResponse], error) {
	return doAPI[DeleteOutputResponse](
		c, ctx,
		"DELETE", "/outputs/{id}",
		map[string]string{"id": outputId},
		nil, nil,
	)
}

type DeleteOutputResponse struct {
	Id string `json:"id"`
}

// ============================================================================

type OutputUnion struct {
	raw json.RawMessage
}

func (t OutputUnion) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"type"`
	}
	err := json.Unmarshal(t.raw, &discriminator)
	return discriminator.Discriminator, err
}

func (t OutputUnion) AsOutputElasticsearch() (OutputElasticsearch, error) {
	var body OutputElasticsearch
	err := json.Unmarshal(t.raw, &body)
	return body, err
}

func (t OutputUnion) AsOutputRemoteElasticsearch() (OutputRemoteElasticsearch, error) {
	var body OutputRemoteElasticsearch
	err := json.Unmarshal(t.raw, &body)
	return body, err
}

func (t *OutputUnion) FromOutputElasticsearch(v OutputElasticsearch) error {
	v.Type = OutputTypeElasticsearch
	b, err := json.Marshal(v)
	t.raw = b
	return err
}

func (t *OutputUnion) FromOutputRemoteElasticsearch(v OutputElasticsearch) error {
	v.Type = OutputTypeRemoteElasticsearch
	b, err := json.Marshal(v)
	t.raw = b
	return err
}

func (t OutputUnion) MarshalJSON() ([]byte, error) {
	return t.raw.MarshalJSON()
}

func (t *OutputUnion) UnmarshalJSON(b []byte) error {
	return t.raw.UnmarshalJSON(b)
}

// ============================================================================

type OutputElasticsearch struct {
	CaSha256             *string                     `json:"ca_sha256,omitempty"`
	CaTrustedFingerprint *string                     `json:"ca_trusted_fingerprint,omitempty"`
	Config               *map[string]any             `json:"config,omitempty"`
	ConfigYaml           *string                     `json:"config_yaml,omitempty"`
	Hosts                []string                    `json:"hosts"`
	Id                   *string                     `json:"id,omitempty"`
	IsDefault            *bool                       `json:"is_default,omitempty"`
	IsDefaultMonitoring  *bool                       `json:"is_default_monitoring,omitempty"`
	Name                 string                      `json:"name"`
	Preset               *OutputPreset               `json:"preset,omitempty"`
	ProxyId              *string                     `json:"proxy_id,omitempty"`
	Shipper              *OutputElasticsearchShipper `json:"shipper,omitempty"`
	Ssl                  *OutputElasticsearchSsl     `json:"ssl,omitempty"`
	Type                 OutputType                  `json:"type"`
}

func (t OutputElasticsearch) AsUnion() (*OutputUnion, error) {
	bytes, err := json.Marshal(t)
	return &OutputUnion{bytes}, err
}

type OutputElasticsearchShipper struct {
	CompressionLevel            *float32 `json:"compression_level,omitempty"`
	DiskQueueCompressionEnabled *bool    `json:"disk_queue_compression_enabled,omitempty"`
	DiskQueueEnabled            *bool    `json:"disk_queue_enabled,omitempty"`
	DiskQueueEncryptionEnabled  *bool    `json:"disk_queue_encryption_enabled,omitempty"`
	DiskQueueMaxSize            *float32 `json:"disk_queue_max_size,omitempty"`
	DiskQueuePath               *string  `json:"disk_queue_path,omitempty"`
	Loadbalance                 *bool    `json:"loadbalance,omitempty"`
}

type OutputElasticsearchSsl struct {
	Certificate            *string   `json:"certificate,omitempty"`
	CertificateAuthorities *[]string `json:"certificate_authorities,omitempty"`
	Key                    *string   `json:"key,omitempty"`
}

// ============================================================================

type OutputRemoteElasticsearch struct {
	ConfigYaml          *string                           `json:"config_yaml,omitempty"`
	Hosts               []string                          `json:"hosts"`
	Id                  *string                           `json:"id,omitempty"`
	IsDefault           *bool                             `json:"is_default,omitempty"`
	IsDefaultMonitoring *bool                             `json:"is_default_monitoring,omitempty"`
	Name                string                            `json:"name"`
	Preset              *OutputPreset                     `json:"preset,omitempty"`
	ServiceToken        *string                           `json:"service_token,omitempty"`
	Secrets             *OutputRemoteElasticsearchSecrets `json:"secrets,omitempty"`
	Type                OutputType                        `json:"type"`
}

func (t OutputRemoteElasticsearch) AsUnion() (*OutputUnion, error) {
	bytes, err := json.Marshal(t)
	return &OutputUnion{bytes}, err
}

type OutputRemoteElasticsearchSecrets struct {
	ServiceToken OutputRemoteElasticsearchSecretUnion `json:"service_token"`
}

type OutputRemoteElasticsearchSecretUnion struct {
	raw json.RawMessage
}

func (t OutputRemoteElasticsearchSecretUnion) AsString() (OutputRemoteElasticsearchSecretString, error) {
	var body OutputRemoteElasticsearchSecretString
	err := json.Unmarshal(t.raw, &body)
	return body, err
}

func (t OutputRemoteElasticsearchSecretUnion) AsSecret() (OutputRemoteElasticsearchSecretId, error) {
	var body OutputRemoteElasticsearchSecretId
	err := json.Unmarshal(t.raw, &body)
	return body, err
}

func (t *OutputRemoteElasticsearchSecretUnion) FromString(v OutputRemoteElasticsearchSecretString) error {
	b, err := json.Marshal(v)
	t.raw = b
	return err
}

func (t *OutputRemoteElasticsearchSecretUnion) FromSecret(v OutputRemoteElasticsearchSecretId) error {
	b, err := json.Marshal(v)
	t.raw = b
	return err
}

func (t OutputRemoteElasticsearchSecretUnion) MarshalJSON() ([]byte, error) {
	return t.raw.MarshalJSON()
}

func (t *OutputRemoteElasticsearchSecretUnion) UnmarshalJSON(b []byte) error {
	return t.raw.UnmarshalJSON(b)
}

type OutputRemoteElasticsearchSecretString string

func (t OutputRemoteElasticsearchSecretString) AsUnion() (*OutputRemoteElasticsearchSecretUnion, error) {
	bytes, err := json.Marshal(t)
	return &OutputRemoteElasticsearchSecretUnion{bytes}, err
}

type OutputRemoteElasticsearchSecretId struct {
	Id string `json:"id"`
}

func (t OutputRemoteElasticsearchSecretId) AsUnion() (*OutputRemoteElasticsearchSecretUnion, error) {
	bytes, err := json.Marshal(t)
	return &OutputRemoteElasticsearchSecretUnion{bytes}, err
}

// ============================================================================

type OutputType string

const (
	OutputTypeElasticsearch       OutputType = "elasticsearch"
	OutputTypeRemoteElasticsearch OutputType = "remote_elasticsearch"
)

type OutputPreset string

const (
	OutputPresetBalanced   OutputType = "balanced"
	OutputPresetCustom     OutputType = "custom"
	OutputPresetThroughput OutputType = "throughput"
	OutputPresetScale      OutputType = "scale"
	OutputPresetLatency    OutputType = "latency"
)
