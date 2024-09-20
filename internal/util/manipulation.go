package util

import (
	"strings"
)

// Pointer returns a pointer to a variable holding the supplied T constant
func Pointer[T any](x T) *T {
	return &x
}

// Flatten takes a hierarchy and flattens it using a dot "."
func Flatten(m map[string]any) map[string]any {
	var out = make(map[string]any)

	var flattenRecursive func(m map[string]any, keys []string)
	flattenRecursive = func(m map[string]any, keys []string) {
		for key, val := range m {
			newKeys := append(keys, key)
			if newMap, ok := val.(map[string]any); ok {
				flattenRecursive(newMap, newKeys)
			} else {
				key := strings.Join(newKeys, ".")
				out[key] = val
			}
		}
	}

	flattenRecursive(m, []string{})
	return out
}

// Unflatten will unflatten a map with keys which are comprised of multiple tokens segmented by dots "."
func Unflatten(m map[string]any) map[string]any {
	var out = make(map[string]any)
	for path, val := range m {
		keys := strings.Split(path, ".")

		tree := out
		for i, key := range keys {
			if i == len(keys)-1 {
				tree[key] = val
			} else {
				leaf, ok := tree[key]
				if !ok {
					leaf = make(map[string]any)
					tree[key] = leaf
				}
				tree = leaf.(map[string]any)
			}
		}
	}
	return out
}

func StructToMap(val any) (map[string]any, error) {
	bytes, err := JsonMarshal(val)
	if err != nil {
		return nil, err
	}
	m, err := JsonUnmarshal[map[string]any](bytes)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func MapToStruct[T any](val map[string]any) (*T, error) {
	bytes, err := JsonMarshal(val)
	if err != nil {
		return nil, err
	}
	s, err := JsonUnmarshal[T](bytes)
	if err != nil {
		return nil, err
	}
	return &s, err
}
