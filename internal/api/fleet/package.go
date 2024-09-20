package fleet

import (
	"context"
)

func (c *Client) ListPackages(ctx context.Context, params *ListPackagesParams) (*ApiResponse[ListPackagesResponse], error) {
	return doAPI[ListPackagesResponse](
		c, ctx,
		"GET", "/epm/packages",
		nil, nil, params,
	)
}

type ListPackagesParams struct {
	ExcludeInstallStatus *bool   `url:"excludeInstallStatus,omitempty"`
	Prerelease           *bool   `url:"prerelease,omitempty"`
	Experimental         *bool   `url:"experimental,omitempty"`
	Category             *string `url:"category,omitempty"`
}

type ListPackagesResponse struct {
	Items SearchResults `json:"items"`
}

type SearchResults []SearchResult

type SearchResult struct {
	Description string          `json:"description"`
	Download    string          `json:"download"`
	Name        string          `json:"name"`
	Path        string          `json:"path"`
	SavedObject *map[string]any `json:"savedObject,omitempty"`
	Status      string          `json:"status"`
	Title       string          `json:"title"`
	Type        string          `json:"type"`
	Version     string          `json:"version"`
}

// ============================================================================

func (c *Client) ReadPackage(ctx context.Context, pkgName string, pkgVersion string, params *ReadPackageParams) (*ApiResponse[ReadPackageResponse], error) {
	return doAPI[ReadPackageResponse](
		c, ctx,
		"GET", "/epm/packages/{name}/{version}",
		map[string]string{"name": pkgName, "version": pkgVersion},
		nil, params,
	)
}

type ReadPackageParams struct {
	IgnoreUnverified *bool `url:"ignoreUnverified,omitempty"`
	Full             *bool `url:"full,omitempty"`
	Prerelease       *bool `url:"prerelease,omitempty"`
}

type ReadPackageResponse struct {
	Item                 *PackageInfo   `json:"item,omitempty"`
	KeepPoliciesUpToDate *bool          `json:"keepPoliciesUpToDate,omitempty"`
	LatestVersion        *string        `json:"latestVersion,omitempty"`
	LicensePath          *string        `json:"licensePath,omitempty"`
	Notice               *string        `json:"notice,omitempty"`
	SavedObject          map[string]any `json:"savedObject"`
	Status               PackageStatus  `json:"status"`
}

// ============================================================================

func (c *Client) InstallPackage(ctx context.Context, pkgName string, pkgVersion string, body InstallPackageRequest, params *InstallPackageParams) (*ApiResponse[InstallPackageResponse], error) {
	return doAPI[InstallPackageResponse](
		c, ctx,
		"POST", "/epm/packages/{name}/{version}",
		map[string]string{"name": pkgName, "version": pkgVersion},
		body, params,
	)
}

type InstallPackageRequest struct {
	Force             *bool `json:"force,omitempty"`
	IgnoreConstraints *bool `json:"ignore_constraints,omitempty"`
}

type InstallPackageParams struct {
	IgnoreUnverified *bool `url:"ignoreUnverified,omitempty"`
	Full             *bool `url:"full,omitempty"`
	Prerelease       *bool `url:"prerelease,omitempty"`
}

type InstallPackageResponse struct {
	Meta  *InstallPackageResponseMeta  `json:"_meta,omitempty"`
	Items []InstallPackageResponseItem `json:"items"`
}

type InstallPackageResponseMeta struct {
	InstallSource *PackageInstallSource `json:"install_source,omitempty"`
}

type InstallPackageResponseItem struct {
	Id string `json:"id"`
}

// ============================================================================

func (c *Client) UpdatePackage(ctx context.Context, pkgName string, pkgVersion string, body UpdatePackageRequest, params *UpdatePackageParams) (*ApiResponse[UpdatePackageResponse], error) {
	return doAPI[UpdatePackageResponse](
		c, ctx,
		"PUT", "/epm/packages/{name}/{version}",
		map[string]string{"name": pkgName, "version": pkgVersion},
		body, params,
	)
}

type UpdatePackageRequest struct {
	KeepPoliciesUpToDate *bool `json:"keepPoliciesUpToDate,omitempty"`
}

type UpdatePackageParams struct {
	IgnoreUnverified *bool `url:"ignoreUnverified,omitempty"`
	Full             *bool `url:"full,omitempty"`
	Prerelease       *bool `url:"prerelease,omitempty"`
}

type UpdatePackageResponse struct {
	Items []UpdatePackageResponseItem `json:"items"`
}

type UpdatePackageResponseItem struct {
	Id string `json:"id"`
}

// ============================================================================

func (c *Client) DeletePackage(ctx context.Context, pkgName string, pkgVersion string, body DeletePackageRequest, params *DeletePackageParams) (*ApiResponse[DeletePackageResponse], error) {
	return doAPI[DeletePackageResponse](
		c, ctx,
		"DELETE", "/epm/packages/{name}/{version}",
		map[string]string{"name": pkgName, "version": pkgVersion},
		body, params,
	)
}

type DeletePackageRequest struct {
	Force *bool `json:"force,omitempty"`
}

type DeletePackageParams struct {
	IgnoreUnverified *bool `url:"ignoreUnverified,omitempty"`
	Full             *bool `url:"full,omitempty"`
	Prerelease       *bool `url:"prerelease,omitempty"`
}

type DeletePackageResponse struct {
	Items []DeletePackageResponseItem `json:"items"`
}

type DeletePackageResponseItem struct {
	Id string `json:"id"`
}

// ============================================================================

type PackageInfo struct {
	Assets        map[string]any      `json:"assets"`
	Categories    []string            `json:"categories"`
	Description   string              `json:"description"`
	Download      string              `json:"download"`
	FormatVersion string              `json:"format_version"`
	Internal      *bool               `json:"internal,omitempty"`
	Name          string              `json:"name"`
	Path          string              `json:"path"`
	Readme        *string             `json:"readme,omitempty"`
	Release       *PackageInfoRelease `json:"release,omitempty"`
	Status        string              `json:"status"`
	Title         string              `json:"title"`
	Type          string              `json:"type"`
	Version       string              `json:"version"`
}

type PackageInfoRelease string

const (
	Beta         PackageInfoRelease = "beta"
	Experimental PackageInfoRelease = "experimental"
	Ga           PackageInfoRelease = "ga"
)

type PackageStatus string

const (
	InstallFailed PackageStatus = "install_failed"
	Installed     PackageStatus = "installed"
	Installing    PackageStatus = "installing"
	NotInstalled  PackageStatus = "not_installed"
)

type PackageInstallSource string

const (
	Bundled  PackageInstallSource = "bundled"
	Registry PackageInstallSource = "registry"
	Upload   PackageInstallSource = "upload"
)
