package util

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ validator.Describer = elasticDurationValidator{}
	_ validator.String    = elasticDurationValidator{}
)

// ElasticDurationValidator validates that the value is an Elastic duration.
func ElasticDurationValidator() elasticDurationValidator {
	return elasticDurationValidator{}
}

// elasticDurationValidator validates that the value is an Elastic duration.
type elasticDurationValidator struct{}

func (v elasticDurationValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v elasticDurationValidator) MarkdownDescription(ctx context.Context) string {
	return "value must be a valid Elastic duration. See https://www.elastic.co/guide/en/elasticsearch/reference/current/api-conventions.html#time-units."
}

func (v elasticDurationValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue
	err := v.checkValue(request.ConfigValue.ValueString())
	if err != nil {
		d := validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		)
		response.Diagnostics.Append(d)
	}
}

// Validate that a string is a valid duration using Elastic time units: d, h, m, s, ms, micros, nanos.
// See https://www.elastic.co/guide/en/elasticsearch/reference/current/api-conventions.html#time-units.
func (v elasticDurationValidator) checkValue(value string) error {
	if value == "" {
		return fmt.Errorf("%q contains an invalid duration: [empty]", value)
	}

	rex := regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?(?:d|h|m|s|ms|micros|nanos)$`)

	if !rex.MatchString(value) {
		return fmt.Errorf("%q contains an invalid duration: not conforming to Elastic time-units format", value)
	}

	return nil
}

// ============================================================================

var (
	_ validator.Describer = durationValidator{}
	_ validator.String    = durationValidator{}
)

// DurationValidator validates that the value is an Golang duration.
func DurationValidator() durationValidator {
	return durationValidator{}
}

// durationValidator validates that the value is an Golang duration.
type durationValidator struct{}

func (v durationValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v durationValidator) MarkdownDescription(ctx context.Context) string {
	return "value must be a valid Golang duration. See https://pkg.go.dev/time#ParseDuration."
}

func (v durationValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue
	err := v.checkValue(request.ConfigValue.ValueString())
	if err != nil {
		d := validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		)
		response.Diagnostics.Append(d)
	}
}

// Validate that a string is a valid duration using Golang time units: ns, us, ms, s, m, h.
// See https://pkg.go.dev/time#ParseDuration.
func (v durationValidator) checkValue(value string) error {
	if value == "" {
		return fmt.Errorf("%q contains an invalid duration: [empty]", value)
	}

	if _, err := time.ParseDuration(value); err != nil {
		return fmt.Errorf("%q contains an invalid duration: %s", value, err)
	}

	return nil
}
