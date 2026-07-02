package provider

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func stringSetValues(ctx context.Context, set types.Set, diagnostics *diag.Diagnostics) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	var values []types.String
	diagnostics.Append(set.ElementsAs(ctx, &values, false)...)
	if diagnostics.HasError() {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if !value.IsNull() && !value.IsUnknown() {
			result = append(result, value.ValueString())
		}
	}
	sort.Strings(result)
	return result
}

func setFromStrings(ctx context.Context, values []string, diagnostics *diag.Diagnostics) types.Set {
	sort.Strings(values)
	result, diags := types.SetValueFrom(ctx, types.StringType, values)
	diagnostics.Append(diags...)
	return result
}
