package confluence

import (
	"context"

	"github.com/ctreminiom/go-atlassian/confluence"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableConfluenceContentBodyView() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_body_view",
		Description: "Confluence Content Body in the View Format.",
		List: &plugin.ListConfig{
			ParentHydrate: listContent,
			Hydrate:       listContentBodyView,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the content.",
			},
			{
				Name:        "representation",
				Type:        proto.ColumnType_STRING,
				Description: "The representation type of the content.",
			},
			{
				Name:        "value",
				Type:        proto.ColumnType_STRING,
				Description: "The content body.",
			},
		},
	}
}

//// LIST FUNCTIONS

func listContentBodyView(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContentBody")

	content := h.Item.(*confluence.ContentScheme)
	c := contentBody{
		ID:             content.ID,
		Representation: content.Body.View.Representation,
		Value:          content.Body.View.Value,
	}
	d.StreamListItem(ctx, c)

	return nil, nil
}
