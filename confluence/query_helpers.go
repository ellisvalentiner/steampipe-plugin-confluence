package confluence

import (
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func getStringQual(d *plugin.QueryData, name string) string {
	if d == nil {
		return ""
	}

	qualValue, ok := d.EqualsQuals[name]
	if !ok || qualValue == nil {
		return ""
	}

	value := grpc.GetQualValue(qualValue)
	if str, ok := value.(string); ok {
		return str
	}

	return fmt.Sprintf("%v", value)
}

func getQueryLimit(d *plugin.QueryData) int {
	if d == nil || d.QueryContext == nil || d.QueryContext.Limit == nil {
		return 0
	}

	limit := int(*d.QueryContext.Limit)
	if limit < 0 {
		return 0
	}

	return limit
}

func requestPageSize(maxItems, streamed, defaultPageSize int) int {
	if maxItems <= 0 {
		return defaultPageSize
	}

	remaining := maxItems - streamed
	if remaining <= 0 {
		return 0
	}
	if remaining < defaultPageSize {
		return remaining
	}

	return defaultPageSize
}
