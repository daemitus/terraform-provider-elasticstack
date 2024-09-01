/*
Alerting

OpenAPI schema for alerting endpoints

API version: 0.2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package alerting

import (
	"encoding/json"
)

// checks if the GetRuleTypes200ResponseInnerAlerts type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetRuleTypes200ResponseInnerAlerts{}

// GetRuleTypes200ResponseInnerAlerts Details for writing alerts as data documents for this rule type.
type GetRuleTypes200ResponseInnerAlerts struct {
	// The namespace for this rule type.
	Context *string `json:"context,omitempty"`
	// Indicates whether new fields are added dynamically.
	Dynamic *string `json:"dynamic,omitempty"`
	// Indicates whether the alerts are space-aware. If true, space-specific alert indices are used.
	IsSpaceAware *bool                                       `json:"isSpaceAware,omitempty"`
	Mappings     *GetRuleTypes200ResponseInnerAlertsMappings `json:"mappings,omitempty"`
	// A secondary alias. It is typically used to support the signals alias for detection rules.
	SecondaryAlias *string `json:"secondaryAlias,omitempty"`
	// Indicates whether the rule should write out alerts as data.
	ShouldWrite *bool `json:"shouldWrite,omitempty"`
	// Indicates whether to include the ECS component template for the alerts.
	UseEcs *bool `json:"useEcs,omitempty"`
	// Indicates whether to include the legacy component template for the alerts.
	UseLegacyAlerts *bool `json:"useLegacyAlerts,omitempty"`
}

// NewGetRuleTypes200ResponseInnerAlerts instantiates a new GetRuleTypes200ResponseInnerAlerts object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetRuleTypes200ResponseInnerAlerts() *GetRuleTypes200ResponseInnerAlerts {
	this := GetRuleTypes200ResponseInnerAlerts{}
	var useLegacyAlerts bool = false
	this.UseLegacyAlerts = &useLegacyAlerts
	return &this
}

// NewGetRuleTypes200ResponseInnerAlertsWithDefaults instantiates a new GetRuleTypes200ResponseInnerAlerts object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetRuleTypes200ResponseInnerAlertsWithDefaults() *GetRuleTypes200ResponseInnerAlerts {
	this := GetRuleTypes200ResponseInnerAlerts{}
	var useLegacyAlerts bool = false
	this.UseLegacyAlerts = &useLegacyAlerts
	return &this
}

// GetContext returns the Context field value if set, zero value otherwise.
func (o *GetRuleTypes200ResponseInnerAlerts) GetContext() string {
	if o == nil || IsNil(o.Context) {
		var ret string
		return ret
	}
	return *o.Context
}

// GetContextOk returns a tuple with the Context field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) GetContextOk() (*string, bool) {
	if o == nil || IsNil(o.Context) {
		return nil, false
	}
	return o.Context, true
}

// HasContext returns a boolean if a field has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) HasContext() bool {
	if o != nil && !IsNil(o.Context) {
		return true
	}

	return false
}

// SetContext gets a reference to the given string and assigns it to the Context field.
func (o *GetRuleTypes200ResponseInnerAlerts) SetContext(v string) {
	o.Context = &v
}

// GetDynamic returns the Dynamic field value if set, zero value otherwise.
func (o *GetRuleTypes200ResponseInnerAlerts) GetDynamic() string {
	if o == nil || IsNil(o.Dynamic) {
		var ret string
		return ret
	}
	return *o.Dynamic
}

// GetDynamicOk returns a tuple with the Dynamic field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) GetDynamicOk() (*string, bool) {
	if o == nil || IsNil(o.Dynamic) {
		return nil, false
	}
	return o.Dynamic, true
}

// HasDynamic returns a boolean if a field has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) HasDynamic() bool {
	if o != nil && !IsNil(o.Dynamic) {
		return true
	}

	return false
}

// SetDynamic gets a reference to the given string and assigns it to the Dynamic field.
func (o *GetRuleTypes200ResponseInnerAlerts) SetDynamic(v string) {
	o.Dynamic = &v
}

// GetIsSpaceAware returns the IsSpaceAware field value if set, zero value otherwise.
func (o *GetRuleTypes200ResponseInnerAlerts) GetIsSpaceAware() bool {
	if o == nil || IsNil(o.IsSpaceAware) {
		var ret bool
		return ret
	}
	return *o.IsSpaceAware
}

// GetIsSpaceAwareOk returns a tuple with the IsSpaceAware field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) GetIsSpaceAwareOk() (*bool, bool) {
	if o == nil || IsNil(o.IsSpaceAware) {
		return nil, false
	}
	return o.IsSpaceAware, true
}

// HasIsSpaceAware returns a boolean if a field has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) HasIsSpaceAware() bool {
	if o != nil && !IsNil(o.IsSpaceAware) {
		return true
	}

	return false
}

// SetIsSpaceAware gets a reference to the given bool and assigns it to the IsSpaceAware field.
func (o *GetRuleTypes200ResponseInnerAlerts) SetIsSpaceAware(v bool) {
	o.IsSpaceAware = &v
}

// GetMappings returns the Mappings field value if set, zero value otherwise.
func (o *GetRuleTypes200ResponseInnerAlerts) GetMappings() GetRuleTypes200ResponseInnerAlertsMappings {
	if o == nil || IsNil(o.Mappings) {
		var ret GetRuleTypes200ResponseInnerAlertsMappings
		return ret
	}
	return *o.Mappings
}

// GetMappingsOk returns a tuple with the Mappings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) GetMappingsOk() (*GetRuleTypes200ResponseInnerAlertsMappings, bool) {
	if o == nil || IsNil(o.Mappings) {
		return nil, false
	}
	return o.Mappings, true
}

// HasMappings returns a boolean if a field has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) HasMappings() bool {
	if o != nil && !IsNil(o.Mappings) {
		return true
	}

	return false
}

// SetMappings gets a reference to the given GetRuleTypes200ResponseInnerAlertsMappings and assigns it to the Mappings field.
func (o *GetRuleTypes200ResponseInnerAlerts) SetMappings(v GetRuleTypes200ResponseInnerAlertsMappings) {
	o.Mappings = &v
}

// GetSecondaryAlias returns the SecondaryAlias field value if set, zero value otherwise.
func (o *GetRuleTypes200ResponseInnerAlerts) GetSecondaryAlias() string {
	if o == nil || IsNil(o.SecondaryAlias) {
		var ret string
		return ret
	}
	return *o.SecondaryAlias
}

// GetSecondaryAliasOk returns a tuple with the SecondaryAlias field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) GetSecondaryAliasOk() (*string, bool) {
	if o == nil || IsNil(o.SecondaryAlias) {
		return nil, false
	}
	return o.SecondaryAlias, true
}

// HasSecondaryAlias returns a boolean if a field has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) HasSecondaryAlias() bool {
	if o != nil && !IsNil(o.SecondaryAlias) {
		return true
	}

	return false
}

// SetSecondaryAlias gets a reference to the given string and assigns it to the SecondaryAlias field.
func (o *GetRuleTypes200ResponseInnerAlerts) SetSecondaryAlias(v string) {
	o.SecondaryAlias = &v
}

// GetShouldWrite returns the ShouldWrite field value if set, zero value otherwise.
func (o *GetRuleTypes200ResponseInnerAlerts) GetShouldWrite() bool {
	if o == nil || IsNil(o.ShouldWrite) {
		var ret bool
		return ret
	}
	return *o.ShouldWrite
}

// GetShouldWriteOk returns a tuple with the ShouldWrite field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) GetShouldWriteOk() (*bool, bool) {
	if o == nil || IsNil(o.ShouldWrite) {
		return nil, false
	}
	return o.ShouldWrite, true
}

// HasShouldWrite returns a boolean if a field has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) HasShouldWrite() bool {
	if o != nil && !IsNil(o.ShouldWrite) {
		return true
	}

	return false
}

// SetShouldWrite gets a reference to the given bool and assigns it to the ShouldWrite field.
func (o *GetRuleTypes200ResponseInnerAlerts) SetShouldWrite(v bool) {
	o.ShouldWrite = &v
}

// GetUseEcs returns the UseEcs field value if set, zero value otherwise.
func (o *GetRuleTypes200ResponseInnerAlerts) GetUseEcs() bool {
	if o == nil || IsNil(o.UseEcs) {
		var ret bool
		return ret
	}
	return *o.UseEcs
}

// GetUseEcsOk returns a tuple with the UseEcs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) GetUseEcsOk() (*bool, bool) {
	if o == nil || IsNil(o.UseEcs) {
		return nil, false
	}
	return o.UseEcs, true
}

// HasUseEcs returns a boolean if a field has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) HasUseEcs() bool {
	if o != nil && !IsNil(o.UseEcs) {
		return true
	}

	return false
}

// SetUseEcs gets a reference to the given bool and assigns it to the UseEcs field.
func (o *GetRuleTypes200ResponseInnerAlerts) SetUseEcs(v bool) {
	o.UseEcs = &v
}

// GetUseLegacyAlerts returns the UseLegacyAlerts field value if set, zero value otherwise.
func (o *GetRuleTypes200ResponseInnerAlerts) GetUseLegacyAlerts() bool {
	if o == nil || IsNil(o.UseLegacyAlerts) {
		var ret bool
		return ret
	}
	return *o.UseLegacyAlerts
}

// GetUseLegacyAlertsOk returns a tuple with the UseLegacyAlerts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) GetUseLegacyAlertsOk() (*bool, bool) {
	if o == nil || IsNil(o.UseLegacyAlerts) {
		return nil, false
	}
	return o.UseLegacyAlerts, true
}

// HasUseLegacyAlerts returns a boolean if a field has been set.
func (o *GetRuleTypes200ResponseInnerAlerts) HasUseLegacyAlerts() bool {
	if o != nil && !IsNil(o.UseLegacyAlerts) {
		return true
	}

	return false
}

// SetUseLegacyAlerts gets a reference to the given bool and assigns it to the UseLegacyAlerts field.
func (o *GetRuleTypes200ResponseInnerAlerts) SetUseLegacyAlerts(v bool) {
	o.UseLegacyAlerts = &v
}

func (o GetRuleTypes200ResponseInnerAlerts) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetRuleTypes200ResponseInnerAlerts) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Context) {
		toSerialize["context"] = o.Context
	}
	if !IsNil(o.Dynamic) {
		toSerialize["dynamic"] = o.Dynamic
	}
	if !IsNil(o.IsSpaceAware) {
		toSerialize["isSpaceAware"] = o.IsSpaceAware
	}
	if !IsNil(o.Mappings) {
		toSerialize["mappings"] = o.Mappings
	}
	if !IsNil(o.SecondaryAlias) {
		toSerialize["secondaryAlias"] = o.SecondaryAlias
	}
	if !IsNil(o.ShouldWrite) {
		toSerialize["shouldWrite"] = o.ShouldWrite
	}
	if !IsNil(o.UseEcs) {
		toSerialize["useEcs"] = o.UseEcs
	}
	if !IsNil(o.UseLegacyAlerts) {
		toSerialize["useLegacyAlerts"] = o.UseLegacyAlerts
	}
	return toSerialize, nil
}

type NullableGetRuleTypes200ResponseInnerAlerts struct {
	value *GetRuleTypes200ResponseInnerAlerts
	isSet bool
}

func (v NullableGetRuleTypes200ResponseInnerAlerts) Get() *GetRuleTypes200ResponseInnerAlerts {
	return v.value
}

func (v *NullableGetRuleTypes200ResponseInnerAlerts) Set(val *GetRuleTypes200ResponseInnerAlerts) {
	v.value = val
	v.isSet = true
}

func (v NullableGetRuleTypes200ResponseInnerAlerts) IsSet() bool {
	return v.isSet
}

func (v *NullableGetRuleTypes200ResponseInnerAlerts) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetRuleTypes200ResponseInnerAlerts(val *GetRuleTypes200ResponseInnerAlerts) *NullableGetRuleTypes200ResponseInnerAlerts {
	return &NullableGetRuleTypes200ResponseInnerAlerts{value: val, isSet: true}
}

func (v NullableGetRuleTypes200ResponseInnerAlerts) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetRuleTypes200ResponseInnerAlerts) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}