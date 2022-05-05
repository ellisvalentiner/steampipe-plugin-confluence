package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableConfluenceContentBodyStorage() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_body_storage",
		Description: "Confluence Content Body in the Storage Format.",
		List: &plugin.ListConfig{
			ParentHydrate: listContent,
			Hydrate:       listContentBodyStorage,
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

// structs
type contentBody struct {
	ID             string `json:"id,omitempty"`
	Representation string `json:"representation,omitempty"`
	Value          string `json:"value,omitempty"`
}

//// LIST FUNCTIONS

func listContentBodyStorage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContentBody")

	content := h.Item.(*model.ContentScheme)
	row := contentBody{
		ID:             content.ID,
		Representation: content.Body.Storage.Representation,
		Value:          content.Body.Storage.Value,
	}
	d.StreamListItem(ctx, row)

	return nil, nil
}
