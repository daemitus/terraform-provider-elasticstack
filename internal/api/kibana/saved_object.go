package kibana

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

func (c *Client) ImportSavedObjects(ctx context.Context, space string, body []byte, params *ImportSavedObjectsParams) (*ApiResponse[ImportSavedObjectsResponse], error) {
	req := c.http.R().SetContext(ctx)
	req.SetPathParams(map[string]string{"space": space})
	req.SetHeader("Content-Type", "multipart/form-data")
	req.SetFileReader("file", "file.ndjson", bytes.NewReader(body))

	if params != nil {
		queryParams, err := query.Values(params)
		if err != nil {
			return nil, err
		}
		req.SetQueryParamsFromValues(queryParams)
	}

	resp, err := req.Execute("POST", "/s/{space}/api/saved_objects/_import")
	if err != nil {
		return nil, err
	}

	response := &ApiResponse[ImportSavedObjectsResponse]{
		Contents: resp.Body(),
		Response: resp,
	}

	if strings.Contains(resp.Header().Get("Content-Type"), "application/json") {
		switch resp.StatusCode() {
		case http.StatusOK, http.StatusNoContent:
			var dest ImportSavedObjectsResponse
			if err := json.Unmarshal(response.Contents, &dest); err != nil {
				return nil, err
			}
			response.Output = dest
		default:
			var dest Error
			if err := json.Unmarshal(response.Contents, &dest); err != nil {
				return nil, err
			}
			response.Error = &dest
		}
	}

	return response, nil
}

// ============================================================================

type ImportSavedObjectsParams struct {
	CompatibilityMode *bool `url:"compatibilityMode,omitempty"`
	CreateNewCopies   *bool `url:"createNewCopies,omitempty"`
	Overwrite         *bool `url:"overwrite,omitempty"`
}

type ImportSavedObjectsResponse struct {
	Success        bool                        `json:"success"`
	SuccessCount   int64                       `json:"successCount"`
	SuccessResults []ImportSavedObjectsSuccess `json:"successResults"`
	Errors         []ImportSavedObjectsError   `json:"errors"`
}

type ImportSavedObjectsSuccess struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	DestinationID string                 `json:"destinationId"`
	Meta          ImportSavedObjectsMeta `json:"meta"`
}

type ImportSavedObjectsMeta struct {
	Icon  string `json:"icon"`
	Title string `json:"title"`
}

type ImportSavedObjectsError struct {
	ID    string                      `json:"id"`
	Type  string                      `json:"type"`
	Title string                      `json:"title"`
	Error ImportSavedObjectsErrorType `json:"error"`
	Meta  ImportSavedObjectsMeta      `json:"meta"`
}

type ImportSavedObjectsErrorType struct {
	Type string `json:"type"`
}
