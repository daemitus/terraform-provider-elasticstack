package output

import (
	"context"
	"fmt"

	fleetapi "github.com/elastic/terraform-provider-elasticstack/generated/kibana"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type outputModel struct {
	ID                   types.String `tfsdk:"id"`
	OutputID             types.String `tfsdk:"output_id"`
	Name                 types.String `tfsdk:"name"`
	Type                 types.String `tfsdk:"type"`
	Hosts                types.List   `tfsdk:"hosts"` //> string
	CaSha256             types.String `tfsdk:"ca_sha256"`
	CaTrustedFingerprint types.String `tfsdk:"ca_trusted_fingerprint"`
	DefaultIntegrations  types.Bool   `tfsdk:"default_integrations"`
	DefaultMonitoring    types.Bool   `tfsdk:"default_monitoring"`
	Ssl                  types.List   `tfsdk:"ssl"` //> outputSslModel
	ConfigYaml           types.String `tfsdk:"config_yaml"`
}

type outputSslModel struct {
	CertificateAuthorities types.List   `tfsdk:"certificate_authorities"` //> string
	Certificate            types.String `tfsdk:"certificate"`
	Key                    types.String `tfsdk:"key"`
}

func (model *outputModel) populateFromAPICreate(ctx context.Context, data *fleetapi.FleetOutputCreateRequest) (diags diag.Diagnostics) {
	if data == nil {
		return
	}

	union, err := data.ValueByDiscriminator()
	if err != nil {
		diags.AddError(err.Error(), "")
		return
	}

	var nd diag.Diagnostics
	switch data := union.(type) {
	case fleetapi.FleetOutputCreateRequestElasticsearch:
		model.ID = types.StringPointerValue(data.Id)
		model.OutputID = types.StringPointerValue(data.Id)
		model.Name = types.StringValue(data.Name)
		model.Type = types.StringValue(string(data.Type))
		model.Hosts = utils.SliceToListType_String(ctx, data.Hosts, path.Root("hosts"), diags)
		model.CaSha256 = types.StringPointerValue(data.CaSha256)
		model.CaTrustedFingerprint = types.StringPointerValue(data.CaTrustedFingerprint)
		model.DefaultIntegrations = types.BoolPointerValue(data.IsDefault)
		model.DefaultMonitoring = types.BoolPointerValue(data.IsDefaultMonitoring)
		model.ConfigYaml = types.StringPointerValue(data.ConfigYaml)

		if data.Ssl != nil {
			p := path.Root("ssl")
			sslModels := []outputSslModel{{
				CertificateAuthorities: utils.SliceToListType_String(ctx, data.Ssl.CertificateAuthorities, p.AtName("certificate_authorities"), diags),
				Certificate:            types.StringPointerValue(data.Ssl.Certificate),
				Key:                    types.StringPointerValue(data.Ssl.Key),
			}}
			model.Ssl, nd = types.ListValueFrom(ctx, getSslAttrTypes(), sslModels)
			diags.Append(nd...)
		} else {
			model.Ssl = types.ListNull(getSslAttrTypes())
		}

	case fleetapi.FleetOutputCreateRequestLogstash:
		model.ID = types.StringPointerValue(data.Id)
		model.OutputID = types.StringPointerValue(data.Id)
		model.Name = types.StringValue(data.Name)
		model.Type = types.StringValue(string(data.Type))
		model.Hosts = utils.SliceToListType_String(ctx, data.Hosts, path.Root("hosts"), diags)
		model.CaSha256 = types.StringPointerValue(data.CaSha256)
		model.CaTrustedFingerprint = types.StringPointerValue(data.CaTrustedFingerprint)
		model.DefaultIntegrations = types.BoolPointerValue(data.IsDefault)
		model.DefaultMonitoring = types.BoolPointerValue(data.IsDefaultMonitoring)
		model.ConfigYaml = types.StringPointerValue(data.ConfigYaml)

		if data.Ssl != nil {
			p := path.Root("ssl")
			sslModels := []outputSslModel{{
				CertificateAuthorities: utils.SliceToListType_String(ctx, data.Ssl.CertificateAuthorities, p.AtName("certificate_authorities"), diags),
				Certificate:            types.StringPointerValue(data.Ssl.Certificate),
				Key:                    types.StringPointerValue(data.Ssl.Key),
			}}
			model.Ssl, nd = types.ListValueFrom(ctx, getSslAttrTypes(), sslModels)
			diags.Append(nd...)
		} else {
			model.Ssl = types.ListNull(getSslAttrTypes())
		}

	case fleetapi.FleetOutputCreateRequestKafka:
		model.ID = types.StringPointerValue(data.Id)
		model.OutputID = types.StringPointerValue(data.Id)
		model.Name = types.StringValue(data.Name)
		model.Type = types.StringValue(string(data.Type))
		model.Hosts = utils.SliceToListType_String(ctx, data.Hosts, path.Root("hosts"), diags)
		model.CaSha256 = types.StringPointerValue(data.CaSha256)
		model.CaTrustedFingerprint = types.StringPointerValue(data.CaTrustedFingerprint)
		model.DefaultIntegrations = types.BoolPointerValue(data.IsDefault)
		model.DefaultMonitoring = types.BoolPointerValue(data.IsDefaultMonitoring)
		model.ConfigYaml = types.StringPointerValue(data.ConfigYaml)

		if data.Ssl != nil {
			p := path.Root("ssl")
			sslModels := []outputSslModel{{
				CertificateAuthorities: utils.SliceToListType_String(ctx, data.Ssl.CertificateAuthorities, p.AtName("certificate_authorities"), diags),
				Certificate:            types.StringPointerValue(data.Ssl.Certificate),
				Key:                    types.StringPointerValue(data.Ssl.Key),
			}}
			model.Ssl, nd = types.ListValueFrom(ctx, getSslAttrTypes(), sslModels)
			diags.Append(nd...)
		} else {
			model.Ssl = types.ListNull(getSslAttrTypes())
		}

	default:
		diags.AddError(fmt.Sprintf("unhandled output type: %T", data), "")
	}

	return
}

func (model *outputModel) populateFromAPIUpdate(ctx context.Context, data *fleetapi.FleetOutputUpdateRequest) (diags diag.Diagnostics) {
	if data == nil {
		return
	}

	union, err := data.ValueByDiscriminator()
	if err != nil {
		diags.AddError(err.Error(), "")
		return
	}

	var nd diag.Diagnostics
	switch data := union.(type) {
	case fleetapi.FleetOutputUpdateRequestElasticsearch:
		model.ID = types.StringPointerValue(data.Id)
		model.OutputID = types.StringPointerValue(data.Id)
		model.Name = types.StringValue(data.Name)
		model.Type = types.StringValue(string(data.Type))
		model.Hosts = utils.SliceToListType_String(ctx, data.Hosts, path.Root("hosts"), diags)
		model.CaSha256 = types.StringPointerValue(data.CaSha256)
		model.CaTrustedFingerprint = types.StringPointerValue(data.CaTrustedFingerprint)
		model.DefaultIntegrations = types.BoolPointerValue(data.IsDefault)
		model.DefaultMonitoring = types.BoolPointerValue(data.IsDefaultMonitoring)
		model.ConfigYaml = types.StringPointerValue(data.ConfigYaml)

		if data.Ssl != nil {
			p := path.Root("ssl")
			sslModel := []outputSslModel{{
				CertificateAuthorities: utils.SliceToListType_String(ctx, data.Ssl.CertificateAuthorities, p.AtName("certificate_authorities"), diags),
				Certificate:            types.StringPointerValue(data.Ssl.Certificate),
				Key:                    types.StringPointerValue(data.Ssl.Key),
			}}
			model.Ssl, nd = types.ListValueFrom(ctx, getSslAttrTypes(), sslModel)
			diags.Append(nd...)
		} else {
			model.Ssl = types.ListNull(getSslAttrTypes())
		}

	case fleetapi.FleetOutputUpdateRequestLogstash:
		model.ID = types.StringPointerValue(data.Id)
		model.OutputID = types.StringPointerValue(data.Id)
		model.Name = types.StringValue(data.Name)
		model.Type = types.StringValue(string(data.Type))
		model.Hosts = utils.SliceToListType_String(ctx, data.Hosts, path.Root("hosts"), diags)
		model.CaSha256 = types.StringPointerValue(data.CaSha256)
		model.CaTrustedFingerprint = types.StringPointerValue(data.CaTrustedFingerprint)
		model.DefaultIntegrations = types.BoolPointerValue(data.IsDefault)
		model.DefaultMonitoring = types.BoolPointerValue(data.IsDefaultMonitoring)
		model.ConfigYaml = types.StringPointerValue(data.ConfigYaml)

		if data.Ssl != nil {
			p := path.Root("ssl")
			sslModel := []outputSslModel{{
				CertificateAuthorities: utils.SliceToListType_String(ctx, data.Ssl.CertificateAuthorities, p.AtName("certificate_authorities"), diags),
				Certificate:            types.StringPointerValue(data.Ssl.Certificate),
				Key:                    types.StringPointerValue(data.Ssl.Key),
			}}
			model.Ssl, nd = types.ListValueFrom(ctx, getSslAttrTypes(), sslModel)
			diags.Append(nd...)
		} else {
			model.Ssl = types.ListNull(getSslAttrTypes())
		}

	case fleetapi.FleetOutputUpdateRequestKafka:
		model.ID = types.StringPointerValue(data.Id)
		model.OutputID = types.StringPointerValue(data.Id)
		model.Name = types.StringValue(data.Name)
		model.Type = types.StringValue(string(data.Type))
		model.Hosts = utils.SliceToListType_String(ctx, data.Hosts, path.Root("hosts"), diags)
		model.CaSha256 = types.StringPointerValue(data.CaSha256)
		model.CaTrustedFingerprint = types.StringPointerValue(data.CaTrustedFingerprint)
		model.DefaultIntegrations = types.BoolPointerValue(data.IsDefault)
		model.DefaultMonitoring = types.BoolPointerValue(data.IsDefaultMonitoring)
		model.ConfigYaml = types.StringPointerValue(data.ConfigYaml)

		if data.Ssl != nil {
			p := path.Root("ssl")
			sslModel := []outputSslModel{{
				CertificateAuthorities: utils.SliceToListType_String(ctx, data.Ssl.CertificateAuthorities, p.AtName("certificate_authorities"), diags),
				Certificate:            types.StringPointerValue(data.Ssl.Certificate),
				Key:                    types.StringPointerValue(data.Ssl.Key),
			}}
			model.Ssl, nd = types.ListValueFrom(ctx, getSslAttrTypes(), sslModel)
			diags.Append(nd...)
		} else {
			model.Ssl = types.ListNull(getSslAttrTypes())
		}

	default:
		diags.AddError(fmt.Sprintf("unhandled output type: %T", data), "")
	}

	return
}

func (model outputModel) toAPICreateModel(ctx context.Context) (union fleetapi.FleetOutputCreateRequest, diags diag.Diagnostics) {
	outputType := model.Type.ValueString()
	switch outputType {
	case "elasticsearch":
		body := fleetapi.FleetOutputCreateRequestElasticsearch{
			CaSha256:             model.CaSha256.ValueStringPointer(),
			CaTrustedFingerprint: model.CaTrustedFingerprint.ValueStringPointer(),
			ConfigYaml:           model.ConfigYaml.ValueStringPointer(),
			Hosts:                utils.ListTypeToSlice_String(ctx, model.Hosts, path.Root("hosts"), diags),
			Id:                   model.OutputID.ValueStringPointer(),
			IsDefault:            model.DefaultIntegrations.ValueBoolPointer(),
			IsDefaultMonitoring:  model.DefaultMonitoring.ValueBoolPointer(),
			Name:                 model.Name.ValueString(),
			Type:                 fleetapi.FleetOutputCreateRequestElasticsearchTypeElasticsearch,
		}

		// Can't use helpers for anonymous structs
		if utils.IsKnown(model.Ssl) {
			sslModels := utils.ListTypeAs[outputSslModel](ctx, model.Ssl, path.Root("ssl"), diags)
			if len(sslModels) > 0 {
				body.Ssl = &struct {
					Certificate            *string  `json:"certificate,omitempty"`
					CertificateAuthorities []string `json:"certificate_authorities,omitempty"`
					Key                    *string  `json:"key,omitempty"`
				}{
					Certificate:            sslModels[0].Certificate.ValueStringPointer(),
					CertificateAuthorities: utils.ListTypeToSlice_String(ctx, sslModels[0].CertificateAuthorities, path.Root("certificate_authorities"), diags),
					Key:                    sslModels[0].Key.ValueStringPointer(),
				}
			}
		}

		err := union.FromFleetOutputCreateRequestElasticsearch(body)
		if err != nil {
			diags.AddError(err.Error(), "")
		}

	case "logstash":
		body := fleetapi.FleetOutputCreateRequestLogstash{
			CaSha256:             model.CaSha256.ValueStringPointer(),
			CaTrustedFingerprint: model.CaTrustedFingerprint.ValueStringPointer(),
			ConfigYaml:           model.ConfigYaml.ValueStringPointer(),
			Hosts:                utils.ListTypeToSlice_String(ctx, model.Hosts, path.Root("hosts"), diags),
			Id:                   model.OutputID.ValueStringPointer(),
			IsDefault:            model.DefaultIntegrations.ValueBoolPointer(),
			IsDefaultMonitoring:  model.DefaultMonitoring.ValueBoolPointer(),
			Name:                 model.Name.ValueString(),
			Type:                 fleetapi.FleetOutputCreateRequestLogstashTypeLogstash,
		}

		// Can't use helpers for anonymous structs
		if utils.IsKnown(model.Ssl) {
			sslModels := utils.ListTypeAs[outputSslModel](ctx, model.Ssl, path.Root("ssl"), diags)
			if len(sslModels) > 0 {
				body.Ssl = &struct {
					Certificate            *string  `json:"certificate,omitempty"`
					CertificateAuthorities []string `json:"certificate_authorities,omitempty"`
					Key                    *string  `json:"key,omitempty"`
				}{
					Certificate:            sslModels[0].Certificate.ValueStringPointer(),
					CertificateAuthorities: utils.ListTypeToSlice_String(ctx, sslModels[0].CertificateAuthorities, path.Root("certificate_authorities"), diags),
					Key:                    sslModels[0].Key.ValueStringPointer(),
				}
			}
		}

		err := union.FromFleetOutputCreateRequestLogstash(body)
		if err != nil {
			diags.AddError(err.Error(), "")
		}

	case "kafka":
		body := fleetapi.FleetOutputCreateRequestKafka{
			Type:                 fleetapi.FleetOutputCreateRequestKafkaTypeKafka,
			CaSha256:             model.CaSha256.ValueStringPointer(),
			CaTrustedFingerprint: model.CaTrustedFingerprint.ValueStringPointer(),
			ConfigYaml:           model.ConfigYaml.ValueStringPointer(),
			Hosts:                utils.ListTypeToSlice_String(ctx, model.Hosts, path.Root("hosts"), diags),
			Id:                   model.OutputID.ValueStringPointer(),
			IsDefault:            model.DefaultIntegrations.ValueBoolPointer(),
			IsDefaultMonitoring:  model.DefaultMonitoring.ValueBoolPointer(),
			Name:                 model.Name.ValueString(),
		}

		// Can't use helpers for anonymous structs
		if utils.IsKnown(model.Ssl) {
			sslModels := utils.ListTypeAs[outputSslModel](ctx, model.Ssl, path.Root("ssl"), diags)
			if len(sslModels) > 0 {
				body.Ssl = &struct {
					Certificate            *string                                                    `json:"certificate,omitempty"`
					CertificateAuthorities []string                                                   `json:"certificate_authorities,omitempty"`
					Key                    *string                                                    `json:"key,omitempty"`
					VerificationMode       *fleetapi.FleetOutputCreateRequestKafkaSslVerificationMode `json:"verification_mode,omitempty"`
				}{
					Certificate:            sslModels[0].Certificate.ValueStringPointer(),
					CertificateAuthorities: utils.ListTypeToSlice_String(ctx, sslModels[0].CertificateAuthorities, path.Root("certificate_authorities"), diags),
					Key:                    sslModels[0].Key.ValueStringPointer(),
				}
			}
		}

		err := union.FromFleetOutputCreateRequestKafka(body)
		if err != nil {
			diags.AddError(err.Error(), "")
		}

	default:
		diags.AddError(fmt.Sprintf("unhandled output type: %s", outputType), "")
	}

	return
}

func (model outputModel) toAPIUpdateModel(ctx context.Context) (union fleetapi.FleetOutputUpdateRequest, diags diag.Diagnostics) {
	outputType := model.Type.ValueString()
	switch outputType {
	case "elasticsearch":
		body := fleetapi.FleetOutputUpdateRequestElasticsearch{
			Type:                 fleetapi.FleetOutputUpdateRequestElasticsearchTypeElasticsearch,
			CaSha256:             model.CaSha256.ValueStringPointer(),
			CaTrustedFingerprint: model.CaTrustedFingerprint.ValueStringPointer(),
			ConfigYaml:           model.ConfigYaml.ValueStringPointer(),
			Hosts:                utils.ListTypeToSlice_String(ctx, model.Hosts, path.Root("hosts"), diags),
			IsDefault:            model.DefaultIntegrations.ValueBoolPointer(),
			IsDefaultMonitoring:  model.DefaultMonitoring.ValueBoolPointer(),
			Name:                 model.Name.ValueString(),
		}

		// Can't use helpers for anonymous structs
		if utils.IsKnown(model.Ssl) {
			sslModels := utils.ListTypeAs[outputSslModel](ctx, model.Ssl, path.Root("ssl"), diags)
			if len(sslModels) > 0 {
				body.Ssl = &struct {
					Certificate            *string  `json:"certificate,omitempty"`
					CertificateAuthorities []string `json:"certificate_authorities,omitempty"`
					Key                    *string  `json:"key,omitempty"`
				}{
					Certificate:            sslModels[0].Certificate.ValueStringPointer(),
					CertificateAuthorities: utils.ListTypeToSlice_String(ctx, sslModels[0].CertificateAuthorities, path.Root("certificate_authorities"), diags),
					Key:                    sslModels[0].Key.ValueStringPointer(),
				}
			}
		}

		err := union.FromFleetOutputUpdateRequestElasticsearch(body)
		if err != nil {
			diags.AddError(err.Error(), "")
		}

	case "logstash":
		body := fleetapi.FleetOutputUpdateRequestLogstash{
			CaSha256:             model.CaSha256.ValueStringPointer(),
			CaTrustedFingerprint: model.CaTrustedFingerprint.ValueStringPointer(),
			ConfigYaml:           model.ConfigYaml.ValueStringPointer(),
			Hosts:                utils.ListTypeToSlice_String(ctx, model.Hosts, path.Root("hosts"), diags),
			IsDefault:            model.DefaultIntegrations.ValueBoolPointer(),
			IsDefaultMonitoring:  model.DefaultMonitoring.ValueBoolPointer(),
			Name:                 model.Name.ValueString(),
			Type:                 fleetapi.FleetOutputUpdateRequestLogstashTypeLogstash,
		}

		// Can't use helpers for anonymous structs
		if utils.IsKnown(model.Ssl) {
			sslModels := utils.ListTypeAs[outputSslModel](ctx, model.Ssl, path.Root("ssl"), diags)
			if len(sslModels) > 0 {
				body.Ssl = &struct {
					Certificate            *string  `json:"certificate,omitempty"`
					CertificateAuthorities []string `json:"certificate_authorities,omitempty"`
					Key                    *string  `json:"key,omitempty"`
				}{
					Certificate:            sslModels[0].Certificate.ValueStringPointer(),
					CertificateAuthorities: utils.ListTypeToSlice_String(ctx, sslModels[0].CertificateAuthorities, path.Root("certificate_authorities"), diags),
					Key:                    sslModels[0].Key.ValueStringPointer(),
				}
			}
		}

		err := union.FromFleetOutputUpdateRequestLogstash(body)
		if err != nil {
			diags.AddError(err.Error(), "")
		}

	case "kafka":
		body := fleetapi.FleetOutputUpdateRequestKafka{
			Type:                 fleetapi.FleetOutputUpdateRequestKafkaTypeKafka,
			CaSha256:             model.CaSha256.ValueStringPointer(),
			CaTrustedFingerprint: model.CaTrustedFingerprint.ValueStringPointer(),
			ConfigYaml:           model.ConfigYaml.ValueStringPointer(),
			Hosts:                utils.ListTypeToSlice_String(ctx, model.Hosts, path.Root("hosts"), diags),
			IsDefault:            model.DefaultIntegrations.ValueBoolPointer(),
			IsDefaultMonitoring:  model.DefaultMonitoring.ValueBoolPointer(),
			Name:                 model.Name.ValueString(),
		}

		// Can't use helpers for anonymous structs
		if utils.IsKnown(model.Ssl) {
			sslModels := utils.ListTypeAs[outputSslModel](ctx, model.Ssl, path.Root("ssl"), diags)
			if len(sslModels) > 0 {
				body.Ssl = &struct {
					Certificate            *string                                                    `json:"certificate,omitempty"`
					CertificateAuthorities []string                                                   `json:"certificate_authorities,omitempty"`
					Key                    *string                                                    `json:"key,omitempty"`
					VerificationMode       *fleetapi.FleetOutputUpdateRequestKafkaSslVerificationMode `json:"verification_mode,omitempty"`
				}{
					Certificate:            sslModels[0].Certificate.ValueStringPointer(),
					CertificateAuthorities: utils.ListTypeToSlice_String(ctx, sslModels[0].CertificateAuthorities, path.Root("certificate_authorities"), diags),
					Key:                    sslModels[0].Key.ValueStringPointer(),
				}
			}
		}

		err := union.FromFleetOutputUpdateRequestKafka(body)
		if err != nil {
			diags.AddError(err.Error(), "")
		}

	default:
		diags.AddError(fmt.Sprintf("unhandled output type: %s", outputType), "")
	}

	return
}
