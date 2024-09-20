package xstringvalidator

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = regexDoesNotMatchValidator{}

// regexDoesNotMatchValidator validates that a string Attribute's value does not match the specified regular expression.
type regexDoesNotMatchValidator struct {
	regexp  *regexp.Regexp
	message string
}

// Description describes the validation in plain text formatting.
func (validator regexDoesNotMatchValidator) Description(_ context.Context) string {
	if validator.message != "" {
		return validator.message
	}
	return fmt.Sprintf("value must not match regular expression '%s'", validator.regexp)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator regexDoesNotMatchValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v regexDoesNotMatchValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if v.regexp.MatchString(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// RegexDoesNotMatch returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a string.
//   - Does not match the given regular expression https://github.com/google/re2/wiki/Syntax.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
// Optionally an error message can be provided to return something friendlier
// than "value must match regular expression 'regexp'".
func RegexDoesNotMatch(regexp *regexp.Regexp, message string) validator.String {
	return regexDoesNotMatchValidator{
		regexp:  regexp,
		message: message,
	}
}
