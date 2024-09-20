package util

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// ConvertToAttrDiags wraps an existing collection of diagnostics with an attribute path.
func ConvertToAttrDiags(diags diag.Diagnostics, path path.Path) diag.Diagnostics {
	var nd diag.Diagnostics
	for _, d := range diags {
		nd.AddAttributeError(path, d.Summary(), d.Detail())
	}
	return nd
}

func DiagsAsError(diags diag.Diagnostics) error {
	diag := diags.Errors()[0]
	return fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
}
