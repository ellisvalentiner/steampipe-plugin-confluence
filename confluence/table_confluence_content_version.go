package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableConfluenceContentVersion() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_version",
		Description: "Confluence Content Version.",
		List: &plugin.ListConfig{
			ParentHydrate: listContent,
			Hydrate:       listContentVersion,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the content.",
			},
			{
				Name:        "by",
				Type:        proto.ColumnType_STRING,
				Description: "The representation type of the content.",
			},
			{
				Name:        "number",
				Type:        proto.ColumnType_INT,
				Description: "The content body.",
			},
			{
				Name:        "when",
				Type:        proto.ColumnType_STRING,
				Description: "The content body.",
			},
			{
				Name:        "message",
				Type:        proto.ColumnType_STRING,
				Description: "The content body.",
			},
			{
				Name:        "minor_edit",
				Type:        proto.ColumnType_BOOL,
				Description: "The content body.",
			},
		},
	}
}

type contentVersion struct {
	ID        string                   `json:"id,omitempty"`
	By        *model.ContentUserScheme `json:"by,omitempty"`
	Number    int                      `json:"number,omitempty"`
	When      string                   `json:"when,omitempty"`
	Message   string                   `json:"message,omitempty"`
	MinorEdit bool                     `json:"minorEdit,omitempty"`
}

//// LIST FUNCTIONS

func listContentVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContentVersion")

	content := h.Item.(*model.ContentScheme)
	c := contentVersion{
		ID:        content.ID,
		By:        content.Version.By,
		Number:    content.Version.Number,
		When:      content.Version.When,
		Message:   content.Version.Message,
		MinorEdit: content.Version.MinorEdit,
	}
	d.StreamListItem(ctx, c)

	return nil, nil
}
