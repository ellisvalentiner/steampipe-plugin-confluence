package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION
func tableConfluenceContentMetadata() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_metadata",
		Description: "Information about the current user in relation to the content, including when they last viewed it, modified it, contributed to it, or added it as a favorite.",
		List: &plugin.ListConfig{
			ParentHydrate: listContent,
			Hydrate:       listContentMetadata,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the content.",
			},
			{
				Name:        "metadata",
				Type:        proto.ColumnType_JSON,
				Description: "The representation type of the content.",
			},
		},
	}
}

// structs
type metadata struct {
	ID       string                `json:"id,omitempty"`
	Metadata *model.MetadataScheme `json:"metadata,omitempty"`
}

//// LIST FUNCTIONS
func listContentMetadata(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContentMetadata")

	content := h.Item.(*model.ContentScheme)
	c := metadata{
		ID:       content.ID,
		Metadata: content.Metadata,
	}
	d.StreamListItem(ctx, c)

	return nil, nil
}
