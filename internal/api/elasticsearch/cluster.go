package elasticsearch

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8/typedapi/cluster/getsettings"
	"github.com/elastic/go-elasticsearch/v8/typedapi/cluster/health"
	"github.com/elastic/go-elasticsearch/v8/typedapi/cluster/info"
	"github.com/elastic/go-elasticsearch/v8/typedapi/cluster/putsettings"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/deletescript"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/getscript"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/putscript"
	"github.com/elastic/go-elasticsearch/v8/typedapi/slm/deletelifecycle"
	"github.com/elastic/go-elasticsearch/v8/typedapi/slm/getlifecycle"
	"github.com/elastic/go-elasticsearch/v8/typedapi/slm/putlifecycle"
	createsnapshot "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/createrepository"
	deletesnapshot "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/delete"
	"github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/deleterepository"
	getsnapshot "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/getrepository"
)

func (c *Client) GetClusterHealth(ctx context.Context) (*ApiResponse[GetClusterHealthResponse], error) {
	return doApiPtr(c.http.Cluster.Health().Do(ctx))
}

type (
	GetClusterHealthResponse = health.Response
)

// ============================================================================

func (c *Client) GetClusterInfo(ctx context.Context, target string) (*ApiResponse[GetClusterInfoResponse], error) {
	return doApiPtr(c.http.Cluster.Info(target).Do(ctx))
}

type (
	GetClusterInfoResponse = info.Response
)

// ============================================================================

func (c *Client) GetSlmPolicy(ctx context.Context, id string) (*ApiResponse[GetSlmPolicyResponse], error) {
	return doApi(c.http.Slm.GetLifecycle().PolicyId(id).Do(ctx))
}

func (c *Client) PutSlmPolicy(ctx context.Context, id string, req PutSlmPolicyRequest) (*ApiResponse[PutSlmPolicyResponse], error) {
	return doApiPtr(c.http.Slm.PutLifecycle(id).Request(&req).Do(ctx))
}

func (c *Client) DeleteSlmPolicy(ctx context.Context, id string) (*ApiResponse[DeleteSlmPolicyResponse], error) {
	return doApiPtr(c.http.Slm.DeleteLifecycle(id).Do(ctx))
}

type (
	GetSlmPolicyResponse    = getlifecycle.Response
	PutSlmPolicyRequest     = putlifecycle.Request
	PutSlmPolicyResponse    = putlifecycle.Response
	DeleteSlmPolicyResponse = deletelifecycle.Response
)

// ============================================================================

func (c *Client) GetSnapshot(ctx context.Context, repository string, snapshot string) (*ApiResponse[GetSnapshotResponse], error) {
	return doApiPtr(c.http.Snapshot.Get(repository, snapshot).Do(ctx))
}

func (c *Client) CreateSnapshot(ctx context.Context, repository string, snapshot string, req CreateSnapshotRequest) (*ApiResponse[CreateSnapshotResponse], error) {
	return doApiPtr(c.http.Snapshot.Create(repository, snapshot).Request(&req).Do(ctx))
}

func (c *Client) DeleteSnapshot(ctx context.Context, repository string, snapshot string) (*ApiResponse[DeleteSnapshotResponse], error) {
	return doApiPtr(c.http.Snapshot.Delete(repository, snapshot).Do(ctx))
}

type (
	GetSnapshotResponse    = getsnapshot.Response
	CreateSnapshotRequest  = createsnapshot.Request
	CreateSnapshotResponse = createsnapshot.Response
	DeleteSnapshotResponse = deletesnapshot.Response
)

// ============================================================================

func (c *Client) GetRepository(ctx context.Context, repository string) (*ApiResponse[GetRepositoryResponse], error) {
	return doApi(c.http.Snapshot.GetRepository().Repository(repository).Do(ctx))
}

func (c *Client) CreateRepository(ctx context.Context, repository string, req CreateRepositoryRequest) (*ApiResponse[CreateRepositoryResponse], error) {
	return doApiPtr(c.http.Snapshot.CreateRepository(repository).Request(&req).Do(ctx))
}

func (c *Client) DeleteRepository(ctx context.Context, repository string) (*ApiResponse[DeleteRepositoryResponse], error) {
	return doApiPtr(c.http.Snapshot.DeleteRepository(repository).Do(ctx))
}

type (
	GetRepositoryResponse    = getrepository.Response
	CreateRepositoryRequest  = createrepository.Request
	CreateRepositoryResponse = createrepository.Response
	DeleteRepositoryResponse = deleterepository.Response
)

// ============================================================================

func (c *Client) GetClusterSettings(ctx context.Context) (*ApiResponse[GetClusterSettingsResponse], error) {
	return doApiPtr(c.http.Cluster.GetSettings().Do(ctx))
}

func (c *Client) UpdateClusterSettings(ctx context.Context, req PutClusterSettingsRequest) (*ApiResponse[PutClusterSettingsResponse], error) {
	return doApiPtr(c.http.Cluster.PutSettings().Request(&req).Do(ctx))
}

type (
	GetClusterSettingsResponse = getsettings.Response
	PutClusterSettingsRequest  = putsettings.Request
	PutClusterSettingsResponse = putsettings.Response
)

// ============================================================================

func (c *Client) ReadScript(ctx context.Context, id string) (*ApiResponse[GetScriptResponse], error) {
	return doApiPtr(c.http.GetScript(id).Do(ctx))
}

func (c *Client) PutScript(ctx context.Context, id string, req PutScriptRequest) (*ApiResponse[PutScriptResponse], error) {
	return doApiPtr(c.http.PutScript(id).Request(&req).Do(ctx))
}

func (c *Client) DeleteScript(ctx context.Context, id string) (*ApiResponse[DeleteScriptResponse], error) {
	return doApiPtr(c.http.DeleteScript(id).Do(ctx))
}

type (
	GetScriptResponse    = getscript.Response
	PutScriptRequest     = putscript.Request
	PutScriptResponse    = putscript.Response
	DeleteScriptResponse = deletescript.Response
)
