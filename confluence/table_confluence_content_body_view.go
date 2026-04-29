package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableConfluenceContentBodyView() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_body_view",
		Description: "Confluence Content Body in the View Format.",
		List: &plugin.ListConfig{
			ParentHydrate: listContentForBodyView,
			Hydrate:       listContentBodyView,
			KeyColumns:    plugin.OptionalColumns([]string{"id"}),
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

	content := h.Item.(*model.ContentScheme)
	if content == nil || content.Body == nil || content.Body.View == nil {
		return nil, nil
	}
	row := contentBody{
		ID:             content.ID,
		Representation: content.Body.View.Representation,
		Value:          content.Body.View.Value,
	}
	d.StreamListItem(ctx, row)

	return nil, nil
}

func listContentForBodyView(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listContentWithExpand(ctx, d, []string{"body.view"})
}
