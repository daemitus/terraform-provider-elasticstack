package kibana2

import (
	"fmt"

	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
)

func reportUnknownError(statusCode int, body []byte) fwdiag.Diagnostics {
	return fwdiag.Diagnostics{
		fwdiag.NewErrorDiagnostic(
			fmt.Sprintf("Unexpected status code from server: got HTTP %d", statusCode),
			string(body),
		),
	}
}
