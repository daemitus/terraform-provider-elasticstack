package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// An existing null map should should be semantically equal to an empty map.
func SemanticEqualEmptyMap(ctx context.Context, existing types.Map, incoming types.Map) types.Map {
	if !IsKnown(existing) && len(incoming.Elements()) == 0 {
		return types.MapNull(incoming.ElementType(ctx))
	}
	return incoming
}

// An existing null list should be semantically equal to an empty list.
func SemanticEqualEmptyList(ctx context.Context, existing types.List, incoming types.List) types.List {
	if !IsKnown(existing) && len(incoming.Elements()) == 0 {
		return types.ListNull(incoming.ElementType(ctx))
	}
	return incoming
}

// An existing null set should be semantically equal to an empty set.
func SemanticEqualEmptySet(ctx context.Context, existing types.Set, incoming types.Set) types.Set {
	if !IsKnown(existing) && len(incoming.Elements()) == 0 {
		return types.SetNull(incoming.ElementType(ctx))
	}
	return incoming
}
