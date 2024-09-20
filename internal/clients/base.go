package clients

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type baseClient struct {
}

type ApiResponse interface {
	StatusCode() int
	Body() string
}

func (c *baseClient) reportFromErr(err error) diag.Diagnostics {
	return diag.Diagnostics{
		diag.NewErrorDiagnostic("Unexpected error from API", err.Error()),
	}
}

func (c *baseClient) reportUnknownError(resp ApiResponse) diag.Diagnostics {
	return diag.Diagnostics{
		diag.NewErrorDiagnostic(
			fmt.Sprintf("Unexpected status code from server: got HTTP %d", resp.StatusCode()),
			resp.Body(),
		),
	}
}
