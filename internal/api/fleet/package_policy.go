package fleet

import (
	"context"
)

func (c *Client) ReadPackagePolicy(ctx context.Context, packagePolicyId string) (*ApiResponse[ReadPackagePolicyResponse], error) {
	return doAPI[ReadPackagePolicyResponse](
		c, ctx,
		"GET", "/package_policies/{id}",
		map[string]string{"id": packagePolicyId},
		nil, nil,
	)
}

type ReadPackagePolicyResponse struct {
	Item PackagePolicy `json:"item"`
}

// ============================================================================

func (c *Client) CreatePackagePolicy(ctx context.Context, body PackagePolicyRequest) (*ApiResponse[CreatePackagePolicyResponse], error) {
	return doAPI[CreatePackagePolicyResponse](
		c, ctx,
		"POST", "/package_policies",
		nil, body, nil,
	)
}

type CreatePackagePolicyResponse struct {
	Item PackagePolicy `json:"item"`
}

// ============================================================================

func (c *Client) UpdatePackagePolicy(ctx context.Context, packagePolicyId string, body PackagePolicyRequest) (*ApiResponse[UpdatePackagePolicyResponse], error) {
	return doAPI[UpdatePackagePolicyResponse](
		c, ctx,
		"PUT", "/package_policies/{id}",
		map[string]string{"id": packagePolicyId},
		body, nil,
	)
}

type UpdatePackagePolicyResponse struct {
	Item   PackagePolicy `json:"item"`
	Sucess bool          `json:"sucess"`
}

// ============================================================================

func (c *Client) DeletePackagePolicy(ctx context.Context, packagePolicyId string, params *DeletePackagePolicyParams) (*ApiResponse[DeletePackagePolicyResponse], error) {
	return doAPI[DeletePackagePolicyResponse](
		c, ctx,
		"DELETE", "/package_policies/{id}",
		map[string]string{"id": packagePolicyId},
		nil, params,
	)
}

type DeletePackagePolicyParams struct {
	Force *bool `url:"force,omitempty"`
}

type DeletePackagePolicyResponse struct {
	Id string `json:"id"`
}

// ============================================================================

type PackagePolicy struct {
	Description      *string                     `json:"description,omitempty"`
	Enabled          *bool                       `json:"enabled,omitempty"`
	Id               string                      `json:"id"`
	Inputs           []PackagePolicyInput        `json:"inputs"`
	Name             string                      `json:"name"`
	Namespace        *string                     `json:"namespace,omitempty"`
	OutputId         *string                     `json:"output_id,omitempty"`
	Package          *PackagePolicyPackage       `json:"package,omitempty"`
	PolicyId         *string                     `json:"policy_id,omitempty"`
	Revision         float32                     `json:"revision"`
	SecretReferences []PackagePolicySecretRef    `json:"secret_references"`
	Vars             map[string]PackagePolicyVar `json:"vars,omitempty"`
}

type PackagePolicyInput struct {
	Config         map[string]any              `json:"config,omitempty"`
	Enabled        bool                        `json:"enabled"`
	PolicyTemplate string                      `json:"policy_template"`
	Processors     []string                    `json:"processors,omitempty"`
	Streams        []PackagePolicyInputStream  `json:"streams,omitempty"`
	Type           string                      `json:"type"`
	Vars           map[string]PackagePolicyVar `json:"vars,omitempty"`
}

type PackagePolicyInputStream struct {
	DataStream PackagePolicyInputStreamDataStream `json:"data_stream"`
	Enabled    bool                               `json:"enabled,omitempty"`
	Vars       map[string]PackagePolicyVar        `json:"vars,omitempty"`
}

type PackagePolicyInputStreamDataStream struct {
	Type    string `json:"type"`
	Dataset string `json:"dataset"`
}

type PackagePolicyPackage struct {
	Name    string  `json:"name"`
	Title   *string `json:"title,omitempty"`
	Version string  `json:"version"`
}

type PackagePolicySecretRef struct {
	Id string `json:"id"`
}

type PackagePolicyVar struct {
	Value any    `json:"value,omitempty"`
	Type  string `json:"type"`
}

type PackagePolicyVarSecretValue struct {
	Id          string `json:"id"`
	IsSecretRef bool   `json:"isSecretRef"`
}

type PackagePolicyRequest struct {
	Description *string                              `json:"description,omitempty"`
	Force       *bool                                `json:"force,omitempty"`
	Id          *string                              `json:"id,omitempty"`
	Inputs      map[string]PackagePolicyRequestInput `json:"inputs,omitempty"`
	Name        string                               `json:"name"`
	Namespace   *string                              `json:"namespace,omitempty"`
	Package     PackagePolicyRequestPackage          `json:"package"`
	PolicyId    string                               `json:"policy_id"`
	Vars        map[string]any                       `json:"vars,omitempty"`
}

type PackagePolicyRequestInput struct {
	Enabled *bool                                      `json:"enabled,omitempty"`
	Streams map[string]PackagePolicyRequestInputStream `json:"streams,omitempty"`
	Vars    map[string]any                             `json:"vars,omitempty"`
}

type PackagePolicyRequestInputStream struct {
	Enabled *bool          `json:"enabled,omitempty"`
	Vars    map[string]any `json:"vars,omitempty"`
}

type PackagePolicyRequestPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
