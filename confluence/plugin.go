package confluence

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-confluence",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"confluence_content":              tableConfluenceContent(),
			"confluence_content_body_storage": tableConfluenceContentBodyStorage(),
			"confluence_content_body_view":    tableConfluenceContentBodyView(),
			"confluence_space":                tableConfluenceSpace(),
		},
	}
	return p
}
