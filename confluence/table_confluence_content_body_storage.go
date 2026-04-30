package confluence

import (
	"context"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableConfluenceContentBodyStorage() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_body_storage",
		Description: "Confluence Content Body in the Storage Format.",
		List: &plugin.ListConfig{
			ParentHydrate: listContentForBodyStorage,
			Hydrate:       listContentBodyStorage,
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

func listContentBodyStorage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContentBody")

	content := h.Item.(*models.ContentScheme)
	if content == nil || content.Body == nil || content.Body.Storage == nil {
		return nil, nil
	}
	row := contentBody{
		ID:             content.ID,
		Representation: content.Body.Storage.Representation,
		Value:          content.Body.Storage.Value,
	}
	d.StreamListItem(ctx, row)

	return nil, nil
}

func listContentForBodyStorage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listContentWithExpand(ctx, d, []string{"body.storage"})
}
