package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func IntToInt64Type(value int) types.Int64 {
	return types.Int64Value(int64(value))
}

func IntPointerToInt64Type(value *int) types.Int64 {
	if value == nil {
		return types.Int64Null()
	} else {
		return types.Int64Value(int64(*value))
	}
}

func SliceToListType[T any](
	value []T,
	elemType attr.Type,
	path path.Path,
	diags diag.Diagnostics,
	iteratee func(item T, index int) attr.Value,
) types.List {
	if value == nil {
		return types.ListNull(elemType)
	}

	elems := make([]attr.Value, 0, len(value))
	for i, v := range value {
		elems = append(elems, iteratee(v, i))
	}

	val, d := types.ListValue(elemType, elems)
	diags.Append(ConvertToAttrDiags(d, path)...)
	return val
}

func SliceToListType_String(
	value []string,
	path path.Path,
	diags diag.Diagnostics,
) types.List {
	return SliceToListType(value, types.StringType, path, diags, func(item string, index int) attr.Value {
		return types.StringValue(item)
	})
}

func SliceToSetType[T any](
	value []T,
	elemType attr.Type,
	path path.Path,
	diags diag.Diagnostics,
	iteratee func(item T, index int) attr.Value,
) types.Set {
	if value == nil {
		return types.SetNull(elemType)
	}

	elems := make([]attr.Value, 0, len(value))
	for i, v := range value {
		elems = append(elems, iteratee(v, i))
	}

	val, d := types.SetValue(elemType, elems)
	diags.Append(ConvertToAttrDiags(d, path)...)
	return val
}

func SliceToSetType_String(
	value []string,
	path path.Path,
	diags diag.Diagnostics,
) types.Set {
	return SliceToSetType(value, types.StringType, path, diags, func(item string, index int) attr.Value {
		return types.StringValue(item)
	})
}

func MapToMapType[T any](
	value map[string]T,
	elemType attr.Type,
	path path.Path,
	diags diag.Diagnostics,
	iteratee func(key string, value T) attr.Value,
) types.Map {
	if value == nil {
		return types.MapNull(elemType)
	}

	elems := make(map[string]attr.Value, len(value))
	for k, v := range value {
		elems[k] = iteratee(k, v)
	}

	val, d := types.MapValue(elemType, elems)
	diags.Append(ConvertToAttrDiags(d, path)...)
	return val
}

func MapToObjectType(
	attrTypes map[string]attr.Type,
	value map[string]attr.Value,
	path path.Path,
	diags diag.Diagnostics,
) types.Object {
	if value == nil {
		return types.ObjectNull(attrTypes)
	}

	val, d := types.ObjectValue(attrTypes, value)
	diags.Append(ConvertToAttrDiags(d, path)...)

	return val
}

func MapToNormalizedType[T any](
	value map[string]T,
	path path.Path,
	diags diag.Diagnostics,
) jsontypes.Normalized {
	if value == nil {
		return jsontypes.NewNormalizedNull()
	}

	val, err := json.Marshal(value)
	if err != nil {
		diags.AddAttributeError(path, "marshal failure", err.Error())
	}

	return jsontypes.NewNormalizedValue(string(val))
}

func StructToNormalizedType[T any](
	value *T,
	path path.Path,
	diags diag.Diagnostics,
) jsontypes.Normalized {
	if value == nil {
		return jsontypes.NewNormalizedNull()
	}

	val, err := json.Marshal(value)
	if err != nil {
		diags.AddAttributeError(path, "marshal failure", err.Error())
	}

	return jsontypes.NewNormalizedValue(string(val))
}

// ============================================================================

func Int64TypeToIntPointer(value types.Int64) *int {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := int(value.ValueInt64())
	return &v
}

func ObjectTypeToStruct[T1 any, T2 any](
	ctx context.Context,
	value types.Object,
	path path.Path,
	diags diag.Diagnostics,
	transformer func(item T1) T2,
) *T2 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var dest T1
	opts := basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	}
	d := value.As(ctx, &dest, opts)
	diags.Append(ConvertToAttrDiags(d, path)...)
	out := transformer(dest)
	return &out
}

func ListTypeToSliceBasic[T any](
	ctx context.Context,
	value types.List,
	path path.Path,
	diags diag.Diagnostics,
) []T {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var items []T
	d := value.ElementsAs(ctx, &items, false)
	diags.Append(ConvertToAttrDiags(d, path)...)
	return items
}

func ListTypeToSlice[T1 any, T2 any](
	ctx context.Context,
	value types.List,
	path path.Path,
	diags diag.Diagnostics,
	iteratee func(value T1, index int) T2,
) []T2 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	elems := ListTypeToSliceBasic[T1](ctx, value, path, diags)
	if diags.HasError() {
		return nil
	}

	items := make([]T2, 0, len(elems))
	for i, v := range elems {
		items = append(items, iteratee(v, i))
	}

	return items
}

func SetTypeToSliceBasic[T any](
	ctx context.Context,
	value types.Set,
	path path.Path,
	diags diag.Diagnostics,
) []T {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var items []T
	d := value.ElementsAs(ctx, &items, false)
	diags.Append(ConvertToAttrDiags(d, path)...)
	return items
}

func SetTypeToSlice[T1 any, T2 any](
	ctx context.Context,
	value types.Set,
	path path.Path,
	diags diag.Diagnostics,
	iteratee func(value T1, index int) T2,
) []T2 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	elems := SetTypeToSliceBasic[T1](ctx, value, path, diags)
	items := make([]T2, 0, len(elems))
	for i, v := range elems {
		items = append(items, iteratee(v, i))
	}

	return items
}

func MapTypeToMap[T1 any, T2 any](
	ctx context.Context,
	value types.Map,
	path path.Path,
	diags diag.Diagnostics,
	iteratee func(key string, value T1) T2,
) map[string]T2 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var items map[string]T2
	d := value.ElementsAs(ctx, &items, false)
	diags.Append(ConvertToAttrDiags(d, path)...)

	return items
}

func NormalizedTypeToMap[T any](
	value jsontypes.Normalized,
	path path.Path,
	diags diag.Diagnostics,
) map[string]T {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	dest, err := JsonUnmarshalS[map[string]T](value.ValueString())
	if err != nil {
		diags.AddAttributeError(path, "normalized unmarshal failure", err.Error())
	}
	return dest
}

func NormalizedTypeToStruct[T any](
	value jsontypes.Normalized,
	path path.Path,
	diags diag.Diagnostics,
) *T {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	dest, err := JsonUnmarshalS[T](value.ValueString())
	if err != nil {
		diags.AddAttributeError(path, "normalized unmarshal failure", err.Error())
	}
	return &dest
}

// ============================================================================

func TransformMap[T1 any, T2 any](
	value map[string]T1,
	iteratee func(key string, value T1) T2,
) map[string]T2 {
	if value == nil {
		return nil
	}

	elems := make(map[string]T2, len(value))
	for k, v := range value {
		elems[k] = iteratee(k, v)
	}

	return elems
}

func TransformSlice[T1 any, T2 any](
	value []T1,
	iteratee func(item T1, index int) T2,
) []T2 {
	if value == nil {
		return nil
	}

	elems := make([]T2, 0, len(value))
	for i, v := range value {
		elems = append(elems, iteratee(v, i))
	}

	return elems
}

func TransformStruct[T1 any, T2 any](
	value *T1,
	transformer func(item T1) T2,
) *T2 {
	if value == nil {
		return nil
	}

	res := transformer(*value)
	return &res
}

// ============================================================================

func Int64TypeToString(val types.Int64) *string {
	v := val.ValueInt64Pointer()
	if v == nil {
		return nil
	}
	i := strconv.FormatInt(*v, 10)
	return &i
}

func StringToInt64Type(val *string, path path.Path, diags diag.Diagnostics) types.Int64 {
	if val == nil {
		return types.Int64Null()
	}
	i, err := strconv.ParseInt(*val, 10, 0)
	if err != nil {
		diags.AddAttributeError(path, "atoi failure", err.Error())
		return types.Int64Null()
	}
	return types.Int64Value(i)
}

func DurationToInt64Type(val estypes.Duration, path path.Path, diags diag.Diagnostics) types.Int64 {
	if val == nil {
		return types.Int64Null()
	}
	if v, ok := val.(string); ok {
		return StringToInt64Type(&v, path, diags)
	} else {
		diags.AddAttributeError(path, "not a string", fmt.Sprintf("actual type: %T", val))
	}
	return types.Int64Unknown()
}

func Int64TypeToDuration(val types.Int64) estypes.Duration {
	v := val.ValueInt64Pointer()
	if v == nil {
		return estypes.Duration(nil)
	}
	out := strconv.FormatInt(*v, 10)
	return estypes.Duration(out)
}

func Ptr[T any](val T) *T {
	return &val
}
